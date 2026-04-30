import { useEffect, useState } from "react";
import { Star } from "lucide-react";
import { FieldMask, Timestamp } from "@bufbuild/protobuf";

import { CenterDialog } from "./ui/sheet";
import { Button } from "./ui/button";
import { TextField } from "./ui/text-field";
import { pianoPostClient } from "../lib/api-client";
import {
  CreatePianoPostRequest,
  PianoPost,
  PostVisibility,
  UpdatePianoPostRequest,
} from "../lib/gen/piano_post/v1/piano_post_pb";
import { formatPiano } from "../lib/resource-name";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  pianoId: string;
  pianoName: string;
  /** 編集モード時に渡す既存投稿。未指定なら新規作成。 */
  editingPost?: PianoPost;
  /** 作成または更新成功時に親側で再フェッチさせるためのハンドラ。 */
  onSaved?: () => void;
};

type Attr = {
  key: "ambientNoise" | "footTraffic" | "resonance" | "keyTouchWeight" | "tuningQuality";
  leftLabel: string;
  rightLabel: string;
};

const ATTRS: Attr[] = [
  { key: "ambientNoise", leftLabel: "静か", rightLabel: "賑やか" },
  { key: "footTraffic", leftLabel: "人通り少", rightLabel: "人通り多" },
  { key: "resonance", leftLabel: "響き弱", rightLabel: "響き豊か" },
  { key: "keyTouchWeight", leftLabel: "鍵盤軽い", rightLabel: "鍵盤重い" },
  { key: "tuningQuality", leftLabel: "調律悪い", rightLabel: "調律良い" },
];

type Attrs = Partial<Record<Attr["key"], number>>;

const todayISO = () => new Date().toISOString().slice(0, 10);

function hydrateAttrs(post: PianoPost | undefined): Attrs {
  if (!post) return {};
  const a: Attrs = {};
  if (post.ambientNoise != null) a.ambientNoise = post.ambientNoise;
  if (post.footTraffic != null) a.footTraffic = post.footTraffic;
  if (post.resonance != null) a.resonance = post.resonance;
  if (post.keyTouchWeight != null) a.keyTouchWeight = post.keyTouchWeight;
  if (post.tuningQuality != null) a.tuningQuality = post.tuningQuality;
  return a;
}

