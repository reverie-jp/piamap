// dev 用: cmd/genjwt で発行した access token を localStorage に保存する。
// Phase 1 で SocialLogin に置き換え。

const KEY = "piamap.dev.access_token";

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
}
