import { Bell, MapIcon, ScrollText, UserRound } from "lucide-react";
import { NavLink } from "react-router";

type Tab = {
  to: string;
  label: string;
  icon: React.ComponentType<{ size?: number; className?: string }>;
};

const TABS: Tab[] = [
  { to: "/map", label: "マップ", icon: MapIcon },
  { to: "/timeline", label: "タイムライン", icon: ScrollText },
  { to: "/notifications", label: "通知", icon: Bell },
  { to: "/profile/me", label: "プロフィール", icon: UserRound },
];

// max-w-[390px] のシェル内に固定表示。シェル外に出ないよう absolute + bottom-0 + 親の中央寄せに依存。
export function BottomNav() {
  return (
    <nav
      aria-label="メインナビゲーション"
      className="fixed bottom-0 left-1/2 z-30 w-full max-w-97.5 -translate-x-1/2 border-t border-slate-200 bg-white/95 backdrop-blur"
      style={{ paddingBottom: "env(safe-area-inset-bottom, 0)" }}
    >
      <ul className="grid h-16 grid-cols-4">
        {TABS.map((t) => (
          <li key={t.to} className="contents">
            <NavLink
              to={t.to}
              className={({ isActive }) =>
                `flex flex-col items-center justify-center gap-0.5 text-[10px] font-medium ${
                  isActive ? "text-amber-600" : "text-slate-500 hover:text-slate-700"
                }`
              }
            >
              {({ isActive }) => (
                <>
                  <t.icon size={22} className={isActive ? "fill-amber-100" : ""} />
                  <span>{t.label}</span>
                </>
              )}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  );
}
