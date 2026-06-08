import { useEffect, useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIModal } from '../components/tui/TUIModal';

const RESOLUTIONS = ['1280x720','1280x1024','1366x768','1600x900','1920x1080','2560x1440','3840x2160'];
const REFRESH     = ['24Hz','30Hz','60Hz','75Hz','100Hz','120Hz','144Hz','165Hz','240Hz'];
const ORIENTATION = ['Normal (0°)','Rotación 90°','Rotación 180°','Rotación 270°'];
const COLOR_PROF  = ['sRGB','AdobeRGB','P3-D65','Linear'];

const MONITORS = [
  { id: 'DP-1',  name: 'DELL U2722D', res: '2560x1440', hz: '60Hz', connected: true,  primary: true  },
  { id: 'HDMI-1',name: 'LG 24MK600', res: '1920x1080', hz: '60Hz', connected: true,  primary: false },
  { id: 'DP-2',  name: 'Sin señal',   res: '—',        hz: '—',    connected: false, primary: false },
];

function ls(k: string, d: any) { try { const v = localStorage.getItem('disp_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('disp_'+k, JSON.stringify(v)); }

export function DisplayModule({ onModified, onStatusMsg }: { onModified: () => void; onStatusMsg: (m: string) => void }) {
  const [brightness, setBrightness] = useState(() => ls('brightness', 75));
  const [contrast, setContrast]     = useState(() => ls('contrast', 50));
  const [res, setRes]               = useState(() => ls('res', '1920x1080'));
  const [hz, setHz]                 = useState(() => ls('hz', '60Hz'));
  const [orient, setOrient]         = useState(() => ls('orient', ORIENTATION[0]));
  const [colorProf, setColorProf]   = useState(() => ls('color', 'sRGB'));
  const [nightMode, setNightMode]   = useState(() => ls('night', false));
  const [hdr, setHdr]               = useState(() => ls('hdr', false));
  const [selectedMon, setSelectedMon] = useState('DP-1');

  const [pendingRes, setPendingRes]     = useState<string | null>(null);
  const [countdown, setCountdown]       = useState(10);
  const [showCountdown, setShowCountdown] = useState(false);

  function update(key: string, val: any, setter: (v: any) => void) {
    setter(val); ss(key, val); onModified();
  }

  function applyRes(newRes: string) {
    setPendingRes(newRes);
    setCountdown(10);
    setShowCountdown(true);
  }

  useEffect(() => {
    if (!showCountdown) return;
    if (countdown <= 0) {
      // Auto-revert
      setShowCountdown(false);
      setPendingRes(null);
      onStatusMsg(`Pantalla · Resolución revertida automáticamente`);
      return;
    }
    const id = setTimeout(() => setCountdown(c => c - 1), 1000);
    return () => clearTimeout(id);
  }, [showCountdown, countdown]);

  function confirmRes() {
    if (pendingRes) update('res', pendingRes, setRes);
    setShowCountdown(false);
    setPendingRes(null);
    onStatusMsg(`Pantalla · Resolución aplicada: ${pendingRes}`);
  }

  return (
    <div>
      {/* Monitor cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8, padding: '0 0 4px 0' }}>
        {MONITORS.map(mon => (
          <div
            key={mon.id}
            onClick={() => mon.connected && setSelectedMon(mon.id)}
            style={{
              padding: '10px 12px',
              border: `1px solid ${mon.id === selectedMon ? C.accent : C.border}`,
              background: !mon.connected ? 'transparent' : mon.id === selectedMon ? C.accentDim : C.surface,
              opacity: mon.connected ? 1 : 0.4,
              cursor: mon.connected ? 'pointer' : 'not-allowed',
              fontFamily: 'monospace',
            }}
          >
            <div style={{ color: mon.connected ? C.textPrimary : C.textMuted, fontSize: 12 }}>{mon.id}</div>
            <div style={{ color: C.textSecondary, fontSize: 10, marginTop: 2 }}>{mon.name}</div>
            <div style={{ color: C.textMuted, fontSize: 10 }}>{mon.res} · {mon.hz}</div>
            {mon.primary && <div style={{ color: C.accent, fontSize: 10, marginTop: 4 }}>● Primario</div>}
            {!mon.connected && <div style={{ color: C.textMuted, fontSize: 10, marginTop: 4 }}>Desconectado</div>}
          </div>
        ))}
      </div>

      {/* Resolution countdown bar */}
      {showCountdown && (
        <div style={{
          padding: '10px 14px',
          background: C.warnDim,
          borderLeft: `3px solid ${C.warn}`,
          display: 'flex', alignItems: 'center', gap: 12,
          fontFamily: 'monospace',
          marginBottom: 4,
        }}>
          <span style={{ color: C.warn, fontSize: 12 }}>⚠ ¿Mantener resolución {pendingRes}? Revirtiendo en {countdown}s...</span>
          <div style={{ display: 'flex', gap: 8, marginLeft: 'auto' }}>
            <TUIButton label="Confirmar" onClick={async () => confirmRes()} variant="primary" />
            <TUIButton label="Revertir" onClick={async () => { setShowCountdown(false); setPendingRes(null); }} />
          </div>
        </div>
      )}

      <TUISection title="AJUSTES DE IMAGEN">
        <TUISlider label="Brillo" value={brightness} onChange={v => update('brightness', v, setBrightness)} />
        <TUISlider label="Contraste" value={contrast} onChange={v => update('contrast', v, setContrast)} />
        <TUISelect label="Perfil de Color" value={colorProf} options={COLOR_PROF} onChange={v => update('color', v, setColorProf)} />
      </TUISection>

      <TUISection title="RESOLUCIÓN Y FRECUENCIA">
        <TUISelect
          label="Resolución"
          description="Aplicar abre cuenta regresiva de 10s"
          value={res}
          options={RESOLUTIONS}
          onChange={newRes => { if (newRes !== res) applyRes(newRes); }}
        />
        <TUISelect label="Tasa de Refresco" value={hz} options={REFRESH} onChange={v => update('hz', v, setHz)} />
        <TUISelect label="Orientación" value={orient} options={ORIENTATION} onChange={v => update('orient', v, setOrient)} />
      </TUISection>

      <TUISection title="CARACTERÍSTICAS">
        <TUIToggle label="Modo Nocturno" description="Filtro de luz azul activo entre 20:00–07:00" value={nightMode} onChange={v => update('night', v, setNightMode)} />
        <TUIToggle label="HDR" description="Rango dinámico alto (requiere panel compatible)" value={hdr} onChange={v => update('hdr', v, setHdr)} />
      </TUISection>
    </div>
  );
}
