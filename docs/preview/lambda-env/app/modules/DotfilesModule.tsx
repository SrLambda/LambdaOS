// Mirrors: src/lambda-env/internal/modules/dotfiles/main.go
// Actions: run, backup, stow, list-conflicts, list profiles
import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUIInput } from '../components/tui/TUIInput';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIProgress } from '../components/tui/TUIProgress';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIModal } from '../components/tui/TUIModal';

// ── Types ────────────────────────────────────────────────────────────────────

type PackageStatus = 'stowed' | 'unstowed' | 'conflict';

interface DotPackage {
  name: string;
  desc: string;
  files: number;
  status: PackageStatus;
}

interface Conflict {
  pkg: string;
  file: string;
  existing: string;
}

interface BackupEntry {
  id: string;
  date: string;
  size: string;
  note: string;
}

// ── Data ────────────────────────────────────────────────────────────────────

const INITIAL_PACKAGES: DotPackage[] = [
  { name: 'zsh',       desc: '.zshrc, .zsh_aliases, .zprofile',                  files: 4,  status: 'stowed'   },
  { name: 'nvim',      desc: 'init.lua, lua/plugins, lua/config',                 files: 18, status: 'stowed'   },
  { name: 'qtile',     desc: 'config.py, autostart.sh, themes/',                  files: 6,  status: 'stowed'   },
  { name: 'alacritty', desc: 'alacritty.toml',                                    files: 1,  status: 'stowed'   },
  { name: 'tmux',      desc: '.tmux.conf, tmux/plugins/',                         files: 3,  status: 'stowed'   },
  { name: 'git',       desc: '.gitconfig, .gitignore_global',                     files: 2,  status: 'stowed'   },
  { name: 'ssh',       desc: '.ssh/config',                                       files: 1,  status: 'conflict' },
  { name: 'starship',  desc: 'starship.toml',                                     files: 1,  status: 'unstowed' },
  { name: 'dunst',     desc: 'dunstrc',                                           files: 1,  status: 'unstowed' },
  { name: 'picom',     desc: 'picom.conf',                                        files: 1,  status: 'stowed'   },
  { name: 'rofi',      desc: 'config.rasi, themes/',                              files: 5,  status: 'stowed'   },
  { name: 'fonts',     desc: '.local/share/fonts/',                               files: 12, status: 'unstowed' },
];

const INITIAL_CONFLICTS: Conflict[] = [
  { pkg: 'ssh', file: '~/.ssh/config', existing: 'symlink → /etc/ssh/ssh_config (sistema)' },
];

const INITIAL_BACKUPS: BackupEntry[] = [
  { id: 'b001', date: '2026-06-07 08:45', size: '2.4 MB', note: 'Pre-actualización kernel 6.6' },
  { id: 'b002', date: '2026-05-22 14:30', size: '2.1 MB', note: 'Snapshot semanal' },
  { id: 'b003', date: '2026-05-15 09:12', size: '1.9 MB', note: 'Antes de migrar a catppuccin' },
];

const PROFILES = ['personal', 'work', 'minimal', 'gaming'];
const REMOTE_URLS = ['git@github.com:user/dotfiles.git', 'https://github.com/user/dotfiles', '(personalizado)'];

// ── Persistence ──────────────────────────────────────────────────────────────

function ls(k: string, d: any) {
  try { const v = localStorage.getItem('dots_' + k); return v ? JSON.parse(v) : d; } catch { return d; }
}
function ss(k: string, v: any) { localStorage.setItem('dots_' + k, JSON.stringify(v)); }

// ── Component ────────────────────────────────────────────────────────────────

