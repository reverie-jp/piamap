import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import { CalendarClock, Check, Clock, MapPin, Music, Plus, RotateCw, Tag, X } from "lucide-react";
import { FieldMask } from "@bufbuild/protobuf";
import maplibregl, { type LngLatBoundsLike, type Marker } from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

import { Button } from "../components/ui/button";
import { Select } from "../components/ui/select";
import { TextField } from "../components/ui/text-field";
import { Sheet } from "../components/ui/sheet";
import { EMPTY_FILTERS, MapSearchBar, type MapSearchFilters } from "../components/MapSearchBar";
import { PianoAttributeMeters } from "../components/PianoAttributeMeters";
import { SignUpPromptModal } from "../components/SignUpPromptModal";
import { pianoClient } from "../lib/api-client";
import { useAuth } from "../lib/auth";
import { findStyle, getSavedStyleId } from "../lib/map-styles";
import { getSavedMapView, saveMapView } from "../lib/map-view-state";
import {
  CreatePianoRequest,
  GetPianoRequest,
  LatLng,
  LatLngBounds,
  Piano,
  PianoAvailability,
  PianoKind,
  PianoType,
  SearchPianosRequest,
  UpdatePianoRequest,
} from "../lib/gen/piano/v1/piano_pb";
import { formatPiano, parsePiano } from "../lib/resource-name";

import type { Route } from "./+types/map";

export function meta({}: Route.MetaArgs) {
  return [{ title: "マップ — PiaMap" }];
}

const TOKYO: [number, number] = [139.7536, 35.6936];

type CreateDraft = { lng: number; lat: number };
type RelocateState = {
  pianoId: string;
  piano: Piano;
  initialLat: number;
  initialLng: number;
  newLat: number;
  newLng: number;
};

