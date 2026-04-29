# PiaMap プロジェクトガイド

コードを読めば分かることは書かない。**意図・暗黙の前提・運用ルール**のみ。

## プロダクト

ストリートピアノを発見・記録・共有する SNS。ドメイン: piamap.com。リポジトリは GitHub の `reverie-jp` organization 配下で管理。

個人開発として始動。「近くのストリートピアノを探せる適切なツールがない」という開発者本人の不便さが起点。マップ + 投稿(演奏記録 + 評価) + ピアニスト同士の繋がりが核。

MVP は **ストリートピアノのみ**(`piano.kind = 'street'`)。Phase 2 で練習室(`practice_room`)も追加できるよう ENUM だけ用意済み。

## 技術スタック

- Frontend: React Router v7 + React Aria + Tailwind + MapLibre + Connect RPC TS client(Cloudflare Pages)
- Backend: Go 1.25 + Connect RPC + buf + sqlc + pgx/v5 + golang-migrate(Fly.io)。HTTP ルーティングは stdlib `net/http` のみ(Connect ハンドラが `http.Handler` を返すので追加フレームワーク不要)
- DB: PostgreSQL 16 + PostGIS(Neon)
- Map: Protomaps(`.pmtiles` を Cloudflare R2 + Worker で配信)— Google Maps はスケール時のコスト爆発で却下
- Storage: Cloudflare R2(画像/音声/pmtiles) + Stream(動画)
- Auth: Google ID Token 検証 → 自前 JWT。MVP は **Google のみ**(Phase 2 で Apple / LINE 追加)

**純 Go + Connect RPC** を選んだ理由は (1) Web/Flutter 両対応の型安全 API、(2) Supabase ロックイン回避、(3) Phase 2 の Flutter 版が同じ proto から自動生成された gRPC クライアントを使える(connect-go は Connect/gRPC/gRPC-Web を同時に提供)。ハイブリッド構成は却下済み。

## モジュール構造(reverie パターン)

```
proto/<domain>/v1/*.proto       # API 契約(単一の真実)
internal/
  handler/                       # Connect handler(3行、proto変換と認証取り出しのみ)
  usecase/                       # ★ ビジネスロジック + バリデーション + 認可 + Tx境界
  domain/                        # 値オブジェクト(コンストラクタで不変条件強制)
  repository/                    # sqlc 経由で DB 操作、entity のみ返す
  gateway/                       # 他モジュール公開の View 組み立て
  auth/                          # ID Token 検証
  infra/                         # R2 / Stream / Google など外部連携
  platform/                      # ulid / resourcename / xerrors / ratelimit / logger / jwt
  gen/                           # buf 生成物 (Go)
cmd/server/main.go               # DI 配線のみ
db/migrations/                   # golang-migrate
db/queries/                      # sqlc 入力 SQL
web/                             # React Router v7
```

ドメイン分割: `user / piano / piano_post / piano_comment / piano_edit / report / media`(Phase 1d で `follow / notification`)。

## バリデーションは usecase 層に集約

`protovalidate` は **採用しない**。API 層が将来変わってもバリデーションが残るよう、proto はスキーマ定義のみで意味的バリデーションは usecase の `Input.Validate()` に集約。ドメイン層では値オブジェクト(`NewLatLng` 等)のコンストラクタで不変条件を強制する。

## handler は3行

```go
func (h *Handler) GetPiano(ctx, req) (resp, err) {
    input, err := adapter.FromGetPianoRequest(ctx, req)
    if err != nil { return nil, err }
    output, err := h.getPiano.Execute(ctx, input)
    if err != nil { return nil, err }
    return adapter.ToGetPianoResponse(output), nil
}
```

認証コンテキスト抽出(`interceptor.UserIDFromContext`)は adapter 内で行い handler に漏らさない。

## 命名規約(AIP)

将来公開 API を想定して AIP 準拠:

- **AIP-142**: タイムスタンプは `_time` 末尾(`create_time`, `visit_time`)。`_at` 禁止
- **AIP-122/131**: リソースは `name`(`pianos/01HXY...`)が第一識別子。Get/Update/Delete は `name` パラメータ
- **AIP-158**: List は `page_size` / `page_token` / `next_page_token`、ULID をカーソルとして opaque 使用
- **AIP-136**: カスタム動詞は `:verb`(`/v1/pianos/{name=pianos/*}:report`)

## マイグレーション

本番リリース前は `migrations/000001_init.{up,down}.sql` を **直接編集**。追加ファイルは作らない。up/down は対称に保つ。本番後にルール変更予定。

スキーマ変更時は必ず up/down 双方を検証してからコミット:
```bash
docker compose down -v && docker compose up -d piamap-db
docker exec -i piamap-db psql -U piamap -d piamap_db < migrations/000001_init.up.sql
```

## 開発コマンド(host を汚さない)

