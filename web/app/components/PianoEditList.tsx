import { useEffect, useState } from "react";
import { Link } from "react-router";
import { ArrowRight } from "lucide-react";

import { pianoClient } from "../lib/api-client";
import {
  ListPianoEditsRequest,
  PianoEditOperation,
  type PianoEdit,
} from "../lib/gen/piano/v1/piano_pb";
import { formatPiano } from "../lib/resource-name";

type Props = {
  pianoId: string;
  /** ピアノ更新時にこの値を変えると再フェッチ。 */
  refreshKey?: number;
};

export function PianoEditList({ pianoId, refreshKey }: Props) {
  const [edits, setEdits] = useState<PianoEdit[]>([]);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!pianoId) return;
    setErr(null);
    pianoClient
      .listPianoEdits(
        new ListPianoEditsRequest({ parent: formatPiano(pianoId), pageSize: 20 }),
      )
      .then((res) => setEdits(res.edits))
      .catch((e) => setErr((e as Error)?.message || String(e)));
  }, [pianoId, refreshKey]);

  if (err) {
    return <p className="rounded bg-rose-50 p-3 text-sm text-rose-700">{err}</p>;
  }
  if (edits.length === 0) {
    return (
      <p className="rounded-2xl bg-slate-50 p-4 text-center text-sm text-slate-500">
        編集履歴がありません
      </p>
    );
  }
  return (
    <ul className="space-y-2">
      {edits.map((e) => (
        <li
          key={e.name}
          className="rounded-xl border border-slate-200 bg-white p-3 text-sm"
        >
          <div className="flex flex-wrap items-center gap-2 text-xs text-slate-500">
            <span
              className={`rounded-full px-2 py-0.5 text-[10px] font-bold ${opStyle(
                e.operation,
              )}`}
            >
              {opLabel(e.operation)}
            </span>
            {e.editor ? (
              <Link
                to={`/profile/${e.editor.split("/").pop()}`}
                className="font-medium text-slate-700 hover:underline"
              >
                {e.editorDisplayName || `@${e.editor.split("/").pop()}`}
              </Link>
            ) : (
              <span className="text-slate-400">削除済みユーザー</span>
            )}
            <span>{formatTime(e.createTime)}</span>
          </div>
          {e.summary ? (
            <p className="mt-1 text-slate-800 whitespace-pre-line">{e.summary}</p>
          ) : null}
          {e.operation === PianoEditOperation.CREATE ? (
            <p className="mt-2 text-xs text-slate-500">ピアノを新規登録しました</p>
          ) : e.changesJson ? (
            <ChangesDetail json={e.changesJson} />
          ) : null}
        </li>
      ))}
    </ul>
  );
}

type Diff = {
  field: string;
  old: unknown;
  new: unknown;
};

function ChangesDetail({ json }: { json: string }) {
  const diffs = parseDiffs(json);
  if (diffs.length === 0) return null;
  return (
    <ul className="mt-2 space-y-1 text-xs">
      {diffs.map((d) => (
        <li
          key={d.field}
          className="flex flex-wrap items-center gap-1.5 rounded-md bg-slate-50 px-2 py-1"
        >
          <span className="font-medium text-slate-600">{fieldLabel(d.field)}:</span>
          <ValueChip value={formatValue(d.field, d.old)} variant="old" />
          <ArrowRight size={12} className="text-slate-400" aria-hidden />
          <ValueChip value={formatValue(d.field, d.new)} variant="new" />
        </li>
      ))}
    </ul>
  );
}

function ValueChip({ value, variant }: { value: string; variant: "old" | "new" }) {
  if (!value) {
    return <span className="text-slate-400 italic">(なし)</span>;
  }
  const cls =
    variant === "new"
      ? "bg-amber-100 text-amber-800"
      : "bg-slate-200 text-slate-700 line-through";
  return (
    <span className={`max-w-[14rem] truncate rounded px-1.5 py-0.5 ${cls}`} title={value}>
      {value}
    </span>
  );
}

function parseDiffs(json: string): Diff[] {
  try {
    const v = JSON.parse(json);
    if (!v || typeof v !== "object") return [];
    return Object.entries(v as Record<string, unknown>)
      .filter(([, change]) => change && typeof change === "object")
      .map(([field, change]) => {
        const c = change as Record<string, unknown>;
        return { field, old: c.old, new: c.new };
      });
  } catch {
    return [];
  }
}

