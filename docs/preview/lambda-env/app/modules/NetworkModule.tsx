import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { ICONS } from '../data/icon-map';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUIInput } from '../components/tui/TUIInput';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';

interface Network { ssid: string; signal: number; security: string; saved: boolean; }

const NETWORKS: Network[] = [
  { ssid: 'LambdaNet-5G',   signal: 92, security: 'WPA3', saved: true  },
  { ssid: 'Office-Corp',    signal: 74, security: 'WPA2', saved: true  },
  { ssid: 'CaféPublico',    signal: 61, security: 'WPA2', saved: false },
  { ssid: 'IoT_Devices',    signal: 48, security: 'WPA2', saved: false },
  { ssid: 'Hotspot-Mobile', signal: 35, security: 'WPA3', saved: false },
  { ssid: 'Guest_Open',     signal: 28, security: 'OPEN', saved: false },
];

function SignalBars({ pct }: { pct: number }) {
  const filled = Math.ceil(pct / 20);
  return (
    <span style={{ display: 'inline-flex', alignItems: 'flex-end', gap: 1, height: 12 }}>
      {[1,2,3,4,5].map(i => (
        <span key={i} style={{
          width: 3, display: 'inline-block',
          height: `${i * 20}%`,
          background: i <= filled
            ? (pct > 60 ? C.success : pct > 30 ? C.warn : C.error)
            : C.surface2,
        }} />
      ))}
    </span>
  );
}