ホストには `docker` だけあれば十分。Go / Node / sqlc / buf / migrate / air は **コンテナ内のみ**:

```bash
docker compose up -d                                  # 全サービス起動
docker compose exec -T piamap-api make migrate-up
docker compose exec -T piamap-api make sqlc
docker compose exec -T piamap-api make proto
docker compose exec -T piamap-api make dev-up         # air で auto-rebuild
docker compose exec -T piamap-api make psql
```

API サーバーは `make dev-up`(air)で起動。手動でバイナリを起動しないこと(ポート衝突 / 挙動の不一致)。

## 主要設計判断

### 投稿モデル(`piano_posts`)
performances + reviews を統合した「投稿」エンティティ。1訪問 = 1投稿、同一 user × piano に複数行可。

- **`rating` 必須(1-5)** — Google Maps 口コミ仕様。最小投稿は ★だけでも OK、`body` / メディア / 5属性は任意
- 5属性(`ambient_noise` / `foot_traffic` / `resonance` / `key_touch_weight` / `tuning_quality`)はピアノの「特徴メーター」表示用、**rating とは独立に集計**(総合評価には混ぜない)
- 平均評価: `rating_sum / NULLIF(post_count, 0)`(rating 必須なので post_count = rating 母数)
- `visibility` は `public` / `private`(将来 `unlisted` 追加可能なため ENUM)。private でも数値集計には反映する

### ピアノ登録は自動公開(approve フローなし)
個人運営で承認滞留を避けるため `status='active'` で即時公開。`pending` 値は ENUM に温存(admin が問題投稿を一時非表示する用途)。

### コミュニティ編集の多層防御(MVP 全部入り)
1. `piano_edits` で編集ログ公開(JSONB diff + summary)
2. 編集者の即時可視化(最終編集表示)
3. Revert(`operation='restore'` で記録)
4. レート制限(MVP は in-memory、Phase1d で Redis 化)
5. 信頼ライン保護: 削除 / 座標を 500m 以上動かす / 名前全文置換 は信頼ユーザーのみ
6. Watch + 通知(MVP は DB 蓄積 + クライアントポーリング、Phase1d で Redis pub/sub 化)
7. 通報3件で自動非公開(usecase 層で実装、仕様変更容易)

**信頼ユーザー判定**: `(NOW() - users.create_time) >= 7d AND (post_count + edit_count) >= 3` を Go 側で評価。

### ユーザーリスト(行ってみたい/行ったことある/お気に入り)
`piano_user_lists` 単一テーブル + `list_kind` ENUM。同一 (user, piano) でも `list_kind` が違えば共存可能(複合主キー)。

**visited は `piano_post` 作成 usecase で UPSERT**。投稿削除しても visited マークは残す(行った事実は変わらない)。

### EXIF strip 必須(プライバシー)
写真アップロード時の処理:
1. EXIF を読み取り(訪問日時・GPS は **自動入力の参考に活用**)
2. ストレージ保存前に **完全削除**
3. 保存後の写真には EXIF が一切残らない

### アバター画像
初回サインインで Google プロフィール画像 → R2 にコピーして安定 URL 化(`avatar_url`)。設定画面で R2 への独自アップロードで上書き可能。

### handle(@username)
`users.handle VARCHAR(20) UNIQUE CHECK (handle ~ '^[a-z0-9_]{3,20}$')`。URL 識別子・将来のメンション対象。`handle_change_time` を別途持つ。

### マップタイル: Protomaps + R2
日本の OSM データを `planetiler` で `.pmtiles` 化して R2 配信、Cloudflare Worker が z/x/y タイル URL に変換。スケール時もコストはストレージ代(月 $0.075)のみで、Google Maps API のような従量課金で死なない。

### 著作権 / JASRAC
投稿動画の 99% はカバー曲。MVP は利用規約 + 通報削除フローで対応、ユーザー増加後に JASRAC 包括契約検討。

## MVP スコープ

**含める**: マップ + ピアノ情報 + 投稿(rating + メディア) + 検索 + プロフィール + 編集ログ + Watch/通知蓄積 + 3 リスト

**含めない**: フォロー / ブロック / チャット / メンション / フィード / オーナーアカウント / リアルタイム通知配信 / 訪問報告 / 多言語データ層

## リリース前タスク(未着手、要決定)

- 利用規約・プライバシーポリシー草案
- 監視・アラート(Sentry + UptimeRobot 想定、未設定)
- DB バックアップ戦略(Neon Free → Pro 検討、R2 別バケット同期検討)

## 参照

- `references/reverie-project/CLAUDE.md` — モジュール構造・命名規約・パターンの第一参考
- `references/reverie-project/migrations/000001_init.up.sql` — スキーマ書き方の元ネタ

## コメントポリシー

原則書かない。コードで意図が伝わる名前を選ぶ。非自明な why が必要な時だけ短く。
