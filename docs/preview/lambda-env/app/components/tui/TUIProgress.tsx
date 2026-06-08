import { C } from './tokens';

interface TUIProgressProps {
  value: number; // 0-100
  label?: string;
  subLabel?: string;
  color?: string;
}

export function TUIProgress({ value, label, subLabel, color = C.accent }: TUIProgressProps) {
  return (
    <div style={{ padding: '10px 12px', fontFamily: 'monospace' }}>
      {label && (
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
          <span style={{ color: C.textPrimary, fontSize: 12 }}>{label}</span>
          <span style={{ color: color, fontSize: 12 }}>{value}%</span>
        </div>
      )}
      <div style={{ position: 'relative', height: 8, background: C.surface2, border: `1px solid ${C.border}` }}>
        <div style={{
          position: 'absolute', left: 0, top: 0, bottom: 0,
          width: `${value}%`,
          background: color,
          transition: 'width 0.3s ease',
        }} />
        {/* Thumb marker */}
        <div style={{
          position: 'absolute',
          left: `calc(${value}% - 5px)`,
          top: -2, bottom: -2,
          width: 3,
          background: color,
          boxShadow: `0 0 6px ${color}`,
        }} />
      </div>
      {subLabel && <div style={{ color: C.textSecondary, fontSize: 11, marginTop: 4 }}>{subLabel}</div>}
    </div>
  );
}
