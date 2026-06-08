import { C } from './tokens';
import { ICONS } from '../../data/icon-map';

interface TUIToggleProps {
  label: string;
  description?: string;
  value: boolean;
  onChange: (v: boolean) => void;
  disabled?: boolean;
  requireConfirm?: boolean;
  onConfirmRequest?: (cb: () => void) => void;
  focused?: boolean;
}

export function TUIToggle({ label, description, value, onChange, disabled, focused }: TUIToggleProps) {
  return (
    <div
      onClick={() => !disabled && onChange(!value)}
      style={{
        display: 'flex',
        alignItems: 'flex-start',
        justifyContent: 'space-between',
        padding: '10px 12px',
        borderBottom: `1px solid ${C.border}`,
        background: focused ? C.accentDim : 'transparent',
        cursor: disabled ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.45 : 1,
        transition: 'background 0.1s',
      }}
    >
      <div style={{ flex: 1 }}>
        <div style={{ color: C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>{label}</div>
        {description && (
          <div style={{ color: C.textSecondary, fontSize: 11, fontFamily: 'monospace', marginTop: 2 }}>{description}</div>
        )}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginLeft: 16 }}>
        <span style={{
          fontSize: 13,
          color: value ? C.success : C.textSecondary,
          fontFamily: 'monospace',
        }}>
          {value ? ICONS.widgets.toggle_on.nerd : ICONS.widgets.toggle_off.nerd}
        </span>
        <span style={{
          fontSize: 11,
          color: value ? C.success : C.textSecondary,
          fontFamily: 'monospace',
          minWidth: 74,
        }}>
          {value ? 'Activado' : 'Desactivado'}
        </span>
      </div>
    </div>
  );
}
