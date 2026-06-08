import { useEffect, useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUISection } from '../components/tui/TUISection';
import { TUIProgress } from '../components/tui/TUIProgress';
import { TUIButton } from '../components/tui/TUIButton';

const SCREEN_T = ['1 min','5 min','10 min','15 min','30 min','Nunca'];
const SLEEP_T  = ['5 min','15 min','30 min','1 hora','2 horas','Nunca'];
const PROFILES = ['Máximo Ahorro','Ahorro Equilibrado','Balanceado','Alto Rendimiento','Máximo (Turbo)'];
const LID_ACT  = ['Suspender','Hibernar','Nada','Apagar pantalla'];

function ls(k: string, d: any) { try { const v = localStorage.getItem('pwr_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('pwr_'+k, JSON.stringify(v)); }

export function PowerModule({ onModified }: { onModified: () => void }) {
  const [profile, setProfile]   = useState(() => ls('profile', 'Balanceado'));
  const [screenT, setScreenT]   = useState(() => ls('screenT', '10 min'));
  const [sleepT, setSleepT]     = useState(() => ls('sleepT', '30 min'));
  const [hibernate, setHibernate] = useState(() => ls('hibernate', true));
  const [turbo, setTurbo]       = useState(() => ls('turbo', true));
  const [lidAct, setLidAct]     = useState(() => ls('lid', 'Suspender'));
  const [wol, setWol]           = useState(() => ls('wol', false));
  const [chargeLimit, setChargeLimit] = useState(() => ls('charge', 80));
  const [battBrightness, setBattBrightness] = useState(() => ls('battBright', 50));
  const [battPct, setBattPct]   = useState(85);

  useEffect(() => {
    const id = setInterval(() => setBattPct(p => Math.min(100, p + (Math.random() > 0.8 ? 1 : 0))), 4000);
    return () => clearInterval(id);
  }, []);

  function update(key: string, val: any, setter: (v: any) => void) {
    setter(val); ss(key, val); onModified();
  }

  const battColor = battPct > 60 ? C.success : battPct > 20 ? C.warn : C.error;

  return (
    <div>
      {/* Battery widget */}
      <div style={{
        padding: '12px 14px',
        background: C.surface,
        borderLeft: `3px solid ${battColor}`,
        fontFamily: 'monospace',
        marginBottom: 4,
      }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
          <span style={{ color: C.textPrimary, fontSize: 13 }}>Batería</span>
          <span style={{ color: C.success, fontSize: 11 }}>⚡ Cargando · ~47 min completo</span>
        </div>
        <TUIProgress
          value={battPct}
          label={`${battPct}%`}
          color={battColor}
          subLabel={`Ciclos: 142  ·  Salud: 96%  ·  Temp: 31.2°C  ·  Potencia: +45W`}
        />
      </div>

      <TUISection title="PERFIL DE ENERGÍA">
        <TUISelect label="Perfil activo" description="TLP / power-profiles-daemon" value={profile} options={PROFILES} onChange={v => update('profile', v, setProfile)} />
        <TUIToggle label="Turbo Boost CPU" description="Frecuencias boost (consume más energía)" value={turbo} onChange={v => update('turbo', v, setTurbo)} />
      </TUISection>

      <TUISection title="TIMEOUTS">
        <TUISelect
          label="Apagar pantalla después de"
          description="DPMS standby timeout"
          value={screenT}
          options={SCREEN_T}
          onChange={v => update('screenT', v, setScreenT)}
        />
        <TUISelect
          label="Suspender sistema después de"
          description="systemd-suspend.service"
          value={sleepT}
          options={SLEEP_T}
          onChange={v => update('sleepT', v, setSleepT)}
        />
        <TUISelect
          label="Al cerrar la tapa"
          description="Acción al cerrar el portátil"
          value={lidAct}
          options={LID_ACT}
          onChange={v => update('lid', v, setLidAct)}
        />
        <TUIToggle label="Hibernación" description="Guardar RAM en swap al hibernar (swapfile ≥ RAM)" value={hibernate} onChange={v => update('hibernate', v, setHibernate)} />
      </TUISection>

      <TUISection title="BATERÍA">
        <TUISlider label="Brillo con batería" description="Reduce brillo automáticamente al desconectar" value={battBrightness} onChange={v => update('battBright', v, setBattBrightness)} />
        <TUISlider label="Límite de carga (%)" description="Conservar salud de la batería (TLP threshold)" value={chargeLimit} min={60} max={100} onChange={v => update('charge', v, setChargeLimit)} />
        <TUIToggle label="Wake on LAN" description="Despertar equipo por red (ethtool wol g)" value={wol} onChange={v => update('wol', v, setWol)} />
      </TUISection>

      <TUISection title="ACCIONES">
        <div style={{ padding: '10px 12px', display: 'flex', gap: 10 }}>
          <TUIButton label="Suspender" onClick={async () => { onModified(); }} icon="◷" />
          <TUIButton label="Hibernar" onClick={async () => { onModified(); }} icon="◌" />
          <TUIButton label="Reiniciar" onClick={async () => { onModified(); }} variant="danger" icon="↻" />
          <TUIButton label="Apagar" onClick={async () => { onModified(); }} variant="danger" icon="✕" />
        </div>
      </TUISection>
    </div>
  );
}
