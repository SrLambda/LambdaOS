// Mirrors: src/lambda-env/internal/modules/almacenamiento/main.go
// Actions: run, mount, unmount, list-partitions, disk-info
import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUIButton } from '../components/tui/TUIButton';
import { TUIProgress } from '../components/tui/TUIProgress';
import { TUIModal } from '../components/tui/TUIModal';

// ── Types ─────────────────────────────────────────────────────────────────────

interface Disk {
  id: string;
  model: string;
  device: string;
  size: string;
  sizeBytes: number;
  type: 'SSD' | 'HDD' | 'NVMe' | 'USB';
  health: 'OK' | 'WARN' | 'FAIL';
  temp?: number;
}

interface Partition {
  id: string;
  diskId: string;
  device: string;
  label: string;
  fstype: string;
  size: string;
  used: string;
  usedPct: number;
  mountpoint: string | null;
  uuid: string;
}

// ── Data ──────────────────────────────────────────────────────────────────────

const DISKS: Disk[] = [
  { id: 'nvme0', model: 'Samsung 980 Pro',    device: '/dev/nvme0n1', size: '1 TB',   sizeBytes: 1000, type: 'NVMe', health: 'OK',   temp: 38 },
  { id: 'sda',   model: 'WD Blue 2TB',        device: '/dev/sda',     size: '2 TB',   sizeBytes: 2000, type: 'HDD',  health: 'OK',   temp: 31 },
  { id: 'sdb',   model: 'SanDisk Cruzer USB', device: '/dev/sdb',     size: '32 GB',  sizeBytes: 32,   type: 'USB',  health: 'OK' },
];

const INITIAL_PARTITIONS: Partition[] = [
  { id: 'p1', diskId: 'nvme0', device: '/dev/nvme0n1p1', label: 'EFI',      fstype: 'vfat',  size: '512 MB', used: '35 MB',  usedPct: 7,  mountpoint: '/boot/efi',  uuid: 'A1B2-C3D4' },
  { id: 'p2', diskId: 'nvme0', device: '/dev/nvme0n1p2', label: 'root',     fstype: 'ext4',  size: '200 GB', used: '68 GB',  usedPct: 34, mountpoint: '/',          uuid: 'a1b2c3d4-...' },
  { id: 'p3', diskId: 'nvme0', device: '/dev/nvme0n1p3', label: 'home',     fstype: 'btrfs', size: '780 GB', used: '210 GB', usedPct: 27, mountpoint: '/home',      uuid: 'b2c3d4e5-...' },
  { id: 'p4', diskId: 'sda',   device: '/dev/sda1',      label: 'data',     fstype: 'ext4',  size: '1.8 TB', used: '900 GB', usedPct: 50, mountpoint: '/mnt/data',  uuid: 'c3d4e5f6-...' },
  { id: 'p5', diskId: 'sda',   device: '/dev/sda2',      label: 'backup',   fstype: 'ext4',  size: '200 GB', used: '120 GB', usedPct: 60, mountpoint: null,         uuid: 'd4e5f6g7-...' },
  { id: 'p6', diskId: 'sdb',   device: '/dev/sdb1',      label: 'USB_DATA', fstype: 'exfat', size: '32 GB',  used: '8 GB',   usedPct: 25, mountpoint: null,         uuid: 'E5F6-G7H8' },
];

const typeColors: Record<string, string> = {
  NVMe: C.accent, SSD: C.success, HDD: C.textSecondary, USB: C.warn,
};
const healthColor = (h: Disk['health']) =>
  h === 'OK' ? C.success : h === 'WARN' ? C.warn : C.error;

// ── Component ─────────────────────────────────────────────────────────────────

