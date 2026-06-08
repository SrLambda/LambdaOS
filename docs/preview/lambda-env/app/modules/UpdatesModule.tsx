import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUIProgress } from '../components/tui/TUIProgress';

interface Pkg { name: string; current: string; next: string; size: string; selected: boolean; }

const PKGS: Pkg[] = [
  { name: 'linux',        current: '6.5.0-1',  next: '6.6.2-1',  size: '128 MB', selected: true  },
  { name: 'firefox',      current: '120.0-1',  next: '121.0.1-2', size: '64 MB',  selected: true  },
  { name: 'mesa',         current: '23.3.0-1', next: '23.3.2-1', size: '48 MB',  selected: true  },
  { name: 'systemd',      current: '254.5-1',  next: '254.7-1',  size: '12 MB',  selected: true  },
  { name: 'openssl',      current: '3.1.4-1',  next: '3.2.0-1',  size: '4 MB',   selected: true  },
  { name: 'neovim',       current: '0.9.4-1',  next: '0.9.5-1',  size: '8 MB',   selected: false },
  { name: 'git',          current: '2.42.0-1', next: '2.43.0-1', size: '2 MB',   selected: true  },
  { name: 'python',       current: '3.11.5-1', next: '3.12.1-1', size: '20 MB',  selected: false },
];

function ls(k: string, d: any) { try { const v = localStorage.getItem('upd_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('upd_'+k, JSON.stringify(v)); }

export function UpdatesModule({ onModified }: { onModified: () => void }) {
  const [pkgs, setPkgs]         = useState(PKGS);
  const [autoUpd, setAutoUpd]   = useState(() => ls('auto', false));
  const [checking, setChecking] = useState(false);
  const [updating, setUpdating] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentPkg, setCurrentPkg] = useState('');

  const selected  = pkgs.filter(p => p.selected);
  const totalSize = selected.reduce((acc, p) => acc + parseFloat(p.size), 0);

  function togglePkg(name: string) {
    setPkgs(ps => ps.map(p => p.name === name ? { ...p, selected: !p.selected } : p));
  }

  async function checkUpdates() {
    setChecking(true);
    await new Promise(r => setTimeout(r, 1500));
    setChecking(false);
  }

  async function runUpdate() {
    setUpdating(true);
    setProgress(0);
    for (let i = 0; i < selected.length; i++) {
      setCurrentPkg(selected[i].name);
      setProgress(Math.round(((i + 1) / selected.length) * 100));
      await new Promise(r => setTimeout(r, 600));
    }
    setPkgs(ps => ps.map(p => selected.find(s => s.name === p.name) ? { ...p, current: p.next } : p));
    setUpdating(false);
    setCurrentPkg('');
    onModified();
  }

  return (
    <div>
      <TUISection title="ESTADO">
        <div style={{ padding: '10px 12px', borderBottom: `1px solid ${C.border}`, fontFamily: 'monospace' }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <div>
              <span style={{ color: C.warn, fontSize: 12 }}>↻ {pkgs.length} actualizaciones disponibles</span>
              <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>Última verificación: hace 2h · pacman</div>
            </div>
            <TUIButton label={checking ? 'Verificando...' : 'Verificar ahora'} onClick={checkUpdates} />
          </div>
        </div>
        <TUIToggle label="Actualizaciones automáticas" description="Instalar automáticamente las actualizaciones de seguridad" value={autoUpd} onChange={v => { setAutoUpd(v); ss('auto', v); onModified(); }} />
      </TUISection>

      {updating && (
        <div style={{ padding: '10px 14px', background: C.accentDim, borderLeft: `3px solid ${C.accent}`, marginBottom: 4 }}>
          <TUIProgress value={progress} label="Instalando actualizaciones..." subLabel={`Instalando: ${currentPkg}`} />
        </div>
      )}

      <TUISection title={`PAQUETES DISPONIBLES (${pkgs.length})`}>
        {/* Header */}
        <div style={{
          display: 'grid', gridTemplateColumns: '20px 160px 1fr 1fr 70px',
          padding: '5px 12px', background: C.surface, gap: 8,
          borderBottom: `1px solid ${C.border}`, fontFamily: 'monospace',
        }}>
          {['','Paquete','Actual','Nuevo','Tamaño'].map(h => (
            <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
          ))}
        </div>
        {pkgs.map(pkg => (
          <div
            key={pkg.name}
            onClick={() => togglePkg(pkg.name)}
            style={{
              display: 'grid', gridTemplateColumns: '20px 160px 1fr 1fr 70px',
              padding: '7px 12px', borderBottom: `1px solid ${C.border}`,
              gap: 8, fontFamily: 'monospace', cursor: 'pointer',
              background: pkg.selected ? C.accentDim : 'transparent',
              alignItems: 'center',
            }}
          >
            <span style={{ color: pkg.selected ? C.accent : C.textMuted, fontSize: 12 }}>
              {pkg.selected ? '●' : '○'}
            </span>
            <span style={{ color: C.textPrimary, fontSize: 12 }}>{pkg.name}</span>
            <span style={{ color: C.textMuted, fontSize: 11 }}>{pkg.current}</span>
            <span style={{ color: C.success, fontSize: 11 }}>→ {pkg.next}</span>
            <span style={{ color: C.textMuted, fontSize: 11 }}>{pkg.size}</span>
          </div>
        ))}

        {/* Footer */}
        <div style={{
          display: 'flex', alignItems: 'center', justifyContent: 'space-between',
          padding: '10px 12px', background: C.surface,
          fontFamily: 'monospace',
        }}>
          <span style={{ color: C.textSecondary, fontSize: 11 }}>
            {selected.length} paquetes seleccionados · {totalSize.toFixed(0)} MB
          </span>
          <div style={{ display: 'flex', gap: 8 }}>
            <TUIButton label="Seleccionar todos" onClick={async () => setPkgs(ps => ps.map(p => ({ ...p, selected: true })))} />
            <TUIButton label={`Instalar (${selected.length})`} onClick={runUpdate} variant="primary" disabled={selected.length === 0 || updating} />
          </div>
        </div>
      </TUISection>
    </div>
  );
}