export function CreatePianoPostModal({
  isOpen,
  onOpenChange,
  pianoId,
  pianoName,
  editingPost,
  onSaved,
}: Props) {
  const isEditing = Boolean(editingPost);

  const [rating, setRating] = useState<number>(0);
  const [body, setBody] = useState("");
  const [visitDate, setVisitDate] = useState(todayISO);
  const [attrs, setAttrs] = useState<Attrs>({});
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  // モーダルを開いたタイミングで初期値を入れる (編集なら既存値、新規ならクリア)。
  useEffect(() => {
    if (!isOpen) return;
    if (editingPost) {
      setRating(editingPost.rating);
      setBody(editingPost.body ?? "");
      setVisitDate(
        editingPost.visitTime
          ? editingPost.visitTime.toDate().toISOString().slice(0, 10)
          : todayISO(),
      );
      setAttrs(hydrateAttrs(editingPost));
    } else {
      setRating(0);
      setBody("");
      setVisitDate(todayISO());
      setAttrs({});
    }
    setErr(null);
  }, [isOpen, editingPost]);

  const handleSubmit = async () => {
    if (rating < 1 || rating > 5) {
      setErr("評価を選択してください");
      return;
    }
    setSubmitting(true);
    setErr(null);
    try {
      const visit = new Date(visitDate);
      if (Number.isNaN(visit.getTime())) {
        setErr("訪問日が不正です");
        return;
      }
      if (isEditing && editingPost) {
        // 全フィールドを mask に載せ、未入力は undefined のままサーバーに送る。
        // サーバー側は SetX=true + 値=null を NULL として保存するため、これで完全クリアできる。
        const post = new PianoPost({
          name: editingPost.name,
          rating,
          visitTime: Timestamp.fromDate(visit),
        });
        const trimmed = body.trim();
        if (trimmed) post.body = trimmed;
        if (attrs.ambientNoise != null) post.ambientNoise = attrs.ambientNoise;
        if (attrs.footTraffic != null) post.footTraffic = attrs.footTraffic;
        if (attrs.resonance != null) post.resonance = attrs.resonance;
        if (attrs.keyTouchWeight != null) post.keyTouchWeight = attrs.keyTouchWeight;
        if (attrs.tuningQuality != null) post.tuningQuality = attrs.tuningQuality;
        await pianoPostClient.updatePianoPost(
          new UpdatePianoPostRequest({
            pianoPost: post,
            updateMask: new FieldMask({
              paths: [
                "rating",
                "visit_time",
                "body",
                "ambient_noise",
                "foot_traffic",
                "resonance",
                "key_touch_weight",
                "tuning_quality",
              ],
            }),
          }),
        );
      } else {
        const post = new PianoPost({
          rating,
          visitTime: Timestamp.fromDate(visit),
          body: body.trim() || undefined,
          visibility: PostVisibility.PUBLIC,
          ambientNoise: attrs.ambientNoise,
          footTraffic: attrs.footTraffic,
          resonance: attrs.resonance,
          keyTouchWeight: attrs.keyTouchWeight,
          tuningQuality: attrs.tuningQuality,
        });
        await pianoPostClient.createPianoPost(
          new CreatePianoPostRequest({ parent: formatPiano(pianoId), pianoPost: post }),
        );
      }
      onOpenChange(false);
      onSaved?.();
    } catch (e) {
      setErr((e as Error)?.message || String(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <CenterDialog isOpen={isOpen} onOpenChange={onOpenChange} title={pianoName}>
      <div className="flex max-h-[70vh] flex-col">
        <div className="scrollbar-subtle flex-1 space-y-4 overflow-y-auto pr-1">
          <div>
            <label className="text-sm font-medium text-slate-700">
              評価<span className="ml-0.5 text-rose-600">*</span>
            </label>
            <div className="mt-2 flex items-center gap-1">
              {[1, 2, 3, 4, 5].map((n) => (
                <button
                  key={n}
                  type="button"
                  aria-label={`${n}つ星`}
                  onClick={() => setRating(n)}
                  className="flex flex-1 cursor-pointer items-center justify-center rounded p-1 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                >
                  <Star
                    size={36}
                    className={n <= rating ? "fill-amber-500 text-amber-500" : "text-slate-300"}
                  />
                </button>
              ))}
            </div>
          </div>

          <TextField
            label="訪問日"
            isRequired
            value={visitDate}
            onChange={setVisitDate}
            type="date"
          />

          <TextField
            label="感想"
            multiline
            rows={4}
            value={body}
            onChange={setBody}
            placeholder="どんな雰囲気だった?演奏したらどうだった?"
          />

          <div>
            <label className="text-sm font-medium text-slate-700">ピアノの特徴</label>
            <ul className="mt-2 space-y-3">
              {ATTRS.map((a) => {
                const value = attrs[a.key];
                return (
                  <li key={a.key}>
                    <div className="flex justify-between text-xs text-slate-600">
                      <span>{a.leftLabel}</span>
                      <span>{a.rightLabel}</span>
                    </div>
                    <div className="mt-1 flex items-center gap-1">
                      {[1, 2, 3, 4, 5].map((n) => {
                        const selected = value === n;
                        return (
                          <button
                            key={n}
                            type="button"
                            aria-label={`${a.leftLabel}—${a.rightLabel} ${n}`}
                            aria-pressed={selected}
                            onClick={() =>
                              setAttrs((s) => {
                                if (s[a.key] === n) {
                                  const next = { ...s };
                                  delete next[a.key];
                                  return next;
                                }
                                return { ...s, [a.key]: n };
                              })
                            }
                            className={
                              "h-8 flex-1 cursor-pointer rounded-md border text-xs font-semibold outline-none " +
                              "focus-visible:ring-2 focus-visible:ring-amber-500 " +
                              (selected
                                ? "border-amber-500 bg-amber-500 text-white"
                                : "border-slate-300 bg-white text-slate-600 hover:border-slate-400")
                            }
                          >
                            {n}
                          </button>
                        );
                      })}
                    </div>
                  </li>
                );
              })}
            </ul>
          </div>

          {err ? <p className="text-xs text-rose-600">{err}</p> : null}
        </div>

        <div className="mt-3 flex flex-col gap-2 border-slate-200 bg-white pt-3">
          <Button onPress={handleSubmit} isPending={submitting}>
            {isEditing ? "保存する" : "投稿する"}
          </Button>
          <Button variant="ghost" size="sm" onPress={() => onOpenChange(false)}>
            キャンセル
          </Button>
        </div>
      </div>
    </CenterDialog>
  );
}