function formatValue(field: string, value: unknown): string {
  if (value == null || value === "") return "";
  if (field === "location" && typeof value === "object") {
    const v = value as { lat?: number; lng?: number };
    if (typeof v.lat === "number" && typeof v.lng === "number") {
      return `${v.lat.toFixed(5)}, ${v.lng.toFixed(5)}`;
    }
  }
  switch (field) {
    case "kind":
      return kindLabel(String(value));
    case "piano_type":
      return pianoTypeLabel(String(value));
    case "availability":
      return availabilityLabel(String(value));
    case "status":
      return statusLabel(String(value));
  }
  if (typeof value === "string") return value;
  if (typeof value === "number") return String(value);
  return JSON.stringify(value);
}

function kindLabel(s: string): string {
  switch (s) {
    case "street":
      return "ストリート";
    case "practice_room":
      return "練習室";
    case "other":
      return "その他";
    default:
      return s;
  }
}

function pianoTypeLabel(s: string): string {
  switch (s) {
    case "grand":
      return "グランド";
    case "upright":
      return "アップライト";
    case "electronic":
      return "電子";
    case "unknown":
      return "不明";
    default:
      return s;
  }
}

function availabilityLabel(s: string): string {
  switch (s) {
    case "regular":
      return "通年";
    case "irregular":
      return "不定期";
    case "event_only":
      return "イベント時のみ";
    case "weather_dependent":
      return "天候次第";
    default:
      return s;
  }
}

function statusLabel(s: string): string {
  switch (s) {
    case "pending":
      return "申請中";
    case "active":
      return "公開中";
    case "temporary":
      return "一時公開";
    case "removed":
      return "削除";
    default:
      return s;
  }
}

function opLabel(op: PianoEditOperation): string {
  switch (op) {
    case PianoEditOperation.CREATE:
      return "新規登録";
    case PianoEditOperation.UPDATE:
      return "編集";
    case PianoEditOperation.PHOTO_ADD:
      return "写真追加";
    case PianoEditOperation.PHOTO_REMOVE:
      return "写真削除";
    case PianoEditOperation.STATUS_CHANGE:
      return "状態変更";
    case PianoEditOperation.KIND_CHANGE:
      return "種別変更";
    case PianoEditOperation.RESTORE:
      return "復元";
    default:
      return "—";
  }
}

function opStyle(op: PianoEditOperation): string {
  switch (op) {
    case PianoEditOperation.CREATE:
      return "bg-emerald-100 text-emerald-700";
    case PianoEditOperation.STATUS_CHANGE:
    case PianoEditOperation.KIND_CHANGE:
      return "bg-amber-100 text-amber-700";
    case PianoEditOperation.PHOTO_REMOVE:
      return "bg-rose-100 text-rose-700";
    case PianoEditOperation.RESTORE:
      return "bg-sky-100 text-sky-700";
    default:
      return "bg-slate-100 text-slate-700";
  }
}

function fieldLabel(k: string): string {
  switch (k) {
    case "name":
      return "名前";
    case "description":
      return "説明";
    case "address":
      return "住所";
    case "location":
      return "位置";
    case "kind":
      return "種別";
    case "status":
      return "状態";
    case "piano_brand":
      return "メーカー";
    case "piano_model":
      return "モデル";
    case "piano_type":
      return "ピアノ種類";
    case "manufacture_year":
      return "製造年";
    case "hours":
      return "営業時間";
    case "availability":
      return "営業状況";
    case "availability_note":
      return "営業ノート";
    case "venue_type":
      return "会場種別";
    case "prefecture":
      return "都道府県";
    case "city":
      return "市区町村";
    case "install_time":
      return "設置日時";
    case "remove_time":
      return "撤去日時";
    default:
      return k;
  }
}

function formatTime(ts: PianoEdit["createTime"]): string {
  if (!ts) return "";
  const d = ts.toDate();
  const yyyy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  const hh = String(d.getHours()).padStart(2, "0");
  const mi = String(d.getMinutes()).padStart(2, "0");
  return `${yyyy}/${mm}/${dd} ${hh}:${mi}`;
}
