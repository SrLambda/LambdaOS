// Mirrors: src/lambda-env/internal/modules/logs/main.go
// Actions: run, filter by unit/priority/time, export
import { useState, useRef, useEffect } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIInput } from '../components/tui/TUIInput';

// ── Types ─────────────────────────────────────────────────────────────────────

type Priority = 'emerg' | 'alert' | 'crit' | 'err' | 'warning' | 'notice' | 'info' | 'debug';

interface LogEntry {
  id: string;
  ts: string;
  unit: string;
  priority: Priority;
  message: string;
}

// ── Data ──────────────────────────────────────────────────────────────────────

const UNITS = ['(todos)', 'kernel', 'systemd', 'NetworkManager', 'sshd', 'bluetooth', 'pipewire', 'qtile', 'nvim'];
const PRIORITIES: Priority[] = ['emerg', 'alert', 'crit', 'err', 'warning', 'notice', 'info', 'debug'];
const PRIORITY_FILTER_OPTS = ['(todas)', ...PRIORITIES];
const TIME_RANGES = ['última hora', 'últimas 6h', 'hoy', 'ayer', 'esta semana', 'todo'];

const SAMPLE_LOGS: LogEntry[] = [
  { id: 'l01', ts: '2026-06-07 09:43:21', unit: 'kernel',         priority: 'info',    message: 'Linux version 6.6.32-arch1 (builduser@buildhost) (gcc 13.2.1)' },
  { id: 'l02', ts: '2026-06-07 09:43:22', unit: 'systemd',        priority: 'info',    message: 'systemd 255 running in system mode (+PAM +AUDIT +SELINUX)' },
  { id: 'l03', ts: '2026-06-07 09:43:25', unit: 'NetworkManager', priority: 'notice',  message: 'NetworkManager (version 1.46.0) is starting...' },
  { id: 'l04', ts: '2026-06-07 09:43:26', unit: 'NetworkManager', priority: 'info',    message: 'device (wlan0): driver supports Access Point (AP) mode' },
  { id: 'l05', ts: '2026-06-07 09:43:28', unit: 'pipewire',       priority: 'info',    message: 'module-rt: RTKit not available: org.freedesktop.DBus.Error.ServiceUnknown' },
  { id: 'l06', ts: '2026-06-07 09:44:01', unit: 'sshd',           priority: 'notice',  message: 'Server listening on 0.0.0.0 port 22' },
  { id: 'l07', ts: '2026-06-07 09:44:10', unit: 'sshd',           priority: 'warning', message: 'Failed password for invalid user admin from 192.168.1.100 port 52814 ssh2' },
  { id: 'l08', ts: '2026-06-07 09:44:12', unit: 'sshd',           priority: 'warning', message: 'Failed password for invalid user root from 192.168.1.100 port 52815 ssh2' },
  { id: 'l09', ts: '2026-06-07 09:44:14', unit: 'sshd',           priority: 'err',     message: 'error: maximum authentication attempts exceeded for invalid user root from 192.168.1.100' },
  { id: 'l10', ts: '2026-06-07 09:44:15', unit: 'kernel',         priority: 'warning', message: 'audit: type=1130 audit(1717749855.123:45): pid=1 uid=0 auid=4294967295 ses=4294967295' },
  { id: 'l11', ts: '2026-06-07 09:45:00', unit: 'bluetooth',      priority: 'info',    message: 'Bluetooth: hci0: unexpected event for opcode 0x2042' },
  { id: 'l12', ts: '2026-06-07 09:46:30', unit: 'NetworkManager', priority: 'notice',  message: 'device (wlan0): Activation: successful, device is now active' },
  { id: 'l13', ts: '2026-06-07 10:00:00', unit: 'systemd',        priority: 'notice',  message: 'Starting Daily man-db regeneration...' },
  { id: 'l14', ts: '2026-06-07 10:00:01', unit: 'systemd',        priority: 'notice',  message: 'Finished Daily man-db regeneration.' },
  { id: 'l15', ts: '2026-06-07 10:12:44', unit: 'kernel',         priority: 'crit',    message: 'EDAC MC0: 1 CE memory read error on CPU_SrcID#0_Ha#0_Chan#0_DIMM#0' },
  { id: 'l16', ts: '2026-06-07 10:14:05', unit: 'pipewire',       priority: 'info',    message: 'spa.alsa: opened hw:0 capture node 43' },
  { id: 'l17', ts: '2026-06-07 10:22:10', unit: 'qtile',          priority: 'info',    message: 'Qtile started successfully' },
  { id: 'l18', ts: '2026-06-07 10:22:15', unit: 'qtile',          priority: 'warning', message: 'Unhandled exception in window event hook' },
];

