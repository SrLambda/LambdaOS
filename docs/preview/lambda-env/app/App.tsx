import { useEffect, useRef, useState } from 'react';
import { CATEGORIES, type CategoryDef, type ModuleDef } from './data/categories';
import { C } from './components/tui/tokens';
import { ICONS } from './data/icon-map';

// Modules
import { AudioModule }       from './modules/AudioModule';
import { NetworkModule }     from './modules/NetworkModule';
import { DisplayModule }     from './modules/DisplayModule';
import { PowerModule }       from './modules/PowerModule';
import { BluetoothModule }   from './modules/BluetoothModule';
import { SecurityModule }    from './modules/SecurityModule';
import { KeyboardModule }    from './modules/KeyboardModule';
import { AppearanceModule }  from './modules/AppearanceModule';
import { ServicesModule }    from './modules/ServicesModule';
import { MonitorModule }     from './modules/MonitorModule';
import { UpdatesModule }     from './modules/UpdatesModule';
import { DateTimeModule, UsersModule, DefaultsModule, AutostartModule, NotificationsModule, FontsModule } from './modules/SimpleModules';
import { NeovimModule }      from './modules/NeovimModule';
import { QtileModule }       from './modules/QtileModule';
import { DotfilesModule }    from './modules/DotfilesModule';
import { StorageModule }     from './modules/StorageModule';
import { LogsModule }        from './modules/LogsModule';
import { StubModule }        from './modules/StubModule';

// ── Navigation types ──────────────────────────────────────────────────────────

type NavLevel =
  | { level: 'categories' }
  | { level: 'modules'; catId: string }
  | { level: 'detail'; catId: string; moduleId: string };

// ── Status bar state ──────────────────────────────────────────────────────────

interface StatusState {
  context: string;
  message: string;
  type: 'normal' | 'success' | 'error' | 'warn' | 'loading';
  shortcuts: string;
}

// ── Search ────────────────────────────────────────────────────────────────────

function searchModules(q: string): { cat: CategoryDef; mod: ModuleDef }[] {
  if (!q.trim()) return [];
  const lq = q.toLowerCase();
  const results: { cat: CategoryDef; mod: ModuleDef }[] = [];
  for (const cat of CATEGORIES) {
    for (const mod of cat.modules) {
      if (mod.name.toLowerCase().includes(lq) || mod.description.toLowerCase().includes(lq)) {
        results.push({ cat, mod });
      }
    }
  }
  return results;
}

// ── Module renderer ───────────────────────────────────────────────────────────

function renderModule(moduleId: string, onModified: () => void, onStatusMsg: (m: string) => void) {
  switch (moduleId) {
    case 'audio':         return <AudioModule onModified={onModified} />;
    case 'red':           return <NetworkModule onModified={onModified} />;
    case 'pantalla':      return <DisplayModule onModified={onModified} onStatusMsg={onStatusMsg} />;
    case 'energia':       return <PowerModule onModified={onModified} />;
    case 'bluetooth':     return <BluetoothModule onModified={onModified} />;
    case 'seguridad':     return <SecurityModule onModified={onModified} />;
    case 'teclado':       return <KeyboardModule onModified={onModified} />;
    case 'apariencia':    return <AppearanceModule onModified={onModified} />;
    case 'servicios':     return <ServicesModule onModified={onModified} />;
    case 'monitor':       return <MonitorModule />;
    case 'actualizaciones': return <UpdatesModule onModified={onModified} />;
    case 'fecha':         return <DateTimeModule onModified={onModified} />;
    case 'usuarios':      return <UsersModule onModified={onModified} />;
    case 'defaults':      return <DefaultsModule onModified={onModified} />;
    case 'autostart':     return <AutostartModule onModified={onModified} />;
    case 'notificaciones':return <NotificationsModule onModified={onModified} />;
    case 'fuentes':       return <FontsModule onModified={onModified} />;
    case 'neovim':        return <NeovimModule onModified={onModified} />;
    case 'qtile':         return <QtileModule onModified={onModified} />;
    case 'dotfiles':      return <DotfilesModule onModified={onModified} />;
    case 'almacenamiento':return <StorageModule onModified={onModified} />;
    case 'logs':          return <LogsModule onModified={onModified} />;
    default: {
      const mod = CATEGORIES.flatMap(c => c.modules).find(m => m.id === moduleId);
      return <StubModule name={mod?.name ?? moduleId} />;
    }
  }
}

