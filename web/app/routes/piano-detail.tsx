import { useCallback, useEffect, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router";
import {
  ArrowLeft,
  Bookmark,
  Building2,
  CalendarClock,
  CalendarDays,
  Clock,
  Heart,
  MapPin,
  MapPinned,
  MessageSquarePlus,
  Music,
  Pencil,
  Star,
  Tag,
  UserRound,
} from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { Button } from "../components/ui/button";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { CreatePianoPostModal } from "../components/CreatePianoPostModal";
import { EditPianoModal } from "../components/EditPianoModal";
import { MessageDialog } from "../components/MessageDialog";
import { PianoAttributeMeters } from "../components/PianoAttributeMeters";
import { PianoEditList } from "../components/PianoEditList";
import { PianoPostCard } from "../components/PianoPostCard";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { useAuth } from "../lib/auth";
import { useMe } from "../lib/use-me";
import { pianoClient, pianoPostClient, pianoUserListClient } from "../lib/api-client";
import { GetPianoRequest, Piano } from "../lib/gen/piano/v1/piano_pb";
import {
  DeletePianoPostRequest,
  ListPianoPostsRequest,
  type PianoPost,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import {
  AddPianoToUserListRequest,
  GetMyPianoUserListsRequest,
  PianoListKind,
  RemovePianoFromUserListRequest,
} from "../lib/gen/piano_user_list/v1/piano_user_list_pb";
import { formatPiano } from "../lib/resource-name";

import type { Route } from "./+types/piano-detail";

export function meta({}: Route.MetaArgs) {
  return [{ title: "ピアノ — PiaMap" }];
}

export default function PianoDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const { authed } = useAuth();
  const me = useMe();
  const [piano, setPiano] = useState<Piano | null>(null);
  const [posts, setPosts] = useState<PianoPost[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [reviewOpen, setReviewOpen] = useState(false);
  const [signUpOpen, setSignUpOpen] = useState(false);
  const [signUpAction, setSignUpAction] = useState<string | undefined>();
  const [editingPost, setEditingPost] = useState<PianoPost | null>(null);
  const [deletingPost, setDeletingPost] = useState<PianoPost | null>(null);
  const [editPianoOpen, setEditPianoOpen] = useState(false);
  const [savedMessage, setSavedMessage] = useState<string | null>(null);
  const [editsRefreshKey, setEditsRefreshKey] = useState(0);
  const [myListKinds, setMyListKinds] = useState<Set<PianoListKind>>(new Set());
  const [listKindBusy, setListKindBusy] = useState<Set<PianoListKind>>(new Set());

  const reload = useCallback(() => {
    if (!id) return;
    setErr(null);
    pianoClient
      .getPiano(new GetPianoRequest({ name: formatPiano(id) }))
      .then((res) => setPiano(res.piano ?? null))
      .catch((e) => setErr(e?.message || String(e)));
    pianoPostClient
      .listPianoPosts(new ListPianoPostsRequest({ parent: formatPiano(id), pageSize: 20 }))
      .then((res) => setPosts(res.pianoPosts))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  }, [id]);

  useEffect(() => {
    reload();
  }, [reload]);

  // 自分のリスト所属状態をフェッチ。
  useEffect(() => {
    if (!authed || !id) {
      setMyListKinds(new Set());
      return;
    }
    let cancelled = false;
    pianoUserListClient
      .getMyPianoUserLists(new GetMyPianoUserListsRequest({ parent: formatPiano(id) }))
      .then((res) => {
        if (!cancelled) setMyListKinds(new Set(res.listKinds));
      })
      .catch(() => {
        if (!cancelled) setMyListKinds(new Set());
      });
    return () => {
      cancelled = true;
    };
  }, [authed, id]);

  // 他ページ (/map?relocate=*) から saved メッセージを引き継ぐ。
  useEffect(() => {
    const msg = (location.state as { savedMessage?: string } | null)?.savedMessage;
    if (msg) {
      setSavedMessage(msg);
      setEditsRefreshKey((k) => k + 1);
      navigate(location.pathname, { replace: true, state: null });
    }
  }, [location, navigate]);

  const handleReviewPress = () => {
    if (!authed) {
      setSignUpAction("レビュー投稿");
      setSignUpOpen(true);
      return;
    }
    setEditingPost(null);
    setReviewOpen(true);
  };

  const handleEditPiano = () => {
    if (!authed) {
      setSignUpAction("ピアノ情報の編集");
      setSignUpOpen(true);
      return;
    }
    setEditPianoOpen(true);
  };

  const handleRelocate = () => {
    if (!authed) {
      setSignUpAction("ピアノの位置変更");
      setSignUpOpen(true);
      return;
    }
    if (id) navigate(`/map?relocate=${id}`);
  };

  const toggleListKind = async (kind: PianoListKind, actionLabel: string) => {
    if (!authed) {
      setSignUpAction(actionLabel);
      setSignUpOpen(true);
      return;
    }
    if (!id) return;
    if (listKindBusy.has(kind)) return;
    setListKindBusy((s) => new Set(s).add(kind));
    const willAdd = !myListKinds.has(kind);
    // 楽観的更新
    setMyListKinds((prev) => {
      const next = new Set(prev);
      if (willAdd) next.add(kind);
      else next.delete(kind);
      return next;
    });
    try {
      if (willAdd) {
        await pianoUserListClient.addPianoToUserList(
          new AddPianoToUserListRequest({ parent: formatPiano(id), listKind: kind }),
        );
      } else {
        await pianoUserListClient.removePianoFromUserList(
          new RemovePianoFromUserListRequest({ parent: formatPiano(id), listKind: kind }),
        );
      }
      reload();
    } catch (e) {
      // ロールバック
      setMyListKinds((prev) => {
        const next = new Set(prev);
        if (willAdd) next.delete(kind);
        else next.add(kind);
        return next;
      });
      setErr((e as Error)?.message || String(e));
    } finally {
      setListKindBusy((s) => {
        const next = new Set(s);
        next.delete(kind);
        return next;
      });
    }
  };

  const handleEdit = (post: PianoPost) => {
    setEditingPost(post);
    setReviewOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!deletingPost) return;
    try {
      await pianoPostClient.deletePianoPost(
        new DeletePianoPostRequest({ name: deletingPost.name }),
      );
      setDeletingPost(null);
      reload();
    } catch (e) {
      setErr((e as Error)?.message || String(e));
      setDeletingPost(null);
    }
  };

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-3">
        <Link to="/map" aria-label="戻る" className="text-slate-500 hover:text-slate-700">
          <ArrowLeft size={22} />
        </Link>
        <h1 className="flex-1 truncate text-base font-bold text-slate-900">
          {piano?.displayName ?? "読み込み中..."}
        </h1>
        {piano ? (
          <div className="flex items-center gap-1">
            <button
              type="button"
              aria-label="位置を変更"
              onClick={handleRelocate}
              className="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-full text-slate-500 hover:bg-slate-100 hover:text-slate-700 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
            >
              <MapPinned size={18} />
            </button>
            <button
              type="button"
              aria-label="編集"
              onClick={handleEditPiano}
              className="inline-flex h-9 w-9 cursor-pointer items-center justify-center rounded-full text-slate-500 hover:bg-slate-100 hover:text-slate-700 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
            >
              <Pencil size={18} />
            </button>
          </div>
        ) : null}
      </header>

      {err ? <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p> : null}

      {piano ? (
        <section className="space-y-4">
          <RatingBlock piano={piano} />

          <div className="flex gap-2">
            <ListToggleButton
              icon={<Bookmark size={16} />}
              label="行ってみたい"
              active={myListKinds.has(PianoListKind.WISHLIST)}
              busy={listKindBusy.has(PianoListKind.WISHLIST)}
              onPress={() => toggleListKind(PianoListKind.WISHLIST, "行ってみたいに追加")}
              count={piano.wishlistCount}
            />
            <ListToggleButton
              icon={<Heart size={16} />}
              label="お気に入り"
              active={myListKinds.has(PianoListKind.FAVORITE)}
              busy={listKindBusy.has(PianoListKind.FAVORITE)}
              onPress={() => toggleListKind(PianoListKind.FAVORITE, "お気に入りに追加")}
              count={piano.favoriteCount}
              activeColor="rose"
            />
          </div>

          {piano.description ? (
            <p className="rounded-2xl bg-slate-50 p-4 text-sm text-slate-700 whitespace-pre-line">
              {piano.description}
            </p>
          ) : null}

          <PianoAttributeMeters piano={piano} />

          <dl className="space-y-2 rounded-2xl border border-slate-200 p-4 text-sm">
            {piano.address ? (
              <Row icon={<MapPin size={16} className="text-slate-400" />} label="住所" value={piano.address} />
            ) : null}
            {piano.venueType ? (
              <Row
                icon={<Building2 size={16} className="text-slate-400" />}
                label="会場種別"
                value={piano.venueType}
              />
            ) : null}
            <Row
              icon={<Tag size={16} className="text-slate-400" />}
              label="種別"
              value={pianoKindLabel(piano.kind)}
            />
            <Row
              icon={<Music size={16} className="text-slate-400" />}
              label="ピアノ"
              value={`${pianoTypeLabel(piano.pianoType)} / ${piano.pianoBrand}${
                piano.pianoModel ? ` ${piano.pianoModel}` : ""
              }`}
            />
            {piano.manufactureYear ? (
              <Row
                icon={<CalendarDays size={16} className="text-slate-400" />}
                label="製造年"
                value={`${piano.manufactureYear}年`}
              />
            ) : null}
            {piano.hours ? (
              <Row
                icon={<Clock size={16} className="text-slate-400" />}
                label="営業時間"
                value={piano.hours}
              />
            ) : null}
            <Row
              icon={<CalendarClock size={16} className="text-slate-400" />}
              label="営業状況"
              value={availabilityLabel(piano.availability)}
            />
            {piano.availabilityNote ? (
              <Row
                icon={<CalendarClock size={16} className="text-slate-400" />}
                label="営業ノート"
                value={
                  <span className="whitespace-pre-line">{piano.availabilityNote}</span>
                }
              />
            ) : null}
            {piano.creator ? (
              <Row
                icon={<UserRound size={16} className="text-slate-400" />}
                label="登録"
                value={
                  <Link to={`/${piano.creator}`} className="text-amber-600 hover:underline">
                    @{piano.creator.split("/").pop()}
                  </Link>
                }
              />
            ) : null}
          </dl>

          <div className="flex items-center justify-between pt-2">
            <h2 className="text-sm font-bold text-slate-900">レビュー ({posts.length})</h2>
            <Button size="sm" onPress={handleReviewPress}>
              <MessageSquarePlus size={16} aria-hidden />
              レビュー
            </Button>
          </div>
          {posts.length === 0 ? (
            <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
              まだレビューがありません。最初の1件を投稿しよう。
            </p>
          ) : (
            <ul className="space-y-3">
              {posts.map((p) => (
                <li key={p.name}>
                  <PianoPostCard
                    post={p}
                    currentUserCustomId={me?.customId}
                    canLike={authed}
                    onLikeUnauthorized={() => {
                      setSignUpAction("いいね");
                      setSignUpOpen(true);
                    }}
                    onEdit={handleEdit}
                    onDelete={(post) => setDeletingPost(post)}
                  />
                </li>
              ))}
            </ul>
          )}

          <div className="pt-4">
            <h2 className="mb-2 text-sm font-bold text-slate-900">編集履歴</h2>
            {id ? <PianoEditList pianoId={id} refreshKey={editsRefreshKey} /> : null}
          </div>
        </section>
      ) : null}

      {piano ? (
        <CreatePianoPostModal
          isOpen={reviewOpen}
          onOpenChange={(open) => {
            setReviewOpen(open);
            if (!open) setEditingPost(null);
          }}
          pianoId={id ?? ""}
          pianoName={piano.displayName}
          editingPost={editingPost ?? undefined}
          onSaved={reload}
        />
      ) : null}
      <ConfirmDialog
        isOpen={deletingPost !== null}
        onOpenChange={(open) => {
          if (!open) setDeletingPost(null);
        }}
        title="レビューを削除"
        message="この操作は取り消せません。本当に削除しますか?"
        confirmLabel="削除する"
        destructive
        onConfirm={handleConfirmDelete}
      />
      {piano ? (
        <EditPianoModal
          isOpen={editPianoOpen}
          onOpenChange={setEditPianoOpen}
          piano={piano}
          onSaved={(updated) => {
            setPiano(updated);
            setSavedMessage("ピアノ情報を更新しました");
            setEditsRefreshKey((k) => k + 1);
          }}
        />
      ) : null}
      <MessageDialog
        isOpen={savedMessage !== null}
        onOpenChange={(open) => {
          if (!open) setSavedMessage(null);
        }}
        title="保存しました"
        message={savedMessage ?? undefined}
      />
      <SignUpPromptModal
        isOpen={signUpOpen}
        onOpenChange={setSignUpOpen}
        action={signUpAction}
      />
    </MobileShell>
  );
}

