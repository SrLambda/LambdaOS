import { useEffect, useRef, useState } from 'react';
import { C } from './tokens';

interface TUISelectProps {
  label: string;
  description?: string;
  value: string;
  options: string[];
  onChange: (v: string) => void;
  disabled?: boolean;
  focused?: boolean;
}

export function TUISelect({ label, description, value, options, onChange, disabled, focused }: TUISelectProps) {
  const [open, setOpen] = useState(false);
  const [highlightIdx, setHighlightIdx] = useState(() => options.indexOf(value));
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const idx = options.indexOf(value);
    setHighlightIdx(idx >= 0 ? idx : 0);
  }, [open]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'ArrowDown') { setHighlightIdx(i => Math.min(i + 1, options.length - 1)); e.stopPropagation(); e.preventDefault(); }
      if (e.key === 'ArrowUp') { setHighlightIdx(i => Math.max(i - 1, 0)); e.stopPropagation(); e.preventDefault(); }
      if (e.key === 'Enter') { onChange(options[highlightIdx]); setOpen(false); e.stopPropagation(); }
      if (e.key === 'Escape') { setOpen(false); e.stopPropagation(); }
    };
    window.addEventListener('keydown', onKey, true);
    return () => window.removeEventListener('keydown', onKey, true);
  }, [open, highlightIdx, options]);

  useEffect(() => {
    if (!open) return;
    const onClickOut = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', onClickOut);
    return () => document.removeEventListener('mousedown', onClickOut);
  }, [open]);

  return (
    <div ref={ref} style={{ position: 'relative', borderBottom: `1px solid ${C.border}`, background: focused ? C.accentDim : 'transparent' }}>
      <div
        onClick={() => !disabled && setOpen(o => !o)}
        style={{
          display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between',
          padding: '10px 12px',
          cursor: disabled ? 'not-allowed' : 'pointer',
          opacity: disabled ? 0.45 : 1,
        }}
      >
        <div style={{ flex: 1 }}>
          <div style={{ color: C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>{label}</div>
          {description && <div style={{ color: C.textSecondary, fontSize: 11, fontFamily: 'monospace', marginTop: 2 }}>{description}</div>}
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginLeft: 16 }}>
          <span style={{ color: C.accent, fontSize: 11, fontFamily: 'monospace', maxWidth: 200, textAlign: 'right' }}>{value}</span>
          <span style={{ color: C.textSecondary, fontSize: 10, fontFamily: 'monospace' }}>{open ? '▼' : '▶'}</span>
        </div>
      </div>

      {open && (
        <div style={{
          position: 'absolute', right: 12, top: '100%', zIndex: 100,
          background: C.surface, border: `1px solid ${C.accentBorder}`,
          minWidth: 220, maxHeight: 200, overflowY: 'auto',
          boxShadow: `0 4px 20px rgba(0,0,0,0.6), 0 0 10px ${C.accent}20`,
        }}>
          {options.map((opt, i) => (
            <div
              key={opt}
              onMouseEnter={() => setHighlightIdx(i)}
              onClick={() => { onChange(opt); setOpen(false); }}
              style={{
                padding: '7px 12px',
                fontFamily: 'monospace',
                fontSize: 12,
                cursor: 'pointer',
                background: i === highlightIdx ? C.accent : opt === value ? C.accentDim : 'transparent',
                color: i === highlightIdx ? '#000' : opt === value ? C.accent : C.textPrimary,
                borderLeft: opt === value && i !== highlightIdx ? `2px solid ${C.accent}` : '2px solid transparent',
              }}
            >
              {opt === value && i !== highlightIdx ? `● ${opt}` : opt}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