const PRIO_COLOR: Record<Priority, string> = {
  emerg:   C.error,
  alert:   C.error,
  crit:    C.error,
  err:     '#FF6B6B',
  warning: C.warn,
  notice:  C.accent,
  info:    C.textSecondary,
  debug:   C.textMuted,
};

const PRIO_LABEL: Record<Priority, string> = {
  emerg:   'EMERG',
  alert:   'ALERT',
  crit:    ' CRIT',
  err:     '  ERR',
  warning: ' WARN',
  notice:  ' NOTE',
  info:    ' INFO',
  debug:   'DEBUG',
};

// ── Component ─────────────────────────────────────────────────────────────────

export function LogsModule({ onModified }: { onModified: () => void }) {
  const [unit, setUnit]       = useState('(todos)');
  const [prio, setPrio]       = useState('(todas)');
  const [timeRange, setTime]  = useState('hoy');
  const [search, setSearch]   = useState('');
  const [follow, setFollow]   = useState(false);
  const [logs, setLogs]       = useState<LogEntry[]>(SAMPLE_LOGS);
  const bottomRef = useRef<HTMLDivElement>(null);

  const filtered = logs.filter(l => {
    if (unit !== '(todos)' && l.unit !== unit) return false;
    if (prio !== '(todas)' && l.priority !== prio) return false;
    if (search && !l.message.toLowerCase().includes(search.toLowerCase()) &&
        !l.unit.toLowerCase().includes(search.toLowerCase())) return false;
    return true;
  });

  // Simulate new log lines when following
  useEffect(() => {
    if (!follow) return;
    const interval = setInterval(() => {
      const units = ['kernel', 'systemd', 'NetworkManager', 'pipewire'];
      const msgs = [
        'heartbeat check ok',
        'scheduled task completed',
        'connection keepalive sent',
        'buffer flushed (0 bytes pending)',
      ];
      const entry: LogEntry = {
        id: 'live-' + Date.now(),
        ts: new Date().toLocaleString('es-ES', { hour12: false }).replace(',', ''),
        unit: units[Math.floor(Math.random() * units.length)],
        priority: Math.random() > 0.85 ? 'warning' : 'info',
        message: msgs[Math.floor(Math.random() * msgs.length)],
      };
      setLogs(prev => [...prev.slice(-200), entry]);
    }, 2000);
    return () => clearInterval(interval);
  }, [follow]);

  useEffect(() => {
    if (follow && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, follow]);

  const critCount = filtered.filter(l => ['emerg','alert','crit','err'].includes(l.priority)).length;
  const warnCount = filtered.filter(l => l.priority === 'warning').length;

  return (
    <div>
      {/* Stats bar */}
      <div style={{
        padding: '10px 14px', background: C.surface, marginBottom: 4,
        fontFamily: 'monospace', display: 'flex', gap: 24, alignItems: 'center',
        borderLeft: `3px solid ${critCount > 0 ? C.error : warnCount > 0 ? C.warn : C.success}`,
      }}>
        <div>
          <div style={{ color: C.textPrimary, fontSize: 12 }}>{filtered.length} entradas</div>
          <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>
            journalctl · {timeRange}
          </div>
        </div>
        <div style={{ display: 'flex', gap: 20, marginLeft: 'auto' }}>
          {critCount > 0 && <span style={{ color: C.error,   fontSize: 11 }}>✗ {critCount} errores</span>}
          {warnCount > 0 && <span style={{ color: C.warn,    fontSize: 11 }}>⚠ {warnCount} avisos</span>}
          <span style={{ color: C.success, fontSize: 11 }}>● {follow ? 'siguiendo' : 'estático'}</span>
        </div>
      </div>

      {/* Filters */}
      <TUISection title="FILTROS">
        <TUISelect
          label="Unidad"
          description="Servicio systemd o componente de kernel"
          value={unit}
          options={UNITS}
          onChange={setUnit}
        />
        <TUISelect
          label="Prioridad mínima"
          description="Nivel de severidad (journald RFC 5424)"
          value={prio}
          options={PRIORITY_FILTER_OPTS}
          onChange={setPrio}
        />
        <TUISelect
          label="Rango de tiempo"
          description="Ventana temporal de los logs"
          value={timeRange}
          options={TIME_RANGES}
          onChange={v => { setTime(v); onModified(); }}
        />
        <TUIInput
          label="Buscar"
          description="Filtrar por texto en mensaje o unidad"
          value={search}
          onChange={setSearch}
        />
      </TUISection>

      {/* Log output */}
      <TUISection title={`JOURNAL (${filtered.length})`}>
        {/* Toolbar */}
        <div style={{
          padding: '8px 12px', display: 'flex', gap: 8, alignItems: 'center',
          borderBottom: `1px solid ${C.border}`,
        }}>
          <TUIButton
            label={follow ? '■ Detener follow' : '▶ Seguir log'}
            onClick={async () => setFollow(f => !f)}
            variant={follow ? 'danger' : 'primary'}
          />
          <TUIButton
            label="Limpiar filtros"
            onClick={async () => { setUnit('(todos)'); setPrio('(todas)'); setSearch(''); }}
          />
          <TUIButton
            label="Exportar"
            onClick={async () => { onModified(); }}
            icon="↑"
          />
          <span style={{ marginLeft: 'auto', color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>
            {follow && <span style={{ color: C.success }}>● LIVE  </span>}
            {filtered.length} líneas
          </span>
        </div>

        {/* Log lines */}
        <div
          className="tui-scrollbar"
          style={{
            height: 360, overflowY: 'auto',
            background: C.bg, fontFamily: 'monospace', fontSize: 11,
          }}
        >
          {filtered.length === 0 ? (
            <div style={{ padding: '20px 14px', color: C.textMuted }}>
              Sin resultados para los filtros actuales.
            </div>
          ) : filtered.map(l => (
            <div
              key={l.id}
              style={{
                display: 'flex', gap: 8, padding: '3px 12px',
                borderBottom: `1px solid ${C.border}22`,
                background: ['emerg','alert','crit','err'].includes(l.priority)
                  ? `${C.error}08` : 'transparent',
              }}
            >
              <span style={{ color: C.textMuted, flexShrink: 0, fontSize: 10 }}>{l.ts}</span>
              <span style={{
                color: PRIO_COLOR[l.priority], flexShrink: 0, width: 40, textAlign: 'right',
              }}>
                {PRIO_LABEL[l.priority]}
              </span>
              <span style={{ color: C.accent, flexShrink: 0, width: 90, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {l.unit}
              </span>
              <span style={{ color: C.textSecondary, flex: 1 }}>{l.message}</span>
            </div>
          ))}
          <div ref={bottomRef} />
        </div>
      </TUISection>

      {/* Priority legend */}
      <TUISection title="LEYENDA" collapsible defaultOpen={false}>
        <div style={{ padding: '10px 14px', display: 'flex', flexWrap: 'wrap', gap: 16, fontFamily: 'monospace' }}>
          {PRIORITIES.map(p => (
            <span key={p} style={{ color: PRIO_COLOR[p], fontSize: 11 }}>
              {PRIO_LABEL[p].trim()} — {p}
            </span>
          ))}
        </div>
      </TUISection>
    </div>
  );
}
