import { useCallback, useEffect, useRef, useState } from "react";
import maplibregl, { type LngLatBoundsLike, type Marker } from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

import { pianoClient } from "../lib/api-client";
import { setAccessToken, getAccessToken } from "../lib/auth";
import {
  DEFAULT_STYLE_ID,
  MAP_STYLES,
  findStyle,
  getSavedStyleId,
  saveStyleId,
  type MapStyleDef,
} from "../lib/map-styles";
import {
  LatLng,
  LatLngBounds,
  Piano,
  PianoAvailability,
  PianoKind,
  PianoType,
  SearchPianosRequest,
  CreatePianoRequest,
} from "../lib/gen/piano/v1/piano_pb";

import type { Route } from "./+types/map";

export function meta({}: Route.MetaArgs) {
  return [{ title: "PiaMap — マップ" }];
}

const TOKYO: [number, number] = [139.7536, 35.6936]; // [lng, lat]

type CreateDraft = { lng: number; lat: number };

export default function MapPage() {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<Map<string, Marker>>(new Map());

  const [pianos, setPianos] = useState<Piano[]>([]);
  const [selected, setSelected] = useState<Piano | null>(null);
  const [createMode, setCreateMode] = useState(false);
  const [draft, setDraft] = useState<CreateDraft | null>(null);
  const [authed, setAuthed] = useState<boolean>(false);
  const [styleId, setStyleId] = useState<string>(DEFAULT_STYLE_ID);

  useEffect(() => {
    setAuthed(Boolean(getAccessToken()));
    setStyleId(getSavedStyleId());
  }, []);

  const fetchPianos = useCallback(async (bounds: LngLatBoundsLike) => {
    const b = maplibregl.LngLatBounds.convert(bounds);
    const req = new SearchPianosRequest({
      bounds: new LatLngBounds({
        southwest: new LatLng({ latitude: b.getSouth(), longitude: b.getWest() }),
        northeast: new LatLng({ latitude: b.getNorth(), longitude: b.getEast() }),
      }),
      limit: 200,
    });
    try {
      const res = await pianoClient.searchPianos(req);
      setPianos(res.pianos);
    } catch (e) {
      console.error("searchPianos failed", e);
    }
  }, []);

  // map 初期化 (client-only)。初期 style は localStorage から復元。
  useEffect(() => {
    if (!containerRef.current) return;
    const initialStyle = findStyle(getSavedStyleId());
    const map = new maplibregl.Map({
      container: containerRef.current,
      style: initialStyle.style,
      center: TOKYO,
      zoom: 13,
    });
    map.addControl(new maplibregl.NavigationControl(), "top-right");
    mapRef.current = map;

    const onLoadOrMove = () => fetchPianos(map.getBounds());
    map.on("load", onLoadOrMove);
    map.on("moveend", onLoadOrMove);

    return () => {
      map.remove();
      mapRef.current = null;
      markersRef.current.forEach((m) => m.remove());
      markersRef.current.clear();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // styleId が変わったら setStyle で差し替え。HTML マーカーは DOM に直接付くので生存する。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    const def = findStyle(styleId);
    map.setStyle(def.style);
  }, [styleId]);

  // create mode のクリックハンドラ。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    if (!createMode) {
      map.getCanvas().style.cursor = "";
      return;
    }
    map.getCanvas().style.cursor = "crosshair";
    const onClick = (e: maplibregl.MapMouseEvent) => {
      setDraft({ lng: e.lngLat.lng, lat: e.lngLat.lat });
      setCreateMode(false);
    };
    map.on("click", onClick);
    return () => {
      map.off("click", onClick);
      map.getCanvas().style.cursor = "";
    };
  }, [createMode]);

  // pianos 配列の変化に応じて marker を再描画。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;

    const next = new Map<string, Marker>();
    for (const p of pianos) {
      if (!p.location) continue;
      let marker = markersRef.current.get(p.name);
      if (!marker) {
        const el = document.createElement("button");
        el.className =
          "flex h-7 w-7 items-center justify-center rounded-full border-2 border-white bg-amber-500 text-xs font-bold text-white shadow hover:bg-amber-600";
        el.textContent = "♪";
        el.title = p.displayName;
        el.onclick = (ev) => {
          ev.stopPropagation();
          setSelected(p);
        };
        marker = new maplibregl.Marker({ element: el })
          .setLngLat([p.location.longitude, p.location.latitude])
          .addTo(map);
      } else {
        marker.setLngLat([p.location.longitude, p.location.latitude]);
      }
      next.set(p.name, marker);
    }
    // 消えたマーカーを除去。
    for (const [name, marker] of markersRef.current) {
      if (!next.has(name)) marker.remove();
    }
    markersRef.current = next;
  }, [pianos]);

  return (
    <div className="relative h-screen w-screen overflow-hidden">
      <DevTokenBar
        authed={authed}
        onChange={(t) => {
          setAccessToken(t);
          setAuthed(Boolean(t));
        }}
      />

      <div ref={containerRef} className="h-full w-full" />

      <div className="absolute bottom-6 right-6 flex flex-col items-end gap-2">
        <StyleSwitcher
          value={styleId}
          onChange={(id) => {
            setStyleId(id);
            saveStyleId(id);
          }}
        />
        <button
          onClick={() => {
            if (!authed) {
              alert("ピアノを登録するには上部の入力欄に dev token を入れてください");
              return;
            }
            setCreateMode((v) => !v);
          }}
          className={`rounded-full px-4 py-2 font-semibold shadow ${
            createMode ? "bg-rose-600 text-white" : "bg-emerald-600 text-white hover:bg-emerald-700"
          }`}
        >
          {createMode ? "登録キャンセル" : "+ ピアノを登録"}
        </button>
        <button
          onClick={() => {
            const map = mapRef.current;
            if (map) fetchPianos(map.getBounds());
          }}
          className="rounded-full bg-white px-4 py-2 font-semibold shadow hover:bg-slate-100"
        >
          再読込
        </button>
      </div>

      {selected ? (
        <PianoDetail piano={selected} onClose={() => setSelected(null)} />
      ) : null}

      {draft ? (
        <CreatePianoForm
          location={draft}
          onCancel={() => setDraft(null)}
          onCreated={(p) => {
            setDraft(null);
            setPianos((prev) => [p, ...prev]);
            setSelected(p);
          }}
        />
      ) : null}
    </div>
  );
}

