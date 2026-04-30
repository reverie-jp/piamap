import { useEffect, useRef, useState } from "react";
import { Link } from "react-router";
import { Clock, MapPin, Search, SlidersHorizontal, X } from "lucide-react";

import { pianoClient } from "../lib/api-client";
import {
  GetPianoRequest,
  Piano,
  PianoType,
  SearchPianosRequest,
} from "../lib/gen/piano/v1/piano_pb";
import { formatPiano, parsePiano } from "../lib/resource-name";
import {
  clearSearchHistory,
  getSearchHistory,
  pushSearchHistory,
  removeSearchHistoryItem,
  type SearchHistoryItem,
} from "../lib/search-history";
import { Button } from "./ui/button";
import { Select } from "./ui/select";

// 5 環境属性の最低平均値フィルタのキー (proto の min_*_average に対応)。
const ATTRIBUTE_KEYS = [
  "minAmbientNoiseAverage",
  "minFootTrafficAverage",
  "minResonanceAverage",
  "minKeyTouchWeightAverage",
  "minTuningQualityAverage",
] as const;
type AttributeKey = (typeof ATTRIBUTE_KEYS)[number];

export type MapSearchFilters = {
  pianoType?: PianoType;
  pianoBrand?: string;
  minRatingAverage?: number;
} & Partial<Record<AttributeKey, number>>;

export const EMPTY_FILTERS: MapSearchFilters = {};

export function activeFilterCount(f: MapSearchFilters): number {
  let n = 0;
  if (f.pianoType !== undefined) n += 1;
  if (f.pianoBrand !== undefined) n += 1;
  if (f.minRatingAverage !== undefined) n += 1;
  for (const k of ATTRIBUTE_KEYS) {
    if (f[k] !== undefined) n += 1;
  }
  return n;
}

const PIANO_TYPE_ITEMS = [
  { id: "any", label: "すべて" },
  { id: String(PianoType.GRAND), label: "グランド" },
  { id: String(PianoType.UPRIGHT), label: "アップライト" },
  { id: String(PianoType.ELECTRONIC), label: "電子" },
  { id: String(PianoType.UNKNOWN), label: "不明" },
];

const MIN_RATING_ITEMS = [
  { id: "any", label: "すべて" },
  { id: "3", label: "★3 以上" },
  { id: "4", label: "★4 以上" },
  { id: "4.5", label: "★4.5 以上" },
];

// 5 環境属性の絞り込み選択肢 (1=控えめ ... 5=強い)。
// 例: 響きなら "響き 4以上" は resonance_average >= 4.0 の意。
const ATTRIBUTE_SCORE_ITEMS = [
  { id: "any", label: "すべて" },
  { id: "2", label: "2 以上" },
  { id: "3", label: "3 以上" },
  { id: "4", label: "4 以上" },
  { id: "4.5", label: "4.5 以上" },
];

// 各属性のラベルとスコアの方向性 (1 ↔ 5)。
const ATTRIBUTE_LABELS: { key: AttributeKey; label: string; lo: string; hi: string }[] = [
  { key: "minResonanceAverage", label: "響き", lo: "弱い", hi: "豊か" },
  { key: "minTuningQualityAverage", label: "調律", lo: "悪い", hi: "良い" },
  { key: "minKeyTouchWeightAverage", label: "鍵盤の重さ", lo: "軽い", hi: "重い" },
  { key: "minAmbientNoiseAverage", label: "周囲音", lo: "静か", hi: "賑やか" },
  { key: "minFootTrafficAverage", label: "人通り", lo: "少ない", hi: "多い" },
];

// よく登録される主要メーカー。値は piano_brand カラムに格納される文字列。
// マッチは大文字小文字無視 (ILIKE) なので "yamaha"/"YAMAHA" などはどれでも当たる。
const PIANO_BRAND_ITEMS = [
  { id: "any", label: "すべて" },
  { id: "YAMAHA", label: "YAMAHA" },
  { id: "KAWAI", label: "KAWAI" },
  { id: "STEINWAY", label: "STEINWAY" },
  { id: "BOSENDORFER", label: "BÖSENDORFER" },
  { id: "BECHSTEIN", label: "BECHSTEIN" },
  { id: "BLUTHNER", label: "BLÜTHNER" },
  { id: "FAZIOLI", label: "FAZIOLI" },
  { id: "BOSTON", label: "BOSTON" },
  { id: "DIAPASON", label: "DIAPASON" },
  { id: "ROLAND", label: "ROLAND" },
  { id: "KORG", label: "KORG" },
  { id: "CASIO", label: "CASIO" },
  { id: "unknown", label: "不明" },
];

