import type { ReactNode } from "react";
import { CheckCircle2 } from "lucide-react";

import { CenterDialog } from "./ui/sheet";
import { Button } from "./ui/button";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  message?: string;
  okLabel?: string;
  /** タイトル左に出す装飾アイコン (デフォルト: 緑のチェック)。 */
  icon?: ReactNode;
};

export function MessageDialog({
  isOpen,
  onOpenChange,
  title,
  message,
  okLabel = "OK",
  icon,
}: Props) {
  return (
    <CenterDialog isOpen={isOpen} onOpenChange={onOpenChange} title={title}>
      <div className="flex items-start gap-3">
        <div className="mt-0.5 shrink-0">
          {icon ?? <CheckCircle2 size={24} className="text-emerald-500" />}
        </div>
        {message ? <p className="text-sm text-slate-700">{message}</p> : null}
      </div>
      <div className="mt-4">
        <Button className="w-full" onPress={() => onOpenChange(false)}>
          {okLabel}
        </Button>
      </div>
    </CenterDialog>
  );
}
