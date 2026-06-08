import { useEffect, useState } from 'react';

interface TUIWindowProps {
  title: string;
  children: React.ReactNode;
  width?: string;
}

function useTime() {
  const [time, setTime] = useState(new Date());
  useEffect(() => {
    const id = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(id);
  }, []);
  return time;
}

export function TUIWindow({ title, children, width = 'max-w-5xl' }: TUIWindowProps) {
  const time = useTime();
  const [blink, setBlink] = useState(true);

  useEffect(() => {
    const id = setInterval(() => setBlink(b => !b), 500);
    return () => clearInterval(id);
  }, []);

  const timeStr = time.toLocaleTimeString('es-ES', { hour12: false });
  const dateStr = time.toLocaleDateString('es-ES', { day: '2-digit', month: '2-digit', year: 'numeric' });

  return (
    <div className={`bg-black text-[#6D40FF] font-mono w-full ${width} shadow-[0_0_40px_rgba(109,64,255,0.4),0_0_80px_rgba(109,64,255,0.15)] border border-[#6D40FF]`}>
      {/* Top chrome bar */}
      <div className="bg-[#0a0018] border-b border-[#6D40FF] px-2 py-1 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <span className="text-[10px] text-[#6D40FF]/50 select-none">[ SYS ]</span>
          <span className="text-[10px] text-[#6D40FF] tracking-widest">{title}</span>
        </div>
        <div className="flex items-center gap-4 text-[10px]">
          <span className="text-[#6D40FF]/60">{dateStr}</span>
          <span className="text-[#6D40FF] font-bold tracking-wider">{timeStr}</span>
          <span className={`text-[#6D40FF] ${blink ? 'opacity-100' : 'opacity-0'}`}>█</span>
        </div>
      </div>

      {/* Double border title */}
      <div className="px-3 pt-2 pb-1 border-b border-[#6D40FF]/30">
        <div className="text-[#6D40FF]/40 text-[10px] leading-none select-none overflow-hidden">
          {'─'.repeat(120)}
        </div>
        <div className="flex items-center justify-between py-1 px-1">
          <div className="flex items-center gap-2">
            <span className="text-[#6D40FF] text-xs">◈</span>
            <span className="text-[#6D40FF] text-xs tracking-wider">LAMBDA-OS</span>
            <span className="text-[#6D40FF]/40 text-[10px]">│</span>
            <span className="text-[#6D40FF]/70 text-[10px]">System Preferences</span>
          </div>
          <div className="flex items-center gap-3 text-[10px] text-[#6D40FF]/50">
            <span>F1:Help</span>
            <span>F5:Refresh</span>
            <span>F10:Exit</span>
            <span>Tab:Next</span>
          </div>
        </div>
        <div className="text-[#6D40FF]/40 text-[10px] leading-none select-none overflow-hidden">
          {'─'.repeat(120)}
        </div>
      </div>

      {/* Content */}
      <div className="p-3">
        {children}
      </div>

      {/* Status bar */}
      <div className="border-t border-[#6D40FF] bg-[#6D40FF] text-black px-3 py-0.5 flex items-center justify-between text-[10px]">
        <div className="flex items-center gap-3">
          <span className="font-bold">TERMINUS-OS v2.4.1</span>
          <span>│</span>
          <span>root@localhost</span>
          <span>│</span>
          <span>TTY1</span>
        </div>
        <div className="flex items-center gap-3">
          <span>CPU: 12%</span>
          <span>MEM: 4.2G/16G</span>
          <span>SWAP: 0K</span>
          <span>│</span>
          <span>uptime: 7d 14h 32m</span>
        </div>
      </div>
    </div>
  );
}
