import { useEffect, useState } from "react";
import { Link, useParams } from "react-router";
import { ArrowLeft, UserRound } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { userClient } from "../lib/api-client";
import { GetUserRequest, User } from "../lib/gen/user/v1/user_pb";
import { formatUser } from "../lib/resource-name";

import type { Route } from "./+types/profile-other";

export function meta({ params }: Route.MetaArgs) {
  return [{ title: `@${params.customId} — PiaMap` }];
}

export default function ProfileOther() {
  const { customId } = useParams();
  const [user, setUser] = useState<User | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!customId) return;
    setErr(null);
    userClient
      .getUser(new GetUserRequest({ name: formatUser(customId) }))
      .then((res) => setUser(res.user ?? null))
      .catch((e) => setErr(e?.message || String(e)));
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
        </section>
      ) : null}
    </MobileShell>
  );
}
