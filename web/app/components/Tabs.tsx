type TabId = string | number;

type Tab<T extends TabId> = {
  id: T;
  label: string;
};

type Props<T extends TabId> = {
  tabs: Tab<T>[];
  active: T;
  onChange: (id: T) => void;
};

export function Tabs<T extends TabId>({ tabs, active, onChange }: Props<T>) {
  return (
    <div className="scrollbar-none flex gap-1 overflow-x-auto border-b border-slate-200">
      {tabs.map((t) => {
        const isActive = t.id === active;
        return (
          <button
            key={String(t.id)}
            type="button"
            onClick={() => onChange(t.id)}
            className={
              "shrink-0 cursor-pointer border-b-2 px-3 py-2 text-sm font-medium outline-none transition focus-visible:ring-2 focus-visible:ring-amber-500 " +
              (isActive
                ? "border-amber-500 text-amber-700"
                : "border-transparent text-slate-500 hover:text-slate-800")
            }
            role="tab"
            aria-selected={isActive}
          >
            {t.label}
          </button>
        );
      })}
    </div>
  );
}
