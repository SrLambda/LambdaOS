import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIModal } from '../components/tui/TUIModal';

type ServiceState = 'running' | 'stopped' | 'failed';

interface Service {
  name: string; description: string;
  state: ServiceState; enabled: boolean; pid?: number;
}

const INITIAL_SERVICES: Service[] = [
  { name: 'NetworkManager',  description: 'Gestión de red',               state: 'running', enabled: true,  pid: 1203 },
  { name: 'bluetooth',       description: 'Protocolo Bluetooth',           state: 'running', enabled: true,  pid: 1418 },
  { name: 'sshd',            description: 'Servidor OpenSSH',             state: 'running', enabled: true,  pid: 2041 },
  { name: 'docker',          description: 'Motor de contenedores Docker',  state: 'stopped', enabled: false },
  { name: 'postgresql',      description: 'Base de datos PostgreSQL',      state: 'failed',  enabled: true  },
  { name: 'nginx',           description: 'Servidor web / proxy inverso',  state: 'stopped', enabled: false },
  { name: 'cups',            description: 'Sistema de impresión',          state: 'running', enabled: true,  pid: 1876 },
  { name: 'avahi-daemon',    description: 'Descubrimiento de red local',   state: 'running', enabled: true,  pid: 1654 },
  { name: 'cron',            description: 'Planificador de tareas',        state: 'running', enabled: true,  pid: 1203 },
  { name: 'ufw',             description: 'Uncomplicated Firewall',        state: 'stopped', enabled: false },
];

const stateColor = (s: ServiceState) => s === 'running' ? C.success : s === 'failed' ? C.error : C.textMuted;
const stateIcon  = (s: ServiceState) => s === 'running' ? '●' : s === 'failed' ? '⚠' : '○';
const stateLabel = (s: ServiceState) => s === 'running' ? 'running' : s === 'failed' ? 'failed' : 'stopped';

export function ServicesModule({ onModified }: { onModified: () => void }) {
  const [services, setServices] = useState(INITIAL_SERVICES);
  const [selected, setSelected] = useState<string | null>(null);
  const [restartTarget, setRestartTarget] = useState<string | null>(null);

  function updateState(name: string, state: ServiceState) {
    setServices(s => s.map(sv => sv.name === name ? { ...sv, state, pid: state === 'running' ? Math.floor(Math.random() * 9000 + 1000) : undefined } : sv));
    onModified();
  }

  function toggleEnabled(name: string) {
    setServices(s => s.map(sv => sv.name === name ? { ...sv, enabled: !sv.enabled } : sv));
    onModified();
  }

  const running = services.filter(s => s.state === 'running').length;
  const failed  = services.filter(s => s.state === 'failed').length;

  return (
    <div>
      {restartTarget && (
        <TUIModal
          title={`¿Reiniciar ${restartTarget}?`}
          description={`El servicio ${restartTarget} se detendrá y volverá a iniciar. Las conexiones activas se interrumpirán brevemente.`}
          confirmLabel="Reiniciar"
          onConfirm={async () => {
            updateState(restartTarget, 'stopped');
            setTimeout(() => updateState(restartTarget, 'running'), 600);
            setRestartTarget(null);
          }}
          onCancel={() => setRestartTarget(null)}
        />
      )}

      {/* Summary */}
      <div style={{
        padding: '8px 14px', background: C.surface, marginBottom: 4,
        fontFamily: 'monospace', display: 'flex', gap: 20,
      }}>
        <span style={{ color: C.success, fontSize: 11 }}>● {running} activos</span>
        {failed > 0 && <span style={{ color: C.error, fontSize: 11 }}>⚠ {failed} fallidos</span>}
        <span style={{ color: C.textMuted, fontSize: 11 }}>○ {services.length - running - failed} detenidos</span>
      </div>

      <TUISection title="SERVICIOS DEL SISTEMA" rootRequired>
        {/* Header */}
        <div style={{
          display: 'grid', gridTemplateColumns: '160px 1fr 80px 50px auto',
          padding: '5px 12px', background: C.surface,
          gap: 8, fontFamily: 'monospace',
          borderBottom: `1px solid ${C.border}`,
        }}>
          {['Servicio','Descripción','Estado','PID',''].map(h => (
            <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
          ))}
        </div>

        {services.map(svc => {
          const isSel = selected === svc.name;
          return (
            <div key={svc.name}>
              <div
                onClick={() => setSelected(isSel ? null : svc.name)}
                style={{
                  display: 'grid', gridTemplateColumns: '160px 1fr 80px 50px auto',
                  padding: '8px 12px',
                  borderBottom: `1px solid ${C.border}`,
                  gap: 8, fontFamily: 'monospace',
                  cursor: 'pointer',
                  background: isSel ? C.accentDim : 'transparent',
                  alignItems: 'center',
                }}
              >
                <span style={{ color: svc.enabled ? C.textPrimary : C.textSecondary, fontSize: 12 }}>{svc.name}</span>
                <span style={{ color: C.textMuted, fontSize: 11, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{svc.description}</span>
                <span style={{ color: stateColor(svc.state), fontSize: 11 }}>
                  {stateIcon(svc.state)} {stateLabel(svc.state)}
                </span>
                <span style={{ color: C.textMuted, fontSize: 11 }}>{svc.pid ?? '—'}</span>
                <span style={{ color: C.accent, fontSize: 10 }}>{isSel ? '▼' : '▶'}</span>
              </div>

              {/* Contextual actions */}
              {isSel && (
                <div style={{
                  padding: '8px 12px 10px 20px',
                  borderBottom: `1px solid ${C.border}`,
                  background: C.surface,
                  display: 'flex', gap: 8, alignItems: 'center',
                  flexWrap: 'wrap',
                }}>
                  {svc.state !== 'running' && (
                    <TUIButton label="Iniciar" onClick={async () => { updateState(svc.name, 'running'); setSelected(null); }} variant="primary" />
                  )}
                  {svc.state === 'running' && (
                    <TUIButton label="Detener" onClick={async () => { updateState(svc.name, 'stopped'); setSelected(null); }} />
                  )}
                  {svc.state === 'running' && (
                    <TUIButton label="Reiniciar" onClick={async () => setRestartTarget(svc.name)} />
                  )}
                  <TUIButton
                    label={svc.enabled ? 'Deshabilitar' : 'Habilitar'}
                    onClick={async () => { toggleEnabled(svc.name); }}
                    variant={svc.enabled ? 'danger' : 'secondary'}
                  />
                  <span style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', marginLeft: 8 }}>
                    {svc.enabled ? '● Habilitado al inicio' : '○ Deshabilitado al inicio'}
                  </span>
                </div>
              )}
            </div>
          );
        })}
      </TUISection>
    </div>
  );
}
