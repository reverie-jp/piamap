import { ChevronDown } from "lucide-react";
import {
  Select as RACSelect,
  SelectValue,
  Label,
  Button,
  Popover,
  ListBox,
  ListBoxItem,
  type SelectProps,
  type ListBoxItemProps,
} from "react-aria-components";

type Props<T extends object> = Omit<SelectProps<T>, "children"> & {
  label?: string;
  items: { id: string | number; label: string }[];
  /** 未選択時に表示する文言。 */
  placeholder?: string;
  /** フィールド直下に赤字で表示するエラーメッセージ。 */
  errorMessage?: string;
  className?: string;
};

const baseTriggerCls =
  "mt-1 inline-flex h-10 w-full items-center justify-between gap-2 rounded-lg border bg-white px-3 text-sm text-slate-900 " +
  "outline-none data-[focus-visible]:ring-2";

export function Select<T extends object>({
  label,
  items,
  placeholder,
  errorMessage,
  className = "",
  ...props
}: Props<T>) {
  const isInvalid = Boolean(errorMessage);
  const triggerCls =
    `${baseTriggerCls} ` +
    (isInvalid
      ? "border-rose-500 data-[focus-visible]:ring-rose-500"
      : "border-slate-300 data-[focus-visible]:ring-amber-500");
  return (
    <RACSelect {...props} isInvalid={isInvalid} className={`block ${className}`}>
      {label ? (
        <Label className="text-sm font-medium text-slate-700">
          {label}
          {props.isRequired ? <span className="ml-0.5 text-rose-600">*</span> : null}
        </Label>
      ) : null}
      <Button className={triggerCls}>
        <SelectValue className="truncate">
          {({ isPlaceholder, defaultChildren }) =>
            isPlaceholder ? (
              <span className="text-slate-400">{placeholder ?? "選択してください"}</span>
            ) : (
              defaultChildren
            )
          }
        </SelectValue>
        <ChevronDown size={16} className="shrink-0 text-slate-500" />
      </Button>
      <Popover className="rounded-lg border border-slate-200 bg-white shadow-lg outline-none">
        <ListBox className="max-h-64 overflow-auto p-1 outline-none">
          {items.map((it) => (
            <ListBoxItem
              key={it.id}
              id={it.id}
              className={({ isFocused, isSelected }) =>
                `cursor-default rounded px-3 py-2 text-sm outline-none ${
                  isFocused ? "bg-amber-50" : ""
                } ${isSelected ? "font-semibold text-amber-700" : "text-slate-800"}`
              }
            >
              {it.label}
            </ListBoxItem>
          ))}
        </ListBox>
      </Popover>
      {errorMessage ? (
        <p className="mt-1 text-xs text-rose-600" role="alert">
          {errorMessage}
        </p>
      ) : null}
    </RACSelect>
  );
}

export type { ListBoxItemProps };
