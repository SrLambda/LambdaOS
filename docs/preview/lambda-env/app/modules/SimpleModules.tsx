// Módulos simples implementados en un solo archivo
import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { ICONS } from '../data/icon-map';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIInput } from '../components/tui/TUIInput';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIModal } from '../components/tui/TUIModal';

// ── FECHA Y HORA ──────────────────────────────────────────────────────────────

const TIMEZONES = ['UTC-08:00 (Los Angeles)','UTC-05:00 (Nueva York)','UTC+00:00 (Londres/UTC)','UTC+01:00 (Madrid/CET)','UTC+02:00 (Helsinki)','UTC+05:30 (Mumbai)','UTC+09:00 (Tokio)'];
const DATE_FORMATS = ['DD/MM/YYYY','MM/DD/YYYY','YYYY-MM-DD','D de MMMM de YYYY'];
const TIME_FORMATS = ['HH:MM:SS (24h)','h:MM:SS AM/PM (12h)'];

function ls(ns: string, k: string, d: any) { try { const v = localStorage.getItem(`${ns}_${k}`); return v ? JSON.parse(v) : d; } catch { return d; } }
function ss(ns: string, k: string, v: any) { localStorage.setItem(`${ns}_${k}`, JSON.stringify(v)); }

export function DateTimeModule({ onModified }: { onModified: () => void }) {
  const [ntp, setNtp]   = useState(() => ls('dt','ntp',true));
  const [tz, setTz]     = useState(() => ls('dt','tz','UTC+01:00 (Madrid/CET)'));
  const [df, setDf]     = useState(() => ls('dt','df','DD/MM/YYYY'));
  const [tf, setTf]     = useState(() => ls('dt','tf','HH:MM:SS (24h)'));
  const now = new Date();

  function update(k: string, v: any, s: (x: any) => void) { s(v); ss('dt', k, v); onModified(); }

  return (
    <div>
      <div style={{ padding: '10px 14px', background: C.surface, marginBottom: 4, fontFamily: 'monospace', textAlign: 'center' }}>
        <div style={{ color: C.accent, fontSize: 28 }}>{now.toLocaleTimeString('es-ES', { hour12: false })}</div>
        <div style={{ color: C.textSecondary, fontSize: 13, marginTop: 4 }}>{now.toLocaleDateString('es-ES', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</div>
        <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>{tz}</div>
      </div>
      <TUISection title="ZONA HORARIA Y SINCRONIZACIÓN">
        <TUIToggle label="NTP automático" description="Sincronizar hora con servidores NTP (timedatectl)" value={ntp} onChange={v => update('ntp', v, setNtp)} />
        <TUISelect label="Zona horaria" description="timedatectl set-timezone" value={tz} options={TIMEZONES} onChange={v => update('tz', v, setTz)} disabled={!ntp} />
      </TUISection>
      <TUISection title="FORMATO">
        <TUISelect label="Formato de fecha" value={df} options={DATE_FORMATS} onChange={v => update('df', v, setDf)} />
        <TUISelect label="Formato de hora" value={tf} options={TIME_FORMATS} onChange={v => update('tf', v, setTf)} />
      </TUISection>
    </div>
  );
}

// ── USUARIOS ──────────────────────────────────────────────────────────────────

interface User { name: string; groups: string[]; shell: string; last: string; locked: boolean; }
const USERS_LIST: User[] = [
  { name: 'root',   groups: ['root'],                shell: '/bin/zsh',  last: '2026-06-07 09:12', locked: false },
  { name: 'lambda', groups: ['wheel','audio','video'],shell: '/bin/zsh',  last: '2026-06-07 08:45', locked: false },
  { name: 'guest',  groups: ['users'],               shell: '/bin/bash', last: '2026-05-20 15:30', locked: true  },
];

export function UsersModule({ onModified }: { onModified: () => void }) {
  const [users, setUsers] = useState(USERS_LIST);
  const [selected, setSelected] = useState<string | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null);

  function toggleLock(name: string) {
    setUsers(u => u.map(usr => usr.name === name ? { ...usr, locked: !usr.locked } : usr));
    onModified();
  }

  return (
    <div>
      {deleteTarget && (
        <TUIModal
          title={`¿Eliminar usuario "${deleteTarget}"?`}
          description="Se eliminará la cuenta y todos los archivos del usuario. Esta acción es irreversible."
          confirmLabel="Eliminar usuario"
          variant="danger"
          onConfirm={() => { setUsers(u => u.filter(usr => usr.name !== deleteTarget)); setDeleteTarget(null); onModified(); }}
          onCancel={() => setDeleteTarget(null)}
        />
      )}
      <TUISection title="CUENTAS DE USUARIO" rootRequired>
        {users.map(usr => {
          const isSel = selected === usr.name;
          return (
            <div key={usr.name}>
              <div
                onClick={() => setSelected(isSel ? null : usr.name)}
                style={{
                  display: 'flex', alignItems: 'center', gap: 12, padding: '10px 14px',
                  borderBottom: `1px solid ${C.border}`, cursor: 'pointer',
                  background: isSel ? C.accentDim : 'transparent',
                }}
              >
                <div style={{ width: 32, height: 32, border: `1px solid ${C.border}`, display: 'flex', alignItems: 'center', justifyContent: 'center', background: C.surface, color: C.accent, fontFamily: 'monospace', fontSize: 14 }}>
                  {usr.name[0].toUpperCase()}
                </div>
                <div style={{ flex: 1 }}>
                  <div style={{ color: usr.locked ? C.textMuted : C.textPrimary, fontFamily: 'monospace', fontSize: 13 }}>
                    {usr.name} {usr.locked && <span style={{ color: C.error, fontSize: 10 }}>{ICONS.widgets.lock.nerd} bloqueado</span>}
                  </div>
                  <div style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>
                    {usr.groups.join(', ')} · {usr.shell} · Última sesión: {usr.last}
                  </div>
                </div>
                <span style={{ color: C.accent, fontSize: 10, fontFamily: 'monospace' }}>{isSel ? '▼' : '▶'}</span>
              </div>
              {isSel && (
                <div style={{ padding: '8px 14px 10px 58px', borderBottom: `1px solid ${C.border}`, background: C.surface, display: 'flex', gap: 8 }}>
                  <TUIButton label={usr.locked ? 'Desbloquear' : 'Bloquear'} onClick={async () => toggleLock(usr.name)} />
                  <TUIButton label="Cambiar contraseña" onClick={async () => { onModified(); }} />
                  {usr.name !== 'root' && <TUIButton label="Eliminar" onClick={async () => setDeleteTarget(usr.name)} variant="danger" />}
                </div>
              )}
            </div>
          );
        })}
        <div style={{ padding: '10px 14px' }}>
          <TUIButton label="Crear nuevo usuario" onClick={async () => { onModified(); }} variant="primary" icon="+" />
        </div>
      </TUISection>
    </div>
  );
}

// ── DEFAULTS ──────────────────────────────────────────────────────────────────

const DEFAULTS_MAP = [
  { category: 'Navegador web',    current: 'Firefox',         options: ['Firefox','Chromium','LibreWolf','qutebrowser'] },
  { category: 'Editor de texto',  current: 'Neovim',          options: ['Neovim','Vim','Nano','Helix','VSCode'] },
  { category: 'Terminal',         current: 'Alacritty',       options: ['Alacritty','Kitty','WezTerm','Foot','xterm'] },
  { category: 'Reproductor video',current: 'mpv',             options: ['mpv','VLC','celluloid'] },
  { category: 'Reproductor audio',current: 'Spotify',         options: ['Spotify','cmus','ncmpcpp','Rhythmbox'] },
  { category: 'Visor de imágenes',current: 'imv',             options: ['imv','feh','eog','nomacs'] },
  { category: 'File manager',     current: 'Thunar',          options: ['Thunar','Nautilus','Ranger','nnn','lf'] },
  { category: 'PDF viewer',       current: 'zathura',         options: ['zathura','Evince','Okular','Mupdf'] },
  { category: 'Correo',           current: 'Thunderbird',     options: ['Thunderbird','Evolution','aerc','neomutt'] },
];

export function DefaultsModule({ onModified }: { onModified: () => void }) {
  const [values, setValues] = useState(() => Object.fromEntries(DEFAULTS_MAP.map(d => [d.category, ls('def', d.category, d.current)])));

  return (
    <div>
      <TUISection title="APLICACIONES POR DEFECTO">
        {DEFAULTS_MAP.map(d => (
          <TUISelect
            key={d.category}
            label={d.category}
            value={values[d.category]}
            options={d.options}
            onChange={v => { setValues(prev => ({ ...prev, [d.category]: v })); ss('def', d.category, v); onModified(); }}
          />
        ))}
      </TUISection>
    </div>
  );
}

// ── AUTOSTART ─────────────────────────────────────────────────────────────────

const AUTOSTART_ITEMS = [
  { name: 'NetworkManager-applet', description: 'Bandeja de red', enabled: true  },
  { name: 'blueman-applet',        description: 'Bandeja Bluetooth', enabled: true  },
  { name: 'dunst',                 description: 'Daemon de notificaciones', enabled: true  },
  { name: 'picom',                 description: 'Compositor X11', enabled: true  },
  { name: 'flameshot',             description: 'Capturas de pantalla', enabled: true  },
  { name: 'spotify',               description: 'Reproductor de música', enabled: false },
  { name: 'discord',               description: 'Cliente de chat', enabled: false },
  { name: 'nextcloud-client',      description: 'Sincronización en la nube', enabled: false },
];

export function AutostartModule({ onModified }: { onModified: () => void }) {
  const [items, setItems] = useState(AUTOSTART_ITEMS);

  return (
    <div>
      <TUISection title="PROGRAMAS AL INICIAR SESIÓN">
        {items.map(item => (
          <TUIToggle
            key={item.name}
            label={item.name}
            description={item.description}
            value={item.enabled}
            onChange={v => { setItems(it => it.map(i => i.name === item.name ? { ...i, enabled: v } : i)); onModified(); }}
          />
        ))}
        <div style={{ padding: '10px 12px' }}>
          <TUIButton label="Añadir programa" onClick={async () => { onModified(); }} icon="+" />
        </div>
      </TUISection>
    </div>
  );
}

// ── NOTIFICACIONES ────────────────────────────────────────────────────────────

const POSITIONS = ['Arriba derecha','Arriba izquierda','Abajo derecha','Abajo izquierda','Arriba centro'];
const NOTIF_APPS = ['Firefox','Discord','Spotify','Thunderbird','Sistema','Updates'];

export function NotificationsModule({ onModified }: { onModified: () => void }) {
  const [pos, setPos]     = useState(() => ls('notif','pos','Arriba derecha'));
  const [timeout, setTimeout_] = useState(() => ls('notif','timeout',5));
  const [dnd, setDnd]     = useState(() => ls('notif','dnd',false));
  const [appToggles, setAppToggles] = useState<Record<string, boolean>>(() =>
    Object.fromEntries(NOTIF_APPS.map(a => [a, ls('notif', a, true)]))
  );

  function update(k: string, v: any, s: (x: any) => void) { s(v); ss('notif', k, v); onModified(); }

  return (
    <div>
      <TUISection title="CONFIGURACIÓN GENERAL">
        <TUIToggle label="No molestar" description="Suprimir todas las notificaciones" value={dnd} onChange={v => update('dnd', v, setDnd)} />
        <TUISelect label="Posición" description="Esquina de la pantalla donde aparecen" value={pos} options={POSITIONS} onChange={v => update('pos', v, setPos)} disabled={dnd} />
        <TUISlider label="Duración" description="Segundos antes de que desaparezca" value={timeout} min={1} max={30} unit="s" onChange={v => update('timeout', v, setTimeout_)} disabled={dnd} />
      </TUISection>
      <TUISection title="POR APLICACIÓN" collapsible defaultOpen={false}>
        {NOTIF_APPS.map(app => (
          <TUIToggle
            key={app}
            label={app}
            value={appToggles[app]}
            disabled={dnd}
            onChange={v => { setAppToggles(t => ({ ...t, [app]: v })); ss('notif', app, v); onModified(); }}
          />
        ))}
      </TUISection>
    </div>
  );
}

// ── FUENTES ───────────────────────────────────────────────────────────────────

const FONTS_LIST = [
  { name: 'JetBrains Mono',   style: 'JetBrains Mono',   installed: true  },
  { name: 'Fira Code',        style: 'Fira Code',         installed: true  },
  { name: 'Cascadia Code',    style: 'Cascadia Code',     installed: true  },
  { name: 'Iosevka',          style: 'Iosevka',           installed: true  },
  { name: 'Hack',             style: 'Hack',              installed: true  },
  { name: 'Victor Mono',      style: 'Victor Mono',       installed: false },
  { name: 'Monaspace Neon',   style: 'monospace',         installed: false },
  { name: 'Berkeley Mono',    style: 'monospace',         installed: false },
];

export function FontsModule({ onModified }: { onModified: () => void }) {
  const [selected, setSelected] = useState('JetBrains Mono');
  const font = FONTS_LIST.find(f => f.name === selected);

  return (
    <div>
      <TUISection title="FUENTES INSTALADAS">
        {FONTS_LIST.map(f => (
          <div
            key={f.name}
            onClick={() => setSelected(f.name)}
            style={{
              display: 'flex', alignItems: 'center', justifyContent: 'space-between',
              padding: '8px 14px', borderBottom: `1px solid ${C.border}`,
              cursor: 'pointer',
              background: f.name === selected ? C.accentDim : 'transparent',
            }}
          >
            <div style={{ flex: 1 }}>
              <span style={{ color: f.name === selected ? C.accent : C.textPrimary, fontFamily: f.style, fontSize: 13 }}>{f.name}</span>
            </div>
            {f.installed
              ? <span style={{ color: C.success, fontSize: 10, fontFamily: 'monospace' }}>● instalada</span>
              : <TUIButton label="Instalar" onClick={async () => { onModified(); }} />
            }
          </div>
        ))}
      </TUISection>

      {/* Preview */}
      {font && font.installed && (
        <div style={{ padding: '12px 14px', background: C.surface, marginTop: 4 }}>
          <div style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace', marginBottom: 8 }}>VISTA PREVIA — {selected}</div>
          {['13px','16px','20px'].map(size => (
            <div key={size} style={{ fontFamily: font.style, fontSize: size, color: C.textPrimary, lineHeight: 1.4, marginBottom: 6 }}>
              El veloz murciélago hindú comía feliz cardillo y kiwi. — <span style={{ color: C.accent }}>0123456789</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
