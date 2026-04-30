import { Link } from "react-router";

import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "PiaMap" },
    { name: "description", content: "ストリートピアノを発見・記録・共有する SNS" },
  ];
}

export default function Home() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-slate-900">PiaMap</h1>
        <p className="mt-2 text-slate-600">ストリートピアノを発見・記録・共有する SNS</p>
        <Link
          to="/map"
          className="mt-6 inline-block rounded-full bg-amber-500 px-6 py-3 font-semibold text-white shadow hover:bg-amber-600"
        >
          マップを開く →
        </Link>
      </div>
    </main>
  );
}
