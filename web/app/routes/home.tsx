import { useEffect } from "react";
import { Link, useNavigate } from "react-router";
import { MapPinned, Music, Star } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { Button } from "../components/ui/button";
import { useAuth } from "../lib/auth";

import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "PiaMap — ストリートピアノを発見・記録・共有" },
    { name: "description", content: "近くのストリートピアノを地図で見つけて、演奏を記録するSNS" },
  ];
}

export default function Landing() {
  const { authed } = useAuth();
  const nav = useNavigate();

  // 認証済なら /map にリダイレクト。
  useEffect(() => {
    if (authed) nav("/map", { replace: true });
  }, [authed, nav]);

  return (
    <MobileShell>
      <div className="flex flex-col items-center pt-10 pb-6 text-center">
        <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-amber-500 text-white shadow-lg">
          <Music size={32} />
        </div>
        <h1 className="mt-4 text-2xl font-bold text-slate-900">PiaMap</h1>
        <p className="mt-1 text-sm text-slate-600">
          ストリートピアノを発見・記録・共有する SNS
        </p>
      </div>

      <ul className="mt-2 space-y-3 text-sm text-slate-700">
        <li className="flex items-start gap-3 rounded-2xl bg-slate-50 p-4">
          <MapPinned size={20} className="mt-0.5 shrink-0 text-amber-600" />
          <div>
            <strong className="block text-slate-900">マップで探す</strong>
            <span className="text-xs text-slate-600">近くのピアノを地図から発見</span>
          </div>
        </li>
        <li className="flex items-start gap-3 rounded-2xl bg-slate-50 p-4">
          <Star size={20} className="mt-0.5 shrink-0 text-amber-600" />
          <div>
            <strong className="block text-slate-900">演奏を記録・評価</strong>
            <span className="text-xs text-slate-600">★評価 + 動画/音声 + 環境メモ</span>
          </div>
        </li>
        <li className="flex items-start gap-3 rounded-2xl bg-slate-50 p-4">
          <Music size={20} className="mt-0.5 shrink-0 text-amber-600" />
          <div>
            <strong className="block text-slate-900">ピアニストと繋がる</strong>
            <span className="text-xs text-slate-600">同じピアノを訪れた人のレビューを読む</span>
          </div>
        </li>
      </ul>

      <div className="mt-8 space-y-2">
        <Button
          size="lg"
          className="w-full"
          onPress={() => alert("Google ログインは Phase 1 で実装予定")}
        >
          Google でログイン
        </Button>
        <Link
          to="/map"
          className="block text-center text-sm font-medium text-slate-600 underline-offset-2 hover:underline"
        >
          ログインせずにマップを見る
        </Link>
        <Link
          to="/settings"
          className="block text-center text-xs text-slate-400 underline-offset-2 hover:underline"
        >
          dev token を入力 (開発用)
        </Link>
      </div>
    </MobileShell>
  );
}
