import { useEffect, useState } from "react";
import { Link } from "react-router";
import { Star } from "lucide-react";

import { pianoUserListClient } from "../lib/api-client";
import type { Piano } from "../lib/gen/piano/v1/piano_pb";
import {
  ListUserListPianosRequest,
  type PianoListKind,
} from "../lib/gen/piano_user_list/v1/piano_user_list_pb";
import { formatUser, parsePiano } from "../lib/resource-name";

type Props = {
  customId: string;
  listKind: PianoListKind;
  /** 親が更新トリガを変えると再フェッチする。 */
  refreshKey?: number;
  emptyMessage?: string;
};

export function UserPianoList({ customId, listKind, refreshKey, emptyMessage }: Props) {
  const [pianos, setPianos] = useState<Piano[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!customId) return;
    setLoading(true);
    setErr(null);
    pianoUserListClient
      .listUserListPianos(
        new ListUserListPianosRequest({
          parent: formatUser(customId),
          listKind,
          pageSize: 20,
        }),
      )
      .then((res) => setPianos(res.pianos))
      .catch((e) => setErr((e as Error)?.message || String(e)))
      .finally(() => setLoading(false));
  }, [customId, listKind, refreshKey]);

  if (err) {
    return <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p>;
  }
  if (loading) {
    return <p className="text-center text-xs text-slate-400">読み込み中...</p>;
  }
  if (pianos.length === 0) {
    return (
      <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        {emptyMessage ?? "ピアノがありません"}
      </p>
    );
  }
  return (
    <ul className="space-y-2">
      {pianos.map((p) => {
        const id = parsePianoId(p.name);
        return (
          <li key={p.name}>
            <Link
              to={id ? `/pianos/${id}` : "#"}
              className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white p-3 hover:border-slate-300"
            >
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-bold text-slate-900">{p.displayName}</p>
                {p.address ? (
                  <p className="truncate text-xs text-slate-500">{p.address}</p>
                ) : null}
              </div>
              {p.ratingCount > 0 ? (
                <div className="flex shrink-0 items-center gap-1 text-xs text-amber-600">
                  <Star size={12} className="fill-amber-500 text-amber-500" />
                  <span className="font-bold">{p.ratingAverage.toFixed(1)}</span>
                  <span className="text-slate-400">({p.ratingCount})</span>
                </div>
              ) : null}
            </Link>
          </li>
        );
      })}
    </ul>
  );
}

function parsePianoId(name: string): string | null {
  try {
    return parsePiano(name);
  } catch {
    return null;
  }
}
