import { useState } from 'react';
import { TUISettingItem } from '../TUISettingItem';

const NETWORKS = [
  { ssid: 'LambdaNet-5G', signal: 92, security: 'WPA3', band: '5GHz' },
  { ssid: 'Office-Corp', signal: 74, security: 'WPA2', band: '5GHz' },
  { ssid: 'Guest_Public', signal: 61, security: 'OPEN', band: '2.4GHz' },
  { ssid: 'IoT_Devices', signal: 48, security: 'WPA2', band: '2.4GHz' },
  { ssid: 'Hotspot-Mobile', signal: 35, security: 'WPA3', band: '5GHz' },
];

function SignalBar({ strength }: { strength: number }) {
  const bars = Math.ceil(strength / 20);
  return (
    <span className="flex items-end gap-px h-3">
      {[1,2,3,4,5].map(i => (
        <span
          key={i}
          style={{ height: `${i * 20}%` }}
          className={`w-1 inline-block ${i <= bars ? 'bg-[#6D40FF]' : 'bg-[#6D40FF]/20'}`}
        />
      ))}
    </span>
  );
}

export function NetworkSettings() {
  const [scanning, setScanning] = useState(false);
  const [connected, setConnected] = useState('LambdaNet-5G');

  function handleScan() {
    setScanning(true);
    setTimeout(() => setScanning(false), 2000);
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ RED Y CONECTIVIDAD ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">ifconfig / iwconfig</span>
      </div>

      <TUISettingItem label="Wi-Fi" type="toggle" value={true} storageKey="wifi_enabled"
        description="Interfaz: wlan0  │  MAC: 3C:A9:F4:12:7B:E2" />
      <TUISettingItem label="Bluetooth" type="toggle" value={false} storageKey="bt_enabled"
        description="Interfaz: hci0  │  Estado: desconectado" />
      <TUISettingItem label="Ethernet" type="toggle" value={true} storageKey="eth_enabled"
        description="eth0  │  100/1000 Mbps Full Duplex  │  IP: 192.168.1.105" />
      <TUISettingItem label="DNS Primario" type="input" value="1.1.1.1" storageKey="dns_primary"
        description="Servidor DNS principal (Cloudflare / Google / Custom)" />
      <TUISettingItem label="DNS Secundario" type="input" value="8.8.8.8" storageKey="dns_secondary"
        description="Servidor DNS de respaldo" />
      <TUISettingItem label="Proxy HTTP/S" type="input" value="" storageKey="proxy_url"
        description="Formato: http://host:puerto  (vacío = sin proxy)" />
      <TUISettingItem label="Firewall" type="toggle" value={true} storageKey="firewall_enabled"
        description="iptables activo  │  Política: DROP  │  Reglas: 24" />
      <TUISettingItem label="VPN" type="toggle" value={false} storageKey="vpn_enabled"
        description="WireGuard / OpenVPN  │  No conectado" />
      <TUISettingItem label="IPv6" type="toggle" value={true} storageKey="ipv6_enabled"
        description="Protocolo de Internet versión 6" />

      {/* Networks table */}
      <div className="mt-4">
        <div className="flex items-center justify-between mb-1">
          <span className="text-[#6D40FF]/70 text-[10px] tracking-wider">REDES DISPONIBLES</span>
          <button
            onClick={handleScan}
            className="text-[10px] border border-[#6D40FF]/60 text-[#6D40FF] px-2 py-0.5 hover:bg-[#6D40FF]/10 transition-all"
          >
            {scanning ? '[ ESCANEANDO... ]' : '[ ESCANEAR ]'}
          </button>
        </div>
        <div className="border border-[#6D40FF]/30">
          <div className="grid grid-cols-[1fr_auto_auto_auto] gap-2 px-2 py-1 border-b border-[#6D40FF]/30 bg-[#6D40FF]/10">
            <span className="text-[#6D40FF] text-[10px]">SSID</span>
            <span className="text-[#6D40FF] text-[10px]">SEÑAL</span>
            <span className="text-[#6D40FF] text-[10px]">SEGURIDAD</span>
            <span className="text-[#6D40FF] text-[10px]">BANDA</span>
          </div>
          {NETWORKS.map((net) => (
            <div
              key={net.ssid}
              onClick={() => setConnected(net.ssid)}
              className={`grid grid-cols-[1fr_auto_auto_auto] gap-2 px-2 py-1 border-b border-[#6D40FF]/10 last:border-0 cursor-pointer transition-all
                ${connected === net.ssid ? 'bg-[#6D40FF]/20' : 'hover:bg-[#6D40FF]/5'}`}
            >
              <span className="text-[10px] flex items-center gap-1">
                {connected === net.ssid && <span className="text-[#6D40FF]">●</span>}
                {connected !== net.ssid && <span className="text-[#6D40FF]/20">○</span>}
                <span className={connected === net.ssid ? 'text-[#6D40FF]' : 'text-[#6D40FF]/70'}>{net.ssid}</span>
              </span>
              <span className="flex items-center"><SignalBar strength={net.signal} /></span>
              <span className={`text-[10px] ${net.security === 'OPEN' ? 'text-[#FF4040]/70' : 'text-[#6D40FF]/60'}`}>{net.security}</span>
              <span className="text-[#6D40FF]/50 text-[10px]">{net.band}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
