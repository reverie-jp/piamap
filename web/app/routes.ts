import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  // 未認証時の LP (認証済みは中で /map にリダイレクト)。
  index("routes/home.tsx"),

  // ボトムナビ付きシェルを共有するページ群。
  layout("routes/_app.tsx", [
    route("map", "routes/map.tsx"),
    route("timeline", "routes/timeline.tsx"),
    route("notifications", "routes/notifications.tsx"),
    route("profile/me", "routes/profile-me.tsx"),
  ]),

  // ボトムナビ無し、プッシュ遷移のページ。
  route("pianos/:id", "routes/piano-detail.tsx"),
  route("pianos/:id/posts/:postId", "routes/piano-post-detail.tsx"),
  route("profile/:customId", "routes/profile-other.tsx"),
  route("profile/:customId/saved/:kind", "routes/saved-list.tsx"),
  route("settings", "routes/settings.tsx"),
] satisfies RouteConfig;
