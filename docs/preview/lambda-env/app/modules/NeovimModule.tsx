// Mirrors: src/lambda-env/internal/modules/neovim/main.go
// Actions: run, set-theme, set-font, set-lines, set-columns, enable/disable features
import { useState } from 'react';
import { C } from '../components/tui/tokens';
import { TUISection } from '../components/tui/TUISection';
import { TUISelect } from '../components/tui/TUISelect';
import { TUIToggle } from '../components/tui/TUIToggle';
import { TUISlider } from '../components/tui/TUISlider';
import { TUIInput } from '../components/tui/TUIInput';
import { TUIButton } from '../components/tui/TUIButton';

// ── Data ────────────────────────────────────────────────────────────────────

const COLORSCHEMES = [
  'catppuccin-mocha', 'catppuccin-latte', 'gruvbox', 'gruvbox-light',
  'nord', 'tokyonight', 'tokyonight-day', 'kanagawa', 'rose-pine',
  'dracula', 'onedark', 'nightfox', 'everforest',
];

const FONTS = [
  'JetBrains Mono', 'Fira Code', 'Cascadia Code',
  'Iosevka', 'Hack', 'Victor Mono', 'Monaspace Neon',
];

const FONT_SIZES = ['10', '11', '12', '13', '14', '15', '16', '18', '20'];

const LSP_SERVERS = [
  { id: 'lua_ls',       name: 'lua_ls',       desc: 'Lua Language Server' },
  { id: 'gopls',        name: 'gopls',         desc: 'Go language server' },
  { id: 'tsserver',     name: 'tsserver',      desc: 'TypeScript / JavaScript' },
  { id: 'rust_analyzer',name: 'rust_analyzer', desc: 'Rust Analyzer' },
  { id: 'pyright',      name: 'pyright',       desc: 'Python static type checker' },
  { id: 'clangd',       name: 'clangd',        desc: 'C / C++ / ObjC' },
  { id: 'bashls',       name: 'bashls',        desc: 'Bash Language Server' },
  { id: 'jsonls',       name: 'jsonls',        desc: 'JSON Language Server' },
];

const PLUGINS = [
  { id: 'treesitter',  name: 'nvim-treesitter',  desc: 'Highlighting y parseo semántico', category: 'core' },
  { id: 'telescope',   name: 'telescope.nvim',   desc: 'Fuzzy finder universal',           category: 'core' },
  { id: 'lspconfig',   name: 'nvim-lspconfig',   desc: 'Configuración de servidores LSP',  category: 'core' },
  { id: 'cmp',         name: 'nvim-cmp',         desc: 'Motor de autocompletado',           category: 'core' },
  { id: 'gitsigns',    name: 'gitsigns.nvim',    desc: 'Indicadores Git en el gutter',      category: 'ui' },
  { id: 'lualine',     name: 'lualine.nvim',     desc: 'Statusline configurada',            category: 'ui' },
  { id: 'neo-tree',    name: 'neo-tree.nvim',    desc: 'Árbol de archivos lateral',         category: 'ui' },
  { id: 'which-key',   name: 'which-key.nvim',   desc: 'Popup de atajos pendientes',        category: 'ui' },
  { id: 'copilot',     name: 'copilot.vim',      desc: 'GitHub Copilot',                   category: 'ai' },
  { id: 'gp',          name: 'gp.nvim',          desc: 'GPT prompt integrado',              category: 'ai' },
  { id: 'conform',     name: 'conform.nvim',     desc: 'Formateador multi-lenguaje',        category: 'fmt' },
  { id: 'lint',        name: 'nvim-lint',        desc: 'Linters asíncronos',               category: 'fmt' },
];

// ── Persistence ──────────────────────────────────────────────────────────────

function ls(k: string, d: any) {
  try { const v = localStorage.getItem('nvim_' + k); return v ? JSON.parse(v) : d; } catch { return d; }
}
function ss(k: string, v: any) { localStorage.setItem('nvim_' + k, JSON.stringify(v)); }

// ── Component ────────────────────────────────────────────────────────────────