export default function MapPage() {
  const { authed } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const relocateId = searchParams.get("relocate");

  const containerRef = useRef<HTMLDivElement | null>(null);
  const mapRef = useRef<maplibregl.Map | null>(null);
  const markersRef = useRef<Map<string, Marker>>(new Map());
  const relocateMarkerRef = useRef<Marker | null>(null);

  const [pianos, setPianos] = useState<Piano[]>([]);
  const [selected, setSelected] = useState<Piano | null>(null);
  const [createMode, setCreateMode] = useState(false);
  const [draft, setDraft] = useState<CreateDraft | null>(null);
  const [signupOpen, setSignupOpen] = useState(false);
  const [signupAction, setSignupAction] = useState<string | undefined>();
  const [isFetching, setIsFetching] = useState(false);
  const [relocate, setRelocate] = useState<RelocateState | null>(null);
  const [relocateSummary, setRelocateSummary] = useState("");
  const [relocateSubmitting, setRelocateSubmitting] = useState(false);
  const [relocateErr, setRelocateErr] = useState<string | null>(null);
  const [filters, setFilters] = useState<MapSearchFilters>(EMPTY_FILTERS);
  const [searchClearSignal, setSearchClearSignal] = useState(0);

  const inRelocate = relocate !== null;

  // showLoading=true のときだけ画面のグレーアウトとボタンのスピナーを出す。
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
      if (filters.pianoType !== undefined) req.pianoType = filters.pianoType;
      if (filters.pianoBrand !== undefined) req.pianoBrand = filters.pianoBrand;
      if (filters.minRatingAverage !== undefined) req.minRatingAverage = filters.minRatingAverage;
      if (filters.minAmbientNoiseAverage !== undefined) req.minAmbientNoiseAverage = filters.minAmbientNoiseAverage;
      if (filters.minFootTrafficAverage !== undefined) req.minFootTrafficAverage = filters.minFootTrafficAverage;
      if (filters.minResonanceAverage !== undefined) req.minResonanceAverage = filters.minResonanceAverage;
      if (filters.minKeyTouchWeightAverage !== undefined) req.minKeyTouchWeightAverage = filters.minKeyTouchWeightAverage;
      if (filters.minTuningQualityAverage !== undefined) req.minTuningQualityAverage = filters.minTuningQualityAverage;
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
    [
      filters.pianoType,
      filters.pianoBrand,
      filters.minRatingAverage,
      filters.minAmbientNoiseAverage,
      filters.minFootTrafficAverage,
      filters.minResonanceAverage,
      filters.minKeyTouchWeightAverage,
      filters.minTuningQualityAverage,
    ],
  );

  // map.on("moveend", ...) のクロージャを毎度作り直さないため、最新の fetchPianos を ref で参照する。
  const fetchPianosRef = useRef(fetchPianos);
  useEffect(() => {
    fetchPianosRef.current = fetchPianos;
  }, [fetchPianos]);

  // フィルタ変更時に現在のbboxで再検索 (マーカーに反映)。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    fetchPianos(map.getBounds());
  }, [fetchPianos]);

  // map 初期化。
  useEffect(() => {
    if (!containerRef.current) return;
    const initial = findStyle(getSavedStyleId());
    // 詳細画面から戻った時に位置を復元するため、sessionStorage の前回ビューを優先。
    const saved = getSavedMapView();
    const map = new maplibregl.Map({
      container: containerRef.current,
      style: initial.style,
      center: saved ? [saved.lng, saved.lat] : TOKYO,
      zoom: saved ? saved.zoom : 13,
      attributionControl: { compact: true },
    });
    // モバイルは pinch zoom が基本。NavigationControl は省略してヘッダーとの重なりを避ける。
    mapRef.current = map;

    const onLoadOrMove = () => {
      const c = map.getCenter();
      saveMapView({ lng: c.lng, lat: c.lat, zoom: map.getZoom() });
      fetchPianosRef.current(map.getBounds());
    };
    map.on("load", onLoadOrMove);
    map.on("moveend", onLoadOrMove);

    return () => {
      map.remove();
      mapRef.current = null;
      markersRef.current.forEach((m) => m.remove());
      markersRef.current.clear();
      relocateMarkerRef.current?.remove();
      relocateMarkerRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // create mode のクリックハンドラ。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    if (!createMode || inRelocate) {
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
  }, [createMode, inRelocate]);

  // relocate=<id> 起動: 該当ピアノを取得し、ドラッグ可能なマーカーを置く。
  useEffect(() => {
    if (!relocateId) {
      setRelocate(null);
      return;
    }
    if (!authed) {
      setSignupAction("ピアノの位置変更");
      setSignupOpen(true);
      return;
    }
    let cancelled = false;
    pianoClient
      .getPiano(new GetPianoRequest({ name: formatPiano(relocateId) }))
      .then((res) => {
        if (cancelled) return;
        const piano = res.piano;
        if (!piano || !piano.location) return;
        setRelocate({
          pianoId: relocateId,
          piano,
          initialLat: piano.location.latitude,
          initialLng: piano.location.longitude,
          newLat: piano.location.latitude,
          newLng: piano.location.longitude,
        });
        setRelocateSummary("");
        setRelocateErr(null);
      })
      .catch((e) => {
        if (!cancelled) setRelocateErr((e as Error)?.message || String(e));
      });
    return () => {
      cancelled = true;
    };
  }, [relocateId, authed]);

  // relocate のドラッグ可能マーカー。relocate state がある間だけ表示。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    if (!relocate) {
      relocateMarkerRef.current?.remove();
      relocateMarkerRef.current = null;
      return;
    }
    if (!relocateMarkerRef.current) {
      const el = document.createElement("div");
      el.className =
        "flex h-9 w-9 items-center justify-center rounded-full border-2 border-white bg-rose-500 text-base font-bold text-white shadow ring-4 ring-rose-500/30";
      el.textContent = "♪";
      el.title = "ドラッグして移動";
      const marker = new maplibregl.Marker({ element: el, draggable: true })
        .setLngLat([relocate.newLng, relocate.newLat])
        .addTo(map);
      marker.on("dragend", () => {
        const ll = marker.getLngLat();
        setRelocate((s) => (s ? { ...s, newLat: ll.lat, newLng: ll.lng } : s));
      });
      relocateMarkerRef.current = marker;
      map.flyTo({ center: [relocate.newLng, relocate.newLat], zoom: Math.max(map.getZoom(), 16) });
    } else {
      relocateMarkerRef.current.setLngLat([relocate.newLng, relocate.newLat]);
    }
  }, [relocate]);

  // pianos 配列の変化に応じて marker を再描画。
  useEffect(() => {
    const map = mapRef.current;
    if (!map) return;
    const next = new Map<string, Marker>();
    for (const p of pianos) {
      if (!p.location) continue;
      // relocate 中のピアノは draggable マーカーで表しているので通常マーカーは出さない。
      if (relocate && p.name === relocate.piano.name) continue;
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
  }, [pianos, relocate]);

  const handleCreate = () => {
    if (!authed) {
      setSignupAction("ピアノの登録");
      setSignupOpen(true);
      return;
    }
    setCreateMode((v) => !v);
  };

  const cancelRelocate = useCallback(() => {
    if (relocate) {
      navigate(`/pianos/${relocate.pianoId}`);
    } else {
      navigate("/map", { replace: true });
    }
  }, [relocate, navigate]);

  const submitRelocate = useCallback(async () => {
    if (!relocate) return;
    setRelocateSubmitting(true);
    setRelocateErr(null);
    try {
      await pianoClient.updatePiano(
        new UpdatePianoRequest({
          piano: new Piano({
            name: relocate.piano.name,
            location: new LatLng({ latitude: relocate.newLat, longitude: relocate.newLng }),
          }),
          updateMask: new FieldMask({ paths: ["location"] }),
          editSummary: relocateSummary.trim() || undefined,
        }),
      );
      navigate(`/pianos/${relocate.pianoId}`, {
        state: { savedMessage: "位置を更新しました" },
      });
    } catch (e) {
      setRelocateErr((e as Error)?.message || String(e));
    } finally {
      setRelocateSubmitting(false);
    }
  }, [relocate, relocateSummary, navigate]);

  const movedDistanceM = useMemo(() => {
    if (!relocate) return 0;
    return haversineM(
      relocate.initialLat,
      relocate.initialLng,
      relocate.newLat,
      relocate.newLng,
    );
  }, [relocate]);

  return (
    <div className="relative h-[calc(100dvh-4rem)] w-full overflow-hidden">
      <div ref={containerRef} className="h-full w-full" />

      {!inRelocate && !createMode ? (
        <MapSearchBar
          filters={filters}
          onFiltersChange={setFilters}
          clearSignal={searchClearSignal}
          onSelect={(p) => {
            const map = mapRef.current;
            if (map && p.location) {
              map.flyTo({
                center: [p.location.longitude, p.location.latitude],
                zoom: Math.max(map.getZoom(), 15),
              });
            }
            setSelected(p);
          }}
        />
      ) : null}

      <div
        aria-hidden
        className={`pointer-events-none absolute inset-0 bg-slate-900/40 transition-opacity duration-200 ${
          isFetching ? "opacity-100" : "opacity-0"
        }`}
      />

      {/* relocate モード以外のフローティングアクション */}
      {!inRelocate ? (
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
      ) : null}

      {createMode ? (
        <div className="pointer-events-none absolute left-1/2 top-3 z-30 -translate-x-1/2 rounded-full bg-slate-900/85 px-4 py-2 text-xs font-medium text-white shadow">
          地図をタップして登録地点を選択
        </div>
      ) : null}

      {/* relocate モードの上部ヒント + 下部アクションバー */}
      {inRelocate && relocate ? (
        <>
          <div className="pointer-events-none absolute left-1/2 top-4 z-10 -translate-x-1/2 rounded-full bg-slate-900/85 px-4 py-2 text-xs font-medium text-white shadow">
            ピンをドラッグして新しい位置に移動
          </div>
          <div className="absolute inset-x-0 bottom-0 z-20 border-t border-slate-200 bg-white p-3 shadow-[0_-8px_20px_rgba(15,23,42,0.08)]">
            <div className="mb-1 text-sm font-bold text-slate-900 truncate">
              {relocate.piano.displayName}
            </div>
            <div className="mb-2 text-xs text-slate-500">
              移動距離: {movedDistanceM < 1 ? "0" : movedDistanceM.toFixed(0)} m
              {movedDistanceM > 500 ? (
                <span className="ml-2 text-rose-600">(500m 超は信頼ユーザーのみ)</span>
              ) : null}
            </div>
            <TextField
              label="編集メモ"
              value={relocateSummary}
              onChange={setRelocateSummary}
              placeholder="位置がずれていたので修正、など (任意)"
            />
            {relocateErr ? <p className="mt-1 text-xs text-rose-600">{relocateErr}</p> : null}
            <div className="mt-3 flex gap-2">
              <Button
                variant="secondary"
                className="flex-1"
                onPress={cancelRelocate}
                isDisabled={relocateSubmitting}
              >
                キャンセル
              </Button>
              <Button
                className="flex-1"
                onPress={submitRelocate}
                isPending={relocateSubmitting}
                isDisabled={relocateSubmitting || movedDistanceM < 1}
              >
                <Check size={16} aria-hidden /> 確定
              </Button>
            </div>
          </div>
        </>
      ) : null}

      <Sheet
        isOpen={!inRelocate && selected !== null}
        onOpenChange={(open) => {
          if (!open) {
            setSelected(null);
            setSearchClearSignal((n) => n + 1);
          }
        }}
        title={selected?.displayName}
      >
        {selected ? <PianoSummary piano={selected} /> : null}
      </Sheet>

      <Sheet
        isOpen={!inRelocate && draft !== null}
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

function haversineM(lat1: number, lng1: number, lat2: number, lng2: number): number {
  const R = 6371000;
  const toRad = (d: number) => (d * Math.PI) / 180;
  const dLat = toRad(lat2 - lat1);
  const dLng = toRad(lng2 - lng1);
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLng / 2) ** 2;
  return 2 * R * Math.asin(Math.min(1, Math.sqrt(a)));
}

function PianoSummary({ piano }: { piano: Piano }) {
  const id = parsePianoSafe(piano.name);
  const avg = piano.ratingAverage;
  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2 text-sm text-slate-600">
        <span>
          {piano.ratingCount > 0 ? `★ ${avg.toFixed(1)} (${piano.ratingCount}件)` : "未評価"}
        </span>
      </div>
      {piano.description ? (
        <p className="line-clamp-3 text-sm text-slate-700">{piano.description}</p>
      ) : null}
      <PianoAttributeMeters piano={piano} />
      <ul className="space-y-1.5 text-xs text-slate-700">
        {piano.address ? (
          <SummaryRow icon={<MapPin size={14} />} text={piano.address} />
        ) : null}
        <SummaryRow
          icon={<Tag size={14} />}
          text={`${pianoKindLabel(piano.kind)} · ${pianoTypeLabel(piano.pianoType)} / ${piano.pianoBrand}`}
        />
        {piano.pianoModel ? (
          <SummaryRow icon={<Music size={14} />} text={piano.pianoModel} />
        ) : null}
        <SummaryRow
          icon={<CalendarClock size={14} />}
          text={availabilityLabel(piano.availability)}
        />
        {piano.hours ? <SummaryRow icon={<Clock size={14} />} text={piano.hours} /> : null}
      </ul>
      <Link
        to={id ? `/pianos/${id}` : "#"}
        className="block rounded-full bg-amber-500 py-2.5 text-center text-sm font-semibold text-white hover:bg-amber-600"
      >
        詳細を見る →
      </Link>
    </div>
  );
}

function SummaryRow({ icon, text }: { icon: React.ReactNode; text: string }) {
  return (
    <li className="flex items-start gap-2">
      <span className="mt-0.5 shrink-0 text-slate-400">{icon}</span>
      <span className="min-w-0 flex-1 truncate">{text}</span>
    </li>
  );
}

function pianoKindLabel(k: number): string {
  switch (k) {
    case 1:
      return "ストリート";
    case 2:
      return "練習室";
    case 3:
      return "その他";
    default:
      return "—";
  }
}

function pianoTypeLabel(t: number): string {
  switch (t) {
    case 1:
      return "グランド";
    case 2:
      return "アップライト";
    case 3:
      return "電子";
    case 4:
      return "不明";
    default:
      return "—";
  }
}

function availabilityLabel(a: number): string {
  switch (a) {
    case 1:
      return "通年";
    case 2:
      return "不定期";
    case 3:
      return "イベント時のみ";
    case 4:
      return "天候次第";
    default:
      return "—";
  }
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
  const [pianoType, setPianoType] = useState<PianoType | null>(null);
  const [kind, setKind] = useState<PianoKind>(PianoKind.STREET);
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);
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
