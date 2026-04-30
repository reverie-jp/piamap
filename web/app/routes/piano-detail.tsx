import { useEffect, useState } from "react";
import { Link, useParams } from "react-router";
import { ArrowLeft, MapPin, Music, Star, UserRound } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { pianoClient } from "../lib/api-client";
import { GetPianoRequest, Piano } from "../lib/gen/piano/v1/piano_pb";
import { formatPiano } from "../lib/resource-name";

import type { Route } from "./+types/piano-detail";

export function meta({}: Route.MetaArgs) {
  return [{ title: "ピアノ — PiaMap" }];
}

export default function PianoDetail() {
  const { id } = useParams();
  const [piano, setPiano] = useState<Piano | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    setErr(null);
    pianoClient
      .getPiano(new GetPianoRequest({ name: formatPiano(id) }))
      .then((res) => setPiano(res.piano ?? null))
      .catch((e) => setErr(e?.message || String(e)));
  }, [id]);

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-3">
        <Link to="/map" aria-label="戻る" className="text-slate-500 hover:text-slate-700">
          <ArrowLeft size={22} />
        </Link>
        <h1 className="truncate text-base font-bold text-slate-900">
          {piano?.displayName ?? "読み込み中..."}
        </h1>
      </header>

      {err ? <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p> : null}

      {piano ? (
        <section className="space-y-4">
          <RatingBlock piano={piano} />

          {piano.description ? (
            <p className="rounded-2xl bg-slate-50 p-4 text-sm text-slate-700 whitespace-pre-line">
              {piano.description}
            </p>
          ) : null}

          <dl className="space-y-2 rounded-2xl border border-slate-200 p-4 text-sm">
            {piano.address ? (
              <Row icon={<MapPin size={16} className="text-slate-400" />} label="住所" value={piano.address} />
            ) : null}
            <Row
              icon={<Music size={16} className="text-slate-400" />}
              label="ピアノ"
              value={`${pianoTypeLabel(piano.pianoType)} / ${piano.pianoBrand}${
                piano.pianoModel ? ` ${piano.pianoModel}` : ""
              }`}
            />
            {piano.creator ? (
              <Row
                icon={<UserRound size={16} className="text-slate-400" />}
                label="登録"
                value={
                  <Link to={`/${piano.creator}`} className="text-amber-600 hover:underline">
                    @{piano.creator.split("/").pop()}
                  </Link>
                }
              />
            ) : null}
          </dl>

          <AttributeMeters piano={piano} />

          <p className="text-center text-xs text-slate-400">
            piano_post (投稿/レビュー) ができるとここに一覧が出ます。
          </p>
        </section>
      ) : null}
    </MobileShell>
  );
}

function RatingBlock({ piano }: { piano: Piano }) {
  if (piano.postCount === 0) {
    return (
      <div className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        まだレビューがありません
      </div>
    );
  }
  return (
    <div className="flex items-center gap-3 rounded-2xl bg-amber-50 p-4">
      <Star size={28} className="fill-amber-500 text-amber-500" />
      <div>
        <span className="text-2xl font-bold text-amber-700">{piano.ratingAverage.toFixed(1)}</span>
        <span className="ml-2 text-xs text-slate-600">({piano.postCount}件のレビュー)</span>
      </div>
    </div>
  );
}

function Row({ icon, label, value }: { icon: React.ReactNode; label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5">{icon}</div>
      <div className="flex-1">
        <dt className="text-xs text-slate-500">{label}</dt>
        <dd className="text-slate-800">{value}</dd>
      </div>
    </div>
  );
}

function AttributeMeters({ piano }: { piano: Piano }) {
  const items: { label: string; value: number }[] = [
    { label: "静か—賑やか", value: piano.ambientNoiseAverage },
    { label: "人通り", value: piano.footTrafficAverage },
    { label: "響き", value: piano.resonanceAverage },
    { label: "鍵盤の重さ", value: piano.keyTouchWeightAverage },
    { label: "調律状態", value: piano.tuningQualityAverage },
  ].filter((i) => i.value > 0);
  if (items.length === 0) return null;
  return (
    <div className="rounded-2xl border border-slate-200 p-4">
      <h3 className="mb-3 text-xs font-bold uppercase tracking-wide text-slate-500">特徴</h3>
      <ul className="space-y-2">
        {items.map((it) => (
          <li key={it.label} className="text-sm">
            <div className="flex justify-between text-xs text-slate-600">
              <span>{it.label}</span>
              <span>{it.value.toFixed(1)} / 5</span>
            </div>
            <div className="mt-1 h-1.5 rounded-full bg-slate-100">
              <div
                className="h-full rounded-full bg-amber-500"
                style={{ width: `${(it.value / 5) * 100}%` }}
              />
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

function pianoTypeLabel(t: number): string {
  switch (t) {
    case 1:
      return "グランド";
    case 2:
      return "アップライト";
    case 3:
      return "電子";
    case 4:
      return "不明";
    default:
      return "—";
  }
}
