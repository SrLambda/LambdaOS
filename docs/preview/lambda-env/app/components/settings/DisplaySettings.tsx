import { TUISettingItem } from '../TUISettingItem';

const RESOLUTIONS = ['800x600', '1024x768', '1280x720', '1280x1024', '1366x768', '1920x1080', '2560x1440', '3840x2160'];
const REFRESH_RATES = ['24Hz', '30Hz', '60Hz', '75Hz', '100Hz', '120Hz', '144Hz', '165Hz', '240Hz'];
const COLOR_DEPTHS = ['16-bit (65K)', '24-bit (16.7M)', '32-bit (4.3B)'];
const GAMMA_PROFILES = ['sRGB', 'AdobeRGB', 'P3-D65', 'Linear', 'Custom'];
const ORIENT = ['Normal (0°)', 'Rotación 90°', 'Rotación 180°', 'Rotación 270°'];

export function DisplaySettings() {
  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ PANTALLA Y MONITOR ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">xrandr / drm</span>
      </div>

      <TUISettingItem label="Brillo" type="slider" value={75} storageKey="display_brightness"
        description="Nivel de luminosidad de la pantalla" unit="%" />
      <TUISettingItem label="Contraste" type="slider" value={50} storageKey="display_contrast"
        description="Diferencial entre zonas claras y oscuras" unit="%" />
      <TUISettingItem label="Temperatura Color" type="slider" value={65} storageKey="display_temp"
        min={20} max={100} unit=""
        description="20=Cálido(3200K) → 100=Frío(6500K)" />
      <TUISettingItem label="Gamma" type="slider" value={22} storageKey="display_gamma"
        min={10} max={30} unit=""
        description="Corrección gamma  (10=oscuro, 22=estándar, 30=claro)" />
      <TUISettingItem label="Resolución" type="select" value="1920x1080" storageKey="display_res"
        options={RESOLUTIONS} description="Resolución de pantalla en píxeles" />
      <TUISettingItem label="Tasa de Refresco" type="select" value="60Hz" storageKey="display_hz"
        options={REFRESH_RATES} description="Frecuencia de actualización del panel" />
      <TUISettingItem label="Profundidad de Color" type="select" value="24-bit (16.7M)" storageKey="display_depth"
        options={COLOR_DEPTHS} description="Bits por canal de color" />
      <TUISettingItem label="Perfil de Color" type="select" value="sRGB" storageKey="display_gamma_p"
        options={GAMMA_PROFILES} description="Espacio de color y calibración" />
      <TUISettingItem label="Orientación" type="select" value="Normal (0°)" storageKey="display_orient"
        options={ORIENT} description="Rotación del framebuffer" />
      <TUISettingItem label="Modo Nocturno" type="toggle" value={false} storageKey="display_nightmode"
        description="Filtro luz azul activo entre 20:00–07:00" />
      <TUISettingItem label="Protector de Pantalla" type="toggle" value={true} storageKey="display_ss"
        description="Activar screensaver tras 5 min de inactividad" />
      <TUISettingItem label="HDR" type="toggle" value={false} storageKey="display_hdr"
        description="Alto rango dinámico (requiere panel compatible)" />

      {/* Monitor info panel */}
      <div className="mt-4 border border-[#6D40FF]/30 p-2">
        <div className="text-[#6D40FF]/50 text-[10px] mb-1">INFORMACIÓN DEL MONITOR</div>
        <div className="grid grid-cols-2 gap-x-4 gap-y-0.5 text-[10px]">
          {[
            ['Modelo', 'DELL U2722D IPS'],
            ['Conector', 'DisplayPort 1.4'],
            ['Resolución Nativa', '2560x1440 @ 60Hz'],
            ['Panel', 'IPS 27" 16:9'],
            ['Brillo Máx', '350 cd/m²'],
            ['Contraste', '1000:1'],
          ].map(([k, v]) => (
            <div key={k} className="flex gap-1">
              <span className="text-[#6D40FF]/40">{k}:</span>
              <span className="text-[#6D40FF]/80">{v}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