export function NeovimModule({ onModified }: { onModified: () => void }) {
  const [theme, setTheme]         = useState(() => ls('theme', 'catppuccin-mocha'));
  const [font, setFont]           = useState(() => ls('font', 'JetBrains Mono'));
  const [fontSize, setFontSize]   = useState(() => ls('fontSize', '13'));
  const [lines, setLines]         = useState(() => ls('lines', 40));
  const [columns, setColumns]     = useState(() => ls('columns', 120));
  const [relativeLn, setRelativeLn] = useState(() => ls('relativeLn', true));
  const [wordWrap, setWordWrap]   = useState(() => ls('wordWrap', false));
  const [clipboard, setClipboard] = useState(() => ls('clipboard', true));
  const [undofile, setUndofile]   = useState(() => ls('undofile', true));
  const [lspEnabled, setLspEnabled] = useState<Record<string, boolean>>(
    () => ls('lsp', Object.fromEntries(LSP_SERVERS.map(s => [s.id, ['lua_ls','gopls','tsserver'].includes(s.id)])))
  );
  const [pluginsEnabled, setPluginsEnabled] = useState<Record<string, boolean>>(
    () => ls('plugins', Object.fromEntries(PLUGINS.map(p => [p.id, !['copilot','gp'].includes(p.id)])))
  );
  const [healthOutput, setHealthOutput] = useState<string[]>([]);
  const [checking, setChecking]   = useState(false);

  function update(k: string, v: any, setter: (x: any) => void) {
    setter(v); ss(k, v); onModified();
  }

  function toggleLsp(id: string) {
    const next = { ...lspEnabled, [id]: !lspEnabled[id] };
    setLspEnabled(next); ss('lsp', next); onModified();
  }

  function togglePlugin(id: string) {
    const next = { ...pluginsEnabled, [id]: !pluginsEnabled[id] };
    setPluginsEnabled(next); ss('plugins', next); onModified();
  }

  async function runHealthCheck() {
    setChecking(true);
    setHealthOutput([]);
    const checks = [
      { label: 'nvim', ok: true,  msg: 'NVIM v0.9.5  Build type: Release' },
      { label: 'python3', ok: true,  msg: 'Python 3.12.1 (pynvim 0.5.0)' },
      { label: 'node',   ok: true,  msg: 'Node v21.4.0' },
      { label: 'go',     ok: true,  msg: 'go version go1.21.5 linux/amd64' },
      { label: 'ripgrep',ok: true,  msg: 'ripgrep 14.1.0' },
      { label: 'fd',     ok: true,  msg: 'fd 9.0.0' },
      { label: 'lazygit',ok: false, msg: 'lazygit: not found in PATH' },
    ];
    for (const c of checks) {
      await new Promise(r => setTimeout(r, 300));
      setHealthOutput(prev => [...prev, `${c.ok ? '✓' : '✗'} ${c.label.padEnd(10)} ${c.msg}`]);
    }
    setChecking(false);
  }

  const pluginsByCategory = ['core', 'ui', 'ai', 'fmt'] as const;
  const catLabels: Record<string, string> = {
    core: 'CORE', ui: 'INTERFAZ', ai: 'IA / ASISTENTES', fmt: 'FORMATO Y LINT',
  };

  // Theme preview colors (rough approximation per scheme)
  const themePreview: Record<string, { bg: string; fg: string; accent: string }> = {
    'catppuccin-mocha':   { bg: '#1E1E2E', fg: '#CDD6F4', accent: '#CBA6F7' },
    'catppuccin-latte':   { bg: '#EFF1F5', fg: '#4C4F69', accent: '#8839EF' },
    'gruvbox':            { bg: '#282828', fg: '#EBDBB2', accent: '#FABD2F' },
    'gruvbox-light':      { bg: '#FBF1C7', fg: '#3C3836', accent: '#D79921' },
    'nord':               { bg: '#2E3440', fg: '#D8DEE9', accent: '#88C0D0' },
    'tokyonight':         { bg: '#1A1B26', fg: '#C0CAF5', accent: '#7AA2F7' },
    'tokyonight-day':     { bg: '#E1E2E7', fg: '#3760BF', accent: '#7847BD' },
    'kanagawa':           { bg: '#1F1F28', fg: '#DCD7BA', accent: '#957FB8' },
    'rose-pine':          { bg: '#191724', fg: '#E0DEF4', accent: '#C4A7E7' },
    'dracula':            { bg: '#282A36', fg: '#F8F8F2', accent: '#BD93F9' },
    'onedark':            { bg: '#282C34', fg: '#ABB2BF', accent: '#61AFEF' },
    'nightfox':           { bg: '#192330', fg: '#CDCECF', accent: '#81A1C1' },
    'everforest':         { bg: '#2D353B', fg: '#D3C6AA', accent: '#A7C080' },
  };
  const preview = themePreview[theme] ?? themePreview['catppuccin-mocha'];

  return (
    <div>
      {/* Theme preview card */}
      <div style={{
        margin: '0 0 4px 0',
        border: `1px solid ${C.border}`,
        overflow: 'hidden',
      }}>
        <div style={{ background: preview.bg, padding: '12px 16px', fontFamily: 'monospace' }}>
          <div style={{ color: preview.accent, fontSize: 11, marginBottom: 6 }}>
            -- {theme} --
          </div>
          <div style={{ fontSize: 12, lineHeight: 1.7 }}>
            <span style={{ color: '#569cd6' }}>local </span>
            <span style={{ color: preview.fg }}>M </span>
            <span style={{ color: preview.fg }}> = </span>
            <span style={{ color: '#ce9178' }}>&#123;&#125;</span>
          </div>
          <div style={{ fontSize: 12, lineHeight: 1.7 }}>
            <span style={{ color: '#569cd6' }}>function </span>
            <span style={{ color: preview.accent }}>M.setup</span>
            <span style={{ color: preview.fg }}>(opts)</span>
          </div>
          <div style={{ fontSize: 12, lineHeight: 1.7 }}>
            {'  '}<span style={{ color: preview.fg }}>vim.cmd </span>
            <span style={{ color: '#ce9178' }}>"colorscheme </span>
            <span style={{ color: preview.accent }}>{theme}</span>
            <span style={{ color: '#ce9178' }}>"</span>
          </div>
          <div style={{ fontSize: 12, lineHeight: 1.7 }}>
            <span style={{ color: '#569cd6' }}>end</span>
          </div>
        </div>
      </div>

      <TUISection title="APARIENCIA">
        <TUISelect
          label="Colorscheme"
          description="set-theme — equivalente a vim.cmd 'colorscheme ...'"
          value={theme}
          options={COLORSCHEMES}
          onChange={v => update('theme', v, setTheme)}
        />
        <TUISelect
          label="Fuente"
          description="set-font — guifont en init.lua (GUI: Neovide, Neovim-qt)"
          value={font}
          options={FONTS}
          onChange={v => update('font', v, setFont)}
        />
        <TUISelect
          label="Tamaño de fuente"
          description="set-font — guifont=Nombre:h{tamaño}"
          value={fontSize}
          options={FONT_SIZES}
          onChange={v => update('fontSize', v, setFontSize)}
        />
      </TUISection>

      <TUISection title="DIMENSIONES DE VENTANA">
        <TUISlider
          label="Líneas (height)"
          description="set-lines — vim.o.lines"
          value={lines}
          min={20}
          max={80}
          unit=""
          onChange={v => update('lines', v, setLines)}
        />
        <TUISlider
          label="Columnas (width)"
          description="set-columns — vim.o.columns"
          value={columns}
          min={60}
          max={220}
          unit=""
          onChange={v => update('columns', v, setColumns)}
        />
      </TUISection>

      <TUISection title="COMPORTAMIENTO DEL EDITOR">
        <TUIToggle
          label="Números relativos"
          description="enable relativenumber — vim.o.relativenumber"
          value={relativeLn}
          onChange={v => update('relativeLn', v, setRelativeLn)}
        />
        <TUIToggle
          label="Ajuste de línea (wrap)"
          description="enable wrap — vim.o.wrap"
          value={wordWrap}
          onChange={v => update('wordWrap', v, setWordWrap)}
        />
        <TUIToggle
          label="Portapapeles del sistema"
          description="enable clipboard — vim.o.clipboard = 'unnamedplus'"
          value={clipboard}
          onChange={v => update('clipboard', v, setClipboard)}
        />
        <TUIToggle
          label="Historial persistente (undofile)"
          description="enable undofile — vim.o.undofile"
          value={undofile}
          onChange={v => update('undofile', v, setUndofile)}
        />
      </TUISection>

      <TUISection title="SERVIDORES LSP" collapsible defaultOpen={true}>
        {LSP_SERVERS.map(srv => (
          <TUIToggle
            key={srv.id}
            label={srv.name}
            description={srv.desc}
            value={lspEnabled[srv.id] ?? false}
            onChange={() => toggleLsp(srv.id)}
          />
        ))}
      </TUISection>

      <TUISection title="PLUGINS" collapsible defaultOpen={false}>
        {pluginsByCategory.map(cat => (
          <div key={cat}>
            {/* category sub-header */}
            <div style={{
              padding: '5px 12px',
              background: C.surface,
              borderBottom: `1px solid ${C.border}`,
              color: C.textMuted,
              fontFamily: 'monospace',
              fontSize: 10,
              letterSpacing: '0.06em',
            }}>
              {catLabels[cat]}
            </div>
            {PLUGINS.filter(p => p.category === cat).map(plugin => (
              <TUIToggle
                key={plugin.id}
                label={plugin.name}
                description={plugin.desc}
                value={pluginsEnabled[plugin.id] ?? true}
                onChange={() => togglePlugin(plugin.id)}
              />
            ))}
          </div>
        ))}
      </TUISection>

      <TUISection title="DIAGNÓSTICO" collapsible defaultOpen={false}>
        <div style={{ padding: '10px 12px', borderBottom: `1px solid ${C.border}` }}>
          <TUIButton
            label={checking ? 'Ejecutando :checkhealth...' : 'Ejecutar :checkhealth'}
            onClick={runHealthCheck}
          />
        </div>
        {healthOutput.length > 0 && (
          <div style={{ padding: '10px 14px', fontFamily: 'monospace', background: C.surface }}>
            {healthOutput.map((line, i) => (
              <div key={i} style={{
                fontSize: 11,
                color: line.startsWith('✓') ? C.success : line.startsWith('✗') ? C.error : C.textSecondary,
                lineHeight: 1.7,
              }}>
                {line}
              </div>
            ))}
          </div>
        )}
      </TUISection>
    </div>
  );
}
