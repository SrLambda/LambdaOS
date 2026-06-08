import { C } from '../components/tui/tokens';

export function StubModule({ name }: { name: string }) {
  return (
    <div style={{
      padding: '40px 24px', textAlign: 'center', fontFamily: 'monospace',
    }}>
      <div style={{ color: C.textMuted, fontSize: 24, marginBottom: 16 }}>◇</div>
      <div style={{ color: C.textSecondary, fontSize: 14, marginBottom: 8 }}>{name}</div>
      <div style={{ color: C.textMuted, fontSize: 11 }}>Módulo en desarrollo · Próximamente en lambda-env</div>
      <div style={{ color: C.textMuted, fontSize: 10, marginTop: 16 }}>
        — El módulo completo aparecerá en una futura versión —
      </div>
    </div>
  );
}
