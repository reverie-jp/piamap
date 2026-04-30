// マップの中心+ズームを sessionStorage に保存・復元する。
// 詳細画面から戻った時に位置がリセットされないようにするため、
// pan/zoom 終了時に保存して、初期化時に読み込む。
// sessionStorage はタブ単位で消えるので、長時間後に戻ってきた時のためのものではない。

const STORAGE_KEY = "piamap.map.viewState.v1";

export type MapViewState = {
  lng: number;
  lat: number;
  zoom: number;
};

function isBrowser(): boolean {
  return typeof window !== "undefined" && typeof window.sessionStorage !== "undefined";
}

export function getSavedMapView(): MapViewState | null {
  if (!isBrowser()) return null;
  try {
    const raw = window.sessionStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const v = JSON.parse(raw);
    if (
      v &&
      typeof v.lng === "number" &&
      typeof v.lat === "number" &&
      typeof v.zoom === "number"
    ) {
      return v;
    }
    return null;
  } catch {
    return null;
  }
}

export function saveMapView(v: MapViewState): void {
  if (!isBrowser()) return;
  try {
    window.sessionStorage.setItem(STORAGE_KEY, JSON.stringify(v));
  } catch {
    /* quota など無視 */
  }
}
