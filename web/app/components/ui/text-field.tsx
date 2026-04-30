import {
  TextField as RACTextField,
  Label,
  Input,
  TextArea,
  type TextFieldProps,
} from "react-aria-components";

type Props = TextFieldProps & {
  label?: string;
  placeholder?: string;
  multiline?: boolean;
  rows?: number;
  /** Input の type 属性 (date / time / email など)。 */
  type?: string;
  /** フィールド直下に赤字で表示するエラーメッセージ。 */
  errorMessage?: string;
};

const baseInputCls =
  "mt-1 w-full rounded-lg border bg-white px-3 py-2 text-sm text-slate-900 " +
  "outline-none focus:ring-1 disabled:bg-slate-100";

export function TextField({
  label,
  placeholder,
  multiline,
  rows = 3,
  type,
  errorMessage,
  ...props
}: Props) {
  const isInvalid = Boolean(errorMessage);
  const inputCls =
    `${baseInputCls} ` +
    (isInvalid
      ? "border-rose-500 focus:border-rose-500 focus:ring-rose-500"
      : "border-slate-300 focus:border-amber-500 focus:ring-amber-500");

  return (
    <RACTextField {...props} isInvalid={isInvalid} className="block">
      {label ? (
        <Label className="text-sm font-medium text-slate-700">
          {label}
          {props.isRequired ? <span className="ml-0.5 text-rose-600">*</span> : null}
        </Label>
      ) : null}
      {multiline ? (
        <TextArea placeholder={placeholder} rows={rows} className={`${inputCls} resize-none`} />
      ) : (
        <Input placeholder={placeholder} type={type} className={inputCls} />
      )}
      {errorMessage ? (
        <p className="mt-1 text-xs text-rose-600" role="alert">
          {errorMessage}
        </p>
      ) : null}
    </RACTextField>
  );
}
