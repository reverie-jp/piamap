import { useEffect, useState } from "react";
import { Link } from "react-router";
import { Bookmark, ChevronRight, Heart, MapPin } from "lucide-react";

import { pianoUserListClient } from "../lib/api-client";
import {
  ListUserListPianosRequest,
  PianoListKind,
} from "../lib/gen/piano_user_list/v1/piano_user_list_pb";
import { formatUser } from "../lib/resource-name";

const LIST_DEFS: {
  kind: PianoListKind;
  slug: string;
  label: string;
  description: string;
  icon: React.ComponentType<{ size?: number; className?: string }>;
}[] = [
  {
    kind: PianoListKind.WISHLIST,
    slug: "wishlist",
    label: "行ってみたい",
    description: "気になるピアノをマーク",
    icon: MapPin,
  },
  {
    kind: PianoListKind.VISITED,
    slug: "visited",
    label: "行ったことある",
    description: "投稿で自動的に追加",
    icon: Bookmark,
  },
  {
    kind: PianoListKind.FAVORITE,
    slug: "favorite",
    label: "お気に入り",
    description: "また弾きたいピアノ",
    icon: Heart,
  },
];

type Props = {
  customId: string;
};

export function SavedListsCards({ customId }: Props) {
  const [counts, setCounts] = useState<Partial<Record<PianoListKind, number>>>({});

  useEffect(() => {
    if (!customId) return;
    let cancelled = false;
    const reqs = LIST_DEFS.map((d) =>
      pianoUserListClient
        .listUserListPianos(
          new ListUserListPianosRequest({
            parent: formatUser(customId),
            listKind: d.kind,
            pageSize: 100,
          }),
        )
        .then((res) => [d.kind, res.pianos.length] as const)
        .catch(() => [d.kind, 0] as const),
    );
    Promise.all(reqs).then((rows) => {
      if (cancelled) return;
      const next: Partial<Record<PianoListKind, number>> = {};
      for (const [k, n] of rows) next[k] = n;
      setCounts(next);
    });
    return () => {
      cancelled = true;
    };
  }, [customId]);

  return (
    <ul className="space-y-2">
      {LIST_DEFS.map((d) => {
        const Icon = d.icon;
        const count = counts[d.kind];
        return (
          <li key={d.slug}>
            <Link
              to={`/profile/${customId}/saved/${d.slug}`}
              className="flex items-center gap-3 rounded-2xl border border-slate-200 bg-white p-3 transition hover:border-amber-300 hover:bg-amber-50/30"
            >
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-amber-100 text-amber-600">
                <Icon size={18} />
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-bold text-slate-900">{d.label}</p>
                <p className="text-xs text-slate-500">{d.description}</p>
              </div>
              {count !== undefined ? (
                <span className="shrink-0 text-sm font-bold text-slate-700">{count}</span>
              ) : null}
              <ChevronRight size={18} className="shrink-0 text-slate-400" />
            </Link>
          </li>
        );
      })}
    </ul>
  );
}
