import { useEffect, useState } from 'react';
import { TUISettingItem } from '../TUISettingItem';

const LANGS = ['Español (es_ES)', 'English (en_US)', 'Français (fr_FR)', 'Deutsch (de_DE)', 'Italiano (it_IT)', '日本語 (ja_JP)', '中文 (zh_CN)', 'Português (pt_BR)'];
const TIMEZONES = ['UTC-08:00 (Los Angeles)', 'UTC-05:00 (Nueva York)', 'UTC+00:00 (Londres/UTC)', 'UTC+01:00 (Madrid/París)', 'UTC+02:00 (Helsinki)', 'UTC+05:30 (Mumbai)', 'UTC+09:00 (Tokio)'];
const INIT_SYSTEMS = ['systemd', 'OpenRC', 'runit', 'SysV Init'];
const SCHEDULERS = ['CFQ', 'Deadline', 'NOOP', 'BFQ', 'MQ-Deadline'];

const BOOT_LOG = [
  '[  0.000000] Initializing cgroup subsys cpuset',
  '[  0.000000] Linux version 6.5.0-terminus (gcc 13.2)',
  '[  0.000000] BIOS-provided physical RAM map',
  '[  0.183241] PCI: Using configuration type 1 for base access',
  '[  0.421884] NET: Registered PF_INET6 protocol family',
  '[  1.204571] systemd[1]: Detected architecture x86-64',
  '[  1.871239] Started Network Service.',
  '[  2.441028] terminus: login: root@TERMINUS-WORKSTATION-01',
];

function ProgressBar({ label, value, max, color = '#6D40FF' }: { label: string; value: number; max: number; color?: string }) {
  const pct = (value / max) * 100;
  return (
    <div className="flex items-center gap-2 text-[10px]">
      <span className="text-[#6D40FF]/50 w-14 shrink-0">{label}</span>
      <div className="flex-1 border border-[#6D40FF]/30 h-3 relative bg-black">
        <div className="absolute left-0 top-0 h-full" style={{ width: `${pct}%`, background: color }} />
      </div>
      <span className="text-[#6D40FF]/70 w-16 text-right">{value}/{max}</span>
    </div>
  );
}

export function SystemSettings() {
  const [logIdx, setLogIdx] = useState(0);
  useEffect(() => {
    const id = setInterval(() => setLogIdx(i => (i + 1) % BOOT_LOG.length), 2000);
    return () => clearInterval(id);
  }, []);

  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ SISTEMA GENERAL ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">systemd / uname</span>
      </div>

      {/* System info panel */}
      <div className="border border-[#6D40FF]/30 p-2 mb-3 grid grid-cols-2 gap-x-4 text-[10px]">
        <div className="space-y-0.5">
          {[
            ['OS', 'LambdaOS v2.4.1'],
            ['Kernel', 'Linux 6.5.0-terminus'],
            ['Arch', 'x86_64 (64-bit)'],
            ['Init', 'systemd 254'],
            ['Uptime', '7d 14h 32m 18s'],
          ].map(([k, v]) => (
            <div key={k} className="flex gap-1">
              <span className="text-[#6D40FF]/40 w-12">{k}:</span>
              <span className="text-[#6D40FF]/80">{v}</span>
            </div>
          ))}
        </div>
        <div className="space-y-0.5">
          {[
            ['CPU', 'Intel i7-12700K'],
            ['Cores', '12C / 20T @ 3.6GHz'],
            ['RAM', '16 GB DDR5-4800'],
            ['Disco', 'NVMe 1TB (Samsung 980)'],
            ['GPU', 'NVIDIA RTX 3070'],
          ].map(([k, v]) => (
            <div key={k} className="flex gap-1">
              <span className="text-[#6D40FF]/40 w-12">{k}:</span>
              <span className="text-[#6D40FF]/80">{v}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Resource usage */}
      <div className="border border-[#6D40FF]/30 p-2 mb-3 space-y-1">
        <div className="text-[#6D40FF]/50 text-[10px] mb-1">USO DE RECURSOS</div>
        <ProgressBar label="CPU" value={12} max={100} color="#6D40FF" />
        <ProgressBar label="RAM" value={4200} max={16384} color="#6D40FF" />
        <ProgressBar label="SWAP" value={0} max={8192} color="#FFAA00" />
        <ProgressBar label="Disco /" value={127} max={1000} color="#6D40FF" />
        <ProgressBar label="GPU" value={8} max={100} color="#9966FF" />
      </div>

      <TUISettingItem label="Nombre del Equipo" type="input" value="TERMINUS-WORKSTATION-01"
        storageKey="sys_hostname" description="hostname  │  /etc/hostname" />
      <TUISettingItem label="Idioma del Sistema" type="select" value="Español (es_ES)"
        storageKey="sys_lang" options={LANGS} description="LANG, LC_ALL, LANGUAGE" />
      <TUISettingItem label="Zona Horaria" type="select" value="UTC+01:00 (Madrid/París)"
        storageKey="sys_tz" options={TIMEZONES} description="timedatectl  │  /etc/localtime" />
      <TUISettingItem label="NTP / Reloj en Red" type="toggle" value={true}
        storageKey="sys_ntp" description="Sincronización de hora con servidores NTP" />
      <TUISettingItem label="Init System" type="select" value="systemd"
        storageKey="sys_init" options={INIT_SYSTEMS} description="Sistema de inicialización del OS" />
      <TUISettingItem label="I/O Scheduler" type="select" value="BFQ"
        storageKey="sys_scheduler" options={SCHEDULERS} description="Planificador de E/S del kernel" />
      <TUISettingItem label="Actualizaciones Automáticas" type="toggle" value={true}
        storageKey="sys_autoupd" description="pacman / apt actualización desatendida" />
      <TUISettingItem label="Informes de Errores" type="toggle" value={false}
        storageKey="sys_coredump" description="Habilitar core dumps  │  /proc/sys/kernel/core_pattern" />

      {/* Scrolling boot log */}
      <div className="mt-3 border border-[#6D40FF]/30 p-2 bg-black">
        <div className="text-[#6D40FF]/50 text-[10px] mb-1">KERNEL LOG (dmesg)</div>
        <div className="space-y-0.5">
          {BOOT_LOG.slice(0, logIdx + 1).map((line, i) => (
            <div key={i} className={`text-[10px] ${i === logIdx ? 'text-[#6D40FF]' : 'text-[#6D40FF]/40'}`}>
              {line}
            </div>
          ))}
        </div>
      </div>

      {/* Action buttons */}
      <div className="mt-3 flex gap-2">
        <TUISettingItem label="Reiniciar" type="action" actionLabel="[ REBOOT ]"
          onAction={() => {}} description="" />
        <TUISettingItem label="Apagar" type="action" actionLabel="[ POWEROFF ]"
          onAction={() => {}} description="" danger={true} />
      </div>
    </div>
  );
}
