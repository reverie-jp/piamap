import { CenterDialog } from "./ui/sheet";
import { Button } from "./ui/button";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  message?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  /** 確認ボタンを danger スタイルにする (削除など)。 */
  destructive?: boolean;
  onConfirm: () => void | Promise<void>;
};

export function ConfirmDialog({
  isOpen,
  onOpenChange,
  title,
  message,
  confirmLabel = "OK",
  cancelLabel = "キャンセル",
  destructive,
  onConfirm,
}: Props) {
  return (
    <CenterDialog isOpen={isOpen} onOpenChange={onOpenChange} title={title}>
      {message ? <p className="text-sm text-slate-600">{message}</p> : null}
      <div className="mt-4 flex flex-col gap-2">
        <Button
          variant={destructive ? "danger" : "primary"}
          onPress={async () => {
            await onConfirm();
          }}
        >
          {confirmLabel}
        </Button>
        <Button variant="ghost" size="sm" onPress={() => onOpenChange(false)}>
          {cancelLabel}
        </Button>
      </div>
    </CenterDialog>
  );
}
