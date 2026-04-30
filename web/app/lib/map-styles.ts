// 開発用に使える、API キー不要のタイルスタイル一覧。
// 本番は Protomaps (.pmtiles on R2 + Cloudflare Worker) に切替する。

import type { StyleSpecification } from "maplibre-gl";

export type MapStyleDef = {
  id: string;
  label: string;
  /** ピアノマーカーの背景に対するコントラスト調整用 (light / dark)。 */
  theme: "light" | "dark";
  style: StyleSpecification;
};

const osmAttribution =
  '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors';

const cartoAttribution =
  osmAttribution + ', &copy; <a href="https://carto.com/attributions">CARTO</a>';

const esriAttribution =
  'Tiles &copy; Esri &mdash; Source: Esri, Maxar, Earthstar Geographics, and the GIS User Community';

function rasterStyle(tiles: string[], attribution: string): StyleSpecification {
  return {
    version: 8,
    sources: {
      base: {
        type: "raster",
        tiles,
        tileSize: 256,
        attribution,
      },
    },
    layers: [{ id: "base", type: "raster", source: "base" }],
  };
}

export const MAP_STYLES: MapStyleDef[] = [
  {
    id: "osm",
    label: "OSM Standard",
    theme: "light",
    style: rasterStyle(
      ["https://tile.openstreetmap.org/{z}/{x}/{y}.png"],
      osmAttribution,
    ),
  },
  {
    id: "carto-positron",
    label: "Positron (淡)",
    theme: "light",
    style: rasterStyle(
      [
        "https://a.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png",
        "https://b.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png",
        "https://c.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png",
        "https://d.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png",
      ],
      cartoAttribution,
    ),
  },
  {
    id: "carto-voyager",
    label: "Voyager",
    theme: "light",
    style: rasterStyle(
      [
        "https://a.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png",
        "https://b.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png",
        "https://c.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png",
        "https://d.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png",
      ],
      cartoAttribution,
    ),
  },
  {
    id: "carto-dark",
    label: "Dark Matter",
    theme: "dark",
    style: rasterStyle(
      [
        "https://a.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png",
        "https://b.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png",
        "https://c.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png",
        "https://d.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png",
      ],
      cartoAttribution,
    ),
  },
  {
    id: "esri-satellite",
    label: "衛星写真 (Esri)",
    theme: "dark",
    style: rasterStyle(
      [
        "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
      ],
      esriAttribution,
    ),
  },
  {
    id: "opentopomap",
    label: "地形 (OpenTopoMap)",
    theme: "light",
    style: rasterStyle(
      [
        "https://a.tile.opentopomap.org/{z}/{x}/{y}.png",
        "https://b.tile.opentopomap.org/{z}/{x}/{y}.png",
        "https://c.tile.opentopomap.org/{z}/{x}/{y}.png",
      ],
      osmAttribution +
        ', SRTM | Map style: &copy; <a href="https://opentopomap.org">OpenTopoMap</a> (CC-BY-SA)',
    ),
  },
];

export const DEFAULT_STYLE_ID = "carto-voyager";

const STORAGE_KEY = "piamap.map.style_id";

export function getSavedStyleId(): string {
  if (typeof window === "undefined") return DEFAULT_STYLE_ID;
  return window.localStorage.getItem(STORAGE_KEY) || DEFAULT_STYLE_ID;
}

export function saveStyleId(id: string) {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(STORAGE_KEY, id);
}

export function findStyle(id: string): MapStyleDef {
  return MAP_STYLES.find((s) => s.id === id) || MAP_STYLES[0];
}
