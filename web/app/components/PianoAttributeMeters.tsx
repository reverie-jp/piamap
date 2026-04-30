import type { Piano } from "../lib/gen/piano/v1/piano_pb";

type Props = {
  piano: Piano;
  /** ヘッダ枠やラベルを省く軽量表示 (マップのボトムシート用)。 */
  compact?: boolean;
};

export function PianoAttributeMeters({ piano, compact }: Props) {
  const items: { label: string; value: number }[] = [
    { label: "賑やかさ", value: piano.ambientNoiseAverage },
    { label: "人通り", value: piano.footTrafficAverage },
    { label: "響き", value: piano.resonanceAverage },
    { label: "鍵盤の重さ", value: piano.keyTouchWeightAverage },
    { label: "調律状態", value: piano.tuningQualityAverage },
  ].filter((i) => i.value > 0);
  if (items.length === 0) return null;

  const ul = (
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
  );

  if (compact) return ul;

  return (
    <div className="rounded-2xl border border-slate-200 p-4">
      <h3 className="mb-3 text-xs font-bold uppercase tracking-wide text-slate-500">特徴</h3>
      {ul}
    </div>
  );
}
