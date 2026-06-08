import { useState } from 'react';
import { TUISettingItem } from '../TUISettingItem';

const OUTPUT_DEVICES = ['Speakers (Realtek ALC892)', 'Headphones (3.5mm)', 'HDMI Audio (GPU)', 'USB DAC (FiiO)', 'Bluetooth (Sony WH-1000XM5)'];
const INPUT_DEVICES = ['Micrófono Integrado', 'USB Microphone (Blue Yeti)', 'Headset Mic (Jabra)', 'Line-In'];
const EQ_PRESETS = ['Plano', 'Bass Boost', 'Treble', 'Clásico', 'Rock', 'Jazz', 'Vocal', 'Gaming'];
const SAMPLE_RATES = ['44100 Hz', '48000 Hz', '96000 Hz', '192000 Hz'];

const VU_CHARS = ['▁','▂','▃','▄','▅','▆','▇','█'];

function VUMeter({ label }: { label: string }) {
  const levels = Array.from({ length: 16 }, () => Math.floor(Math.random() * 8));
  return (
    <div className="flex items-end gap-px h-6 mt-1">
      {levels.map((l, i) => (
        <span key={i} className={`text-[8px] leading-none ${l > 5 ? 'text-[#FF4040]' : l > 3 ? 'text-[#FFAA00]' : 'text-[#6D40FF]'}`}>
          {VU_CHARS[l]}
        </span>
      ))}
    </div>
  );
}

export function AudioSettings() {
  const [vuActive, setVuActive] = useState(false);

  return (
    <div>
      <div className="flex items-center justify-between mb-3 pb-1 border-b border-[#6D40FF]/50">
        <span className="text-[#6D40FF] text-xs tracking-wider">[ AUDIO Y SONIDO ]</span>
        <span className="text-[#6D40FF]/50 text-[10px]">alsa / pulseaudio / pipewire</span>
      </div>

      <TUISettingItem label="Volumen Principal" type="slider" value={80} storageKey="audio_master"
        description="Master volume  │  amixer sset Master" unit="%" />
      <TUISettingItem label="Vol. Aplicaciones" type="slider" value={90} storageKey="audio_apps"
        description="PulseAudio sink: sistema y apps" unit="%" />
      <TUISettingItem label="Vol. Notificaciones" type="slider" value={40} storageKey="audio_notif"
        description="Sonidos de alerta del sistema" unit="%" />
      <TUISettingItem label="Dispositivo de Salida" type="select" value="Speakers (Realtek ALC892)"
        storageKey="audio_out_dev" options={OUTPUT_DEVICES}
        description="Sink PulseAudio activo" />
      <TUISettingItem label="Dispositivo de Entrada" type="select" value="Micrófono Integrado"
        storageKey="audio_in_dev" options={INPUT_DEVICES}
        description="Source PulseAudio activo" />
      <TUISettingItem label="Volumen Micrófono" type="slider" value={65} storageKey="audio_mic_vol"
        description="Ganancia del micrófono activo" unit="%" />
      <TUISettingItem label="Tasa de Muestreo" type="select" value="48000 Hz" storageKey="audio_rate"
        options={SAMPLE_RATES} description="Sample rate del daemon de audio" />
      <TUISettingItem label="Perfil EQ" type="select" value="Plano" storageKey="audio_eq"
        options={EQ_PRESETS} description="Ecualizador paramétrico" />
      <TUISettingItem label="Reducción Ruido" type="toggle" value={false} storageKey="audio_nr"
        description="Cancelación de ruido RNNoise (CPU +3%)" />
      <TUISettingItem label="Sonidos del Sistema" type="toggle" value={true} storageKey="audio_beep"
        description="Efectos de sonido para eventos del OS" />
      <TUISettingItem label="Modo Exclusivo" type="toggle" value={false} storageKey="audio_excl"
        description="Permitir apps tomar control exclusivo del HW" />

      {/* VU Meter */}
      <div className="mt-4 border border-[#6D40FF]/30 p-2">
        <div className="flex items-center justify-between mb-1">
          <span className="text-[#6D40FF]/50 text-[10px]">ANALIZADOR DE ESPECTRO</span>
          <button
            onClick={() => setVuActive(!vuActive)}
            className="text-[10px] border border-[#6D40FF]/40 text-[#6D40FF]/70 px-2 py-0.5 hover:bg-[#6D40FF]/10"
          >
            {vuActive ? '[ DETENER ]' : '[ SIMULAR ]'}
          </button>
        </div>
        {vuActive ? (
          <>
            <div className="text-[#6D40FF]/40 text-[10px]">OUT:</div>
            <VUMeter label="out" />
            <div className="text-[#6D40FF]/40 text-[10px] mt-1">IN:</div>
            <VUMeter label="in" />
          </>
        ) : (
          <div className="text-[#6D40FF]/30 text-[10px] text-center py-2">── sin señal ──</div>
        )}
      </div>
    </div>
  );
}
