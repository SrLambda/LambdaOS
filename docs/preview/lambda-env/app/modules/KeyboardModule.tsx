import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUISection } from '../components/tui/TUISection';
import { TUIToggle } from '../components/tui/TUIToggle';

const LAYOUTS: Record<string, string[]> = {
  'Español (es)':    ['Estándar', 'Variante Latinoamericana', 'Dvorak'],
  'English (us)':    ['Estándar', 'Dvorak', 'Colemak', 'Workman', 'Intl'],
  'Français (fr)':   ['Estándar (AZERTY)', 'Bépo'],
  'Deutsch (de)':    ['Estándar', 'Dvorak'],
  'Português (pt)':  ['Brasil','Portugal'],
  '日本語 (jp)':       ['106-key','109-key'],
  'Русский (ru)':    ['Estándar','Phonetic'],
};

const PREVIEW_ROWS = [
  ['q','w','e','r','t','y','u','i','o','p'],
  ['a','s','d','f','g','h','j','k','l'],
  ['z','x','c','v','b','n','m'],
];

function ls(k: string, d: any) { try { const v = localStorage.getItem('kb_'+k); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(k: string, v: any) { localStorage.setItem('kb_'+k, JSON.stringify(v)); }

export function KeyboardModule({ onModified }: { onModified: () => void }) {
  const [layout, setLayout]   = useState(() => ls('layout', 'Español (es)'));
  const [variant, setVariant] = useState(() => ls('variant', 'Estándar'));
  const [delay, setDelay]     = useState(() => ls('delay', 250));
  const [rate, setRate]       = useState(() => ls('rate', 25));
  const [numlock, setNumlock] = useState(() => ls('numlock', true));
  const [compose, setCompose] = useState(() => ls('compose', true));

  function updateLayout(l: string) {
    setLayout(l); ss('layout', l);
    const vars = LAYOUTS[l] ?? ['Estándar'];
    const v = vars[0];
    setVariant(v); ss('variant', v);
    onModified();
  }

  return (
    <div>
      <TUISection title="DISTRIBUCIÓN">
        <TUISelect
          label="Layout de teclado"
          description="Cambio inmediato al seleccionar"
          value={layout}
          options={Object.keys(LAYOUTS)}
          onChange={updateLayout}
        />
        <TUISelect
          label="Variante"
          description="Variante del layout seleccionado"
          value={variant}
          options={LAYOUTS[layout] ?? ['Estándar']}
          onChange={v => { setVariant(v); ss('variant', v); onModified(); }}
        />
      </TUISection>

      {/* Keyboard preview */}
      <div style={{ padding: '12px 14px', background: C.surface, margin: '4px 0' }}>
        <div style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', marginBottom: 8 }}>
          VISTA PREVIA — {layout} · {variant}
        </div>
        {PREVIEW_ROWS.map((row, ri) => (
          <div key={ri} style={{
            display: 'flex', gap: 4, marginBottom: 4,
            paddingLeft: ri * 12,
          }}>
            {row.map(k => (
              <div key={k} style={{
                width: 28, height: 28,
                border: `1px solid ${C.border}`,
                background: C.surface2,
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                fontFamily: 'monospace', fontSize: 11,
                color: C.textPrimary,
              }}>{k}</div>
            ))}
          </div>
        ))}
      </div>

      <TUISection title="VELOCIDAD DE REPETICIÓN">
        <TUISlider
          label="Retardo inicial"
          description="Tiempo antes de que la tecla empiece a repetir"
          value={delay}
          min={100}
          max={1000}
          unit="ms"
          onChange={v => { setDelay(v); ss('delay', v); onModified(); }}
        />
        <TUISlider
          label="Velocidad de repetición"
          description="Caracteres por segundo al mantener una tecla"
          value={rate}
          min={5}
          max={50}
          unit="/s"
          onChange={v => { setRate(v); ss('rate', v); onModified(); }}
        />
      </TUISection>

      <TUISection title="OPCIONES">
        <TUIToggle
          label="NumLock al inicio"
          description="Activar NumLock automáticamente al arrancar"
          value={numlock}
          onChange={v => { setNumlock(v); ss('numlock', v); onModified(); }}
        />
        <TUIToggle
          label="Tecla Compose"
          description="Activar tecla Compose para caracteres especiales (AltGr)"
          value={compose}
          onChange={v => { setCompose(v); ss('compose', v); onModified(); }}
        />
      </TUISection>
    </div>
  );
}
