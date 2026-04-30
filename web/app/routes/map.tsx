import { useCallback, useEffect, useRef, useState } from "react";
import { Link } from "react-router";
import { Plus, RotateCw, X } from "lucide-react";
import maplibregl, { type LngLatBoundsLike, type Marker } from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

import { Button } from "../components/ui/button";
import { Select } from "../components/ui/select";
import { TextField } from "../components/ui/text-field";
import { Sheet } from "../components/ui/sheet";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { pianoClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import { findStyle, getSavedStyleId } from "../lib/map-styles";
import {
  CreatePianoRequest,
  LatLng,
  LatLngBounds,
  Piano,
  PianoAvailability,
  PianoKind,
  PianoType,
  SearchPianosRequest,
} from "../lib/gen/piano/v1/piano_pb";
import { parsePiano } from "../lib/resource-name";

import type { Route } from "./+types/map";

export function meta({}: Route.MetaArgs) {
  return [{ title: "マップ — PiaMap" }];
}

const TOKYO: [number, number] = [139.7536, 35.6936];

type CreateDraft = { lng: number; lat: number };

export default function MapPage() {
  const { authed } = useAuth();
  const containerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<Map<string, Marker>>(new Map());

  const [pianos, setPianos] = useState<Piano[]>([]);
  const [selected, setSelected] = useState<Piano | null>(null);
  const [createMode, setCreateMode] = useState(false);
  const [draft, setDraft] = useState<CreateDraft | null>(null);
  const [signupOpen, setSignupOpen] = useState(false);
  const [signupAction, setSignupAction] = useState<string | undefined>();
  const [isFetching, setIsFetching] = useState(false);

  // showLoading=true のときだけ画面のグレーアウトとボタンのスピナーを出す。
  // map の moveend 等の自動再取得では false (画面が頻繁に暗くなるのを避ける)。
  const fetchPianos = useCallback(
    async (bounds: LngLatBoundsLike, opts?: { showLoading?: boolean }) => {
      const b = maplibregl.LngLatBounds.convert(bounds);
      const req = new SearchPianosRequest({
        bounds: new LatLngBounds({
          southwest: new LatLng({ latitude: b.getSouth(), longitude: b.getWest() }),
          northeast: new LatLng({ latitude: b.getNorth(), longitude: b.getEast() }),
        }),
        limit: 200,
      });
      if (opts?.showLoading) setIsFetching(true);
      try {
        const res = await pianoClient.searchPianos(req);
        setPianos(res.pianos);
      } catch (e) {
        console.error("searchPianos failed", e);
      } finally {
        if (opts?.showLoading) setIsFetching(false);
      }
    },
    [],
  );

  // map 初期化。
  useEffect(() => {
    if (!containerRef.current) return;
    const initial = findStyle(getSavedStyleId());
    const map = new maplibregl.Map({
      container: containerRef.current,
      style: initial.style,
      center: TOKYO,
      zoom: 13,
      attributionControl: { compact: true },
    });
    map.addControl(new maplibregl.NavigationControl({ showCompass: false }), "top-right");
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
    for (const [name, marker] of markersRef.current) {
      if (!next.has(name)) marker.remove();
    }
    markersRef.current = next;
  }, [pianos]);

  const handleCreate = () => {
    if (!authed) {
      setSignupAction("ピアノの登録");
      setSignupOpen(true);
      return;
    }
    setCreateMode((v) => !v);
  };

  return (
    <div className="relative h-[calc(100dvh-4rem)] w-full overflow-hidden">
      {/* maplibre は container に .maplibregl-map (position: relative) を当てるので
          absolute inset-0 で潰れる。明示的な h/w で渡す。 */}
      <div ref={containerRef} className="h-full w-full" />

      {/* ローディング時のグレーアウト (透過、ボタン等は塞がない pointer-events-none) */}
      <div
        aria-hidden
        className={`pointer-events-none absolute inset-0 bg-slate-900/40 transition-opacity duration-200 ${
          isFetching ? "opacity-100" : "opacity-0"
        }`}
      />

      {/* 右下フローティングアクション */}
      <div className="absolute bottom-4 right-4 flex flex-col items-end gap-2">
        <Button
          aria-label={isFetching ? "更新中" : "再読込"}
          variant="secondary"
          size="md"
          isPending={isFetching}
          onPress={() => {
            const map = mapRef.current;
            if (map) fetchPianos(map.getBounds(), { showLoading: true });
          }}
          className="h-12 w-12 rounded-full p-0"
        >
          <RotateCw size={18} />
        </Button>
        <Button
          aria-label={createMode ? "登録キャンセル" : "ピアノを登録"}
          variant={createMode ? "danger" : "primary"}
          size="md"
          onPress={handleCreate}
          className="h-14 w-14 rounded-full p-0"
        >
          {createMode ? <X size={22} /> : <Plus size={26} />}
        </Button>
      </div>

      {/* createMode 中のヒント */}
      {createMode ? (
        <div className="pointer-events-none absolute left-1/2 top-4 z-10 -translate-x-1/2 rounded-full bg-slate-900/85 px-4 py-2 text-xs font-medium text-white shadow">
          地図をタップして登録地点を選択
        </div>
      ) : null}

      <Sheet
        isOpen={selected !== null}
        onOpenChange={(open) => !open && setSelected(null)}
        title={selected?.displayName}
      >
        {selected ? <PianoSummary piano={selected} /> : null}
      </Sheet>

      <Sheet
        isOpen={draft !== null}
        onOpenChange={(open) => !open && setDraft(null)}
        title="ピアノを登録"
      >
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
      </Sheet>

      <SignUpPromptModal isOpen={signupOpen} onOpenChange={setSignupOpen} action={signupAction} />
    </div>
  );
}


function PianoSummary({ piano }: { piano: Piano }) {
  const id = parsePianoSafe(piano.name);
  const avg = piano.ratingAverage;
  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2 text-sm text-slate-600">
        <span>
          {piano.postCount > 0 ? `★ ${avg.toFixed(1)} (${piano.postCount}件)` : "未評価"}
        </span>
        {piano.address ? <span className="truncate text-slate-500">— {piano.address}</span> : null}
      </div>
      {piano.description ? (
        <p className="line-clamp-3 text-sm text-slate-700">{piano.description}</p>
      ) : null}
      <Link
        to={id ? `/pianos/${id}` : "#"}
        className="block rounded-full bg-amber-500 py-2.5 text-center text-sm font-semibold text-white hover:bg-amber-600"
      >
        詳細を見る →
      </Link>
    </div>
  );
}

