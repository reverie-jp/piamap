import type { ReactNode } from "react";

type Props = {
  children: ReactNode;
  /** ボトムナビ分の余白を確保する (4rem 高)。flush 時はページ側で h 計算する想定。 */
  withBottomNav?: boolean;
  /** ヘッダ非表示の全画面ページ (マップ等)。padding を消す。 */
  flush?: boolean;
};

// アプリ全体の最大幅を 390px に固定 (iPhone 12 ポートレート)。
// 画面外の余白には淡いグレーを敷いて、PC でも擬似スマホっぽく表示する。
//
// 注意: 中身が viewport 高さ計算に依存する場合 (マップ等) は、ラッパーで
// flex/min-h を作らないこと。 子側で `h-[calc(100dvh-4rem)]` を直接指定する。
export function MobileShell({ children, withBottomNav, flush }: Props) {
  return (
    <div className="min-h-dvh w-full bg-slate-100">
      <div
        className={`relative mx-auto w-full max-w-[390px] bg-white shadow-xl ${
          flush ? "" : "px-4 py-4"
        }`}
        style={{ minHeight: "100dvh" }}
      >
        {children}
        {withBottomNav && !flush ? <div className="h-16" aria-hidden /> : null}
      </div>
    </div>
  );
}
