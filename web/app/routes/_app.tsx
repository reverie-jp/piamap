import { Outlet } from "react-router";

import { BottomNav } from "../components/BottomNav";
import { MobileShell } from "../components/MobileShell";

// ボトムナビ付きシェル。子ルート (map / timeline / notifications / profile/me) に共有される。
export default function AppLayout() {
  return (
    <MobileShell withBottomNav flush>
      <Outlet />
      <BottomNav />
    </MobileShell>
  );
}
