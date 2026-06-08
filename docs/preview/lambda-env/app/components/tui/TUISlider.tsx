import { useState } from 'react';
import { C } from './tokens';

interface TUISliderProps {
  label: string;
  description?: string;
  value: number;
  min?: number;
  max?: number;
  unit?: string;
  onChange: (v: number) => void;
  disabled?: boolean;
  focused?: boolean;
}

export function TUISlider({ label, description, value, min = 0, max = 100, unit = '%', onChange, disabled, focused }: TUISliderProps) {
  const pct = ((value - min) / (max - min)) * 100;
  const [dragging, setDragging] = useState(false);

  return (
    <div style={{
      padding: '10px 12px',
      borderBottom: `1px solid ${C.border}`,
      background: focused ? C.accentDim : 'transparent',
      opacity: disabled ? 0.45 : 1,
    }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 6 }}>
        <div>
          <div style={{ color: C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>{label}</div>
          {description && <div style={{ color: C.textSecondary, fontSize: 11, fontFamily: 'monospace', marginTop: 2 }}>{description}</div>}
        </div>
        <span style={{
          color: C.accent,
          fontSize: 12,
          fontFamily: 'monospace',
          minWidth: 48,
          textAlign: 'right',
          border: `1px solid ${C.border}`,
          padding: '1px 6px',
        }}>{value}{unit}</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <span style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', width: 20, textAlign: 'right' }}>{min}</span>
        <div style={{ flex: 1, position: 'relative', height: 20, display: 'flex', alignItems: 'center' }}>
          {/* Track */}
          <div style={{ position: 'absolute', left: 0, right: 0, height: 3, background: C.surface2, border: `1px solid ${C.border}` }}>
            <div style={{ width: `${pct}%`, height: '100%', background: C.accent }} />
          </div>
          {/* Thumb */}
          <div style={{
            position: 'absolute',
            left: `calc(${pct}% - 7px)`,
            width: 14,
            height: 14,
            border: `2px solid ${C.accent}`,
            background: dragging ? C.accent : C.bg,
            boxShadow: `0 0 6px ${C.accent}80`,
            cursor: 'pointer',
          }} />
          {/* Native range on top for interaction */}
          <input
            type="range"
            min={min}
            max={max}
            value={value}
            disabled={disabled}
            onMouseDown={() => setDragging(true)}
            onMouseUp={() => setDragging(false)}
            onChange={e => onChange(parseInt(e.target.value))}
            style={{
              position: 'absolute', left: 0, right: 0, width: '100%',
              opacity: 0, cursor: 'pointer', height: 20, margin: 0,
            }}
          />
        </div>
        <span style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', width: 20 }}>{max}</span>
      </div>
    </div>
  );
}
