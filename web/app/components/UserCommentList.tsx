import { useEffect, useState } from "react";
import { Link } from "react-router";
import { MessageSquare, UserRound } from "lucide-react";

import { pianoPostCommentClient } from "../lib/api-client";
import {
  ListPianoPostCommentsRequest,
  type PianoPostComment,
} from "../lib/gen/piano_post_comment/v1/piano_post_comment_pb";
import { formatUser, parsePianoPostComment } from "../lib/resource-name";

type Props = {
  customId: string;
  emptyMessage?: string;
};

export function UserCommentList({ customId, emptyMessage }: Props) {
  const [comments, setComments] = useState<PianoPostComment[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!customId) return;
    setLoading(true);
    setErr(null);
    pianoPostCommentClient
      .listPianoPostComments(
        new ListPianoPostCommentsRequest({ parent: formatUser(customId), pageSize: 50 }),
      )
      .then((res) => setComments(res.pianoPostComments))
      .catch((e) => setErr((e as Error)?.message || String(e)))
      .finally(() => setLoading(false));
  }, [customId]);

  if (err) {
    return <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p>;
  }
  if (loading) {
    return <p className="text-center text-xs text-slate-400">読み込み中...</p>;
  }
  if (comments.length === 0) {
    return (
      <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        {emptyMessage ?? "まだ返信はありません"}
      </p>
    );
  }
  return (
    <ul className="space-y-3">
      {comments.map((c) => (
        <CommentRow key={c.name} comment={c} />
      ))}
    </ul>
  );
}

function CommentRow({ comment }: { comment: PianoPostComment }) {
  let pianoId = "";
  let postId = "";
  try {
    const parsed = parsePianoPostComment(comment.name);
    pianoId = parsed.pianoId;
    postId = parsed.postId;
  } catch {
    /* noop */
  }
  const createdAt = comment.createTime?.toDate();
  const dateLabel = createdAt
    ? `${createdAt.getFullYear()}/${createdAt.getMonth() + 1}/${createdAt.getDate()}`
    : "";
  return (
    <li className="rounded-2xl border border-slate-200 bg-white p-3">
      <header className="flex items-center gap-2 text-xs text-slate-500">
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-slate-200">
          {comment.authorAvatarUrl ? (
            <img
              src={comment.authorAvatarUrl}
              alt=""
              className="h-full w-full rounded-full object-cover"
            />
          ) : (
            <UserRound size={14} className="text-slate-500" />
          )}
        </div>
        <span className="truncate font-medium text-slate-700">
          {comment.authorDisplayName || ""}
        </span>
        {dateLabel ? <span className="text-slate-400">{dateLabel}</span> : null}
      </header>
      <p className="mt-2 whitespace-pre-line text-sm text-slate-800">{comment.body}</p>
      {pianoId && postId ? (
        <Link
          to={`/pianos/${pianoId}/posts/${postId}`}
          className="mt-2 inline-flex items-center gap-1 text-xs text-amber-700 hover:underline"
        >
          <MessageSquare size={12} /> 投稿を見る
        </Link>
      ) : null}
    </li>
  );
}
