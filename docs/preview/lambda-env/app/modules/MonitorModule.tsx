import { useEffect, useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUIProgress } from '../components/tui/TUIProgress';

interface Process { pid: number; name: string; cpu: number; mem: number; state: string; }

function rand(min: number, max: number) { return Math.random() * (max - min) + min; }

const PROCS: Process[] = [
  { pid: 1,    name: 'systemd',        cpu: 0.1, mem: 12.4,  state: 'S' },
  { pid: 847,  name: 'Xorg',          cpu: 2.4, mem: 180.2, state: 'S' },
  { pid: 1204, name: 'pipewire',      cpu: 0.8, mem: 24.1,  state: 'S' },
  { pid: 2341, name: 'qtile',         cpu: 1.2, mem: 96.3,  state: 'S' },
  { pid: 3892, name: 'firefox',       cpu: 8.7, mem: 512.4, state: 'S' },
  { pid: 4012, name: 'discord',       cpu: 3.2, mem: 248.1, state: 'S' },
  { pid: 5001, name: 'spotify',       cpu: 1.9, mem: 184.7, state: 'S' },
  { pid: 6123, name: 'alacritty',     cpu: 0.4, mem: 48.2,  state: 'S' },
  { pid: 7891, name: 'lambda-env',    cpu: 0.3, mem: 32.1,  state: 'R' },
  { pid: 8234, name: 'NetworkManager',cpu: 0.2, mem: 18.4,  state: 'S' },
];

export function MonitorModule() {
  const [cpu, setCpu]       = useState(12);
  const [mem, setMem]       = useState(4200);
  const [swap, setSwap]     = useState(0);
  const [net, setNet]       = useState({ down: 2.4, up: 0.8 });
  const [procs, setProcs]   = useState(PROCS);
  const [tick, setTick]     = useState(0);

  useEffect(() => {
    const id = setInterval(() => {
      setCpu(rand(5, 35));
      setMem(Math.floor(rand(3800, 5000)));
      setNet({ down: rand(0.1, 15), up: rand(0.05, 5) });
      setTick(t => t + 1);
    }, 1500);
    return () => clearInterval(id);
  }, []);

  const cpuColor = cpu > 80 ? C.error : cpu > 50 ? C.warn : C.accent;
  const memColor = mem > 12000 ? C.error : mem > 8000 ? C.warn : C.accent;

  return (
    <div>
      {/* Resources */}
      <TUISection title="RECURSOS DEL SISTEMA">
        <TUIProgress value={Math.round(cpu)} label={`CPU — ${Math.round(cpu)}%`} color={cpuColor} subLabel={`Intel i7-12700K · 12C/20T · Temp: ${Math.floor(rand(45,65))}°C`} />
        <TUIProgress value={Math.round((mem / 16384) * 100)} label={`RAM — ${mem}M / 16384M`} color={memColor} subLabel={`${(mem / 1024).toFixed(1)} GB usados · ${((16384 - mem) / 1024).toFixed(1)} GB libres`} />
        <TUIProgress value={Math.round((swap / 8192) * 100)} label={`SWAP — ${swap}M / 8192M`} color={C.accent} subLabel="Sin uso de swap" />
        <TUIProgress value={28} label="Disco / — 280 GB / 1000 GB" color={C.accent} subLabel="NVMe Samsung 980 · Lectura: 3.2 GB/s · Escritura: 2.1 GB/s" />
        <TUIProgress value={8} label="GPU — 8%" color={C.accent} subLabel={`NVIDIA RTX 3070 · VRAM: 1.2 GB / 8 GB · Temp: ${Math.floor(rand(48,72))}°C`} />
      </TUISection>

      {/* Network stats */}
      <div style={{
        padding: '10px 14px', background: C.surface, margin: '4px 0',
        fontFamily: 'monospace', display: 'grid', gridTemplateColumns: '1fr 1fr 1fr 1fr', gap: 16,
      }}>
        <div>
          <div style={{ color: C.textMuted, fontSize: 10 }}>DESCARGA</div>
          <div style={{ color: C.success, fontSize: 16, marginTop: 2 }}>↓ {net.down.toFixed(1)} MB/s</div>
        </div>
        <div>
          <div style={{ color: C.textMuted, fontSize: 10 }}>SUBIDA</div>
          <div style={{ color: C.accent, fontSize: 16, marginTop: 2 }}>↑ {net.up.toFixed(1)} MB/s</div>
        </div>
        <div>
          <div style={{ color: C.textMuted, fontSize: 10 }}>LATENCIA</div>
          <div style={{ color: C.textPrimary, fontSize: 16, marginTop: 2 }}>{Math.floor(rand(5, 30))} ms</div>
        </div>
        <div>
          <div style={{ color: C.textMuted, fontSize: 10 }}>UPTIME</div>
          <div style={{ color: C.textPrimary, fontSize: 14, marginTop: 2 }}>7d 14h {tick}m</div>
        </div>
      </div>

      {/* Process list */}
      <TUISection title="PROCESOS (TOP 10 por CPU)">
        <div style={{
          display: 'grid', gridTemplateColumns: '55px 1fr 70px 80px 30px',
          padding: '5px 12px', background: C.surface, gap: 8,
          borderBottom: `1px solid ${C.border}`, fontFamily: 'monospace',
        }}>
          {['PID','Nombre','CPU','MEM','ST'].map(h => (
            <span key={h} style={{ color: C.textMuted, fontSize: 10 }}>{h}</span>
          ))}
        </div>
        {procs.sort((a,b) => b.cpu - a.cpu).map(p => (
          <div key={p.pid} style={{
            display: 'grid', gridTemplateColumns: '55px 1fr 70px 80px 30px',
            padding: '5px 12px', borderBottom: `1px solid ${C.border}`,
            gap: 8, fontFamily: 'monospace', alignItems: 'center',
          }}>
            <span style={{ color: C.textMuted, fontSize: 11 }}>{p.pid}</span>
            <span style={{ color: p.name === 'lambda-env' ? C.accent : C.textPrimary, fontSize: 12 }}>{p.name}</span>
            <div style={{ position: 'relative', height: 3, background: C.surface2 }}>
              <div style={{ position: 'absolute', left: 0, top: 0, height: '100%', width: `${Math.min(p.cpu * 10, 100)}%`, background: p.cpu > 8 ? C.warn : C.accent }} />
            </div>
            <span style={{ color: C.textSecondary, fontSize: 11 }}>{p.mem.toFixed(1)} MB</span>
            <span style={{ color: p.state === 'R' ? C.success : C.textMuted, fontSize: 11 }}>{p.state}</span>
          </div>
        ))}
      </TUISection>
    </div>
  );
}
