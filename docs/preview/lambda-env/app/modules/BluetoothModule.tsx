import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIModal } from '../components/tui/TUIModal';

interface BtDevice {
  mac: string; name: string; type: string; icon: string;
  state: 'connected' | 'paired' | 'available';
  battery?: number;
}

const DEVICES: BtDevice[] = [
  { mac: '00:11:22:33:44:55', name: 'Sony WH-1000XM5', type: 'Auriculares', icon: '🎧', state: 'connected', battery: 82 },
  { mac: 'AA:BB:CC:DD:EE:FF', name: 'Logitech MX Keys', type: 'Teclado',    icon: '⌨',  state: 'paired' },
  { mac: '11:22:33:AA:BB:CC', name: 'Apple Magic Mouse', type: 'Ratón',     icon: '◎',  state: 'paired' },
  { mac: 'FF:EE:DD:CC:BB:AA', name: 'Jabra Speak 510',  type: 'Altavoz',   icon: '♪',  state: 'available' },
  { mac: '12:34:56:78:90:AB', name: 'Samsung Galaxy S24',type: 'Teléfono',  icon: '◱',  state: 'available' },
];

function ls(k: string, d: any) { try { const v = localStorage.getItem('bt_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }

export function BluetoothModule({ onModified }: { onModified: () => void }) {
  const [btOn, setBtOn]       = useState(() => ls('on', true));
  const [visible, setVisible] = useState(() => ls('visible', true));
  const [scanning, setScanning] = useState(false);
  const [selected, setSelected] = useState<string | null>(null);
  const [forgetTarget, setForgetTarget] = useState<BtDevice | null>(null);
  const [devices, setDevices] = useState(DEVICES);

  async function scan() {
    setScanning(true);
    await new Promise(r => setTimeout(r, 2000));
    setScanning(false);
  }

  function disconnect(mac: string) {
    setDevices(d => d.map(dev => dev.mac === mac ? { ...dev, state: 'paired' } : dev));
    onModified();
  }

  function connect(mac: string) {
    setDevices(d => d.map(dev => dev.mac === mac ? { ...dev, state: 'connected' } : dev));
    setSelected(null);
    onModified();
  }

  function forget(mac: string) {
    setDevices(d => d.map(dev => dev.mac === mac ? { ...dev, state: 'available' } : dev));
    setForgetTarget(null);
    onModified();
  }

  const stateColor = (s: BtDevice['state']) =>
    s === 'connected' ? C.success : s === 'paired' ? C.accent : C.textMuted;
  const stateLabel = (s: BtDevice['state']) =>
    s === 'connected' ? '● Conectado' : s === 'paired' ? '○ Emparejado' : '· Disponible';

  return (
    <div>
      {forgetTarget && (
        <TUIModal
          title={`¿Olvidar "${forgetTarget.name}"?`}
          description={`Se eliminará "${forgetTarget.name}" de los dispositivos conocidos. Tendrás que emparejarlo de nuevo para usarlo.`}
          confirmLabel="Sí, olvidar"
          variant="danger"
          onConfirm={() => forget(forgetTarget.mac)}
          onCancel={() => setForgetTarget(null)}
        />
      )}

      <TUISection title="BLUETOOTH">
        <TUIToggle
          label="Bluetooth"
          description="Controlador: hci0"
          value={btOn}
          onChange={v => { setBtOn(v); localStorage.setItem('bt_on', JSON.stringify(v)); onModified(); }}
        />
        <TUIToggle
          label="Visible para otros"
          description="Permitir que otros dispositivos detecten este equipo"
          value={visible}
          disabled={!btOn}
          onChange={v => { setVisible(v); localStorage.setItem('bt_visible', JSON.stringify(v)); onModified(); }}
        />
      </TUISection>

      <TUISection title="DISPOSITIVOS">
        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${C.border}` }}>
          <TUIButton
            label={scanning ? 'Buscando dispositivos...' : 'Buscar dispositivos'}
            onClick={scan}
            icon="⊛"
            disabled={!btOn}
          />
        </div>

        {devices.map(dev => {
          const isSel = selected === dev.mac;
          return (
            <div key={dev.mac}>
              <div
                onClick={() => btOn && setSelected(isSel ? null : dev.mac)}
                style={{
                  display: 'flex', alignItems: 'center', gap: 10,
                  padding: '10px 14px',
                  borderBottom: `1px solid ${C.border}`,
                  cursor: btOn ? 'pointer' : 'not-allowed',
                  background: dev.state === 'connected' ? C.successDim : isSel ? C.accentDim : 'transparent',
                  opacity: btOn ? 1 : 0.4,
                }}
              >
                <span style={{ fontSize: 16 }}>{dev.icon}</span>
                <div style={{ flex: 1 }}>
                  <div style={{ color: C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>{dev.name}</div>
                  <div style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace' }}>{dev.type} · {dev.mac}</div>
                </div>
                {dev.battery !== undefined && (
                  <span style={{ color: C.textSecondary, fontSize: 10, fontFamily: 'monospace' }}>{dev.battery}%</span>
                )}
                <span style={{ color: stateColor(dev.state), fontSize: 11, fontFamily: 'monospace' }}>
                  {stateLabel(dev.state)}
                </span>
                <span style={{ color: C.accent, fontSize: 10, fontFamily: 'monospace' }}>{isSel ? '▼' : '▶'}</span>
              </div>

              {/* Contextual actions */}
              {isSel && (
                <div style={{
                  padding: '8px 14px 10px 40px',
                  borderBottom: `1px solid ${C.border}`,
                  background: C.surface,
                  display: 'flex', gap: 8,
                }}>
                  {dev.state === 'connected' && (
                    <TUIButton label="Desconectar" onClick={async () => disconnect(dev.mac)} />
                  )}
                  {dev.state === 'paired' && (
                    <TUIButton label="Conectar" onClick={async () => connect(dev.mac)} variant="primary" />
                  )}
                  {dev.state === 'available' && (
                    <TUIButton label="Emparejar" onClick={async () => connect(dev.mac)} variant="primary" />
                  )}
                  {(dev.state === 'connected' || dev.state === 'paired') && (
                    <TUIButton label="Olvidar" onClick={async () => setForgetTarget(dev)} variant="danger" />
                  )}
                </div>
              )}
            </div>
          );
        })}
      </TUISection>
    </div>
  );
}
