import { useState } from 'react';

interface TUISettingItemProps {
  label: string;
  type: 'toggle' | 'select' | 'slider' | 'input' | 'action';
  value?: string | number | boolean;
  options?: string[];
  onChange?: (value: any) => void;
  description?: string;
  storageKey?: string;
  actionLabel?: string;
  onAction?: () => void;
  danger?: boolean;
  min?: number;
  max?: number;
  unit?: string;
}

function loadValue(key: string | undefined, fallback: any) {
  if (!key) return fallback;
  try {
    const v = localStorage.getItem(`tui_${key}`);
    if (v === null) return fallback;
    const parsed = JSON.parse(v);
    return parsed;
  } catch {
    return fallback;
  }
}

function saveValue(key: string | undefined, value: any) {
  if (!key) return;
  localStorage.setItem(`tui_${key}`, JSON.stringify(value));
}

export function TUISettingItem({
  label, type, value, options, onChange, description,
  storageKey, actionLabel, onAction, danger, min = 0, max = 100, unit = '%'
}: TUISettingItemProps) {
  const [localValue, setLocalValue] = useState(() => loadValue(storageKey, value));

  const handleChange = (newValue: any) => {
    setLocalValue(newValue);
    saveValue(storageKey, newValue);
    onChange?.(newValue);
  };

  const accentColor = danger ? '#FF4040' : '#6D40FF';
  const borderClass = danger ? 'border-[#FF4040]/40' : 'border-[#6D40FF]/30';

  return (
    <div className={`border-b ${borderClass} py-2 flex items-start gap-4 group hover:bg-[#6D40FF]/5 px-1`}>
      {/* Label column */}
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-[#6D40FF]/40 text-[10px] select-none">›</span>
          <span className={`text-xs ${danger ? 'text-[#FF4040]' : 'text-[#6D40FF]'}`}>{label}</span>
        </div>
        {description && (
          <div className="text-[#6D40FF]/40 text-[10px] ml-4 mt-0.5">{description}</div>
        )}
      </div>

      {/* Control column */}
      <div className="shrink-0 flex items-center">
        {type === 'toggle' && (
          <button
            onClick={() => handleChange(!localValue)}
            className={`
              flex items-center gap-1.5 px-2 py-0.5 border text-[10px] transition-all min-w-[72px] justify-center
              ${localValue
                ? `bg-[${accentColor}] text-black border-[${accentColor}] shadow-[0_0_8px_rgba(109,64,255,0.6)]`
                : `bg-black text-[${accentColor}] border-[${accentColor}]/60 hover:border-[${accentColor}]`
              }
            `}
            style={localValue
              ? { background: accentColor, color: '#000', borderColor: accentColor, boxShadow: `0 0 8px ${accentColor}80` }
              : { background: '#000', color: accentColor, borderColor: `${accentColor}60` }
            }
          >
            <span className="font-bold">{localValue ? '●' : '○'}</span>
            <span>{localValue ? ' ON ' : 'OFF'}</span>
          </button>
        )}

        {type === 'select' && options && (
          <div className="relative">
            <select
              value={localValue as string}
              onChange={(e) => handleChange(e.target.value)}
              className="bg-black text-[#6D40FF] border border-[#6D40FF]/60 px-2 py-0.5 text-[10px] outline-none
                focus:border-[#6D40FF] focus:bg-[#6D40FF]/10 appearance-none pr-6 cursor-pointer min-w-[160px]"
            >
              {options.map((option) => (
                <option key={option} value={option} className="bg-[#0a0018] text-[#6D40FF]">
                  {option}
                </option>
              ))}
            </select>
            <span className="absolute right-1.5 top-1/2 -translate-y-1/2 text-[#6D40FF] text-[8px] pointer-events-none">▼</span>
          </div>
        )}

        {type === 'slider' && (
          <div className="flex items-center gap-2">
            <div className="text-[#6D40FF]/40 text-[10px] w-4 text-right">{min}</div>
            <div className="relative flex items-center">
              <input
                type="range"
                min={min}
                max={max}
                value={localValue as number}
                onChange={(e) => handleChange(parseInt(e.target.value))}
                className="tui-slider w-28"
              />
              {/* tick marks */}
              <div className="absolute -bottom-2 left-0 right-0 flex justify-between px-0.5">
                {[0,25,50,75,100].map(t => (
                  <div key={t} className="w-px h-1 bg-[#6D40FF]/30" />
                ))}
              </div>
            </div>
            <div className="text-[#6D40FF]/40 text-[10px] w-4">{max}</div>
            <div className="text-[#6D40FF] text-[10px] w-10 text-right border border-[#6D40FF]/40 px-1">
              {localValue}{unit}
            </div>
          </div>
        )}

        {type === 'input' && (
          <input
            type="text"
            value={localValue as string}
            onChange={(e) => handleChange(e.target.value)}
            className="bg-black text-[#6D40FF] border border-[#6D40FF]/60 px-2 py-0.5 text-[10px]
              outline-none focus:border-[#6D40FF] focus:bg-[#6D40FF]/10 w-48
              caret-[#6D40FF]"
            spellCheck={false}
          />
        )}

        {type === 'action' && (
          <button
            onClick={onAction}
            className={`
              px-3 py-0.5 border text-[10px] transition-all
              ${danger
                ? 'border-[#FF4040]/60 text-[#FF4040] hover:bg-[#FF4040] hover:text-black hover:border-[#FF4040]'
                : 'border-[#6D40FF]/60 text-[#6D40FF] hover:bg-[#6D40FF] hover:text-black hover:border-[#6D40FF]'
              }
            `}
          >
            {actionLabel ?? '[ EJECUTAR ]'}
          </button>
        )}
      </div>
    </div>
  );
}
