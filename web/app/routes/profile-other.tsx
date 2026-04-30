import { useEffect, useState } from "react";
import { Link, useParams } from "react-router";
import { ArrowLeft, UserRound } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { LikedPianoPostList } from "../components/LikedPianoPostList";
import { PianoPostCard } from "../components/PianoPostCard";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { Tabs } from "../components/Tabs";
import { UserCommentList } from "../components/UserCommentList";
import { UserPianoList } from "../components/UserPianoList";
import { useAuth } from "../lib/auth";
import { useMe } from "../lib/use-me";
import { pianoPostClient, userClient } from "../lib/api-client";
import { GetUserRequest, User } from "../lib/gen/user/v1/user_pb";
import {
  ListPianoPostsRequest,
  type PianoPost,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import { PianoListKind } from "../lib/gen/piano_user_list/v1/piano_user_list_pb";
import { formatUser } from "../lib/resource-name";

import type { Route } from "./+types/profile-other";

export function meta({ params }: Route.MetaArgs) {
  return [{ title: `@${params.customId} — PiaMap` }];
}

type ProfileTab = "posts" | "wishlist" | "visited" | "favorite" | "liked" | "comments";

export default function ProfileOther() {
  const { authed } = useAuth();
  const me = useMe();
  const { customId } = useParams();
  const [user, setUser] = useState<User | null>(null);
  const [posts, setPosts] = useState<PianoPost[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [tab, setTab] = useState<ProfileTab>("posts");
  const [signUpOpen, setSignUpOpen] = useState(false);

  useEffect(() => {
    if (!customId) return;
    setErr(null);
    userClient
      .getUser(new GetUserRequest({ name: formatUser(customId) }))
      .then((res) => setUser(res.user ?? null))
      .catch((e) => setErr(e?.message || String(e)));
    pianoPostClient
      .listPianoPosts(
        new ListPianoPostsRequest({ parent: formatUser(customId), pageSize: 20 }),
      )
      .then((res) => setPosts(res.pianoPosts))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  }, [customId]);

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-3">
        <Link to="/map" aria-label="戻る" className="text-slate-500 hover:text-slate-700">
          <ArrowLeft size={22} />
        </Link>
        <h1 className="text-base font-bold text-slate-900">@{customId}</h1>
      </header>

      {err ? (
        <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">読み込みエラー: {err}</p>
      ) : null}

      {user ? (
        <section>
          <div className="mt-2 flex items-center gap-3">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-slate-200">
              {user.avatarUrl ? (
                <img src={user.avatarUrl} alt="" className="h-full w-full rounded-full object-cover" />
              ) : (
                <UserRound size={28} className="text-slate-500" />
              )}
            </div>
            <div className="flex-1">
              <h2 className="text-lg font-bold text-slate-900">{user.displayName}</h2>
              <p className="text-xs text-slate-500">@{user.customId}</p>
            </div>
          </div>
          {user.biography ? (
            <p className="mt-4 text-sm text-slate-700 whitespace-pre-line">{user.biography}</p>
          ) : null}

          <dl className="mt-5 grid grid-cols-3 gap-2 rounded-2xl bg-slate-50 p-3 text-center text-xs">
            <div>
              <dt className="text-slate-500">投稿</dt>
              <dd className="text-base font-bold text-slate-900">{user.postCount}</dd>
            </div>
            <div>
              <dt className="text-slate-500">編集</dt>
              <dd className="text-base font-bold text-slate-900">{user.editCount}</dd>
            </div>
            <div>
              <dt className="text-slate-500">登録年</dt>
              <dd className="text-base font-bold text-slate-900">
                {user.createTime ? new Date(Number(user.createTime.seconds) * 1000).getFullYear() : "—"}
              </dd>
            </div>
          </dl>

          <section className="mt-6">
            <Tabs<ProfileTab>
              tabs={[
                { id: "posts", label: "投稿" },
                { id: "liked", label: "いいね" },
                { id: "comments", label: "コメント" },
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
                          currentUserCustomId={me?.customId}
                          canLike={authed}
                          onLikeUnauthorized={() => setSignUpOpen(true)}
                        />
                      </li>
                    ))}
                  </ul>
                )
              ) : tab === "liked" ? (
                <LikedPianoPostList
                  customId={customId ?? ""}
                  currentUserCustomId={me?.customId}
                  onLikeUnauthorized={() => setSignUpOpen(true)}
                />
              ) : tab === "comments" ? (
                <UserCommentList customId={customId ?? ""} />
              ) : tab === "wishlist" ? (
                <UserPianoList
                  customId={customId ?? ""}
                  listKind={PianoListKind.WISHLIST}
                  emptyMessage="行ってみたいピアノはまだありません"
                />
              ) : tab === "visited" ? (
                <UserPianoList
                  customId={customId ?? ""}
                  listKind={PianoListKind.VISITED}
                  emptyMessage="行ったピアノはまだありません"
                />
              ) : (
                <UserPianoList
                  customId={customId ?? ""}
                  listKind={PianoListKind.FAVORITE}
                  emptyMessage="お気に入りはまだありません"
                />
              )}
            </div>
          </section>
        </section>
      ) : null}
      <SignUpPromptModal isOpen={signUpOpen} onOpenChange={setSignUpOpen} action="いいね" />
    </MobileShell>
  );
}
