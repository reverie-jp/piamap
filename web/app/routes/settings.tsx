import { useEffect, useState } from "react";
import { Link } from "react-router";
import { ArrowLeft } from "lucide-react";

import { MobileShell } from "../components/MobileShell";
import { Button } from "../components/ui/button";
import { TextField } from "../components/ui/text-field";
import { Select } from "../components/ui/select";
import { getAccessToken, setAccessToken, useAuth } from "../lib/auth";
import { MAP_STYLES, getSavedStyleId, saveStyleId } from "../lib/map-styles";

import type { Route } from "./+types/settings";

export function meta({}: Route.MetaArgs) {
  return [{ title: "設定 — PiaMap" }];
}

export default function Settings() {
  const { authed } = useAuth();
  const [token, setToken] = useState("");
  const [styleId, setStyleId] = useState("");

  useEffect(() => {
    setStyleId(getSavedStyleId());
    setToken(getAccessToken() ?? "");
  }, []);

  return (
    <MobileShell>
      <header className="flex items-center gap-2 pb-4">
        <Link to="/profile/me" aria-label="戻る" className="text-slate-500 hover:text-slate-700">
          <ArrowLeft size={22} />
        </Link>
        <h1 className="text-base font-bold text-slate-900">設定</h1>
      </header>

      <Section title="マップ">
        <Select
          label="タイルスタイル"
          selectedKey={styleId}
          onSelectionChange={(key) => {
            const id = String(key);
            setStyleId(id);
            saveStyleId(id);
          }}
          items={MAP_STYLES.map((s) => ({ id: s.id, label: s.label }))}
        />
        <p className="mt-1 text-xs text-slate-500">
          MVP は OSM/Cartocdn/Esri の dev タイル。本番は Protomaps + R2 に切替予定。
        </p>
      </Section>

      <Section title="開発用">
        <p className="text-xs text-slate-500">
          <code className="rounded bg-slate-100 px-1.5 py-0.5">make genjwt &lt;custom_id&gt;</code> で発行した
          access token をここに貼ると認証されます。Phase 1 で Google ログインに置き換え。
        </p>
        <div className="mt-3 space-y-2">
          <TextField
            label="access token"
            value={token}
            onChange={setToken}
            placeholder="eyJhbGciOi..."
          />
          <div className="flex gap-2">
            <Button
              size="sm"
              onPress={() => {
                setAccessToken(token || null);
              }}
              isDisabled={!token}
            >
              保存
            </Button>
            {authed ? (
              <Button
                size="sm"
                variant="secondary"
                onPress={() => {
                  setAccessToken(null);
                  setToken("");
                }}
              >
                ログアウト
              </Button>
            ) : null}
            <span className={`ml-auto self-center text-xs ${authed ? "text-emerald-600" : "text-slate-400"}`}>
              {authed ? "● 認証済み" : "○ 未認証"}
            </span>
          </div>
        </div>
      </Section>
    </MobileShell>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="mb-6">
      <h2 className="mb-2 text-xs font-bold uppercase tracking-wide text-slate-500">{title}</h2>
      <div className="rounded-2xl border border-slate-200 bg-white p-4">{children}</div>
    </section>
  );
}
