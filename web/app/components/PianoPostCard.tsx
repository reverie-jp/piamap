import { useState } from "react";
import { Link } from "react-router";
import { Heart, MessageSquare, Pencil, Star, Trash2 } from "lucide-react";

import { pianoPostLikeClient } from "../lib/api-client";
import {
  LikePianoPostRequest,
  UnlikePianoPostRequest,
} from "../lib/gen/piano_post_like/v1/piano_post_like_pb";
import type { PianoPost } from "../lib/gen/piano_post/v1/piano_post_pb";
import { parsePianoPost } from "../lib/resource-name";

type Props = {
  post: PianoPost;
  /** ピアノ名/リンクを表示するか (タイムラインなど post 単独表示の時)。 */
  showPiano?: boolean;
  /** 投稿者として認識する custom_id (ログインユーザー)。投稿の author と一致する場合だけ編集/削除を出す。 */
  currentUserCustomId?: string;
  /** ❤️ ボタンを動作可能にする (未ログインなら false)。 */
  canLike?: boolean;
  /** 未ログイン時に ❤️ を押したときに親へ通知 (SignUpPromptModal を出す用)。 */
  onLikeUnauthorized?: () => void;
  /** 投稿詳細ページ自身でレンダリングする時は true。コメントボタンを Link でなくただの表示にする。 */
  hideCommentLink?: boolean;
  onEdit?: (post: PianoPost) => void;
  onDelete?: (post: PianoPost) => void;
};

export function PianoPostCard({
  post,
  showPiano,
  currentUserCustomId,
  canLike,
  onLikeUnauthorized,
  hideCommentLink,
  onEdit,
  onDelete,
}: Props) {
  const visitedAt = post.visitTime?.toDate();
  const visitedLabel = visitedAt
    ? `${visitedAt.getFullYear()}/${visitedAt.getMonth() + 1}/${visitedAt.getDate()}`
    : "";
  const authorCustomId = post.author?.split("/").pop() ?? "";
  const isAuthor = Boolean(currentUserCustomId) && currentUserCustomId === authorCustomId;
  let pianoIdParsed = "";
  let postIdParsed = "";
  try {
    const r = parsePianoPost(post.name);
    pianoIdParsed = r.pianoId;
    postIdParsed = r.postId;
  } catch {
    /* noop */
  }
  const detailPath =
    pianoIdParsed && postIdParsed ? `/pianos/${pianoIdParsed}/posts/${postIdParsed}` : "";

  // 楽観的更新でローカル state を持つ。
  const [liked, setLiked] = useState(post.viewerLiked);
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [busy, setBusy] = useState(false);

  const toggleLike = async () => {
    if (!canLike) {
      onLikeUnauthorized?.();
      return;
    }
    if (busy) return;
    setBusy(true);
    const willLike = !liked;
    setLiked(willLike);
    setLikeCount((c) => Math.max(0, c + (willLike ? 1 : -1)));
    try {
      if (willLike) {
        await pianoPostLikeClient.likePianoPost(
          new LikePianoPostRequest({ parent: post.name }),
        );
      } else {
        await pianoPostLikeClient.unlikePianoPost(
          new UnlikePianoPostRequest({ parent: post.name }),
        );
      }
    } catch {
      // ロールバック
      setLiked(!willLike);
      setLikeCount((c) => Math.max(0, c + (willLike ? -1 : 1)));
    } finally {
      setBusy(false);
    }
  };

  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-4">
      <header className="flex items-center justify-between gap-3 text-xs text-slate-500">
        <div className="min-w-0 flex-1 truncate">
          {authorCustomId ? (
            <Link
              to={`/profile/${authorCustomId}`}
              className="font-medium text-slate-700 hover:underline"
            >
              {post.authorDisplayName || `@${authorCustomId}`}
            </Link>
          ) : (
            <span className="text-slate-500">不明なユーザー</span>
          )}
          {visitedLabel ? <span className="ml-2">{visitedLabel} 訪問</span> : null}
        </div>
        {post.rating !== undefined ? (
          <div className="flex items-center gap-1 text-amber-600">
            <Star size={14} className="fill-amber-500 text-amber-500" />
            <span className="font-bold">{post.rating}</span>
          </div>
        ) : null}
      </header>
      {showPiano && pianoIdParsed ? (
        <Link
          to={`/pianos/${pianoIdParsed}`}
          className="mt-2 block text-sm font-medium text-amber-700 hover:underline"
        >
          {post.pianoDisplayName || pianoIdParsed}
        </Link>
      ) : null}
      {post.body ? (
        <p className="mt-2 whitespace-pre-line text-sm text-slate-800">{post.body}</p>
      ) : null}
      <PostAttributes post={post} />
      <div className="mt-2 flex items-center justify-between">
        <div className="flex items-center gap-1">
          <button
            type="button"
            aria-label={liked ? "いいねを解除" : "いいね"}
            aria-pressed={liked}
            onClick={toggleLike}
            disabled={busy}
            className={
              "inline-flex cursor-pointer items-center gap-1 rounded-full px-2 py-1 text-xs outline-none transition focus-visible:ring-2 focus-visible:ring-rose-500 disabled:cursor-progress " +
              (liked
                ? "text-rose-600 hover:bg-rose-50"
                : "text-slate-500 hover:bg-slate-100 hover:text-slate-700")
            }
          >
            <Heart size={14} className={liked ? "fill-rose-500 text-rose-500" : ""} />
            {likeCount > 0 ? <span className="font-semibold">{likeCount}</span> : null}
          </button>
          <CommentLink
            count={post.commentCount}
            to={hideCommentLink ? "" : detailPath}
          />
        </div>
        {isAuthor && (onEdit || onDelete) ? (
          <div className="flex gap-1">
            {onEdit ? (
              <button
                type="button"
                aria-label="編集"
                onClick={() => onEdit(post)}
                className="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-slate-500 hover:bg-slate-100 hover:text-slate-700 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
              >
                <Pencil size={14} />
              </button>
            ) : null}
            {onDelete ? (
              <button
                type="button"
                aria-label="削除"
                onClick={() => onDelete(post)}
                className="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-slate-500 hover:bg-rose-50 hover:text-rose-600 outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
              >
                <Trash2 size={14} />
              </button>
            ) : null}
          </div>
        ) : null}
      </div>
    </article>
  );
}

