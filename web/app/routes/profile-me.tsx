import { useEffect, useState } from "react";
import { Link } from "react-router";
import { LogIn, Settings, UserRound } from "lucide-react";

import { Button } from "../components/ui/button";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { userClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import { User } from "../lib/gen/user/v1/user_pb";
import { GetMyUserRequest } from "../lib/gen/user/v1/user_pb";

import type { Route } from "./+types/profile-me";

export function meta({}: Route.MetaArgs) {
  return [{ title: "プロフィール — PiaMap" }];
}

export default function ProfileMe() {
  const { authed } = useAuth();
  const [me, setMe] = useState<User | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!authed) return;
    setErr(null);
    userClient
      .getMyUser(new GetMyUserRequest())
      .then((res) => setMe(res.user ?? null))
      .catch((e) => setErr(e?.message || String(e)));
  }, [authed]);

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
        <section className="px-4 pb-6">
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
      ) : null}
    </div>
  );
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
