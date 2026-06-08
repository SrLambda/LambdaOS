import { useEffect } from 'react';
import { C } from './tokens';

interface TUIModalProps {
  title: string;
  description: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: 'default' | 'danger';
  onConfirm: () => void;
  onCancel: () => void;
}

export function TUIModal({ title, description, confirmLabel = 'Confirmar', cancelLabel = 'Cancelar', variant = 'default', onConfirm, onCancel }: TUIModalProps) {
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') { onCancel(); e.stopPropagation(); }
      if (e.key === 'Enter') { onConfirm(); e.stopPropagation(); }
    };
    window.addEventListener('keydown', onKey, true);
    return () => window.removeEventListener('keydown', onKey, true);
  }, []);

  return (
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 1000,
        background: 'rgba(0,0,0,0.6)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
      }}
      onClick={onCancel}
    >
      <div
        onClick={e => e.stopPropagation()}
        style={{
          background: C.surface,
          border: `1px solid ${variant === 'danger' ? C.error : C.accentBorder}`,
          padding: 24,
          minWidth: 340,
          maxWidth: 480,
          fontFamily: 'monospace',
          boxShadow: `0 8px 40px rgba(0,0,0,0.7), 0 0 20px ${variant === 'danger' ? C.error : C.accent}20`,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 12 }}>
          <span style={{ color: variant === 'danger' ? C.error : C.warn, fontSize: 16 }}>⚠</span>
          <span style={{ color: C.textPrimary, fontSize: 14 }}>{title}</span>
        </div>
        <div style={{ color: C.textSecondary, fontSize: 12, marginBottom: 20, lineHeight: 1.6 }}>{description}</div>
        <div style={{ display: 'flex', gap: 10, justifyContent: 'flex-end' }}>
          <button
            onClick={onCancel}
            style={{ background: 'transparent', border: `1px solid ${C.border}`, color: C.textSecondary, fontFamily: 'monospace', fontSize: 12, padding: '6px 16px', cursor: 'pointer' }}
          >
            {cancelLabel}
          </button>
          <button
            onClick={onConfirm}
            style={{ background: variant === 'danger' ? C.error : C.accent, border: 'none', color: '#000', fontFamily: 'monospace', fontSize: 12, padding: '6px 16px', cursor: 'pointer', fontWeight: 'bold' }}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
