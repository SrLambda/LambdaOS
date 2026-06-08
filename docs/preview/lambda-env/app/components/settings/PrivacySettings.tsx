import { useState } from 'react';
import { TUISettingItem } from '../TUISettingItem';

const TRACKER_LEVELS = ['Desactivado', 'Bajo', 'Medio', 'Alto', 'Máximo (Paranoico)'];
const CLEAR_OPTIONS = ['Caché', 'Cookies', 'Historial de Comandos', 'Archivos Temporales', 'Logs del Sistema', 'TODO'];

export function PrivacySettings() {
  const [locked, setLocked] = useState(false);
  const [log, setLog] = useState<string[]>([]);

  function addLog(msg: string) {
    setLog(prev => [`[${new Date().toLocaleTimeString('es-ES', { hour12: false })}] ${msg}`, ...prev.slice(0, 4)]);
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ PRIVACIDAD Y SEGURIDAD ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">SELinux / AppArmor</span>
      </div>

      {/* Security level indicator */}
      <div className="border border-[#6D40FF]/30 p-2 mb-3 flex items-center gap-3">
        <span className="text-[10px] text-[#6D40FF]/50">NIVEL DE SEGURIDAD:</span>
        {['CRÍTICO','BAJO','MEDIO','ALTO','MÁXIMO'].map((lvl, i) => (
          <span key={lvl} className={`text-[10px] px-1.5 py-0.5 border
            ${i === 3 ? 'bg-[#6D40FF] text-black border-[#6D40FF]' : 'text-[#6D40FF]/30 border-[#6D40FF]/20'}`}>
            {lvl}
          </span>
        ))}
      </div>

      <TUISettingItem label="Compartir Ubicación" type="toggle" value={false} storageKey="priv_location"
        description="Acceso GPS/triangulación  │  Aplicaciones: 0 activas" />
      <TUISettingItem label="Acceso a Cámara" type="toggle" value={false} storageKey="priv_camera"
        description="LED indicador activo cuando en uso" />
      <TUISettingItem label="Acceso a Micrófono" type="toggle" value={false} storageKey="priv_mic"
        description="Monitor de acceso por proceso habilitado" />
      <TUISettingItem label="Telemetría del OS" type="toggle" value={false} storageKey="priv_telemetry"
        description="Envío de estadísticas de uso anónimas" />
      <TUISettingItem label="Informes de Fallos" type="toggle" value={true} storageKey="priv_crash"
        description="Reportes a desarrolladores (sin datos personales)" />
      <TUISettingItem label="Bloqueo de Rastreadores" type="select" value="Alto" storageKey="priv_tracker"
        options={TRACKER_LEVELS} description="Nivel de protección contra fingerprinting y tracking" />
      <TUISettingItem label="Cookies de Terceros" type="toggle" value={false} storageKey="priv_3rdcookie"
        description="Bloquear cookies cross-site" />
      <TUISettingItem label="Cifrado de Disco (LUKS)" type="toggle" value={true} storageKey="priv_luks"
        description="Full-disk encryption activo  │  Algoritmo: AES-256-XTS" />
      <TUISettingItem label="Firewall AppArmor" type="toggle" value={true} storageKey="priv_apparmor"
        description="Perfiles de confinamiento para 847 aplicaciones" />
      <TUISettingItem label="Sudo sin Contraseña" type="toggle" value={false} storageKey="priv_nopasswd"
        description="WHEEL_NOPASSWD en /etc/sudoers" danger={true} />

      <div className="mt-4 grid grid-cols-2 gap-3">
        {/* Danger zone */}
        <div className="border border-[#FF4040]/30 p-2">
          <div className="text-[#FF4040]/70 text-[10px] mb-2">⚠ ZONA DE PELIGRO</div>
          <TUISettingItem label="Bloquear Sesión" type="action"
            actionLabel="[ BLOQUEAR ]"
            onAction={() => { setLocked(true); addLog('Sesión bloqueada por usuario'); }}
            description="Lock inmediato de pantalla" danger={true} />
          <TUISettingItem label="Limpiar Datos" type="action"
            actionLabel="[ LIMPIAR ]"
            onAction={() => addLog('Limpieza de datos iniciada')}
            description="Caché, cookies, temporales" danger={true} />
        </div>
        {/* Audit log */}
        <div className="border border-[#6D40FF]/30 p-2">
          <div className="text-[#6D40FF]/50 text-[10px] mb-1">LOG DE AUDITORÍA</div>
          <div className="space-y-0.5 min-h-[60px]">
            {log.length === 0 && (
              <div className="text-[#6D40FF]/20 text-[10px]">── sin eventos recientes ──</div>
            )}
            {log.map((entry, i) => (
              <div key={i} className="text-[10px] text-[#6D40FF]/60">{entry}</div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
