import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';

const OUTPUT_DEVS = ['Speakers (Realtek ALC892)', 'Headphones (3.5mm)', 'HDMI Audio (GPU)', 'USB DAC (FiiO K3)', 'Bluetooth (Sony WH-1000XM5)'];
const INPUT_DEVS  = ['Micrófono Integrado', 'USB Microphone (Blue Yeti)', 'Headset Mic (Jabra)', 'Line-In'];
const EQ_PRESETS  = ['Plano', 'Bass Boost', 'Treble', 'Clásico', 'Rock', 'Jazz', 'Vocal', 'Gaming'];
const PROFILES    = ['Salida estéreo', 'Surround 5.1', 'Pro Audio', 'Bajo consumo'];

const PER_APP = [
  { name: 'Firefox', icon: '◉', vol: 90 },
  { name: 'Spotify', icon: '♪', vol: 70 },
  { name: 'Alacritty', icon: '◰', vol: 50 },
  { name: 'Discord', icon: '◎', vol: 80 },
  { name: 'mpv', icon: '▷', vol: 100 },
];

function ls(key: string, def: any) {
  try { const v = localStorage.getItem('audio_' + key); return v ? JSON.parse(v) : def; } catch { return def; }
}
function ss(key: string, val: any) { localStorage.setItem('audio_' + key, JSON.stringify(val)); }

export function AudioModule({ onModified }: { onModified: () => void }) {
  const [master, setMaster]   = useState(() => ls('master', 80));
  const [muted, setMuted]     = useState(() => ls('muted', false));
  const [outDev, setOutDev]   = useState(() => ls('outDev', OUTPUT_DEVS[0]));
  const [inDev, setInDev]     = useState(() => ls('inDev', INPUT_DEVS[0]));
  const [micVol, setMicVol]   = useState(() => ls('micVol', 65));
  const [eq, setEq]           = useState(() => ls('eq', 'Plano'));
  const [profile, setProfile] = useState(() => ls('profile', PROFILES[0]));
  const [sysBeep, setSysBeep] = useState(() => ls('sysBeep', true));
  const [nr, setNr]           = useState(() => ls('nr', false));
  const [perApp, setPerApp]   = useState(() => ls('perApp', PER_APP.map(a => a.vol)));

  function update(key: string, val: any, setter: (v: any) => void) {
    setter(val); ss(key, val); onModified();
  }

  return (
    <div>
      <TUISection title="SALIDA DE AUDIO">
        <TUIToggle
          label="Silenciar"
          description="Silenciar toda la salida de audio"
          value={muted}
          onChange={v => update('muted', v, setMuted)}
        />
        <TUISlider
          label="Volumen Principal"
          description="amixer sset Master"
          value={master}
          onChange={v => update('master', v, setMaster)}
          disabled={muted}
        />
        <TUISelect
          label="Dispositivo de Salida"
          description="Sink PulseAudio / PipeWire activo"
          value={outDev}
          options={OUTPUT_DEVS}
          onChange={v => update('outDev', v, setOutDev)}
        />
        <TUISelect
          label="Perfil"
          description="Perfil de hardware del sink"
          value={profile}
          options={PROFILES}
          onChange={v => update('profile', v, setProfile)}
        />
      </TUISection>

      <TUISection title="ENTRADA DE AUDIO" defaultOpen={true}>
        <TUISelect
          label="Dispositivo de Entrada"
          description="Source PulseAudio / PipeWire activo"
          value={inDev}
          options={INPUT_DEVS}
          onChange={v => update('inDev', v, setInDev)}
        />
        <TUISlider
          label="Volumen Micrófono"
          description="Ganancia del source activo"
          value={micVol}
          onChange={v => update('micVol', v, setMicVol)}
        />
        <TUIToggle
          label="Reducción de Ruido"
          description="RNNoise (requiere pipewire-pulse, +3% CPU)"
          value={nr}
          onChange={v => update('nr', v, setNr)}
        />
      </TUISection>

      <TUISection title="AJUSTES">
        <TUISelect
          label="Perfil EQ"
          description="Ecualizador paramétrico del sistema"
          value={eq}
          options={EQ_PRESETS}
          onChange={v => update('eq', v, setEq)}
        />
        <TUIToggle
          label="Sonidos del Sistema"
          description="Efectos de audio para eventos del OS"
          value={sysBeep}
          onChange={v => update('sysBeep', v, setSysBeep)}
        />
      </TUISection>

      <TUISection title="VOLUMEN POR APLICACIÓN" collapsible defaultOpen={false}>
        {PER_APP.map((app, i) => (
          <TUISlider
            key={app.name}
            label={`${app.icon} ${app.name}`}
            value={perApp[i]}
            onChange={v => {
              const next = [...perApp]; next[i] = v;
              update('perApp', next, setPerApp);
            }}
          />
        ))}
        <div style={{ padding: '8px 12px' }}>
          <TUIButton label="Restablecer todos" onClick={async () => {
            const def = PER_APP.map(a => a.vol);
            update('perApp', def, setPerApp);
          }} />
        </div>
      </TUISection>

      <TUISection title="PERFILES GUARDADOS" collapsible defaultOpen={false}>
        {['Casa', 'Oficina', 'Estudio', 'Gaming'].map(p => (
          <div key={p} style={{
            display: 'flex', alignItems: 'center', justifyContent: 'space-between',
            padding: '8px 12px', borderBottom: `1px solid ${C.border}`,
          }}>
            <span style={{ color: C.textPrimary, fontSize: 12, fontFamily: 'monospace' }}>{p}</span>
            <div style={{ display: 'flex', gap: 8 }}>
              <TUIButton label="Cargar" onClick={async () => { onModified(); }} />
              <TUIButton label="Eliminar" onClick={async () => {}} variant="danger" />
            </div>
          </div>
        ))}
        <div style={{ padding: '8px 12px' }}>
          <TUIButton label="Guardar perfil actual" onClick={async () => { onModified(); }} variant="primary" />
        </div>
      </TUISection>
    </div>
  );
}