// ── Main App ──────────────────────────────────────────────────────────────────

export default function App() {
  const [nav, setNav]           = useState<NavLevel>({ level: 'categories' });
  const [cursorIdx, setCursorIdx] = useState(0);
  const [status, setStatus]     = useState<StatusState>({
    context: 'Lambda-env', message: 'Listo', type: 'normal',
    shortcuts: '↑↓ Navegar  Enter Abrir  ? Ayuda  q Salir',
  });
  const [modified, setModified] = useState<Set<string>>(new Set());
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<{ cat: CategoryDef; mod: ModuleDef }[]>([]);
  const [searchCursor, setSearchCursor] = useState(0);
  const [time, setTime]         = useState(new Date());
  const contentRef = useRef<HTMLDivElement>(null);
  const searchRef  = useRef<HTMLInputElement>(null);

  // Clock
  useEffect(() => {
    const id = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  // Search
  useEffect(() => {
    const results = searchModules(searchQuery);
    setSearchResults(results);
    setSearchCursor(0);
  }, [searchQuery]);

  // Status helpers
  function setStatusMsg(msg: string, type: StatusState['type'] = 'normal') {
    const { context, shortcuts } = getContextInfo();
    setStatus({ context, message: msg, type, shortcuts });
    if (type !== 'normal') {
      setTimeout(() => setStatus(s => ({ ...s, message: '', type: 'normal' })), 3000);
    }
  }

  function getContextInfo(): { context: string; shortcuts: string } {
    if (nav.level === 'categories') return { context: 'Lambda-env', shortcuts: '↑↓ Navegar  Enter Abrir  / Buscar  q Salir' };
    if (nav.level === 'modules') {
      const cat = CATEGORIES.find(c => c.id === nav.catId);
      return { context: `${cat?.name ?? ''}`, shortcuts: '↑↓ Navegar  Enter Abrir  Esc Volver  / Buscar' };
    }
    const cat = CATEGORIES.find(c => c.id === nav.catId);
    const mod = cat?.modules.find(m => m.id === nav.moduleId);
    const isModified = modified.has(nav.moduleId ?? '');
    return {
      context: `${cat?.name} › ${mod?.name}${isModified ? ' *' : ''}`,
      shortcuts: 'Esc Volver  / Buscar  F5 Guardar',
    };
  }

  // Navigate
  function goToModules(catId: string) {
    setNav({ level: 'modules', catId });
    setCursorIdx(0);
    const cat = CATEGORIES.find(c => c.id === catId);
    setStatus({ context: cat?.name ?? '', message: `${cat?.modules.length} módulos`, type: 'normal', shortcuts: '↑↓ Navegar  Enter Abrir  Esc Volver' });
  }

  function goToDetail(catId: string, moduleId: string) {
    setNav({ level: 'detail', catId, moduleId });
    setCursorIdx(0);
    const cat = CATEGORIES.find(c => c.id === catId);
    const mod = cat?.modules.find(m => m.id === moduleId);
    setStatus({ context: `${cat?.name} › ${mod?.name}`, message: '', type: 'normal', shortcuts: 'Esc Volver  / Buscar' });
    setTimeout(() => { contentRef.current?.scrollTo({ top: 0 }); }, 50);
  }

  function goBack() {
    if (nav.level === 'detail') { setNav({ level: 'modules', catId: nav.catId }); setCursorIdx(0); }
    else if (nav.level === 'modules') { setNav({ level: 'categories' }); setCursorIdx(0); }
  }

  function onModified() {
    if (nav.level === 'detail') setModified(m => new Set([...m, nav.moduleId]));
    setStatusMsg('Cambio aplicado', 'success');
  }

  // Current items for keyboard nav
  const currentItems: string[] = nav.level === 'categories'
    ? CATEGORIES.map(c => c.id)
    : nav.level === 'modules'
    ? (CATEGORIES.find(c => c.id === nav.catId)?.modules.map(m => m.id) ?? [])
    : [];

  // Keyboard handler
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const tag = (e.target as HTMLElement).tagName;
      const isSearch = (e.target as HTMLElement) === searchRef.current;
      if (!isSearch && (tag === 'INPUT' || tag === 'SELECT' || tag === 'TEXTAREA')) return;

      if (e.key === '/' && !searchOpen) {
        setSearchOpen(true);
        setTimeout(() => searchRef.current?.focus(), 50);
        e.preventDefault(); return;
      }

      if (searchOpen) {
        if (e.key === 'Escape') { setSearchOpen(false); setSearchQuery(''); e.preventDefault(); }
        if (e.key === 'ArrowDown') { setSearchCursor(i => Math.min(i + 1, searchResults.length - 1)); e.preventDefault(); }
        if (e.key === 'ArrowUp') { setSearchCursor(i => Math.max(i - 1, 0)); e.preventDefault(); }
        if (e.key === 'Enter' && searchResults[searchCursor]) {
          const { cat, mod } = searchResults[searchCursor];
          goToDetail(cat.id, mod.id);
          setSearchOpen(false); setSearchQuery('');
          e.preventDefault();
        }
        return;
      }

      if (e.key === 'Escape') { goBack(); e.preventDefault(); return; }
      if (e.key === 'ArrowUp' || (e.key === 'k' && !isSearch)) { setCursorIdx(i => Math.max(i - 1, 0)); e.preventDefault(); }
      if (e.key === 'ArrowDown' || (e.key === 'j' && !isSearch)) { setCursorIdx(i => Math.min(i + 1, currentItems.length - 1)); e.preventDefault(); }
      if (e.key === 'Enter' || e.key === 'l') {
        if (nav.level === 'categories') { goToModules(currentItems[cursorIdx]); e.preventDefault(); }
        else if (nav.level === 'modules') { goToDetail(nav.catId, currentItems[cursorIdx]); e.preventDefault(); }
      }
      if (e.key === 'h' && nav.level !== 'categories') { goBack(); e.preventDefault(); }
      if (e.key === 'q' && nav.level === 'categories') { /* exit */ e.preventDefault(); }
      if (e.key === 'F5') { setStatusMsg('Configuración recargada', 'success'); e.preventDefault(); }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [nav, cursorIdx, searchOpen, searchResults, searchCursor, currentItems]);

  // Status bar colors
  const statusBg = {
    normal: C.accent, success: C.success, error: C.error,
    warn: C.warn, loading: C.accent,
  }[status.type];

  const { context } = getContextInfo();

  return (
    <div className="crt-screen" style={{
      background: C.bg, minHeight: '100vh', display: 'flex', flexDirection: 'column',
      alignItems: 'center', justifyContent: 'center', padding: 16,
      position: 'relative', overflow: 'hidden',
    }}>
      {/* Vignette */}
      <div style={{ position: 'fixed', inset: 0, pointerEvents: 'none', zIndex: 9997, background: 'radial-gradient(ellipse at center, transparent 55%, rgba(0,0,0,0.65) 100%)' }} />

      {/* Search overlay */}
      {searchOpen && (
        <div style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.55)', display: 'flex', alignItems: 'flex-start', justifyContent: 'center', paddingTop: 80 }}>
          <div style={{ width: 480, background: C.surface, border: `1px solid ${C.accentBorder}`, boxShadow: `0 8px 40px rgba(0,0,0,0.7), 0 0 20px ${C.accent}20` }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, padding: '10px 14px', borderBottom: `1px solid ${C.border}` }}>
              <span style={{ color: C.accent, fontFamily: 'monospace', fontSize: 13 }}>/</span>
              <input
                ref={searchRef}
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder="Buscar módulo o configuración..."
                style={{ flex: 1, background: 'transparent', border: 'none', outline: 'none', color: C.textPrimary, fontFamily: 'monospace', fontSize: 13, caretColor: C.accent }}
              />
              <span style={{ color: C.textMuted, fontSize: 10, fontFamily: 'monospace' }}>Esc para cerrar</span>
            </div>
            {searchQuery && (
              <div>
                {searchResults.length === 0 ? (
                  <div style={{ padding: '14px 14px', color: C.textMuted, fontFamily: 'monospace', fontSize: 12 }}>Sin resultados para "{searchQuery}"</div>
                ) : searchResults.map(({ cat, mod }, i) => (
                  <div
                    key={mod.id}
                    onClick={() => { goToDetail(cat.id, mod.id); setSearchOpen(false); setSearchQuery(''); }}
                    style={{
                      padding: '9px 14px', borderBottom: `1px solid ${C.border}`,
                      background: i === searchCursor ? C.accentDim : 'transparent',
                      cursor: 'pointer',
                    }}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span style={{ color: C.accent, fontFamily: 'monospace', fontSize: 13 }}>{mod.icon}</span>
                      <span style={{ color: C.textPrimary, fontFamily: 'monospace', fontSize: 13 }}>{mod.name}</span>
                      <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10, marginLeft: 'auto' }}>{cat.name}</span>
                    </div>
                    <div style={{ color: C.textSecondary, fontFamily: 'monospace', fontSize: 11, marginTop: 2, marginLeft: 21 }}>{mod.description}</div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Main window */}
      <div style={{
        width: '100%', maxWidth: 900, display: 'flex', flexDirection: 'column',
        border: `1px solid ${C.border}`,
        boxShadow: `0 0 40px rgba(125,86,244,0.2), 0 0 80px rgba(125,86,244,0.08)`,
        background: C.bg,
      }}>
        {/* Header */}
        <div style={{
          height: 44, background: C.surface, borderBottom: `1px solid ${C.border}`,
          display: 'flex', alignItems: 'center', padding: '0 16px', gap: 12,
        }}>
          <span style={{ color: C.accent, fontSize: 16 }}>{ICONS.categories.system.nerd}</span>
          <span style={{ color: C.textPrimary, fontFamily: 'monospace', fontSize: 13, letterSpacing: '0.05em' }}>LambdaOS Settings</span>

          {/* Breadcrumb */}
          {nav.level !== 'categories' && (
            <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginLeft: 8 }}>
              <span style={{ color: C.textMuted, fontSize: 11 }}>—</span>
              <button
                onClick={() => setNav({ level: 'categories' })}
                style={{ background: 'none', border: 'none', color: C.textSecondary, fontFamily: 'monospace', fontSize: 11, cursor: 'pointer', padding: 0 }}
              >
                {CATEGORIES.find(c => c.id === (nav as any).catId)?.name}
              </button>
              {nav.level === 'detail' && (
                <>
                  <span style={{ color: C.textMuted, fontSize: 11 }}>›</span>
                  <span style={{ color: C.accent, fontFamily: 'monospace', fontSize: 11 }}>
                    {CATEGORIES.find(c => c.id === nav.catId)?.modules.find(m => m.id === nav.moduleId)?.name}
                    {modified.has(nav.moduleId) && <span style={{ color: C.error }}> *</span>}
                  </span>
                </>
              )}
            </div>
          )}

          <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 16 }}>
            <button
              onClick={() => { setSearchOpen(true); setTimeout(() => searchRef.current?.focus(), 50); }}
              style={{ background: 'none', border: `1px solid ${C.border}`, color: C.textMuted, fontFamily: 'monospace', fontSize: 10, cursor: 'pointer', padding: '3px 10px' }}
            >
              / Buscar
            </button>
            <span style={{ color: C.textSecondary, fontFamily: 'monospace', fontSize: 11 }}>
              {time.toLocaleTimeString('es-ES', { hour12: false })}
            </span>
          </div>
        </div>

        {/* Content */}
        <div style={{ display: 'flex', flex: 1 }}>

          {/* ── CATEGORIES VIEW ────────────────────────────────────────────── */}
          {nav.level === 'categories' && (
            <div style={{ flex: 1, padding: '8px 0' }}>
              {CATEGORIES.map((cat, i) => (
                <div
                  key={cat.id}
                  onClick={() => goToModules(cat.id)}
                  style={{
                    display: 'flex', alignItems: 'center', gap: 14,
                    padding: '14px 20px',
                    borderBottom: `1px solid ${C.border}`,
                    cursor: 'pointer',
                    background: i === cursorIdx ? C.accentDim : 'transparent',
                    borderLeft: i === cursorIdx ? `3px solid ${C.accent}` : '3px solid transparent',
                    transition: 'all 0.1s',
                  }}
                >
                  <span style={{ color: C.accent, fontSize: 22, width: 28, textAlign: 'center' }}>{cat.icon}</span>
                  <div style={{ flex: 1 }}>
                    <div style={{ color: i === cursorIdx ? C.textPrimary : C.textPrimary, fontFamily: 'monospace', fontSize: 14 }}>
                      {i === cursorIdx ? '▶ ' : '  '}{cat.name}
                    </div>
                    <div style={{ color: C.textSecondary, fontFamily: 'monospace', fontSize: 11, marginTop: 3 }}>{cat.description}</div>
                  </div>
                  <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 11 }}>{cat.modules.length} módulos</span>
                  <span style={{ color: C.accent, fontFamily: 'monospace', fontSize: 12 }}>▶</span>
                </div>
              ))}
            </div>
          )}

          {/* ── MODULES VIEW ───────────────────────────────────────────────── */}
          {nav.level === 'modules' && (() => {
            const cat = CATEGORIES.find(c => c.id === nav.catId)!;
            return (
              <div style={{ flex: 1, padding: '8px 0' }}>
                {cat.modules.map((mod, i) => (
                  <div
                    key={mod.id}
                    onClick={() => goToDetail(cat.id, mod.id)}
                    style={{
                      display: 'flex', alignItems: 'center', gap: 14,
                      padding: '12px 20px',
                      borderBottom: `1px solid ${C.border}`,
                      cursor: 'pointer',
                      background: i === cursorIdx ? C.accentDim : 'transparent',
                      borderLeft: i === cursorIdx ? `3px solid ${C.accent}` : '3px solid transparent',
                      transition: 'all 0.1s',
                      opacity: mod.implemented === false ? 0.5 : 1,
                    }}
                  >
                    <span style={{ color: i === cursorIdx ? C.accent : C.textSecondary, fontSize: 16, width: 20, textAlign: 'center', fontFamily: 'monospace' }}>{mod.icon}</span>
                    <div style={{ flex: 1 }}>
                      <div style={{ fontFamily: 'monospace', fontSize: 13, display: 'flex', alignItems: 'center', gap: 8 }}>
                        <span style={{ color: i === cursorIdx ? C.textPrimary : C.textPrimary }}>
                          {i === cursorIdx ? '▶ ' : '  '}{mod.name}
                        </span>
                        {mod.rootRequired && <span style={{ color: C.error, fontSize: 10 }}>{ICONS.widgets.lock.nerd}</span>}
                        {modified.has(mod.id) && <span style={{ color: C.error, fontSize: 10 }}>*</span>}
                        {!mod.implemented && <span style={{ color: C.textMuted, fontSize: 10 }}>· próximamente</span>}
                      </div>
                      <div style={{ color: C.textSecondary, fontFamily: 'monospace', fontSize: 11, marginTop: 2 }}>{mod.description}</div>
                    </div>
                    <span style={{ color: C.accent, fontFamily: 'monospace', fontSize: 12 }}>▶</span>
                  </div>
                ))}
              </div>
            );
          })()}

          {/* ── DETAIL VIEW ────────────────────────────────────────────────── */}
          {nav.level === 'detail' && (
            <div ref={contentRef} style={{ flex: 1, overflowY: 'auto', maxHeight: 'calc(100vh - 160px)' }}
              className="tui-scrollbar">
              {renderModule(
                nav.moduleId,
                onModified,
                (msg) => setStatusMsg(msg, 'success'),
              )}
            </div>
          )}
        </div>

        {/* Status bar */}
        <div style={{
          height: 30, background: statusBg, borderTop: `1px solid ${C.border}`,
          display: 'flex', alignItems: 'center', padding: '0 14px', gap: 12,
          transition: 'background 0.3s',
        }}>
          <span style={{ color: '#000', fontFamily: 'monospace', fontSize: 11, fontWeight: 'bold' }}>
            {context}
          </span>
          {status.message && (
            <>
              <span style={{ color: '#00000050', fontSize: 10 }}>·</span>
              <span style={{ color: '#000', fontFamily: 'monospace', fontSize: 11 }}>{status.message}</span>
            </>
          )}
          <div style={{ marginLeft: 'auto', display: 'flex', gap: 16 }}>
            <span style={{ color: '#00000070', fontFamily: 'monospace', fontSize: 10 }}>
              {getContextInfo().shortcuts}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
