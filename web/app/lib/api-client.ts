import { createPromiseClient, type Interceptor } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import { AccountService } from "./gen/account/v1/account_connect";
import { PianoService } from "./gen/piano/v1/piano_connect";
import { PianoPostService } from "./gen/piano_post/v1/piano_post_connect";
import { PianoPostCommentService } from "./gen/piano_post_comment/v1/piano_post_comment_connect";
import { PianoPostLikeService } from "./gen/piano_post_like/v1/piano_post_like_connect";
import { PianoUserListService } from "./gen/piano_user_list/v1/piano_user_list_connect";
import { UserService } from "./gen/user/v1/user_connect";
import { getAccessToken } from "./auth";

const baseUrl =
  (typeof window !== "undefined" && (window as any).PIAMAP_API_URL) ||
  "http://localhost:50051";

// dev token があれば Authorization ヘッダに乗せる。
const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getAccessToken();
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  return await next(req);
};

const transport = createConnectTransport({
  baseUrl,
  useBinaryFormat: false,
  interceptors: [authInterceptor],
});

export const accountClient = createPromiseClient(AccountService, transport);
export const pianoClient = createPromiseClient(PianoService, transport);
export const pianoPostClient = createPromiseClient(PianoPostService, transport);
export const pianoPostCommentClient = createPromiseClient(PianoPostCommentService, transport);
export const pianoPostLikeClient = createPromiseClient(PianoPostLikeService, transport);
export const pianoUserListClient = createPromiseClient(PianoUserListService, transport);
export const userClient = createPromiseClient(UserService, transport);
