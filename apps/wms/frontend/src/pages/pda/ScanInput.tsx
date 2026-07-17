interface Props {
  value: string;
  onChange: (v: string) => void;
  placeholder?: string;
}

export function ScanInput({ value, onChange, placeholder = '扫描条码' }: Props) {
  return (
    <div>
      <label className="text-xs font-medium mb-1 block" style={{ color: 'var(--color-muted)' }}>{placeholder}</label>
      <div className="relative">
        <input
          value={value}
          onChange={e => onChange(e.target.value)}
          className="w-full px-4 py-3 border rounded-lg outline-none text-center focus:ring-2 transition-shadow"
          style={{ borderColor: 'var(--border)', '--tw-ring-color': 'var(--ring)' } as React.CSSProperties}
          placeholder={placeholder}
          autoFocus
        />
      </div>
    </div>
  );
}
