import { useEffect, useState } from "react";

import { pianoPostLikeClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import type { PianoPost } from "../lib/gen/piano_post/v1/piano_post_pb";
import { ListLikedPianoPostsRequest } from "../lib/gen/piano_post_like/v1/piano_post_like_pb";
import { formatUser } from "../lib/resource-name";
import { PianoPostCard } from "./PianoPostCard";

type Props = {
  customId: string;
  currentUserCustomId?: string;
  onLikeUnauthorized?: () => void;
  emptyMessage?: string;
};

export function LikedPianoPostList({
  customId,
  currentUserCustomId,
  onLikeUnauthorized,
  emptyMessage,
}: Props) {
  const { authed } = useAuth();
  const [posts, setPosts] = useState<PianoPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!customId) return;
    setLoading(true);
    setErr(null);
    pianoPostLikeClient
      .listLikedPianoPosts(
        new ListLikedPianoPostsRequest({ parent: formatUser(customId), pageSize: 20 }),
      )
      .then((res) => setPosts(res.pianoPosts))
      .catch((e) => setErr((e as Error)?.message || String(e)))
      .finally(() => setLoading(false));
  }, [customId]);

  if (err) {
    return <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p>;
  }
  if (loading) {
    return <p className="text-center text-xs text-slate-400">読み込み中...</p>;
  }
  if (posts.length === 0) {
    return (
      <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        {emptyMessage ?? "いいねした投稿はまだありません"}
      </p>
    );
  }
  return (
    <ul className="space-y-3">
      {posts.map((p) => (
        <li key={p.name}>
          <PianoPostCard
            post={p}
            showPiano
            currentUserCustomId={currentUserCustomId}
            canLike={authed}
            onLikeUnauthorized={onLikeUnauthorized}
          />
        </li>
      ))}
    </ul>
  );
}
