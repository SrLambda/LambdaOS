import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUISection } from '../components/tui/TUISection';
import { TUIToggle } from '../components/tui/TUIToggle';

const THEMES = [
  { id: 'lambda-dark',   name: 'Lambda Dark',   primary: '#7D56F4', bg: '#0D0D0D', desc: 'Oscuro con acento púrpura' },
  { id: 'lambda-rose',   name: 'Lambda Rose',   primary: '#F47F7F', bg: '#1A0A0A', desc: 'Oscuro con acento rosa' },
  { id: 'lambda-teal',   name: 'Lambda Teal',   primary: '#56F4D0', bg: '#0A1A1A', desc: 'Oscuro con acento aguamarina' },
  { id: 'catppuccin',    name: 'Catppuccin Mocha',primary: '#CBA6F7', bg: '#1E1E2E', desc: 'Paleta catppuccin mocha' },
  { id: 'gruvbox',       name: 'Gruvbox Dark',  primary: '#FABD2F', bg: '#282828', desc: 'Paleta gruvbox clásica' },
  { id: 'nord',          name: 'Nord',           primary: '#88C0D0', bg: '#2E3440', desc: 'Paleta nórdica minimalista' },
];

const FONTS = ['JetBrains Mono','Fira Code','Cascadia Code','Iosevka','Hack','Source Code Pro','Inconsolata'];
const WALLPAPERS = ['lambda-nebula.png','dark-grid.png','geometric-dark.png','minimal-dots.png','circuit-board.png','mountain-night.png'];

function ls(k: string, d: any) { try { const v = localStorage.getItem('app_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('app_'+k, JSON.stringify(v)); }

export function AppearanceModule({ onModified }: { onModified: () => void }) {
  const [theme, setTheme]     = useState(() => ls('theme', 'lambda-dark'));
  const [font, setFont]       = useState(() => ls('font', 'JetBrains Mono'));
  const [fontSize, setFontSize] = useState(() => ls('fontSize', 13));
  const [opacity, setOpacity] = useState(() => ls('opacity', 90));
  const [wallpaper, setWallpaper] = useState(() => ls('wallpaper', WALLPAPERS[0]));
  const [animations, setAnimations] = useState(() => ls('anim', true));

  function update(key: string, val: any, setter: (v: any) => void) {
    setter(val); ss(key, val); onModified();
  }

  const activeTheme = THEMES.find(t => t.id === theme) ?? THEMES[0];

  return (
    <div>
      {/* Theme grid */}
      <TUISection title="TEMA">
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 6, padding: '10px 12px' }}>
          {THEMES.map(t => (
            <div
              key={t.id}
              onClick={() => update('theme', t.id, setTheme)}
              style={{
                border: `1px solid ${t.id === theme ? t.primary : C.border}`,
                padding: 10, cursor: 'pointer',
                background: t.id === theme ? `${t.primary}15` : C.surface,
                transition: 'all 0.1s',
              }}
            >
              {/* Mini preview */}
              <div style={{ display: 'flex', gap: 4, marginBottom: 6 }}>
                <div style={{ width: 32, height: 20, background: t.bg, border: `1px solid ${t.primary}40` }}>
                  <div style={{ width: '60%', height: 3, background: t.primary, margin: '4px auto 0' }} />
                </div>
              </div>
              <div style={{ color: t.id === theme ? t.primary : C.textPrimary, fontSize: 11, fontFamily: 'monospace' }}>
                {t.id === theme ? '● ' : '○ '}{t.name}
              </div>
              <div style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', marginTop: 2 }}>{t.desc}</div>
            </div>
          ))}
        </div>
      </TUISection>

      <TUISection title="TIPOGRAFÍA">
        <TUISelect label="Fuente" description="Fuente monoespaciada del sistema" value={font} options={FONTS} onChange={v => update('font', v, setFont)} />
        <TUISlider label="Tamaño de fuente" value={fontSize} min={9} max={20} unit="px" onChange={v => update('fontSize', v, setFontSize)} />

        {/* Live preview */}
        <div style={{ padding: '10px 14px', borderBottom: `1px solid ${C.border}` }}>
          <div style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', marginBottom: 6 }}>VISTA PREVIA</div>
          <div style={{ fontFamily: font + ', monospace', fontSize: fontSize, color: activeTheme.primary, lineHeight: 1.5 }}>
            El veloz murciélago hindú comía feliz cardillo y kiwi.
            <br />
            <span style={{ color: C.textSecondary }}>
              0123456789  !@#$%^&*()  {'{'}{'}'}[]&lt;&gt;
            </span>
          </div>
        </div>
      </TUISection>

      <TUISection title="VENTANAS Y COMPOSICIÓN">
        <TUISlider label="Opacidad de ventanas" value={opacity} min={40} max={100} onChange={v => update('opacity', v, setOpacity)} />
        <TUIToggle label="Animaciones" description="Transiciones y efectos de ventana" value={animations} onChange={v => update('anim', v, setAnimations)} />
      </TUISection>

      <TUISection title="WALLPAPER">
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 6, padding: '10px 12px' }}>
          {WALLPAPERS.map(w => (
            <div
              key={w}
              onClick={() => update('wallpaper', w, setWallpaper)}
              style={{
                height: 48,
                border: `1px solid ${w === wallpaper ? C.accent : C.border}`,
                background: C.surface2,
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                cursor: 'pointer',
                position: 'relative',
                overflow: 'hidden',
              }}
            >
              <span style={{ color: C.textMuted, fontSize: 9, fontFamily: 'monospace', textAlign: 'center', padding: 4 }}>
                {w === wallpaper && <span style={{ color: C.accent }}>●</span>}
                {'\n'}{w.replace('.png', '')}
              </span>
            </div>
          ))}
        </div>
      </TUISection>
    </div>
  );
}