function parsePianoSafe(name: string): string | null {
  try {
    return parsePiano(name);
  } catch {
    return null;
  }
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
  const [pianoBrand, setPianoBrand] = useState("");
  // 未選択を表現するため null 許容。submit 時に UNKNOWN にフォールバック。
  const [pianoType, setPianoType] = useState<PianoType | null>(null);
  const [kind, setKind] = useState<PianoKind>(PianoKind.STREET);
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  // 必須項目の per-field エラー (submit 時のみセット、入力で自動クリア)。
  const [fieldErrors, setFieldErrors] = useState<{ displayName?: string }>({});

  const submit = async () => {
    const next: typeof fieldErrors = {};
    if (!displayName.trim()) next.displayName = "名前を入力してください";
    if (Object.keys(next).length > 0) {
      setFieldErrors(next);
      return;
    }
    setFieldErrors({});
    setSubmitting(true);
    setErr(null);
    try {
      const res = await pianoClient.createPiano(
        new CreatePianoRequest({
          piano: new Piano({
            displayName,
            description: description || undefined,
            location: new LatLng({ latitude: location.lat, longitude: location.lng }),
            kind,
            pianoType: pianoType ?? PianoType.UNKNOWN,
            pianoBrand: pianoBrand.trim() || "unknown",
            availability: PianoAvailability.REGULAR,
          }),
        }),
      );
      if (res.piano) onCreated(res.piano);
    } catch (e: any) {
      setErr(e?.message || String(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="space-y-3">
      <p className="text-xs text-slate-500">
        座標: {location.lat.toFixed(6)}, {location.lng.toFixed(6)}
      </p>

      <TextField
        label="名前"
        value={displayName}
        onChange={(v) => {
          setDisplayName(v);
          if (fieldErrors.displayName) setFieldErrors((e) => ({ ...e, displayName: undefined }));
        }}
        placeholder="〇〇駅前ストリートピアノ"
        autoFocus
        isRequired
        errorMessage={fieldErrors.displayName}
      />
      <TextField
        label="説明"
        value={description}
        onChange={setDescription}
        multiline
        rows={2}
      />

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
          placeholder="種類を選択"
          selectedKey={pianoType !== null ? String(pianoType) : null}
          onSelectionChange={(k) => setPianoType(k === null ? null : (Number(k) as PianoType))}
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
        placeholder="ヤマハ、カワイ など (不明なら空欄)"
      />

      {err ? <p className="text-sm text-rose-600">{err}</p> : null}

      <div className="flex gap-2 pt-2">
        <Button variant="secondary" size="md" className="flex-1" onPress={onCancel} isDisabled={submitting}>
          キャンセル
        </Button>
        <Button size="md" className="flex-1" onPress={submit} isPending={submitting} isDisabled={submitting}>
          {submitting ? "登録中..." : "登録"}
        </Button>
      </div>
    </div>
  );
}
