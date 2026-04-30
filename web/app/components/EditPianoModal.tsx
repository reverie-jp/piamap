import { useEffect, useState } from "react";
import { FieldMask } from "@bufbuild/protobuf";

import { CenterDialog } from "./ui/sheet";
import { Button } from "./ui/button";
import { Select } from "./ui/select";
import { TextField } from "./ui/text-field";
import { pianoClient } from "../lib/api-client";
import {
  Piano,
  PianoAvailability,
  PianoKind,
  PianoType,
  UpdatePianoRequest,
} from "../lib/gen/piano/v1/piano_pb";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  piano: Piano;
  onSaved?: (updated: Piano) => void;
};

export function EditPianoModal({ isOpen, onOpenChange, piano, onSaved }: Props) {
  const [displayName, setDisplayName] = useState(piano.displayName);
  const [description, setDescription] = useState(piano.description ?? "");
  const [address, setAddress] = useState(piano.address ?? "");
  const [hours, setHours] = useState(piano.hours ?? "");
  const [pianoBrand, setPianoBrand] = useState(piano.pianoBrand);
  const [pianoModel, setPianoModel] = useState(piano.pianoModel ?? "");
  const [manufactureYear, setManufactureYear] = useState(
    piano.manufactureYear != null ? String(piano.manufactureYear) : "",
  );
  const [pianoType, setPianoType] = useState<PianoType>(piano.pianoType);
  const [kind, setKind] = useState<PianoKind>(piano.kind);
  const [availability, setAvailability] = useState<PianoAvailability>(piano.availability);
  const [availabilityNote, setAvailabilityNote] = useState(piano.availabilityNote ?? "");
  const [venueType, setVenueType] = useState(piano.venueType ?? "");
  const [editSummary, setEditSummary] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [fieldErr, setFieldErr] = useState<{ displayName?: string; manufactureYear?: string }>({});

  // 開くたびに渡された piano で初期化し直す。
  useEffect(() => {
    if (!isOpen) return;
    setDisplayName(piano.displayName);
    setDescription(piano.description ?? "");
    setAddress(piano.address ?? "");
    setHours(piano.hours ?? "");
    setPianoBrand(piano.pianoBrand);
    setPianoModel(piano.pianoModel ?? "");
    setManufactureYear(piano.manufactureYear != null ? String(piano.manufactureYear) : "");
    setPianoType(piano.pianoType);
    setKind(piano.kind);
    setAvailability(piano.availability);
    setAvailabilityNote(piano.availabilityNote ?? "");
    setVenueType(piano.venueType ?? "");
    setEditSummary("");
    setErr(null);
    setFieldErr({});
  }, [isOpen, piano]);

  const submit = async () => {
    const errors: typeof fieldErr = {};
    if (!displayName.trim()) errors.displayName = "名前を入力してください";
    let yearNum: number | undefined;
    if (manufactureYear.trim()) {
      const n = Number(manufactureYear);
      if (!Number.isFinite(n) || n < 1700 || n > 2100) {
        errors.manufactureYear = "1700–2100 の範囲で";
      } else {
        yearNum = Math.trunc(n);
      }
    }
    if (Object.keys(errors).length > 0) {
      setFieldErr(errors);
      return;
    }
    setFieldErr({});
    setSubmitting(true);
    setErr(null);

    // 差分のみ mask に乗せる。空文字は NULL クリアとして扱い、protobuf を未設定で送る。
    // サーバー側は SetX=true + 値=null を NULL として保存する (UpdatePiano は CASE WHEN ベース)。
    const paths: string[] = [];
    const next = new Piano({ name: piano.name });

    if (displayName.trim() !== piano.displayName) {
      paths.push("display_name");
      next.displayName = displayName.trim();
    }
    // optional 文字列 (空文字 = NULL クリア)
    const optStrChanged = (
      maskPath: string,
      oldV: string | undefined,
      newV: string,
      assign: (v: string | undefined) => void,
    ) => {
      const trimmed = newV.trim();
      const oldNorm = oldV ?? "";
      if (trimmed === oldNorm) return;
      paths.push(maskPath);
      assign(trimmed || undefined);
    };
    optStrChanged("description", piano.description, description, (v) => {
      if (v != null) next.description = v;
    });
    optStrChanged("address", piano.address, address, (v) => {
      if (v != null) next.address = v;
    });
    optStrChanged("hours", piano.hours, hours, (v) => {
      if (v != null) next.hours = v;
    });
    if (pianoBrand.trim() !== piano.pianoBrand) {
      paths.push("piano_brand");
      next.pianoBrand = pianoBrand.trim() || "unknown";
    }
    optStrChanged("piano_model", piano.pianoModel, pianoModel, (v) => {
      if (v != null) next.pianoModel = v;
    });
    const oldYear = piano.manufactureYear ?? undefined;
    if (yearNum !== oldYear) {
      paths.push("manufacture_year");
      if (yearNum != null) next.manufactureYear = yearNum;
    }
    if (pianoType !== piano.pianoType) {
      paths.push("piano_type");
      next.pianoType = pianoType;
    }
    if (kind !== piano.kind) {
      paths.push("kind");
      next.kind = kind;
    }
    if (availability !== piano.availability) {
      paths.push("availability");
      next.availability = availability;
    }
    optStrChanged("availability_note", piano.availabilityNote, availabilityNote, (v) => {
      if (v != null) next.availabilityNote = v;
    });
    optStrChanged("venue_type", piano.venueType, venueType, (v) => {
      if (v != null) next.venueType = v;
    });

    if (paths.length === 0) {
      setErr("変更がありません");
      setSubmitting(false);
      return;
    }

    try {
      const res = await pianoClient.updatePiano(
        new UpdatePianoRequest({
          piano: next,
          updateMask: new FieldMask({ paths }),
          editSummary: editSummary.trim() || undefined,
        }),
      );
      onOpenChange(false);
      if (res.piano) onSaved?.(res.piano);
    } catch (e) {
      setErr((e as Error)?.message || String(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <CenterDialog isOpen={isOpen} onOpenChange={onOpenChange} title="ピアノを編集">
      <div className="flex max-h-[70vh] flex-col">
        <div className="scrollbar-subtle flex-1 space-y-3 overflow-y-auto pr-1">
          <TextField
            label="名前"
            value={displayName}
            onChange={(v) => {
              setDisplayName(v);
              if (fieldErr.displayName) setFieldErr((s) => ({ ...s, displayName: undefined }));
            }}
            isRequired
            errorMessage={fieldErr.displayName}
          />
          <TextField label="説明" value={description} onChange={setDescription} multiline rows={3} />
          <TextField label="住所" value={address} onChange={setAddress} />
          <TextField label="営業時間" value={hours} onChange={setHours} />

          <div className="grid grid-cols-2 gap-2">
            <Select
              label="種別"
              selectedKey={String(kind)}
              onSelectionChange={(k) => setKind(Number(k) as PianoKind)}
              items={[
                { id: String(PianoKind.STREET), label: "ストリート" },
                { id: String(PianoKind.PRACTICE_ROOM), label: "練習室" },
                { id: String(PianoKind.OTHER), label: "その他" },
              ]}
            />
            <Select
              label="ピアノ"
              selectedKey={String(pianoType)}
              onSelectionChange={(k) => setPianoType(Number(k) as PianoType)}
              items={[
                { id: String(PianoType.GRAND), label: "グランド" },
                { id: String(PianoType.UPRIGHT), label: "アップライト" },
                { id: String(PianoType.ELECTRONIC), label: "電子" },
                { id: String(PianoType.UNKNOWN), label: "不明" },
              ]}
            />
          </div>

          <TextField
            label="メーカー"
            value={pianoBrand}
            onChange={setPianoBrand}
            placeholder="ヤマハ、カワイ など"
          />
          <div className="grid grid-cols-2 gap-2">
            <TextField label="モデル" value={pianoModel} onChange={setPianoModel} />
            <TextField
              label="製造年"
              value={manufactureYear}
              onChange={(v) => {
                setManufactureYear(v);
                if (fieldErr.manufactureYear)
                  setFieldErr((s) => ({ ...s, manufactureYear: undefined }));
              }}
              placeholder="2010"
              errorMessage={fieldErr.manufactureYear}
            />
          </div>

          <Select
            label="営業状況"
            selectedKey={String(availability)}
            onSelectionChange={(k) => setAvailability(Number(k) as PianoAvailability)}
            items={[
              { id: String(PianoAvailability.REGULAR), label: "通年" },
              { id: String(PianoAvailability.IRREGULAR), label: "不定期" },
              { id: String(PianoAvailability.EVENT_ONLY), label: "イベント時のみ" },
              { id: String(PianoAvailability.WEATHER_DEPENDENT), label: "天候次第" },
            ]}
          />
          <TextField
            label="営業状況メモ"
            value={availabilityNote}
            onChange={setAvailabilityNote}
            multiline
            rows={2}
          />
          <TextField
            label="会場種別"
            value={venueType}
            onChange={setVenueType}
            placeholder="駅、商業施設、空港 など"
          />

          <TextField
            label="編集メモ"
            value={editSummary}
            onChange={setEditSummary}
            placeholder="どこを更新したか (任意)"
          />

          {err ? <p className="text-xs text-rose-600">{err}</p> : null}
        </div>

        <div className="mt-3 flex flex-col gap-2 border-slate-200 bg-white pt-3">
          <Button onPress={submit} isPending={submitting}>
            保存する
          </Button>
          <Button variant="ghost" size="sm" onPress={() => onOpenChange(false)}>
            キャンセル
          </Button>
        </div>
      </div>
    </CenterDialog>
  );
}
