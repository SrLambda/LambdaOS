import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIInput } from '../components/tui/TUIInput';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIModal } from '../components/tui/TUIModal';

const FW_RULES = [
  { id: 1, rule: 'ALLOW  TCP  22    # SSH',      active: true  },
  { id: 2, rule: 'ALLOW  TCP  80    # HTTP',     active: true  },
  { id: 3, rule: 'ALLOW  TCP  443   # HTTPS',    active: true  },
  { id: 4, rule: 'ALLOW  UDP  53    # DNS',      active: true  },
  { id: 5, rule: 'DROP   TCP  23    # Telnet',   active: true  },
  { id: 6, rule: 'DROP   ALL  0     # Default',  active: true  },
];

const SSH_KEYS = [
  { name: 'id_ed25519', type: 'ED25519', comment: 'user@lambda-ws', added: '2025-11-03' },
  { name: 'id_rsa',     type: 'RSA 4096',comment: 'backup-key',     added: '2024-06-15' },
];

function ls(k: string, d: any) { try { const v = localStorage.getItem('sec_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('sec_'+k, JSON.stringify(v)); }

export function SecurityModule({ onModified }: { onModified: () => void }) {
  const [fw, setFw]           = useState(() => ls('fw', true));
  const [ssh, setSsh]         = useState(() => ls('ssh', true));
  const [luks, setLuks]       = useState(() => ls('luks', true));
  const [apparmor, setApparmor] = useState(() => ls('apparmor', true));
  const [sudo2fa, setSudo2fa] = useState(() => ls('sudo2fa', false));
  const [failLock, setFailLock] = useState(() => ls('faillock', true));
  const [deleteKey, setDeleteKey] = useState<string | null>(null);

  function update(key: string, val: any, setter: (v: any) => void) {
    setter(val); ss(key, val); onModified();
  }

  const levelIdx = [fw, ssh, luks, apparmor].filter(Boolean).length;
  const levels   = ['CRÍTICO','BAJO','MEDIO','ALTO','MÁXIMO'];
  const levelColors = [C.error, C.error, C.warn, C.accent, C.success];

  return (
    <div>
      {deleteKey && (
        <TUIModal
          title={`¿Eliminar clave "${deleteKey}"?`}
          description="Esta acción es irreversible. La clave SSH será eliminada permanentemente del sistema."
          confirmLabel="Eliminar"
          variant="danger"
          onConfirm={() => { setDeleteKey(null); onModified(); }}
          onCancel={() => setDeleteKey(null)}
        />
      )}

      {/* Security level indicator */}
      <div style={{
        padding: '10px 14px', marginBottom: 4,
        background: C.surface,
        borderLeft: `3px solid ${levelColors[levelIdx]}`,
        fontFamily: 'monospace',
        display: 'flex', alignItems: 'center', gap: 12,
      }}>
        <span style={{ color: C.textSecondary, fontSize: 11 }}>NIVEL DE SEGURIDAD</span>
        {levels.map((l, i) => (
          <span key={l} style={{
            fontSize: 10, padding: '2px 8px',
            border: `1px solid ${i === levelIdx ? levelColors[i] : C.border}`,
            background: i === levelIdx ? levelColors[i] : 'transparent',
            color: i === levelIdx ? '#000' : C.textMuted,
          }}>{l}</span>
        ))}
      </div>

      <TUISection title="SISTEMA" collapsible defaultOpen={true}>
        <TUIToggle label="Firewall (iptables/nftables)" description="Filtrado de paquetes de red activo" value={fw} onChange={v => update('fw', v, setFw)} />
        <TUIToggle label="AppArmor" description="Perfiles de confinamiento para aplicaciones" value={apparmor} onChange={v => update('apparmor', v, setApparmor)} />
        <TUIToggle label="Cifrado de Disco (LUKS)" description="Full-disk encryption · AES-256-XTS" value={luks} onChange={v => update('luks', v, setLuks)} />
        <TUIToggle label="Bloqueo por intentos fallidos" description="faillock: bloquear tras 5 fallos en 15 min" value={failLock} onChange={v => update('faillock', v, setFailLock)} />
        <TUIToggle label="2FA para sudo" description="TOTP/FIDO2 al escalar privilegios" value={sudo2fa} onChange={v => update('sudo2fa', v, setSudo2fa)} />
      </TUISection>

      <TUISection title="SSH" collapsible defaultOpen={true} rootRequired>
        <TUIToggle label="Servidor SSH" description="sshd.service · Puerto 22" value={ssh} onChange={v => update('ssh', v, setSsh)} />
        <TUISelect label="Autenticación" description="Método de autenticación permitido"
          value="Solo claves públicas"
          options={['Solo claves públicas','Contraseña + clave','Solo contraseña']}
          onChange={() => { onModified(); }}
        />
        <TUISelect label="Protocolo" description="Versión del protocolo SSH"
          value="SSH-2 únicamente"
          options={['SSH-2 únicamente','SSH-1 + SSH-2']}
          onChange={() => { onModified(); }}
        />

        {/* SSH Keys table */}
        <div style={{ margin: '8px 0 0 0' }}>
          <div style={{
            padding: '4px 12px', background: C.surface,
            display: 'grid', gridTemplateColumns: '1fr 90px 140px 90px auto',
            gap: 8, fontFamily: 'monospace',
          }}>
            {['Nombre','Tipo','Comentario','Añadida',''].map(h => (
              <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
            ))}
          </div>
          {SSH_KEYS.map(key => (
            <div key={key.name} style={{
              padding: '6px 12px', borderBottom: `1px solid ${C.border}`,
              display: 'grid', gridTemplateColumns: '1fr 90px 140px 90px auto',
              gap: 8, fontFamily: 'monospace', alignItems: 'center',
            }}>
              <span style={{ color: C.accent, fontSize: 12 }}>{key.name}</span>
              <span style={{ color: C.textSecondary, fontSize: 11 }}>{key.type}</span>
              <span style={{ color: C.textMuted, fontSize: 11 }}>{key.comment}</span>
              <span style={{ color: C.textMuted, fontSize: 10 }}>{key.added}</span>
              <TUIButton label="Eliminar" onClick={async () => setDeleteKey(key.name)} variant="danger" />
            </div>
          ))}
          <div style={{ padding: '8px 12px' }}>
            <TUIButton label="Generar nueva clave" onClick={async () => { onModified(); }} icon="+" />
          </div>
        </div>
      </TUISection>

      <TUISection title="FIREWALL — REGLAS" collapsible defaultOpen={false} rootRequired>
        {FW_RULES.map(rule => (
          <div key={rule.id} style={{
            display: 'flex', alignItems: 'center', gap: 10,
            padding: '6px 12px', borderBottom: `1px solid ${C.border}`,
            fontFamily: 'monospace',
          }}>
            <span style={{ color: rule.active ? C.success : C.textMuted, fontSize: 10 }}>{rule.active ? '●' : '○'}</span>
            <span style={{ flex: 1, color: C.textSecondary, fontSize: 11 }}>{rule.rule}</span>
          </div>
        ))}
        <div style={{ padding: '8px 12px' }}>
          <TUIButton label="Añadir regla" onClick={async () => { onModified(); }} icon="+" />
        </div>
      </TUISection>
    </div>
  );
}