export function DotfilesModule({ onModified }: { onModified: () => void }) {
  const [packages, setPackages]     = useState<DotPackage[]>(INITIAL_PACKAGES);
  const [conflicts, setConflicts]   = useState<Conflict[]>(INITIAL_CONFLICTS);
  const [backups, setBackups]       = useState<BackupEntry[]>(INITIAL_BACKUPS);
  const [dotfilesDir, setDotfilesDir] = useState(() => ls('dir', '~/dotfiles'));
  const [profile, setProfile]       = useState(() => ls('profile', 'personal'));
  const [remoteUrl, setRemoteUrl]   = useState(() => ls('remote', REMOTE_URLS[0]));
  const [autoSync, setAutoSync]     = useState(() => ls('autoSync', false));
  const [backing, setBacking]       = useState(false);
  const [backupProgress, setBackupProgress] = useState(0);
  const [stowingPkg, setStowingPkg] = useState<string | null>(null);
  const [confirmConflict, setConfirmConflict] = useState<Conflict | null>(null);
  const [pulling, setPulling]       = useState(false);

  function update(k: string, v: any, setter: (x: any) => void) {
    setter(v); ss(k, v); onModified();
  }

  async function stow(pkg: string, action: 'stow' | 'unstow') {
    setStowingPkg(pkg);
    await new Promise(r => setTimeout(r, 800));
    setPackages(ps => ps.map(p => p.name === pkg
      ? { ...p, status: action === 'stow' ? 'stowed' : 'unstowed' }
      : p
    ));
    setStowingPkg(null);
    onModified();
  }

  async function backup() {
    setBacking(true);
    for (let i = 0; i <= 100; i += 10) {
      await new Promise(r => setTimeout(r, 120));
      setBackupProgress(i);
    }
    const entry: BackupEntry = {
      id: 'b' + Date.now(),
      date: new Date().toLocaleString('es-ES', { hour12: false }).replace(',', ''),
      size: '2.5 MB',
      note: 'Backup manual',
    };
    setBackups(b => [entry, ...b]);
    setBacking(false);
    setBackupProgress(0);
    onModified();
  }

  async function pullRemote() {
    setPulling(true);
    await new Promise(r => setTimeout(r, 1500));
    setPulling(false);
    onModified();
  }

  function resolveConflict(c: Conflict, action: 'overwrite' | 'skip') {
    if (action === 'overwrite') {
      setPackages(ps => ps.map(p => p.name === c.pkg ? { ...p, status: 'stowed' } : p));
      setConflicts(cs => cs.filter(x => x.file !== c.file));
    }
    setConfirmConflict(null);
    onModified();
  }

  const stowedCount   = packages.filter(p => p.status === 'stowed').length;
  const conflictCount = packages.filter(p => p.status === 'conflict').length;

  const statusColor: Record<PackageStatus, string> = {
    stowed:   C.success,
    unstowed: C.textMuted,
    conflict: C.error,
  };
  const statusLabel: Record<PackageStatus, string> = {
    stowed:   '● stowed',
    unstowed: '○ unstowed',
    conflict: '⚠ conflicto',
  };

  return (
    <div>
      {/* Conflict override modal */}
      {confirmConflict && (
        <TUIModal
          title={`Conflicto: ${confirmConflict.file}`}
          description={`Ya existe: ${confirmConflict.existing}\n\n¿Sobrescribir con el symlink de dotfiles?`}
          confirmLabel="Sobrescribir"
          cancelLabel="Ignorar"
          variant="danger"
          onConfirm={() => resolveConflict(confirmConflict, 'overwrite')}
          onCancel={() => resolveConflict(confirmConflict, 'skip')}
        />
      )}

      {/* Summary */}
      <div style={{
        padding: '10px 14px', background: C.surface, marginBottom: 4,
        fontFamily: 'monospace', display: 'flex', gap: 24, alignItems: 'center',
        borderLeft: `3px solid ${conflictCount > 0 ? C.error : C.success}`,
      }}>
        <div>
          <div style={{ color: C.textPrimary, fontSize: 12 }}>{dotfilesDir}</div>
          <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>Perfil activo: {profile}</div>
        </div>
        <div style={{ display: 'flex', gap: 20, marginLeft: 'auto' }}>
          <span style={{ color: C.success, fontSize: 11 }}>● {stowedCount} stowed</span>
          {conflictCount > 0 && <span style={{ color: C.error, fontSize: 11 }}>⚠ {conflictCount} conflictos</span>}
          <span style={{ color: C.textMuted, fontSize: 11 }}>○ {packages.length - stowedCount - conflictCount} libres</span>
        </div>
      </div>

      <TUISection title="CONFIGURACIÓN">
        <TUIInput
          label="Directorio de dotfiles"
          description="Ruta local del repositorio de dotfiles"
          value={dotfilesDir}
          onChange={v => update('dir', v, setDotfilesDir)}
        />
        <TUISelect
          label="Perfil activo"
          description="Conjunto de paquetes a gestionar"
          value={profile}
          options={PROFILES}
          onChange={v => update('profile', v, setProfile)}
        />
        <TUIInput
          label="Remote URL"
          description="Repositorio Git remoto para sincronización"
          value={remoteUrl}
          onChange={v => update('remote', v, setRemoteUrl)}
        />
        <TUIToggle
          label="Sincronización automática"
          description="git pull automático al iniciar sesión"
          value={autoSync}
          onChange={v => update('autoSync', v, setAutoSync)}
        />
      </TUISection>

      {/* Conflicts banner */}
      {conflicts.length > 0 && (
        <TUISection title={`⚠ CONFLICTOS (${conflicts.length})`}>
          {conflicts.map(c => (
            <div key={c.file} style={{
              padding: '10px 14px', borderBottom: `1px solid ${C.border}`,
              background: `${C.error}10`,
              borderLeft: `3px solid ${C.error}`,
              fontFamily: 'monospace',
            }}>
              <div style={{ color: C.error, fontSize: 12 }}>{c.file}</div>
              <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>
                {c.existing}
              </div>
              <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
                <TUIButton
                  label="Resolver"
                  onClick={async () => setConfirmConflict(c)}
                  variant="danger"
                />
                <TUIButton
                  label="Ignorar"
                  onClick={async () => setConflicts(cs => cs.filter(x => x.file !== c.file))}
                />
              </div>
            </div>
          ))}
        </TUISection>
      )}

      <TUISection title="PAQUETES">
        {/* Header */}
        <div style={{
          display: 'grid', gridTemplateColumns: '130px 1fr 40px 120px auto',
          padding: '5px 12px', background: C.surface, gap: 8,
          borderBottom: `1px solid ${C.border}`, fontFamily: 'monospace',
        }}>
          {['Paquete', 'Archivos', 'N', 'Estado', ''].map(h => (
            <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
          ))}
        </div>

        {packages.map(pkg => (
          <div key={pkg.name} style={{
            display: 'grid', gridTemplateColumns: '130px 1fr 40px 120px auto',
            padding: '8px 12px', borderBottom: `1px solid ${C.border}`,
            gap: 8, fontFamily: 'monospace', alignItems: 'center',
          }}>
            <span style={{ color: C.textPrimary, fontSize: 12 }}>
              {stowingPkg === pkg.name ? '⟳ ' : ''}{pkg.name}
            </span>
            <span style={{ color: C.textMuted, fontSize: 10, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
              {pkg.desc}
            </span>
            <span style={{ color: C.textMuted, fontSize: 11 }}>{pkg.files}</span>
            <span style={{ color: statusColor[pkg.status], fontSize: 11 }}>
              {statusLabel[pkg.status]}
            </span>
            <div style={{ display: 'flex', gap: 6 }}>
              {pkg.status === 'unstowed' && (
                <TUIButton
                  label="Stow"
                  onClick={async () => stow(pkg.name, 'stow')}
                  variant="primary"
                />
              )}
              {pkg.status === 'stowed' && (
                <TUIButton
                  label="Unstow"
                  onClick={async () => stow(pkg.name, 'unstow')}
                />
              )}
              {pkg.status === 'conflict' && (
                <TUIButton
                  label="Resolver"
                  onClick={async () => {
                    const c = conflicts.find(x => x.pkg === pkg.name);
                    if (c) setConfirmConflict(c);
                  }}
                  variant="danger"
                />
              )}
            </div>
          </div>
        ))}

        {/* Stow all / Unstow all */}
        <div style={{ padding: '10px 12px', display: 'flex', gap: 8 }}>
          <TUIButton
            label="Stow todos"
            onClick={async () => {
              for (const pkg of packages.filter(p => p.status === 'unstowed')) {
                await stow(pkg.name, 'stow');
              }
            }}
            variant="primary"
          />
          <TUIButton
            label="Unstow todos"
            onClick={async () => {
              for (const pkg of packages.filter(p => p.status === 'stowed')) {
                await stow(pkg.name, 'unstow');
              }
            }}
          />
        </div>
      </TUISection>

      <TUISection title="SINCRONIZACIÓN GIT" collapsible defaultOpen={false}>
        <div style={{ padding: '10px 12px', borderBottom: `1px solid ${C.border}`, display: 'flex', gap: 8, flexWrap: 'wrap' }}>
          <TUIButton
            label={pulling ? 'Sincronizando...' : 'git pull'}
            onClick={pullRemote}
            icon="↓"
          />
          <TUIButton
            label="git push"
            onClick={async () => { onModified(); }}
            icon="↑"
          />
          <TUIButton
            label="git status"
            onClick={async () => { onModified(); }}
          />
          <TUIButton
            label="Ver diff"
            onClick={async () => { onModified(); }}
          />
        </div>
        {/* Simulated git status output */}
        <div style={{ padding: '10px 14px', fontFamily: 'monospace', background: C.surface, fontSize: 11 }}>
          <div style={{ color: C.success }}>On branch main</div>
          <div style={{ color: C.textSecondary }}>Your branch is up to date with 'origin/main'.</div>
          <div style={{ color: C.textSecondary, marginTop: 6 }}>nothing to commit, working tree clean</div>
        </div>
      </TUISection>

      <TUISection title="BACKUPS" collapsible defaultOpen={false}>
        {backing && (
          <div style={{ padding: '10px 14px' }}>
            <TUIProgress value={backupProgress} label="Creando backup..." color={C.accent} />
          </div>
        )}
        <div style={{ padding: '8px 12px', borderBottom: `1px solid ${C.border}` }}>
          <TUIButton
            label={backing ? 'Creando backup...' : 'Crear backup ahora'}
            onClick={backup}
            icon="◎"
            disabled={backing}
          />
        </div>

        {/* Backup list */}
        {backups.map(b => (
          <div key={b.id} style={{
            display: 'flex', alignItems: 'center', gap: 12,
            padding: '8px 14px', borderBottom: `1px solid ${C.border}`,
            fontFamily: 'monospace',
          }}>
            <div style={{ flex: 1 }}>
              <div style={{ color: C.textPrimary, fontSize: 12 }}>{b.date}</div>
              <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>{b.note}</div>
            </div>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{b.size}</span>
            <TUIButton
              label="Restaurar"
              onClick={async () => { onModified(); }}
            />
          </div>
        ))}
      </TUISection>
    </div>
  );
}
