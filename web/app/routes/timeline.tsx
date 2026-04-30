import { ScrollText } from "lucide-react";

import type { Route } from "./+types/timeline";

export function meta({}: Route.MetaArgs) {
  return [{ title: "タイムライン — PiaMap" }];
}

export default function Timeline() {
  return (
    <div className="flex h-full flex-col">
      <header className="border-b border-slate-200 px-4 py-3">
        <h1 className="text-base font-bold text-slate-900">タイムライン</h1>
      </header>
      <div className="flex flex-1 flex-col items-center justify-center px-6 text-center text-slate-500">
        <ScrollText size={36} className="mb-3 text-slate-300" />
        <p className="text-sm">投稿(piano_post)ができたら、ここに新着レビューが流れます。</p>
        <p className="mt-1 text-xs text-slate-400">Phase 1 で実装</p>
      </div>
    </div>
  );
}