function ls(k: string, d: any) { try { const v = localStorage.getItem('net_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('net_'+k, JSON.stringify(v)); }

export function NetworkModule({ onModified }: { onModified: () => void }) {
  const [wifiOn, setWifiOn]     = useState(() => ls('wifi', true));
  const [ethOn, setEthOn]       = useState(() => ls('eth', true));
  const [vpnOn, setVpnOn]       = useState(() => ls('vpn', false));
  const [connected, setConnected] = useState('LambdaNet-5G');
  const [selected, setSelected] = useState<string | null>(null);
  const [password, setPassword] = useState('');
  const [scanning, setScanning] = useState(false);
  const [dns, setDns]           = useState(() => ls('dns', '1.1.1.1'));
  const [proxy, setProxy]       = useState(() => ls('proxy', ''));

  function toggle(key: string, val: boolean, setter: (v: boolean) => void) {
    setter(val); ss(key, val); onModified();
  }

  async function scan() {
    setScanning(true);
    await new Promise(r => setTimeout(r, 1800));
    setScanning(false);
  }

  async function connect(net: Network) {
    if (net.security === 'OPEN' || net.saved) {
      setConnected(net.ssid);
      setSelected(null);
      onModified();
    }
    // With password it stays open until user fills it
  }

  return (
    <div>
      {/* Connection info card */}
      <div style={{
        margin: '0 0 4px 0',
        padding: '10px 14px',
        background: C.surface,
        borderLeft: `3px solid ${C.success}`,
        display: 'grid', gridTemplateColumns: '1fr 1fr',
        gap: '4px 24px',
        fontFamily: 'monospace',
      }}>
        <div style={{ color: C.success, fontSize: 11, gridColumn: '1/-1', marginBottom: 4 }}>
          ● CONECTADO · {connected}
        </div>
        {[
          ['IP', '192.168.1.105'],
          ['Gateway', '192.168.1.1'],
          ['DNS', dns],
          ['MAC', '3C:A9:F4:12:7B:E2'],
        ].map(([k, v]) => (
          <div key={k} style={{ display: 'flex', gap: 8 }}>
            <span style={{ color: C.textMuted, fontSize: 11 }}>{k}:</span>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{v}</span>
          </div>
        ))}
      </div>

      <TUISection title="INTERFACES">
        <TUIToggle label="Wi-Fi" description="Interfaz: wlan0" value={wifiOn} onChange={v => toggle('wifi', v, setWifiOn)} />
        <TUIToggle label="Ethernet" description="eth0  ·  100/1000 Mbps  ·  192.168.1.105" value={ethOn} onChange={v => toggle('eth', v, setEthOn)} />
        <TUIToggle label="VPN" description="WireGuard / OpenVPN" value={vpnOn} onChange={v => toggle('vpn', v, setVpnOn)} />
      </TUISection>

      <TUISection title="REDES WI-FI">
        {/* Scan button */}
        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${C.border}` }}>
          <TUIButton
            label={scanning ? 'Escaneando...' : 'Buscar redes Wi-Fi'}
            onClick={scan}
            icon="◉"
          />
        </div>

        {/* Network list */}
        {NETWORKS.map(net => {
          const isConnected = connected === net.ssid;
          const isSelected  = selected === net.ssid;
          return (
            <div key={net.ssid}>
              <div
                onClick={() => !isConnected && setSelected(isSelected ? null : net.ssid)}
                style={{
                  display: 'flex', alignItems: 'center', gap: 10,
                  padding: '9px 14px',
                  borderBottom: `1px solid ${C.border}`,
                  cursor: isConnected ? 'default' : 'pointer',
                  background: isConnected ? C.successDim : isSelected ? C.accentDim : 'transparent',
                  transition: 'background 0.1s',
                }}
              >
                <span style={{ color: isConnected ? C.success : C.textMuted, fontSize: 12, fontFamily: 'monospace', width: 10 }}>
                  {isConnected ? '●' : '○'}
                </span>
                <span style={{ flex: 1, color: isConnected ? C.success : C.textPrimary, fontSize: 13, fontFamily: 'monospace' }}>
                  {net.ssid}
                  {net.saved && !isConnected && <span style={{ color: C.textMuted, fontSize: 10, marginLeft: 8 }}>guardada</span>}
                </span>
                <SignalBars pct={net.signal} />
                <span style={{ color: net.security === 'OPEN' ? C.error : C.textMuted, fontSize: 10, fontFamily: 'monospace', width: 40, textAlign: 'right' }}>
                  {net.security === 'OPEN' ? '○' : ICONS.widgets.lock.nerd} {net.security}
                </span>
                {!isConnected && (
                  <span style={{ color: C.accent, fontSize: 10, fontFamily: 'monospace', marginLeft: 4 }}>
                    {isSelected ? '▼' : '▶'}
                  </span>
                )}
              </div>

              {/* Inline connection panel */}
              {isSelected && !isConnected && (
                <div style={{
                  padding: '10px 14px 12px 34px',
                  borderBottom: `1px solid ${C.border}`,
                  background: C.surface,
                  display: 'flex', flexDirection: 'column', gap: 8,
                }}>
                  {net.security !== 'OPEN' && !net.saved && (
                    <TUIInput
                      label="Contraseña"
                      value={password}
                      onChange={setPassword}
                      type="password"
                      placeholder="Contraseña de la red..."
                    />
                  )}
                  <div style={{ display: 'flex', gap: 8, marginTop: 4 }}>
                    <TUIButton label="Conectar" onClick={() => connect(net)} variant="primary" />
                    {net.saved && <TUIButton label="Olvidar" onClick={async () => { onModified(); }} variant="danger" />}
                    <TUIButton label="Cancelar" onClick={async () => setSelected(null)} />
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </TUISection>

      <TUISection title="CONFIGURACIÓN AVANZADA" collapsible defaultOpen={false}>
        <TUIInput
          label="DNS Primario"
          description="Cloudflare: 1.1.1.1  ·  Google: 8.8.8.8"
          value={dns}
          onChange={v => { ss('dns', v); setDns(v); onModified(); }}
          validate={v => /^(\d{1,3}\.){3}\d{1,3}$/.test(v) ? null : 'Formato inválido (ej: 1.1.1.1)'}
        />
        <TUIInput
          label="Proxy HTTP/S"
          description="Vacío = sin proxy. Formato: http://host:puerto"
          value={proxy}
          onChange={v => { ss('proxy', v); setProxy(v); onModified(); }}
        />
      </TUISection>
    </div>
  );
}
