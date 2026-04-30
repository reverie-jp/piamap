import { useEffect, useState } from "react";
import { ScrollText } from "lucide-react";

import { pianoPostClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import { useMe } from "../lib/use-me";
import {
  ListPianoPostsRequest,
  type PianoPost,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import { PianoPostCard } from "../components/PianoPostCard";
import { SignUpPromptModal } from "../components/SignUpPromptModal";

import type { Route } from "./+types/timeline";

export function meta({}: Route.MetaArgs) {
  return [{ title: "タイムライン — PiaMap" }];
}

export default function Timeline() {
  const { authed } = useAuth();
  const me = useMe();
  const [posts, setPosts] = useState<PianoPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);
  const [signUpOpen, setSignUpOpen] = useState(false);

  useEffect(() => {
    setLoading(true);
    pianoPostClient
      .listPianoPosts(new ListPianoPostsRequest({ parent: "-", pageSize: 20 }))
      .then((res) => setPosts(res.pianoPosts))
      .catch((e) => setErr((e as Error)?.message || String(e)))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="flex h-full flex-col">
      <header className="border-b border-slate-200 px-4 py-3">
        <h1 className="text-base font-bold text-slate-900">タイムライン</h1>
      </header>
      <div className="flex-1 overflow-y-auto p-3">
        {err ? (
          <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p>
        ) : loading ? (
          <p className="text-center text-xs text-slate-400">読み込み中...</p>
        ) : posts.length === 0 ? (
          <div className="flex flex-col items-center justify-center px-6 py-16 text-center text-slate-500">
            <ScrollText size={36} className="mb-3 text-slate-300" />
            <p className="text-sm">まだ投稿がありません</p>
            <p className="mt-1 text-xs text-slate-400">
              地図からピアノを開いて最初の投稿をしてみよう
            </p>
          </div>
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
        )}
      </div>
      <SignUpPromptModal isOpen={signUpOpen} onOpenChange={setSignUpOpen} action="いいね" />
    </div>
  );
}
