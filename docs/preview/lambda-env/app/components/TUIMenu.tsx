interface MenuItem {
  id: string;
  label: string;
  icon?: string;
  shortcut?: string;
}

interface TUIMenuProps {
  items: MenuItem[];
  selectedId: string;
  onSelect: (id: string) => void;
  focusedId?: string;
}

export function TUIMenu({ items, selectedId, onSelect, focusedId }: TUIMenuProps) {
  return (
    <div className="space-y-0 border border-[#6D40FF]/50">
      {/* Menu header */}
      <div className="bg-[#6D40FF]/20 border-b border-[#6D40FF]/50 px-2 py-1">
        <span className="text-[#6D40FF] text-[10px] tracking-widest">┤ MENÚ ├</span>
      </div>
      {items.map((item, index) => {
        const isSelected = selectedId === item.id;
        const isFocused = focusedId === item.id;
        return (
          <div
            key={item.id}
            onClick={() => onSelect(item.id)}
            className={`
              cursor-pointer px-2 py-1.5 flex items-center gap-2 border-b border-[#6D40FF]/20 last:border-b-0 transition-all select-none
              ${isSelected
                ? 'bg-[#6D40FF] text-black'
                : isFocused
                ? 'bg-[#6D40FF]/20 text-[#6D40FF]'
                : 'bg-black text-[#6D40FF]/80 hover:bg-[#6D40FF]/10'
              }
            `}
          >
            <span className="text-[10px] w-4 text-center opacity-60">{String(index + 1)}</span>
            <span className="text-xs w-4 text-center">{item.icon}</span>
            <span className="flex-1 text-xs">{item.label}</span>
            {isSelected && <span className="text-[10px]">◀</span>}
          </div>
        );
      })}
    </div>
  );
}
