// Mirrors: src/lambda-env/internal/modules/qtile/main.go
// Actions: run, set-bar-position, set-bar-size, set-terminal, set-browser, set-layouts
import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUIInput } from '../components/tui/TUIInput';
import { TUIButton } from '../components/tui/TUIButton';

// ── Data ────────────────────────────────────────────────────────────────────

const BAR_POSITIONS  = ['top', 'bottom'];
const TERMINALS      = ['alacritty', 'kitty', 'wezterm', 'foot', 'xterm', 'gnome-terminal'];
const BROWSERS       = ['firefox', 'chromium', 'librewolf', 'qutebrowser', 'brave'];
const FILE_MANAGERS  = ['thunar', 'nautilus', 'dolphin', 'ranger', 'nnn', 'lf'];
const LAYOUT_NAMES   = [
  'Columns', 'Max', 'Stack', 'Bsp', 'Floating',
  'MonadTall', 'MonadWide', 'RatioTile', 'Tile', 'TreeTab',
  'VerticalTile', 'Zoomy',
];
const WALLPAPER_MODES = ['fill', 'fit', 'stretch', 'tile', 'center'];
const WIDGET_SETS = ['minimal', 'standard', 'full'];

// ── Layout preview renderer ───────────────────────────────────────────────

const LAYOUT_PREVIEWS: Record<string, React.ReactNode> = {
  Columns: (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1px 1fr', height: 40, gap: 0 }}>
      <div style={{ background: C.accentDim, border: `1px solid ${C.accentBorder}` }} />
      <div style={{ background: C.border }} />
      <div style={{ background: C.surface2, border: `1px solid ${C.border}` }} />
    </div>
  ),
  Max: (
    <div style={{ height: 40, background: C.accentDim, border: `1px solid ${C.accentBorder}` }} />
  ),
  MonadTall: (
    <div style={{ display: 'grid', gridTemplateColumns: '2fr 1px 1fr', height: 40 }}>
      <div style={{ background: C.accentDim, border: `1px solid ${C.accentBorder}` }} />
      <div style={{ background: C.border }} />
      <div style={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
        <div style={{ flex: 1, background: C.surface2, border: `1px solid ${C.border}` }} />
        <div style={{ flex: 1, background: C.surface2, border: `1px solid ${C.border}` }} />
      </div>
    </div>
  ),
  Bsp: (
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gridTemplateRows: '1fr 1fr', height: 40, gap: 1 }}>
      {[0,1,2,3].map(i => (
        <div key={i} style={{ background: i === 0 ? C.accentDim : C.surface2, border: `1px solid ${i === 0 ? C.accentBorder : C.border}` }} />
      ))}
    </div>
  ),
  Floating: (
    <div style={{ position: 'relative', height: 40, background: C.surface }}>
      {[[2,4,24,20], [12,8,20,18], [20,2,22,22]].map(([l,t,w,h], i) => (
        <div key={i} style={{ position: 'absolute', left: l, top: t, width: w, height: h, background: i === 0 ? C.accentDim : C.surface2, border: `1px solid ${i === 0 ? C.accentBorder : C.border}` }} />
      ))}
    </div>
  ),
};

