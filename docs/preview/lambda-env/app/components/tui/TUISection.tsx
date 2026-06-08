import { useState } from 'react';
import { C } from './tokens';
import { ICONS } from '../../data/icon-map';

interface TUISectionProps {
  title: string;
  collapsible?: boolean;
  defaultOpen?: boolean;
  children: React.ReactNode;
  rootRequired?: boolean;
}

export function TUISection({ title, collapsible = false, defaultOpen = true, children, rootRequired }: TUISectionProps) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <div style={{ marginBottom: 4 }}>
      <div
        onClick={() => collapsible && setOpen(o => !o)}
        style={{
          display: 'flex', alignItems: 'center', gap: 8,
          padding: '6px 12px',
          background: C.surface,
          borderLeft: `3px solid ${C.accent}`,
          cursor: collapsible ? 'pointer' : 'default',
          userSelect: 'none',
        }}
      >
        {collapsible && (
          <span style={{ color: C.accent, fontSize: 10, fontFamily: 'monospace' }}>{open ? '▼' : '▶'}</span>
        )}
        <span style={{ color: C.accent, fontSize: 11, fontFamily: 'monospace', letterSpacing: '0.08em' }}>{title}</span>
        {rootRequired && (
          <span style={{ marginLeft: 'auto', color: C.error, fontSize: 10, fontFamily: 'monospace' }}>{ICONS.widgets.lock.nerd} root</span>
        )}
      </div>
      {(!collapsible || open) && (
        <div>{children}</div>
      )}
    </div>
  );
}
