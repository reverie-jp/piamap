// dev 用: cmd/genjwt で発行した access token を localStorage に保存する。
// Phase 1 で SocialLogin に置き換え。

import { useEffect, useState } from "react";

const KEY = "piamap.dev.access_token";
type Listener = () => void;
const listeners = new Set<Listener>();

export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(KEY);
}

export function setAccessToken(token: string | null) {
  if (typeof window === "undefined") return;
  if (token == null || token === "") {
    window.localStorage.removeItem(KEY);
  } else {
    window.localStorage.setItem(KEY, token);
  }
  listeners.forEach((l) => l());
}

function subscribe(listener: Listener) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

export function useAuth() {
  const [authed, setAuthed] = useState(false);
  useEffect(() => {
    const update = () => setAuthed(Boolean(getAccessToken()));
    update();
    return subscribe(update);
  }, []);
  return { authed };
}
