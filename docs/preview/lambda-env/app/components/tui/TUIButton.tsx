import { useState } from 'react';
import { C } from './tokens';
import { ICONS } from '../../data/icon-map';

interface TUIButtonProps {
  label: string;
  onClick: () => Promise<void> | void;
  variant?: 'primary' | 'secondary' | 'danger';
  icon?: string;
  disabled?: boolean;
  fullWidth?: boolean;
}

export function TUIButton({ label, onClick, variant = 'secondary', icon, disabled, fullWidth }: TUIButtonProps) {
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);

  async function handle() {
    if (loading || disabled) return;
    setLoading(true);
    try { await onClick(); } catch {}
    setLoading(false);
    setDone(true);
    setTimeout(() => setDone(false), 1500);
  }

  const bg = done ? C.success : variant === 'primary' ? C.accent : variant === 'danger' ? 'transparent' : 'transparent';
  const textColor = done ? '#000' : variant === 'primary' ? '#000' : variant === 'danger' ? C.error : C.accent;
  const borderColor = done ? C.success : variant === 'primary' ? C.accent : variant === 'danger' ? C.error : C.accentBorder;

  return (
    <button
      onClick={handle}
      disabled={disabled || loading}
      style={{
        background: bg,
        color: textColor,
        border: `1px solid ${borderColor}`,
        fontFamily: 'monospace',
        fontSize: 12,
        padding: '6px 18px',
        cursor: disabled || loading ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.45 : 1,
        transition: 'all 0.15s',
        width: fullWidth ? '100%' : undefined,
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 6,
      }}
    >
      {done ? `${ICONS.widgets.success.nerd} Listo` : loading ? `${label}...` : `${icon ? icon + ' ' : ''}${label}`}
    </button>
  );
}
