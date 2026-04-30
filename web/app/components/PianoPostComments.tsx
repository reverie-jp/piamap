import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router";
import { Trash2, UserRound } from "lucide-react";

import { pianoPostCommentClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import {
  CreatePianoPostCommentRequest,
  DeletePianoPostCommentRequest,
  ListPianoPostCommentsRequest,
  PianoPostComment,
} from "../lib/gen/piano_post_comment/v1/piano_post_comment_pb";
import { Button } from "./ui/button";
import { TextField } from "./ui/text-field";
import { ConfirmDialog } from "./ConfirmDialog";

type Props = {
  /** 投稿の name = "pianos/{piano_id}/posts/{post_id}" */
  postName: string;
  /** ログインユーザーの custom_id (削除ボタンを出すかの判定)。 */
  currentUserCustomId?: string;
  /** 未ログイン時に投稿しようとしたときの通知 (SignUpPromptModal を出す)。 */
  onCommentUnauthorized?: () => void;
  /** コメント数の変動を親に通知 (count 表示更新用)。 */
  onCountChange?: (delta: number) => void;
};

export function PianoPostComments({
  postName,
  currentUserCustomId,
  onCommentUnauthorized,
  onCountChange,
}: Props) {
  const { authed } = useAuth();
  const [comments, setComments] = useState<PianoPostComment[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);
  const [body, setBody] = useState("");
  const [posting, setPosting] = useState(false);
  const [deletingComment, setDeletingComment] = useState<PianoPostComment | null>(null);

  const reload = useCallback(() => {
    setLoading(true);
    setErr(null);
    return pianoPostCommentClient
      .listPianoPostComments(
        new ListPianoPostCommentsRequest({ parent: postName, pageSize: 100 }),
      )
      .then((res) => setComments(res.pianoPostComments))
      .catch((e) => setErr((e as Error)?.message || String(e)))
      .finally(() => setLoading(false));
  }, [postName]);

  useEffect(() => {
    reload();
  }, [reload]);

  const submit = async () => {
    if (!authed) {
      onCommentUnauthorized?.();
      return;
    }
    const trimmed = body.trim();
    if (!trimmed || posting) return;
    setPosting(true);
    setErr(null);
    try {
      const res = await pianoPostCommentClient.createPianoPostComment(
        new CreatePianoPostCommentRequest({
          parent: postName,
          pianoPostComment: new PianoPostComment({ body: trimmed }),
        }),
      );
      if (res.pianoPostComment) {
        setComments((prev) => [...prev, res.pianoPostComment!]);
        onCountChange?.(1);
      }
      setBody("");
    } catch (e) {
      setErr((e as Error)?.message || String(e));
    } finally {
      setPosting(false);
    }
  };

  const handleConfirmDelete = async () => {
    const target = deletingComment;
    if (!target) return;
    setDeletingComment(null);
    try {
      await pianoPostCommentClient.deletePianoPostComment(
        new DeletePianoPostCommentRequest({ name: target.name }),
      );
      setComments((prev) => prev.filter((c) => c.name !== target.name));
      onCountChange?.(-1);
    } catch (e) {
      setErr((e as Error)?.message || String(e));
    }
  };

  return (
    <section className="mt-3 border-t border-slate-100 pt-3">
      {err ? (
        <p className="mb-2 rounded bg-rose-50 p-2 text-xs text-rose-700">{err}</p>
      ) : null}

      {loading ? (
        <p className="text-center text-xs text-slate-400">読み込み中...</p>
      ) : comments.length === 0 ? (
        <p className="text-xs text-slate-400">まだコメントはありません</p>
      ) : (
        <ul className="space-y-3">
          {comments.map((c) => (
            <CommentItem
              key={c.name}
              comment={c}
              currentUserCustomId={currentUserCustomId}
              onDelete={() => setDeletingComment(c)}
            />
          ))}
        </ul>
      )}

      <div className="mt-3 space-y-2">
        <TextField
          aria-label="コメント"
          placeholder={authed ? "コメントを書く..." : "コメントするにはログインが必要です"}
          multiline
          rows={2}
          value={body}
          onChange={setBody}
          isDisabled={!authed || posting}
        />
        <div className="flex justify-end">
          <Button
            size="sm"
            onPress={submit}
            isDisabled={!body.trim() || posting}
          >
            {posting ? "送信中..." : authed ? "コメントする" : "ログインしてコメント"}
          </Button>
        </div>
      </div>

      <ConfirmDialog
        isOpen={deletingComment !== null}
        onOpenChange={(open) => {
          if (!open) setDeletingComment(null);
        }}
        title="コメントを削除"
        message="この操作は取り消せません。本当に削除しますか?"
        confirmLabel="削除する"
        destructive
        onConfirm={handleConfirmDelete}
      />
    </section>
  );
}

function CommentItem({
  comment,
  currentUserCustomId,
  onDelete,
}: {
  comment: PianoPostComment;
  currentUserCustomId?: string;
  onDelete: () => void;
}) {
  const authorCustomId = comment.author?.split("/").pop() ?? "";
  const isAuthor = Boolean(currentUserCustomId) && currentUserCustomId === authorCustomId;
  const createdAt = comment.createTime?.toDate();
  return (
    <li className="flex gap-2 text-sm">
      <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-slate-200">
        {comment.authorAvatarUrl ? (
          <img
            src={comment.authorAvatarUrl}
            alt=""
            className="h-full w-full rounded-full object-cover"
          />
        ) : (
          <UserRound size={16} className="text-slate-500" />
        )}
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2 text-xs">
          {authorCustomId ? (
            <Link
              to={`/profile/${authorCustomId}`}
              className="truncate font-medium text-slate-700 hover:underline"
            >
              {comment.authorDisplayName || `@${authorCustomId}`}
            </Link>
          ) : (
            <span className="text-slate-500">不明なユーザー</span>
          )}
          {createdAt ? (
            <span className="text-slate-400">
              {createdAt.getFullYear()}/{createdAt.getMonth() + 1}/{createdAt.getDate()}
            </span>
          ) : null}
          {isAuthor ? (
            <button
              type="button"
              aria-label="削除"
              onClick={onDelete}
              className="ml-auto cursor-pointer rounded-full p-1 text-slate-400 hover:bg-rose-50 hover:text-rose-600 outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
            >
              <Trash2 size={12} />
            </button>
          ) : null}
        </div>
        <p className="mt-0.5 whitespace-pre-line text-slate-800">{comment.body}</p>
      </div>
    </li>
  );
}

