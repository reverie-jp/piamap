// マップ検索のローカル履歴。実際にピアノ詳細シートへ遷移したものだけ記録する。
// localStorage キー: piamap.search-history.v1

const STORAGE_KEY = "piamap.search-history.v1";
const MAX_ITEMS = 10;

export type SearchHistoryItem = {
  pianoId: string;
  displayName: string;
  subtitle?: string; // 住所 / 会場種別など補足の1行
  savedAt: number;
};

function isBrowser(): boolean {
  return typeof window !== "undefined" && typeof window.localStorage !== "undefined";
}

export function getSearchHistory(): SearchHistoryItem[] {
  if (!isBrowser()) return [];
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];
    return parsed.filter(
      (it): it is SearchHistoryItem =>
        typeof it === "object" &&
        it !== null &&
        typeof (it as SearchHistoryItem).pianoId === "string" &&
        typeof (it as SearchHistoryItem).displayName === "string",
    );
  } catch {
    return [];
  }
}

export function pushSearchHistory(item: Omit<SearchHistoryItem, "savedAt">): SearchHistoryItem[] {
  if (!isBrowser()) return [];
  const current = getSearchHistory();
  const next: SearchHistoryItem[] = [
    { ...item, savedAt: Date.now() },
    ...current.filter((c) => c.pianoId !== item.pianoId),
  ].slice(0, MAX_ITEMS);
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  } catch {
    /* quota など無視 */
  }
  return next;
}

export function removeSearchHistoryItem(pianoId: string): SearchHistoryItem[] {
  if (!isBrowser()) return [];
  const next = getSearchHistory().filter((c) => c.pianoId !== pianoId);
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  } catch {
    /* noop */
  }
  return next;
}

export function clearSearchHistory(): void {
  if (!isBrowser()) return;
  try {
    window.localStorage.removeItem(STORAGE_KEY);
  } catch {
    /* noop */
  }
}