type Props = {
  filters: MapSearchFilters;
  onFiltersChange: (next: MapSearchFilters) => void;
  /** 検索結果クリック時に呼ぶ。指定なしなら結果は Link で遷移。 */
  onSelect?: (piano: Piano) => void;
  /** 値が変わるたびに query をクリアする (例: ピアノ詳細シートが閉じた時)。 */
  clearSignal?: number;
};

export function MapSearchBar({ filters, onFiltersChange, onSelect, clearSignal }: Props) {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Piano[]>([]);
  const [loading, setLoading] = useState(false);
  const [focused, setFocused] = useState(false);
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [history, setHistory] = useState<SearchHistoryItem[]>([]);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);

  // 初回マウント時に localStorage から履歴を読む (SSR セーフ)。
  useEffect(() => {
    setHistory(getSearchHistory());
  }, []);

  // clearSignal が変わったら query をクリア (シート閉じなど親の合図)。
  // 初回マウントは無視。
  const clearSignalMounted = useRef(false);
  useEffect(() => {
    if (!clearSignalMounted.current) {
      clearSignalMounted.current = true;
      return;
    }
    setQuery("");
  }, [clearSignal]);

  const filterCount = activeFilterCount(filters);

  // デバウンス検索 (250ms)。フィルタ変更でも再検索する。
  useEffect(() => {
    const trimmed = query.trim();
    if (!trimmed) {
      setResults([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    const t = setTimeout(async () => {
      try {
        const req = new SearchPianosRequest({ query: trimmed, limit: 20 });
        if (filters.pianoType !== undefined) req.pianoType = filters.pianoType;
        if (filters.pianoBrand !== undefined) req.pianoBrand = filters.pianoBrand;
        if (filters.minRatingAverage !== undefined) req.minRatingAverage = filters.minRatingAverage;
        for (const k of ATTRIBUTE_KEYS) {
          const v = filters[k];
          if (v !== undefined) req[k] = v;
        }
        const res = await pianoClient.searchPianos(req);
        setResults(res.pianos);
      } catch (e) {
        console.error("searchPianos (text) failed", e);
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 250);
    return () => clearTimeout(t);
    // ESLint exhaustive-deps: filters のキー全部を見る代わりに JSON で簡略化。
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query, JSON.stringify(filters)]);

  // 外側クリックで dropdown / filter パネルを閉じる。
  useEffect(() => {
    if (!focused && !filtersOpen) return;
    const onDocClick = (e: MouseEvent) => {
      if (!containerRef.current) return;
      if (!containerRef.current.contains(e.target as Node)) {
        setFocused(false);
        setFiltersOpen(false);
      }
    };
    document.addEventListener("mousedown", onDocClick);
    return () => document.removeEventListener("mousedown", onDocClick);
  }, [focused, filtersOpen]);

  const showDropdown = focused && query.trim().length > 0 && !filtersOpen;
  const showHistory =
    focused && query.trim().length === 0 && !filtersOpen && history.length > 0;

  // 検索結果クリック時に履歴へ追加。詳細シートやページへ遷移する直前に呼ぶ。
  const recordHistory = (p: Piano, pianoId: string) => {
    if (!pianoId) return;
    const subtitle = p.address || p.venueType || p.prefecture || p.city || "";
    setHistory(
      pushSearchHistory({
        pianoId,
        displayName: p.displayName || pianoId,
        subtitle: subtitle || undefined,
      }),
    );
  };

  // 履歴クリック時: マップが表示されている文脈なら Piano を取得して onSelect に渡し
  // (flyTo + シート表示)、そうでなければ詳細ページへ遷移する。
  const handleHistoryClick = async (item: SearchHistoryItem) => {
    if (!onSelect) return; // Link でカバー
    // タップしたピアノ名を入力欄にフィルしておく (再検索の起点や閉じても残る)。
    setQuery(item.displayName);
    try {
      const res = await pianoClient.getPiano(
        new GetPianoRequest({ name: formatPiano(item.pianoId) }),
      );
      if (res.piano) {
        onSelect(res.piano);
        recordHistory(res.piano, item.pianoId);
      }
    } catch (e) {
      console.error("getPiano (history) failed", e);
    } finally {
      setFocused(false);
      inputRef.current?.blur();
    }
  };

  const handlePianoTypeChange = (id: string) => {
    if (id === "any") {
      const { pianoType: _omit, ...rest } = filters;
      onFiltersChange(rest);
    } else {
      onFiltersChange({ ...filters, pianoType: Number(id) as PianoType });
    }
  };

  const handleMinRatingChange = (id: string) => {
    if (id === "any") {
      const { minRatingAverage: _omit, ...rest } = filters;
      onFiltersChange(rest);
    } else {
      onFiltersChange({ ...filters, minRatingAverage: Number(id) });
    }
  };

  const handlePianoBrandChange = (id: string) => {
    if (id === "any") {
      const { pianoBrand: _omit, ...rest } = filters;
      onFiltersChange(rest);
    } else {
      onFiltersChange({ ...filters, pianoBrand: id });
    }
  };

  const handleAttributeChange = (key: AttributeKey, id: string) => {
    if (id === "any") {
      const next = { ...filters };
      delete next[key];
      onFiltersChange(next);
    } else {
      onFiltersChange({ ...filters, [key]: Number(id) });
    }
  };

  return (
    <div
      ref={containerRef}
      className="absolute left-3 right-3 top-3 z-20 mx-auto max-w-xl"
    >
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Search
            size={16}
            aria-hidden
            className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
          />
          <input
            ref={inputRef}
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onFocus={() => {
              setFocused(true);
              setFiltersOpen(false);
            }}
            placeholder="ピアノ名で検索"
            className="h-11 w-full rounded-full border border-slate-200 bg-white/95 pl-9 pr-9 text-sm text-slate-900 shadow-md outline-none backdrop-blur placeholder:text-slate-400 focus:border-amber-500 focus:ring-1 focus:ring-amber-500"
          />
          {query ? (
            <button
              type="button"
              aria-label="クリア"
              onClick={() => {
                setQuery("");
                setResults([]);
                inputRef.current?.focus();
              }}
              className="absolute right-2 top-1/2 -translate-y-1/2 cursor-pointer rounded-full p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
            >
              <X size={14} />
            </button>
          ) : null}
        </div>
        <button
          type="button"
          aria-label={filtersOpen ? "フィルタを閉じる" : "フィルタを開く"}
          aria-pressed={filtersOpen}
          onClick={() => {
            setFiltersOpen((v) => !v);
            setFocused(false);
          }}
          className={
            "relative inline-flex h-11 w-11 shrink-0 cursor-pointer items-center justify-center rounded-full border bg-white/95 text-slate-700 shadow-md outline-none backdrop-blur transition focus-visible:ring-2 focus-visible:ring-amber-500 " +
            (filtersOpen || filterCount > 0
              ? "border-amber-500 text-amber-700"
              : "border-slate-200 hover:bg-slate-50")
          }
        >
          <SlidersHorizontal size={16} />
          {filterCount > 0 ? (
            <span className="absolute -top-1 -right-1 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-amber-500 px-1 text-[10px] font-bold text-white">
              {filterCount}
            </span>
          ) : null}
        </button>
      </div>

      {filtersOpen ? (
        <div className="mt-2 rounded-2xl border border-slate-200 bg-white shadow-lg">
          <div className="max-h-[60vh] overflow-y-auto p-3">
            <Select
              label="ピアノ種別"
              placeholder="すべて"
              items={PIANO_TYPE_ITEMS}
              selectedKey={
                filters.pianoType !== undefined ? String(filters.pianoType) : "any"
              }
              onSelectionChange={(k) => handlePianoTypeChange(String(k))}
            />
            <div className="mt-2">
              <Select
                label="メーカー"
                placeholder="すべて"
                items={PIANO_BRAND_ITEMS}
                selectedKey={filters.pianoBrand ?? "any"}
                onSelectionChange={(k) => handlePianoBrandChange(String(k))}
              />
            </div>
            <div className="mt-2">
              <Select
                label="最低評価"
                placeholder="すべて"
                items={MIN_RATING_ITEMS}
                selectedKey={
                  filters.minRatingAverage !== undefined
                    ? String(filters.minRatingAverage)
                    : "any"
                }
                onSelectionChange={(k) => handleMinRatingChange(String(k))}
              />
            </div>

            <div className="mt-4 border-t border-slate-100 pt-3">
              <p className="text-xs font-semibold text-slate-500">ピアノの特徴</p>
              <p className="text-[11px] text-slate-400">
                投稿の平均値が指定値以上のピアノに絞り込みます
              </p>
              {ATTRIBUTE_LABELS.map(({ key, label, lo, hi }) => (
                <div key={key} className="mt-2">
                  <Select
                    label={`${label} (1=${lo} / 5=${hi})`}
                    placeholder="すべて"
                    items={ATTRIBUTE_SCORE_ITEMS}
                    selectedKey={
                      filters[key] !== undefined ? String(filters[key]) : "any"
                    }
                    onSelectionChange={(k) => handleAttributeChange(key, String(k))}
                  />
                </div>
              ))}
            </div>
          </div>
          <div className="flex justify-end gap-2 border-t border-slate-100 p-2">
            {filterCount > 0 ? (
              <Button
                variant="ghost"
                size="sm"
                onPress={() => onFiltersChange(EMPTY_FILTERS)}
              >
                クリア
              </Button>
            ) : null}
            <Button size="sm" onPress={() => setFiltersOpen(false)}>
              閉じる
            </Button>
          </div>
        </div>
      ) : null}

      {showHistory ? (
        <div className="mt-2 rounded-2xl border border-slate-200 bg-white shadow-lg">
          <div className="flex items-center justify-between border-b border-slate-100 px-4 py-2">
            <p className="text-xs font-semibold text-slate-500">最近表示したピアノ</p>
            <button
              type="button"
              onClick={() => {
                clearSearchHistory();
                setHistory([]);
              }}
              className="cursor-pointer text-[11px] text-slate-400 hover:text-slate-600 hover:underline outline-none focus-visible:ring-2 focus-visible:ring-amber-500 rounded"
            >
              履歴を全て消す
            </button>
          </div>
          <ul className="max-h-72 divide-y divide-slate-100 overflow-y-auto">
            {history.map((it) => {
              const inner = (
                <>
                  <Clock size={14} className="mt-0.5 shrink-0 text-slate-400" />
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-slate-900">
                      {it.displayName}
                    </p>
                    {it.subtitle ? (
                      <p className="truncate text-xs text-slate-500">{it.subtitle}</p>
                    ) : null}
                  </div>
                </>
              );
              return (
                <li key={it.pianoId} className="relative">
                  {onSelect ? (
                    <button
                      type="button"
                      onClick={() => handleHistoryClick(it)}
                      className="flex w-full cursor-pointer gap-2 px-4 py-2 pr-9 text-left transition hover:bg-slate-50 focus:bg-slate-50 outline-none"
                    >
                      {inner}
                    </button>
                  ) : (
                    <Link
                      to={`/pianos/${it.pianoId}`}
                      className="flex gap-2 px-4 py-2 pr-9 transition hover:bg-slate-50"
                      onClick={() => setFocused(false)}
                    >
                      {inner}
                    </Link>
                  )}
                  <button
                    type="button"
                    aria-label="この履歴を削除"
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      setHistory(removeSearchHistoryItem(it.pianoId));
                    }}
                    className="absolute right-2 top-1/2 -translate-y-1/2 cursor-pointer rounded-full p-1 text-slate-300 hover:bg-slate-100 hover:text-slate-600 outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                  >
                    <X size={12} />
                  </button>
                </li>
              );
            })}
          </ul>
        </div>
      ) : null}

      {showDropdown ? (
        <div className="mt-2 max-h-80 overflow-y-auto rounded-2xl border border-slate-200 bg-white shadow-lg">
          {loading ? (
            <p className="px-4 py-3 text-center text-xs text-slate-400">検索中...</p>
          ) : results.length === 0 ? (
            <p className="px-4 py-3 text-center text-xs text-slate-400">
              該当するピアノがありません
            </p>
          ) : (
            <ul className="divide-y divide-slate-100">
              {results.map((p) => {
                const pianoId = (() => {
                  try {
                    return parsePiano(p.name);
                  } catch {
                    return "";
                  }
                })();
                const subtitle =
                  p.address || p.venueType || p.prefecture || p.city || "";
                const inner = (
                  <>
                    <MapPin size={14} className="mt-0.5 shrink-0 text-amber-500" />
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium text-slate-900">
                        {p.displayName || pianoId}
                      </p>
                      {subtitle ? (
                        <p className="truncate text-xs text-slate-500">{subtitle}</p>
                      ) : null}
                    </div>
                  </>
                );
                if (onSelect) {
                  return (
                    <li key={p.name}>
                      <button
                        type="button"
                        onClick={() => {
                          recordHistory(p, pianoId);
                          onSelect(p);
                          // タップしたピアノ名を入力欄にフィル。
                          setQuery(p.displayName || "");
                          setFocused(false);
                          inputRef.current?.blur();
                        }}
                        className="flex w-full cursor-pointer gap-2 px-4 py-2 text-left transition hover:bg-slate-50 focus:bg-slate-50 outline-none"
                      >
                        {inner}
                      </button>
                    </li>
                  );
                }
                return (
                  <li key={p.name}>
                    <Link
                      to={`/pianos/${pianoId}`}
                      className="flex gap-2 px-4 py-2 transition hover:bg-slate-50"
                      onClick={() => {
                        recordHistory(p, pianoId);
                        setFocused(false);
                      }}
                    >
                      {inner}
                    </Link>
                  </li>
                );
              })}
            </ul>
          )}
        </div>
      ) : null}
    </div>
  );
}