function DevTokenBar({ authed, onChange }: { authed: boolean; onChange: (t: string | null) => void }) {
  const [token, setToken] = useState("");
  return (
    <div className="absolute left-1/2 top-4 z-10 flex -translate-x-1/2 items-center gap-2 rounded-full bg-white/95 px-4 py-2 shadow">
      <span className={`h-2 w-2 rounded-full ${authed ? "bg-emerald-500" : "bg-slate-400"}`} />
      <span className="text-sm text-slate-700">
        {authed ? "認証済 (dev token)" : "未認証"}
      </span>
      <input
        type="password"
        placeholder="dev access token を貼り付け"
        value={token}
        onChange={(e) => setToken(e.target.value)}
        className="ml-2 w-64 rounded border border-slate-300 px-2 py-1 text-xs"
      />
      <button
        onClick={() => {
          onChange(token || null);
          setToken("");
        }}
        className="rounded bg-slate-800 px-3 py-1 text-xs font-medium text-white hover:bg-slate-700"
      >
        保存
      </button>
      {authed ? (
        <button
          onClick={() => onChange(null)}
          className="rounded border border-slate-300 px-3 py-1 text-xs hover:bg-slate-100"
        >
          解除
        </button>
      ) : null}
    </div>
  );
}

function StyleSwitcher({ value, onChange }: { value: string; onChange: (id: string) => void }) {
  return (
    <label className="flex items-center gap-2 rounded-full bg-white px-3 py-1.5 shadow">
      <span className="text-xs text-slate-500">タイル</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="bg-transparent text-sm font-medium text-slate-800 outline-none"
      >
        {MAP_STYLES.map((s: MapStyleDef) => (
          <option key={s.id} value={s.id}>
            {s.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function PianoDetail({ piano, onClose }: { piano: Piano; onClose: () => void }) {
  const avg = piano.ratingAverage;
  return (
    <div className="absolute bottom-6 left-6 z-10 w-80 rounded-lg bg-white p-4 shadow-xl">
      <div className="flex items-start justify-between">
        <h2 className="text-lg font-bold text-slate-900">{piano.displayName}</h2>
        <button onClick={onClose} className="text-slate-400 hover:text-slate-700" aria-label="閉じる">
          ×
        </button>
      </div>
      {piano.description ? (
        <p className="mt-2 text-sm text-slate-700">{piano.description}</p>
      ) : null}
      <dl className="mt-3 space-y-1 text-xs text-slate-600">
        {piano.address ? (
          <div>
            <dt className="inline font-semibold">住所:</dt> <dd className="inline">{piano.address}</dd>
          </div>
        ) : null}
        <div>
          <dt className="inline font-semibold">種別:</dt>{" "}
          <dd className="inline">{kindLabel(piano.kind)} / {pianoTypeLabel(piano.pianoType)}</dd>
        </div>
        {piano.pianoBrand ? (
          <div>
            <dt className="inline font-semibold">メーカー:</dt> <dd className="inline">{piano.pianoBrand}</dd>
          </div>
        ) : null}
        <div>
          <dt className="inline font-semibold">評価:</dt>{" "}
          <dd className="inline">
            {piano.postCount > 0 ? `★${avg.toFixed(1)} (${piano.postCount}件)` : "未評価"}
          </dd>
        </div>
        {piano.creator ? (
          <div>
            <dt className="inline font-semibold">登録:</dt> <dd className="inline">{piano.creator}</dd>
          </div>
        ) : null}
      </dl>
    </div>
  );
}

function CreatePianoForm({
  location,
  onCancel,
  onCreated,
}: {
  location: { lng: number; lat: number };
  onCancel: () => void;
  onCreated: (p: Piano) => void;
}) {
  const [displayName, setDisplayName] = useState("");
  const [description, setDescription] = useState("");
  const [pianoBrand, setPianoBrand] = useState("unknown");
  const [pianoType, setPianoType] = useState<PianoType>(PianoType.UNKNOWN);
  const [kind, setKind] = useState<PianoKind>(PianoKind.STREET);
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    if (!displayName) {
      setErr("名前を入力してください");
      return;
    }
    setSubmitting(true);
    setErr(null);
    try {
      const req = new CreatePianoRequest({
        piano: new Piano({
          displayName,
          description: description || undefined,
          location: new LatLng({ latitude: location.lat, longitude: location.lng }),
          kind,
          pianoType,
          pianoBrand,
          availability: PianoAvailability.REGULAR,
        }),
      });
      const res = await pianoClient.createPiano(req);
      if (res.piano) onCreated(res.piano);
    } catch (e: any) {
      setErr(e?.message || String(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="absolute inset-0 z-20 flex items-center justify-center bg-slate-900/40">
      <div className="w-96 rounded-lg bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-900">ピアノを登録</h2>
        <p className="mt-1 text-xs text-slate-500">
          座標: {location.lat.toFixed(6)}, {location.lng.toFixed(6)}
        </p>

        <label className="mt-4 block text-sm font-medium text-slate-700">名前 *</label>
        <input
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm"
          placeholder="〇〇駅前ストリートピアノ"
        />

        <label className="mt-3 block text-sm font-medium text-slate-700">説明</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm"
          rows={2}
        />

        <div className="mt-3 grid grid-cols-2 gap-2">
          <div>
            <label className="block text-sm font-medium text-slate-700">種別</label>
            <select
              value={kind}
              onChange={(e) => setKind(Number(e.target.value))}
              className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm"
            >
              <option value={PianoKind.STREET}>ストリート</option>
              <option value={PianoKind.PRACTICE_ROOM}>練習室</option>
              <option value={PianoKind.OTHER}>その他</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700">ピアノ</label>
            <select
              value={pianoType}
              onChange={(e) => setPianoType(Number(e.target.value))}
              className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm"
            >
              <option value={PianoType.UNKNOWN}>不明</option>
              <option value={PianoType.GRAND}>グランド</option>
              <option value={PianoType.UPRIGHT}>アップライト</option>
              <option value={PianoType.ELECTRONIC}>電子</option>
            </select>
          </div>
        </div>

        <label className="mt-3 block text-sm font-medium text-slate-700">メーカー</label>
        <input
          value={pianoBrand}
          onChange={(e) => setPianoBrand(e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm"
        />

        {err ? <p className="mt-3 text-sm text-rose-600">{err}</p> : null}

        <div className="mt-4 flex justify-end gap-2">
          <button
            onClick={onCancel}
            disabled={submitting}
            className="rounded px-3 py-1 text-sm hover:bg-slate-100"
          >
            キャンセル
          </button>
          <button
            onClick={submit}
            disabled={submitting}
            className="rounded bg-emerald-600 px-3 py-1 text-sm font-semibold text-white hover:bg-emerald-700 disabled:opacity-50"
          >
            {submitting ? "登録中..." : "登録"}
          </button>
        </div>
      </div>
    </div>
  );
}

function kindLabel(k: PianoKind): string {
  switch (k) {
    case PianoKind.STREET:
      return "ストリート";
    case PianoKind.PRACTICE_ROOM:
      return "練習室";
    case PianoKind.OTHER:
      return "その他";
    default:
      return "不明";
  }
}

function pianoTypeLabel(t: PianoType): string {
  switch (t) {
    case PianoType.GRAND:
      return "グランド";
    case PianoType.UPRIGHT:
      return "アップライト";
    case PianoType.ELECTRONIC:
      return "電子";
    case PianoType.UNKNOWN:
    default:
      return "不明";
  }
}
