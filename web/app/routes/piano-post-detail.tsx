import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { ArrowLeft } from "lucide-react";

import { ConfirmDialog } from "../components/ConfirmDialog";
import { CreatePianoPostModal } from "../components/CreatePianoPostModal";
import { MobileShell } from "../components/MobileShell";
import { PianoPostCard } from "../components/PianoPostCard";
import { PianoPostComments } from "../components/PianoPostComments";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { useAuth } from "../lib/auth";
import { useMe } from "../lib/use-me";
import { pianoPostClient } from "../lib/api-client";
import {
  DeletePianoPostRequest,
  GetPianoPostRequest,
  type PianoPost,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import { formatPianoPost } from "../lib/resource-name";

import type { Route } from "./+types/piano-post-detail";

export function meta({}: Route.MetaArgs) {
  return [{ title: "投稿 — PiaMap" }];
}

export default function PianoPostDetail() {
  const { id, postId } = useParams();
  const navigate = useNavigate();
  const { authed } = useAuth();
  const me = useMe();
  const [post, setPost] = useState<PianoPost | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [editingPost, setEditingPost] = useState<PianoPost | null>(null);
  const [deletingPost, setDeletingPost] = useState<PianoPost | null>(null);
  const [signUpOpen, setSignUpOpen] = useState(false);

  useEffect(() => {
    if (!id || !postId) return;
    setErr(null);
    pianoPostClient
      .getPianoPost(new GetPianoPostRequest({ name: formatPianoPost(id, postId) }))
      .then((res) => setPost(res.pianoPost ?? null))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  }, [id, postId]);

  const reload = () => {
    if (!id || !postId) return;
    pianoPostClient
      .getPianoPost(new GetPianoPostRequest({ name: formatPianoPost(id, postId) }))
      .then((res) => setPost(res.pianoPost ?? null))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  };

  const handleConfirmDelete = async () => {
    if (!deletingPost) return;
    try {
      await pianoPostClient.deletePianoPost(
        new DeletePianoPostRequest({ name: deletingPost.name }),
      );
      setDeletingPost(null);
      navigate(`/pianos/${id}`, { replace: true });
    } catch (e) {
      setErr((e as Error)?.message || String(e));
      setDeletingPost(null);
    }
  };

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-3">
        <button
          type="button"
          onClick={() => navigate(-1)}
          aria-label="戻る"
          className="cursor-pointer text-slate-500 hover:text-slate-700"
        >
          <ArrowLeft size={22} />
        </button>
        <h1 className="text-base font-bold text-slate-900">投稿</h1>
      </header>

      {err ? (
        <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">読み込みエラー: {err}</p>
      ) : null}

      {post ? (
        <>
          {post.pianoDisplayName && id ? (
            <Link
              to={`/pianos/${id}`}
              className="mb-3 block text-sm font-medium text-amber-700 hover:underline"
            >
              {post.pianoDisplayName}
            </Link>
          ) : null}

          <PianoPostCard
            post={post}
            currentUserCustomId={me?.customId}
            canLike={authed}
            hideCommentLink
            onLikeUnauthorized={() => setSignUpOpen(true)}
            onEdit={(p) => setEditingPost(p)}
            onDelete={(p) => setDeletingPost(p)}
          />

          <section className="mt-4 rounded-2xl border border-slate-200 bg-white p-4">
            <h2 className="text-sm font-bold text-slate-900">
              コメント
              {post.commentCount > 0 ? (
                <span className="ml-1 text-slate-500">({post.commentCount})</span>
              ) : null}
            </h2>
            <PianoPostComments
              postName={post.name}
              currentUserCustomId={me?.customId}
              onCommentUnauthorized={() => setSignUpOpen(true)}
              onCountChange={(d) =>
                setPost((p) => {
                  if (!p) return p;
                  const next = p.clone();
                  next.commentCount = Math.max(0, p.commentCount + d);
                  return next;
                })
              }
            />
          </section>
        </>
      ) : !err ? (
        <p className="text-center text-xs text-slate-400">読み込み中...</p>
      ) : null}

      {editingPost && id ? (
        <CreatePianoPostModal
          isOpen={editingPost !== null}
          onOpenChange={(open) => {
            if (!open) setEditingPost(null);
          }}
          pianoId={id}
          pianoName={editingPost.pianoDisplayName || ""}
          editingPost={editingPost}
          onSaved={() => {
            setEditingPost(null);
            reload();
          }}
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

      <SignUpPromptModal
        isOpen={signUpOpen}
        onOpenChange={setSignUpOpen}
        action="この操作"
      />
    </MobileShell>
  );
}
