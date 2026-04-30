import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router";
import { LogIn, Settings, UserRound } from "lucide-react";

import { Button } from "../components/ui/button";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { CreatePianoPostModal } from "../components/CreatePianoPostModal";
import { LikedPianoPostList } from "../components/LikedPianoPostList";
import { PianoPostCard } from "../components/PianoPostCard";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { Tabs } from "../components/Tabs";
import { UserPianoList } from "../components/UserPianoList";
import { PianoListKind } from "../lib/gen/piano_user_list/v1/piano_user_list_pb";
import { pianoPostClient, userClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import {
  DeletePianoPostRequest,
  ListPianoPostsRequest,
  type PianoPost,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import { GetMyUserRequest, type User } from "../lib/gen/user/v1/user_pb";
import { formatUser, parsePianoPost } from "../lib/resource-name";

import type { Route } from "./+types/profile-me";

export function meta({}: Route.MetaArgs) {
  return [{ title: "プロフィール — PiaMap" }];
}

type ProfileTab = "posts" | "wishlist" | "visited" | "favorite" | "liked";

export default function ProfileMe() {
  const { authed } = useAuth();
  const [me, setMe] = useState<User | null>(null);
  const [posts, setPosts] = useState<PianoPost[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [editingPost, setEditingPost] = useState<PianoPost | null>(null);
  const [deletingPost, setDeletingPost] = useState<PianoPost | null>(null);
  const [tab, setTab] = useState<ProfileTab>("posts");

  const reloadPosts = useCallback((customId: string) => {
    return pianoPostClient
      .listPianoPosts(
        new ListPianoPostsRequest({ parent: formatUser(customId), pageSize: 20 }),
      )
      .then((res) => setPosts(res.pianoPosts))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  }, []);

  useEffect(() => {
    if (!authed) return;
    setErr(null);
    userClient
      .getMyUser(new GetMyUserRequest())
      .then((res) => {
        const u = res.user ?? null;
        setMe(u);
        if (u) reloadPosts(u.customId);
      })
      .catch((e) => setErr(e?.message || String(e)));
  }, [authed, reloadPosts]);

  const handleConfirmDelete = async () => {
    if (!deletingPost || !me) return;
    try {
      await pianoPostClient.deletePianoPost(
        new DeletePianoPostRequest({ name: deletingPost.name }),
      );
      setDeletingPost(null);
      await reloadPosts(me.customId);
    } catch (e) {
      setErr((e as Error)?.message || String(e));
      setDeletingPost(null);
    }
  };

  return (
    <div className="flex h-full flex-col">
      <header className="flex items-center justify-between border-b border-slate-200 px-4 py-3">
        <h1 className="text-base font-bold text-slate-900">プロフィール</h1>
        <Link to="/settings" aria-label="設定" className="text-slate-500 hover:text-slate-700">
          <Settings size={20} />
        </Link>
      </header>

      {!authed ? <SignedOutPrompt /> : null}

      {authed && err ? (
        <p className="mx-4 mt-4 rounded bg-rose-50 p-3 text-sm text-rose-700">読み込みエラー: {err}</p>
      ) : null}

      {authed && me ? (
        <div className="flex-1 overflow-y-auto">
          <section className="px-4 pb-4">
            <div className="mt-4 flex items-center gap-3">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-slate-200">
                {me.avatarUrl ? (
                  <img src={me.avatarUrl} alt="" className="h-full w-full rounded-full object-cover" />
                ) : (
                  <UserRound size={28} className="text-slate-500" />
                )}
              </div>
              <div className="flex-1">
                <h2 className="text-lg font-bold text-slate-900">{me.displayName}</h2>
                <p className="text-xs text-slate-500">@{me.customId}</p>
              </div>
            </div>

            {me.biography ? (
              <p className="mt-4 text-sm text-slate-700 whitespace-pre-line">{me.biography}</p>
            ) : null}

            <dl className="mt-5 grid grid-cols-3 gap-2 rounded-2xl bg-slate-50 p-3 text-center text-xs">
              <div>
                <dt className="text-slate-500">投稿</dt>
                <dd className="text-base font-bold text-slate-900">{me.postCount}</dd>
              </div>
              <div>
                <dt className="text-slate-500">編集</dt>
                <dd className="text-base font-bold text-slate-900">{me.editCount}</dd>
              </div>
              <div>
                <dt className="text-slate-500">登録年</dt>
                <dd className="text-base font-bold text-slate-900">
                  {me.createTime ? new Date(Number(me.createTime.seconds) * 1000).getFullYear() : "—"}
                </dd>
              </div>
            </dl>
          </section>

          <section className="px-4 pb-6">
            <Tabs<ProfileTab>
              tabs={[
                { id: "posts", label: "投稿" },
                { id: "liked", label: "いいね" },
                { id: "wishlist", label: "行ってみたい" },
                { id: "visited", label: "行ったことある" },
                { id: "favorite", label: "お気に入り" },
              ]}
              active={tab}
              onChange={setTab}
            />
            <div className="mt-3">
              {tab === "posts" ? (
                posts.length === 0 ? (
                  <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
                    まだ投稿がありません
                  </p>
                ) : (
                  <ul className="space-y-3">
                    {posts.map((p) => (
                      <li key={p.name}>
                        <PianoPostCard
                          post={p}
                          showPiano
                          currentUserCustomId={me.customId}
                          canLike
                          onEdit={(post) => setEditingPost(post)}
                          onDelete={(post) => setDeletingPost(post)}
                        />
                      </li>
                    ))}
                  </ul>
                )
              ) : tab === "liked" ? (
                <LikedPianoPostList customId={me.customId} />
              ) : tab === "wishlist" ? (
                <UserPianoList
                  customId={me.customId}
                  listKind={PianoListKind.WISHLIST}
                  emptyMessage="行ってみたいピアノはまだありません"
                />
              ) : tab === "visited" ? (
                <UserPianoList
                  customId={me.customId}
                  listKind={PianoListKind.VISITED}
                  emptyMessage="行ったピアノはまだありません"
                />
              ) : (
                <UserPianoList
                  customId={me.customId}
                  listKind={PianoListKind.FAVORITE}
                  emptyMessage="お気に入りに登録したピアノはまだありません"
                />
              )}
            </div>
          </section>
        </div>
      ) : null}

      {editingPost ? (
        <CreatePianoPostModal
          isOpen={editingPost !== null}
          onOpenChange={(open) => {
            if (!open) setEditingPost(null);
          }}
          pianoId={parsePostPianoIdSafe(editingPost.name)}
          pianoName={editingPost.pianoDisplayName || ""}
          editingPost={editingPost}
          onSaved={() => {
            if (me) reloadPosts(me.customId);
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
    </div>
  );
}

function parsePostPianoIdSafe(name: string): string {
  try {
    return parsePianoPost(name).pianoId;
  } catch {
    return "";
  }
}

function SignedOutPrompt() {
  const [open, setOpen] = useState(false);
  return (
    <div className="flex flex-1 flex-col items-center justify-center px-6 text-center">
      <UserRound size={48} className="mb-4 text-slate-300" />
      <p className="text-sm text-slate-600">プロフィールを表示するにはログインが必要です。</p>
      <div className="mt-5 w-full max-w-[260px] space-y-2">
        <Button size="lg" className="w-full" onPress={() => setOpen(true)}>
          <LogIn size={16} /> ログインする
        </Button>
        <Link
          to="/settings"
          className="block text-center text-xs text-slate-400 underline-offset-2 hover:underline"
        >
          dev token を入力 (開発用)
        </Link>
      </div>
      <SignUpPromptModal isOpen={open} onOpenChange={setOpen} action="プロフィールの表示" />
    </div>
  );
}
