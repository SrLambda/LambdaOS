import { useEffect, useState } from 'react';
import { TUISettingItem } from '../TUISettingItem';

const SCREEN_TIMEOUTS = ['1 min', '5 min', '10 min', '15 min', '30 min', 'Nunca'];
const SLEEP_TIMEOUTS = ['5 min', '10 min', '15 min', '30 min', '1 hora', '2 horas', 'Nunca'];
const POWER_PROFILES = ['Máximo Ahorro', 'Ahorro Equilibrado', 'Balanceado', 'Alto Rendimiento', 'Máximo (Turbo)'];
const CPU_GOVERNORS = ['powersave', 'conservative', 'ondemand', 'schedutil', 'performance'];

function BatteryBar({ pct }: { pct: number }) {
  const color = pct > 60 ? '#6D40FF' : pct > 20 ? '#FFAA00' : '#FF4040';
  const glowColor = pct > 60 ? 'rgba(109,64,255,0.6)' : pct > 20 ? 'rgba(255,170,0,0.6)' : 'rgba(255,64,64,0.6)';
  return (
    <div className="flex items-center gap-2 mt-1">
      <div className="flex-1 border border-[#6D40FF]/50 h-5 relative bg-black overflow-hidden">
        <div
          className="absolute left-0 top-0 h-full transition-all"
          style={{ width: `${pct}%`, background: color, boxShadow: `inset 0 0 8px ${glowColor}` }}
        />
        {/* Segmented lines */}
        {[25,50,75].map(p => (
          <div key={p} className="absolute top-0 bottom-0 w-px bg-black/40" style={{ left: `${p}%` }} />
        ))}
        <div className="absolute inset-0 flex items-center justify-center text-[10px] font-bold"
          style={{ color: pct > 45 ? '#000' : color, mixBlendMode: pct > 45 ? 'difference' : 'normal' }}>
          {pct}%
        </div>
      </div>
      <div className="text-[#6D40FF]/50 text-[10px] border border-[#6D40FF]/30 px-1 py-0.5">
        ▌
      </div>
    </div>
  );
}

export function PowerSettings() {
  const [pct, setPct] = useState(85);

  useEffect(() => {
    const id = setInterval(() => setPct(p => Math.min(100, p + (Math.random() > 0.7 ? 1 : 0))), 3000);
    return () => clearInterval(id);
  }, []);

  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ ENERGÍA Y BATERÍA ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">upower / tlp / powertop</span>
      </div>

      {/* Battery status widget */}
      <div className="border border-[#6D40FF]/30 p-2 mb-3 grid grid-cols-2 gap-x-4 gap-y-0.5">
        <div>
          <div className="text-[#6D40FF]/50 text-[10px] mb-0.5">ESTADO DE BATERÍA</div>
          <BatteryBar pct={pct} />
        </div>
        <div className="text-[10px] space-y-0.5">
          {[
            ['Estado', '⚡ Cargando'],
            ['Tiempo', '~47 min completo'],
            ['Ciclos', '142 ciclos'],
            ['Salud', '96% capacidad'],
            ['Temp', '31.2°C'],
            ['Potencia', '+45W (carga)'],
          ].map(([k, v]) => (
            <div key={k} className="flex gap-1">
              <span className="text-[#6D40FF]/40 w-14">{k}:</span>
              <span className="text-[#6D40FF]/80">{v}</span>
            </div>
          ))}
        </div>
      </div>

      <TUISettingItem label="Perfil de Energía" type="select" value="Balanceado" storageKey="power_profile"
        options={POWER_PROFILES} description="TLP / power-profiles-daemon" />
      <TUISettingItem label="Gobernador CPU" type="select" value="schedutil" storageKey="power_cpu_gov"
        options={CPU_GOVERNORS} description="cpufreq governor kernel" />
      <TUISettingItem label="Suspender Pantalla" type="select" value="10 min" storageKey="power_screen_off"
        options={SCREEN_TIMEOUTS} description="DPMS standby timeout" />
      <TUISettingItem label="Suspender Sistema" type="select" value="30 min" storageKey="power_sleep"
        options={SLEEP_TIMEOUTS} description="systemd-suspend.service delay" />
      <TUISettingItem label="Brillo en Batería" type="slider" value={50} storageKey="power_batt_bright"
        description="Brillo reducido automático al desconectar AC" unit="%" />
      <TUISettingItem label="Hibernación" type="toggle" value={true} storageKey="power_hibernate"
        description="Guardar RAM en swap e hibernar (swapfile ≥ RAM)" />
      <TUISettingItem label="Turbo Boost CPU" type="toggle" value={true} storageKey="power_turbo"
        description="Frecuencias boost en modo rendimiento" />
      <TUISettingItem label="Límite Carga (%)" type="slider" value={80} storageKey="power_charge_limit"
        min={60} max={100} unit="%"
        description="Conservar salud de batería (TLP threshold)" />
      <TUISettingItem label="Wake on LAN" type="toggle" value={false} storageKey="power_wol"
        description="Despertar equipo por red (ethtool wol g)" />

      {/* CPU freq monitor */}
      <div className="mt-3 border border-[#6D40FF]/30 p-2">
        <div className="text-[#6D40FF]/50 text-[10px] mb-1">FRECUENCIAS CPU (MHz)</div>
        <div className="grid grid-cols-4 gap-1">
          {[3800, 3600, 3200, 4200, 3900, 3700, 3400, 4000].map((f, i) => (
            <div key={i} className="text-[10px]">
              <div className="text-[#6D40FF]/40">C{i}:</div>
              <div className="text-[#6D40FF] font-bold">{f}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
