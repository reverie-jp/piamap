import { Loader2 } from "lucide-react";
import {
  Button as RACButton,
  composeRenderProps,
  type ButtonProps,
} from "react-aria-components";

type Variant = "primary" | "secondary" | "ghost" | "danger";
type Size = "sm" | "md" | "lg";

const base =
  "inline-flex items-center justify-center gap-2 rounded-full font-semibold transition select-none cursor-pointer " +
  "outline-none data-[focus-visible]:ring-2 data-[focus-visible]:ring-amber-500 data-[focus-visible]:ring-offset-2 " +
  // 押下中の視覚フィードバック
  "data-[pressed]:scale-[0.97] data-[pressed]:opacity-90 " +
  // 無効/loading 中はクリック不可
  "data-[disabled]:opacity-50 data-[disabled]:cursor-not-allowed " +
  "data-[pending]:cursor-progress";

const variants: Record<Variant, string> = {
  primary: "bg-amber-500 text-white hover:bg-amber-600 shadow",
  secondary: "bg-white text-slate-900 border border-slate-200 hover:bg-slate-50 shadow-sm",
  ghost: "bg-transparent text-slate-700 hover:bg-slate-100",
  danger: "bg-rose-600 text-white hover:bg-rose-700 shadow",
};

const sizes: Record<Size, string> = {
  sm: "h-8 px-3 text-xs",
  md: "h-10 px-4 text-sm",
  lg: "h-12 px-5 text-base",
};

export function Button({
  variant = "primary",
  size = "md",
  className = "",
  children,
  ...props
}: ButtonProps & { variant?: Variant; size?: Size }) {
  return (
    <RACButton
      {...props}
      className={`${base} ${variants[variant]} ${sizes[size]} ${className}`}
    >
      {composeRenderProps(children, (resolved, { isPending }) =>
        isPending ? <Loader2 size={18} className="animate-spin" aria-hidden /> : resolved,
      )}
    </RACButton>
  );
}