function ListToggleButton({
  icon,
  label,
  active,
  busy,
  onPress,
  count,
  activeColor = "amber",
}: {
  icon: React.ReactNode;
  label: string;
  active: boolean;
  busy: boolean;
  onPress: () => void;
  count: number;
  activeColor?: "amber" | "rose";
}) {
  const base =
    "flex flex-1 cursor-pointer items-center justify-center gap-1.5 rounded-full border px-3 py-2 text-xs font-semibold transition select-none outline-none focus-visible:ring-2 focus-visible:ring-offset-1 disabled:cursor-progress disabled:opacity-60";
  const inactive = "border-slate-300 bg-white text-slate-700 hover:border-slate-400";
  const activeAmber = "border-amber-500 bg-amber-50 text-amber-700 focus-visible:ring-amber-500";
  const activeRose = "border-rose-500 bg-rose-50 text-rose-700 focus-visible:ring-rose-500";
  const cls =
    `${base} ` +
    (active ? (activeColor === "rose" ? activeRose : activeAmber) : inactive);
  return (
    <button type="button" className={cls} onClick={onPress} disabled={busy}>
      {icon}
      <span>{label}</span>
      {count > 0 ? <span className="text-slate-500">{count}</span> : null}
    </button>
  );
}

function RatingBlock({ piano }: { piano: Piano }) {
  if (piano.postCount === 0) {
    return (
      <div className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        まだレビューがありません
      </div>
    );
  }
  return (
    <div className="flex items-center gap-3 rounded-2xl bg-amber-50 p-4">
      <Star size={28} className="fill-amber-500 text-amber-500" />
      <div>
        <span className="text-2xl font-bold text-amber-700">{piano.ratingAverage.toFixed(1)}</span>
        <span className="ml-2 text-xs text-slate-600">({piano.postCount}件のレビュー)</span>
      </div>
    </div>
  );
}

function Row({ icon, label, value }: { icon: React.ReactNode; label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5">{icon}</div>
      <div className="flex-1">
        <dt className="text-xs text-slate-500">{label}</dt>
        <dd className="text-slate-800">{value}</dd>
      </div>
    </div>
  );
}

function pianoTypeLabel(t: number): string {
  switch (t) {
    case 1:
      return "グランド";
    case 2:
      return "アップライト";
    case 3:
      return "電子";
    case 4:
      return "不明";
    default:
      return "—";
  }
}

function pianoKindLabel(k: number): string {
  switch (k) {
    case 1:
      return "ストリート";
    case 2:
      return "練習室";
    case 3:
      return "その他";
    default:
      return "—";
  }
}

function availabilityLabel(a: number): string {
  switch (a) {
    case 1:
      return "通年";
    case 2:
      return "不定期";
    case 3:
      return "イベント時のみ";
    case 4:
      return "天候次第";
    default:
      return "—";
  }
}
