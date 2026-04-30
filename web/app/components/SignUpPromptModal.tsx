import { Link } from "react-router";

import { CenterDialog } from "./ui/sheet";
import { Button } from "./ui/button";

type Props = {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  /** 何のアクションをしようとしたか説明 (任意)。 */
  action?: string;
};

// 未認証ユーザーが認証必須のアクションを試みたときに表示。
export function SignUpPromptModal({ isOpen, onOpenChange, action }: Props) {
  return (
    <CenterDialog isOpen={isOpen} onOpenChange={onOpenChange} title="アカウントが必要です">
      <p className="text-sm text-slate-600">
        {action ?? "この操作"}にはアカウントが必要です。
        <br />
        Google アカウントで数秒で始められます。
      </p>
      <div className="mt-4 flex flex-col gap-2">
        <Button onPress={() => alert("Google ログインは Phase 1 で実装予定")}>
          Google でログイン (Phase 1)
        </Button>
        <Link
          to="/settings"
          onClick={() => onOpenChange(false)}
          className="block text-center text-xs text-slate-500 underline-offset-2 hover:underline"
        >
          dev token を入力 (開発用)
        </Link>
        <Button variant="ghost" size="sm" onPress={() => onOpenChange(false)}>
          あとで
        </Button>
      </div>
    </CenterDialog>
  );
}
