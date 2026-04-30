import { useEffect, useState } from "react";

import { userClient } from "./api-client";
import { useAuth } from "./auth";
import { GetMyUserRequest, type User } from "./gen/user/v1/user_pb";

// ログインユーザーを取得するシンプルなフック。未ログインなら null を返す。
export function useMe() {
  const { authed } = useAuth();
  const [me, setMe] = useState<User | null>(null);

  useEffect(() => {
    if (!authed) {
      setMe(null);
      return;
    }
    let cancelled = false;
    userClient
      .getMyUser(new GetMyUserRequest())
      .then((res) => {
        if (!cancelled) setMe(res.user ?? null);
      })
      .catch(() => {
        if (!cancelled) setMe(null);
      });
    return () => {
      cancelled = true;
    };
  }, [authed]);

  return me;
}
