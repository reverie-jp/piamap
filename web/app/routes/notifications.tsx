import { Bell } from "lucide-react";

import type { Route } from "./+types/notifications";

export function meta({}: Route.MetaArgs) {
  return [{ title: "通知 — PiaMap" }];
}

export default function Notifications() {
  return (
    <div className="flex h-full flex-col">
      <header className="border-b border-slate-200 px-4 py-3">
        <h1 className="text-base font-bold text-slate-900">通知</h1>
      </header>
      <div className="flex flex-1 flex-col items-center justify-center px-6 text-center text-slate-500">
        <Bell size={36} className="mb-3 text-slate-300" />
        <p className="text-sm">編集や投稿への返信が届くとここに表示されます。</p>
        <p className="mt-1 text-xs text-slate-400">Phase 1 で実装</p>
      </div>
    </div>
  );
}
