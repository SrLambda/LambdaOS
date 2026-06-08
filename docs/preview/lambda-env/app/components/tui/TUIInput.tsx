import { useState } from 'react';
import { C } from './tokens';
import { ICONS } from '../../data/icon-map';

interface TUIInputProps {
  label: string;
  description?: string;
  value: string;
  onChange: (v: string) => void;
  placeholder?: string;
  validate?: (v: string) => string | null;
  type?: 'text' | 'password';
  disabled?: boolean;
  focused?: boolean;
}

export function TUIInput({ label, description, value, onChange, placeholder, validate, type = 'text', disabled, focused }: TUIInputProps) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(value);
  const [error, setError] = useState<string | null>(null);
  const [valid, setValid] = useState(false);

  function commit() {
    if (validate) {
      const err = validate(draft);
      if (err) { setError(err); return; }
    }
    setError(null);
    setValid(true);
    onChange(draft);
    setEditing(false);
    setTimeout(() => setValid(false), 1500);
  }

  function cancel() {
    setDraft(value);
    setError(null);
    setEditing(false);
  }

  const borderColor = error ? C.error : valid ? C.success : editing ? C.accent : C.border;

  return (
    <div style={{
      padding: '10px 12px',
      borderBottom: `1px solid ${C.border}`,
      background: focused ? C.accentDim : 'transparent',
      opacity: disabled ? 0.45 : 1,
    }}>
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 16 }}>
        <div style={{ flex: 1 }}>
          <div style={{ color: C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>{label}</div>
          {description && <div style={{ color: C.textSecondary, fontSize: 11, fontFamily: 'monospace', marginTop: 2 }}>{description}</div>}
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 3 }}>
          <div style={{ display: 'flex', alignItems: 'center', border: `1px solid ${borderColor}`, transition: 'border-color 0.15s' }}>
            <span style={{ color: C.textMuted, fontSize: 11, fontFamily: 'monospace', padding: '3px 6px', borderRight: `1px solid ${borderColor}` }}>│</span>
            <input
              type={type}
              value={editing ? draft : value}
              disabled={disabled}
              placeholder={placeholder ?? ''}
              onFocus={() => { setEditing(true); setDraft(value); }}
              onBlur={commit}
              onChange={e => { setDraft(e.target.value); setError(null); }}
              onKeyDown={e => {
                if (e.key === 'Enter') { commit(); e.preventDefault(); }
                if (e.key === 'Escape') { cancel(); e.preventDefault(); }
                e.stopPropagation();
              }}
              style={{
                background: 'transparent',
                border: 'none',
                outline: 'none',
                color: C.textPrimary,
                fontFamily: 'monospace',
                fontSize: 12,
                padding: '3px 8px',
                width: 200,
                caretColor: C.accent,
              }}
            />
            {editing && (
              <span style={{ color: C.textMuted, fontSize: 10, padding: '3px 4px', fontFamily: 'monospace' }}>
                {ICONS.widgets.confirm.nerd}
              </span>
            )}
          </div>
            {error && <div style={{ color: C.error, fontSize: 10, fontFamily: 'monospace' }}>{ICONS.widgets.error.nerd} {error}</div>}
          {valid && <div style={{ color: C.success, fontSize: 10, fontFamily: 'monospace' }}>{ICONS.widgets.success.nerd} Guardado</div>}
        </div>
      </div>
    </div>
  );
}