export function StorageModule({ onModified }: { onModified: () => void }) {
  const [partitions, setPartitions] = useState<Partition[]>(INITIAL_PARTITIONS);
  const [selectedDisk, setSelectedDisk] = useState<string | null>(null);
  const [confirmUnmount, setConfirmUnmount] = useState<Partition | null>(null);
  const [mounting, setMounting] = useState<string | null>(null);

  const filteredPartitions = selectedDisk
    ? partitions.filter(p => p.diskId === selectedDisk)
    : partitions;

  async function mount(p: Partition) {
    setMounting(p.id);
    await new Promise(r => setTimeout(r, 900));
    setPartitions(ps => ps.map(x => x.id === p.id
      ? { ...x, mountpoint: `/mnt/${x.label.toLowerCase()}` }
      : x));
    setMounting(null);
    onModified();
  }

  async function unmount(p: Partition) {
    setConfirmUnmount(null);
    setMounting(p.id);
    await new Promise(r => setTimeout(r, 700));
    setPartitions(ps => ps.map(x => x.id === p.id ? { ...x, mountpoint: null } : x));
    setMounting(null);
    onModified();
  }

  const totalUsedPct = Math.round(
    partitions.reduce((s, p) => s + p.usedPct, 0) / partitions.length
  );

  return (
    <div>
      {confirmUnmount && (
        <TUIModal
          title={`Desmontar ${confirmUnmount.device}`}
          description={`¿Desmontar ${confirmUnmount.label} de ${confirmUnmount.mountpoint}?\nLas operaciones de E/S pendientes se completarán antes.`}
          confirmLabel="Desmontar"
          cancelLabel="Cancelar"
          variant="danger"
          onConfirm={() => unmount(confirmUnmount)}
          onCancel={() => setConfirmUnmount(null)}
        />
      )}

      {/* Summary bar */}
      <div style={{
        padding: '10px 14px', background: C.surface, marginBottom: 4,
        fontFamily: 'monospace', display: 'flex', gap: 24, alignItems: 'center',
        borderLeft: `3px solid ${C.accent}`,
      }}>
        <div>
          <div style={{ color: C.textPrimary, fontSize: 12 }}>{DISKS.length} dispositivos</div>
          <div style={{ color: C.textMuted, fontSize: 10, marginTop: 2 }}>{partitions.length} particiones · uso medio {totalUsedPct}%</div>
        </div>
        <div style={{ display: 'flex', gap: 20, marginLeft: 'auto' }}>
          {DISKS.map(d => (
            <span key={d.id} style={{ color: typeColors[d.type] ?? C.textMuted, fontSize: 11 }}>
              {d.type} {d.size}
            </span>
          ))}
        </div>
      </div>

      {/* Disk cards */}
      <TUISection title="DISPOSITIVOS">
        <div style={{ display: 'flex', gap: 0, flexWrap: 'wrap' }}>
          {DISKS.map(disk => (
            <div
              key={disk.id}
              onClick={() => setSelectedDisk(selectedDisk === disk.id ? null : disk.id)}
              style={{
                flex: '1 1 200px', padding: '12px 14px', cursor: 'pointer',
                borderRight: `1px solid ${C.border}`,
                borderBottom: `1px solid ${C.border}`,
                background: selectedDisk === disk.id ? C.accentDim : 'transparent',
                fontFamily: 'monospace',
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
                <span style={{ color: typeColors[disk.type] ?? C.textMuted, fontSize: 10 }}>
                  [{disk.type}]
                </span>
                <span style={{ color: healthColor(disk.health), fontSize: 10 }}>
                  ● {disk.health}
                </span>
              </div>
              <div style={{ color: C.textPrimary, fontSize: 12, marginBottom: 2 }}>{disk.device}</div>
              <div style={{ color: C.textSecondary, fontSize: 10, marginBottom: 6 }}>{disk.model}</div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: C.textMuted, fontSize: 10 }}>{disk.size}</span>
                {disk.temp !== undefined && (
                  <span style={{ color: disk.temp > 50 ? C.warn : C.textMuted, fontSize: 10 }}>
                    {disk.temp}°C
                  </span>
                )}
              </div>
            </div>
          ))}
        </div>
        {selectedDisk && (
          <div style={{ padding: '6px 12px', borderTop: `1px solid ${C.border}` }}>
            <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>
              Filtrando por: {DISKS.find(d => d.id === selectedDisk)?.device} —{' '}
            </span>
            <span
              style={{ color: C.accent, fontFamily: 'monospace', fontSize: 10, cursor: 'pointer' }}
              onClick={() => setSelectedDisk(null)}
            >
              [mostrar todos]
            </span>
          </div>
        )}
      </TUISection>

      {/* Partition list */}
      <TUISection title={`PARTICIONES${selectedDisk ? ` — ${DISKS.find(d => d.id === selectedDisk)?.device}` : ''}`}>
        {/* Header */}
        <div style={{
          display: 'grid', gridTemplateColumns: '140px 60px 60px 80px 100px 1fr auto',
          padding: '5px 12px', background: C.surface, gap: 8,
          borderBottom: `1px solid ${C.border}`, fontFamily: 'monospace',
        }}>
          {['Dispositivo', 'FS', 'Tamaño', 'Usado', 'Uso', 'Montaje', ''].map(h => (
            <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
          ))}
        </div>

        {filteredPartitions.map(p => (
          <div key={p.id} style={{
            display: 'grid', gridTemplateColumns: '140px 60px 60px 80px 100px 1fr auto',
            padding: '8px 12px', borderBottom: `1px solid ${C.border}`,
            gap: 8, fontFamily: 'monospace', alignItems: 'center',
          }}>
            <div>
              <div style={{ color: C.textPrimary, fontSize: 11 }}>{p.device}</div>
              <div style={{ color: C.textMuted, fontSize: 9, marginTop: 1 }}>{p.label}</div>
            </div>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{p.fstype}</span>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{p.size}</span>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{p.used}</span>
            <div>
              <TUIProgress
                value={p.usedPct}
                color={p.usedPct > 85 ? C.error : p.usedPct > 65 ? C.warn : C.success}
                subLabel={`${p.usedPct}%`}
              />
            </div>
            <span style={{
              color: p.mountpoint ? C.success : C.textMuted,
              fontSize: 10,
              overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
            }}>
              {p.mountpoint ?? '—'}
            </span>
            <div style={{ display: 'flex', gap: 6 }}>
              {mounting === p.id ? (
                <span style={{ color: C.textMuted, fontFamily: 'monospace', fontSize: 10 }}>⟳</span>
              ) : p.mountpoint ? (
                <TUIButton
                  label="Desmontar"
                  onClick={async () => setConfirmUnmount(p)}
                  variant="danger"
                />
              ) : (
                <TUIButton
                  label="Montar"
                  onClick={async () => mount(p)}
                  variant="primary"
                />
              )}
            </div>
          </div>
        ))}
      </TUISection>

      {/* Disk usage overview */}
      <TUISection title="USO POR PARTICIÓN" collapsible defaultOpen={false}>
        <div style={{ padding: '12px 14px', display: 'flex', flexDirection: 'column', gap: 10 }}>
          {partitions.filter(p => p.mountpoint).map(p => (
            <div key={p.id}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4, fontFamily: 'monospace' }}>
                <span style={{ color: C.textPrimary, fontSize: 11 }}>{p.mountpoint}</span>
                <span style={{ color: C.textMuted, fontSize: 10 }}>{p.used} / {p.size}</span>
              </div>
              <TUIProgress
                value={p.usedPct}
                color={p.usedPct > 85 ? C.error : p.usedPct > 65 ? C.warn : C.accent}
                label=""
                subLabel={`${p.usedPct}%`}
              />
            </div>
          ))}
        </div>
      </TUISection>

      {/* Actions */}
      <TUISection title="ACCIONES" collapsible defaultOpen={false} rootRequired>
        <div style={{ padding: '10px 12px', display: 'flex', gap: 8, flexWrap: 'wrap' }}>
          <TUIButton label="Escanear dispositivos" onClick={async () => { onModified(); }} icon="↻" />
          <TUIButton label="Abrir GParted"         onClick={async () => { onModified(); }} icon="✦" />
          <TUIButton label="SMART report"           onClick={async () => { onModified(); }} />
        </div>
      </TUISection>
    </div>
  );
}