function CommentLink({ count, to }: { count: number; to: string }) {
  const cls =
    "inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs outline-none transition text-slate-500 hover:bg-slate-100 hover:text-slate-700 focus-visible:ring-2 focus-visible:ring-amber-500";
  if (!to) {
    return (
      <span className={cls}>
        <MessageSquare size={14} />
        {count > 0 ? <span className="font-semibold">{count}</span> : null}
      </span>
    );
  }
  return (
    <Link to={to} className={`${cls} cursor-pointer`} aria-label="返信を見る">
      <MessageSquare size={14} />
      {count > 0 ? <span className="font-semibold">{count}</span> : null}
    </Link>
  );
}

const ATTR_LABELS: { key: keyof PianoPost; left: string; right: string }[] = [
  { key: "ambientNoise", left: "静か", right: "賑やか" },
  { key: "footTraffic", left: "人通り少", right: "人通り多" },
  { key: "resonance", left: "響き弱", right: "響き豊か" },
  { key: "keyTouchWeight", left: "鍵盤軽い", right: "鍵盤重い" },
  { key: "tuningQuality", left: "調律悪い", right: "調律良い" },
];

function PostAttributes({ post }: { post: PianoPost }) {
  const items = ATTR_LABELS.map((a) => ({ ...a, value: post[a.key] as number | undefined }))
    .filter((it): it is typeof it & { value: number } => it.value != null);
  if (items.length === 0) return null;
  return (
    <ul className="mt-3 space-y-1.5">
      {items.map((it) => (
        <li key={String(it.key)} className="text-xs">
          <div className="flex justify-between text-slate-500">
            <span>{it.left}</span>
            <span>{it.right}</span>
          </div>
          <div className="mt-0.5 h-1 rounded-full bg-slate-100">
            <div
              className="h-full rounded-full bg-amber-500"
              style={{ width: `${(it.value / 5) * 100}%` }}
            />
          </div>
        </li>
      ))}
    </ul>
  );
}
