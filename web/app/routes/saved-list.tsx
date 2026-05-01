import { Link, useNavigate, useParams } from "react-router";
import { ArrowLeft } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { UserPianoList } from "../components/UserPianoList";
import { PianoListKind } from "../lib/gen/piano_user_list/v1/piano_user_list_pb";

import type { Route } from "./+types/saved-list";

export function meta({ params }: Route.MetaArgs) {
  const def = SLUG_TO_DEF[params.kind ?? ""];
  return [{ title: `${def?.label ?? "保存済み"} — PiaMap` }];
}

const SLUG_TO_DEF: Record<
  string,
  { kind: PianoListKind; label: string; emptyMessage: string }
> = {
  wishlist: {
    kind: PianoListKind.WISHLIST,
    label: "行ってみたい",
    emptyMessage: "行ってみたいピアノはまだありません",
  },
  visited: {
    kind: PianoListKind.VISITED,
    label: "行ったことある",
    emptyMessage: "行ったピアノはまだありません",
  },
  favorite: {
    kind: PianoListKind.FAVORITE,
    label: "お気に入り",
    emptyMessage: "お気に入りに登録したピアノはまだありません",
  },
};

export default function SavedList() {
  const { customId = "", kind: slug = "" } = useParams();
  const navigate = useNavigate();
  const def = SLUG_TO_DEF[slug];

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-3">
        <button
          type="button"
          onClick={() => navigate(-1)}
          aria-label="戻る"
          className="cursor-pointer text-slate-500 hover:text-slate-700"
        >
          <ArrowLeft size={22} />
        </button>
        <h1 className="text-base font-bold text-slate-900">{def?.label ?? "保存済み"}</h1>
      </header>

      {!def ? (
        <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
          存在しないリストです。
          <Link to={`/profile/${customId}`} className="ml-2 text-amber-700 underline">
            プロフィールへ
          </Link>
        </p>
      ) : !customId ? (
        <p className="text-center text-xs text-slate-400">読み込み中...</p>
      ) : (
        <UserPianoList
          customId={customId}
          listKind={def.kind}
          emptyMessage={def.emptyMessage}
        />
      )}
    </MobileShell>
  );
}