function LayoutPreview({ name }: { name: string }) {
  const preview = LAYOUT_PREVIEWS[name];
  if (!preview) {
    return (
      <div style={{ height: 40, background: C.surface, border: `1px solid ${C.border}`, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>{name}</span>
      </div>
    );
  }
  return <div>{preview}</div>;
}

// ── Persistence ──────────────────────────────────────────────────────────────

function ls(k: string, d: any) {
  try { const v = localStorage.getItem('qtile_' + k); return v ? JSON.parse(v) : d; } catch { return d; }
}
function ss(k: string, v: any) { localStorage.setItem('qtile_' + k, JSON.stringify(v)); }

// ── Component ────────────────────────────────────────────────────────────────

export function QtileModule({ onModified }: { onModified: () => void }) {
  const [barPos, setBarPos]         = useState(() => ls('barPos', 'top'));
  const [barSize, setBarSize]       = useState(() => ls('barSize', 28));
  const [terminal, setTerminal]     = useState(() => ls('terminal', 'alacritty'));
  const [browser, setBrowser]       = useState(() => ls('browser', 'firefox'));
  const [fileManager, setFm]        = useState(() => ls('fm', 'thunar'));
  const [wpMode, setWpMode]         = useState(() => ls('wpMode', 'fill'));
  const [widgetSet, setWidgetSet]   = useState(() => ls('widgets', 'standard'));
  const [gaps, setGaps]             = useState(() => ls('gaps', 6));
  const [borderWidth, setBorderWidth] = useState(() => ls('borderW', 2));
  const [focusFollowsMouse, setFfm] = useState(() => ls('ffm', true));
  const [autoFullscreen, setAfs]    = useState(() => ls('afs', true));
  const [layouts, setLayouts]       = useState<Record<string, boolean>>(
    () => ls('layouts', Object.fromEntries(LAYOUT_NAMES.map(l => [l, ['Columns','Max','Floating'].includes(l)])))
  );
  const [activeLayout, setActiveLayout] = useState('Columns');
  const [reloading, setReloading]   = useState(false);

  function update(k: string, v: any, setter: (x: any) => void) {
    setter(v); ss(k, v); onModified();
  }

  function toggleLayout(name: string) {
    const next = { ...layouts, [name]: !layouts[name] };
    setLayouts(next); ss('layouts', next); onModified();
  }

  const activeLayouts = LAYOUT_NAMES.filter(l => layouts[l]);

  async function reloadConfig() {
    setReloading(true);
    await new Promise(r => setTimeout(r, 1200));
    setReloading(false);
    onModified();
  }

  return (
    <div>
      {/* Bar preview */}
      <div style={{ margin: '0 0 4px 0', border: `1px solid ${C.border}`, overflow: 'hidden' }}>
        {barPos === 'top' && (
          <div style={{ background: C.surface, borderBottom: `2px solid ${C.accent}`, height: barSize, display: 'flex', alignItems: 'center', padding: '0 12px', gap: 12, fontFamily: 'monospace', fontSize: 10 }}>
            <span style={{ color: C.accent }}>⚙ Qtile</span>
            <span style={{ color: C.textSecondary }}>1:term  2:web  3:code  4:media</span>
            <span style={{ marginLeft: 'auto', color: C.textSecondary }}>♪ 80%  ⚡ 85%  {new Date().toLocaleTimeString('es-ES', { hour12: false })}</span>
          </div>
        )}
        <div style={{ background: '#0a0018', height: 56, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>[ área de trabajo ]</span>
        </div>
        {barPos === 'bottom' && (
          <div style={{ background: C.surface, borderTop: `2px solid ${C.accent}`, height: barSize, display: 'flex', alignItems: 'center', padding: '0 12px', gap: 12, fontFamily: 'monospace', fontSize: 10 }}>
            <span style={{ color: C.accent }}>⚙ Qtile</span>
            <span style={{ color: C.textSecondary }}>1:term  2:web  3:code  4:media</span>
            <span style={{ marginLeft: 'auto', color: C.textSecondary }}>♪ 80%  ⚡ 85%  {new Date().toLocaleTimeString('es-ES', { hour12: false })}</span>
          </div>
        )}
      </div>

      <TUISection title="BARRA DE ESTADO">
        <TUISelect
          label="Posición"
          description="set-bar-position — TOP / BOTTOM de la pantalla"
          value={barPos}
          options={BAR_POSITIONS}
          onChange={v => update('barPos', v, setBarPos)}
        />
        <TUISlider
          label="Altura de la barra"
          description="set-bar-size — píxeles"
          value={barSize}
          min={16}
          max={48}
          unit="px"
          onChange={v => update('barSize', v, setBarSize)}
        />
        <TUISelect
          label="Conjunto de widgets"
          description="Widgets mostrados en la barra"
          value={widgetSet}
          options={WIDGET_SETS}
          onChange={v => update('widgets', v, setWidgetSet)}
        />
      </TUISection>

      <TUISection title="APLICACIONES PREDETERMINADAS">
        <TUISelect
          label="Terminal"
          description="set-terminal — abre con Mod+Return"
          value={terminal}
          options={TERMINALS}
          onChange={v => update('terminal', v, setTerminal)}
        />
        <TUISelect
          label="Navegador"
          description="set-browser — abre con Mod+b"
          value={browser}
          options={BROWSERS}
          onChange={v => update('browser', v, setBrowser)}
        />
        <TUISelect
          label="Gestor de archivos"
          description="Abre con Mod+e"
          value={fileManager}
          options={FILE_MANAGERS}
          onChange={v => update('fm', v, setFm)}
        />
      </TUISection>

      <TUISection title="VENTANAS Y ESPACIADO">
        <TUISlider
          label="Gaps entre ventanas"
          description="margin_gap — píxeles de separación entre ventanas"
          value={gaps}
          min={0}
          max={32}
          unit="px"
          onChange={v => update('gaps', v, setGaps)}
        />
        <TUISlider
          label="Ancho del borde"
          description="border_width — píxeles del borde de la ventana activa"
          value={borderWidth}
          min={0}
          max={6}
          unit="px"
          onChange={v => update('borderW', v, setBorderWidth)}
        />
        <TUIToggle
          label="Foco sigue al ratón"
          description="follow_mouse_focus — foco automático al hover"
          value={focusFollowsMouse}
          onChange={v => update('ffm', v, setFfm)}
        />
        <TUIToggle
          label="Pantalla completa automática"
          description="auto_fullscreen — respeta las hints de fullscreen de X11"
          value={autoFullscreen}
          onChange={v => update('afs', v, setAfs)}
        />
      </TUISection>

      <TUISection title="WALLPAPER">
        <TUISelect
          label="Modo"
          description="Cómo escalar el fondo de pantalla"
          value={wpMode}
          options={WALLPAPER_MODES}
          onChange={v => update('wpMode', v, setWpMode)}
        />
      </TUISection>

      <TUISection title="LAYOUTS" collapsible defaultOpen={true}>
        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${C.border}` }}>
          <div style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10, marginBottom: 8 }}>
            Orden activo: {activeLayouts.join(' → ') || '(ninguno)'}
          </div>
        </div>

        {LAYOUT_NAMES.map(name => (
          <div
            key={name}
            style={{
              display: 'flex', alignItems: 'center', gap: 12,
              padding: '8px 12px',
              borderBottom: `1px solid ${C.border}`,
              background: activeLayout === name ? C.accentDim : 'transparent',
              cursor: 'pointer',
            }}
            onClick={() => setActiveLayout(name)}
          >
            <div style={{ width: 80, flexShrink: 0 }}>
              <LayoutPreview name={name} />
            </div>
            <div style={{ flex: 1 }}>
              <div style={{ fontFamily: 'monospace', fontSize: 13, color: layouts[name] ? C.textPrimary : C.textMuted }}>
                {name}
              </div>
              <div style={{ fontFamily: 'monospace', fontSize: 10, color: C.textMuted, marginTop: 2 }}>
                qtile.layout.{name}
              </div>
            </div>
            <div onClick={e => { e.stopPropagation(); toggleLayout(name); }}>
              <span style={{ fontFamily: 'monospace', fontSize: 13, color: layouts[name] ? C.success : C.textMuted }}>
                {layouts[name] ? '● Activo' : '○ Inactivo'}
              </span>
            </div>
          </div>
        ))}
      </TUISection>

      <TUISection title="ACCIONES">
        <div style={{ padding: '10px 12px', display: 'flex', gap: 10 }}>
          <TUIButton
            label={reloading ? 'Recargando config...' : 'Recargar config'}
            onClick={reloadConfig}
            icon="↻"
          />
          <TUIButton
            label="Abrir config.py"
            onClick={async () => { onModified(); }}
            icon="✦"
          />
          <TUIButton
            label="Reiniciar Qtile"
            onClick={async () => { onModified(); }}
            variant="danger"
            icon="↺"
          />
        </div>
      </TUISection>
    </div>
  );
}
