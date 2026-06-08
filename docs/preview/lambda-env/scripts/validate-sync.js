const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..');
const srcRoot = path.resolve(__dirname, '../../../../src/lambda-env');

let failures = 0;
function assert(condition, message) {
  if (!condition) {
    console.error(`FAIL: ${message}`);
    failures++;
  }
}

// ── Load Go source of truth ──
const iconMapPath = path.join(srcRoot, 'internal/tui/icons/icon-map.json');
const iconMap = JSON.parse(fs.readFileSync(iconMapPath, 'utf8'));

const colorsGoPath = path.join(srcRoot, 'internal/tui/theme/colors.go');
const colorsGo = fs.readFileSync(colorsGoPath, 'utf8');

// ── Read generated icon-map.ts to verify it matches JSON ──
const iconMapTsPath = path.join(root, 'app/data/icon-map.ts');
const iconMapTs = fs.readFileSync(iconMapTsPath, 'utf8');

for (const [section, entries] of Object.entries(iconMap)) {
  for (const [key, vals] of Object.entries(entries)) {
    assert(iconMapTs.includes(vals.nerd), `icon-map.ts: should contain nerd icon for ${section}.${key}`);
    assert(iconMapTs.includes(vals.fallback), `icon-map.ts: should contain fallback icon for ${section}.${key}`);
  }
}

// ── Task 4.1: categories.ts icons match icon-map.json via ICONS import ──
const categoriesPath = path.join(root, 'app/data/categories.ts');
const categoriesSrc = fs.readFileSync(categoriesPath, 'utf8');

assert(categoriesSrc.includes("import { ICONS } from './icon-map'"), 'categories.ts: should import ICONS from icon-map');

const categoryMap = { sistema: 'system', aplicaciones: 'apps', operaciones: 'ops' };
for (const [catId, mapKey] of Object.entries(categoryMap)) {
  const expectedRef = `ICONS.categories.${mapKey}.nerd`;
  const re = new RegExp(`id:\\s*'${catId}'[\\s\\S]*?icon:\\s*${expectedRef.replace(/\./g, '\\.')}`);
  assert(re.test(categoriesSrc), `categories.ts: category ${catId} icon should reference ${expectedRef}`);
}

// Module icon mapping (modules that exist in icon-map.json)
const moduleMap = {
  pantalla: 'display',
  audio: 'audio',
  red: 'network',
  bluetooth: 'bluetooth',
  seguridad: 'security',
  neovim: 'neovim',
  qtile: 'qtile',
  dotfiles: 'dotfiles',
  logs: 'logs',
  almacenamiento: 'storage',
  teclado: 'keyboard',
  apariencia: 'appearance',
  energia: 'power',
  defaults: 'defaults',
  monitor: 'hardware-dashboard',
};

for (const [modId, mapKey] of Object.entries(moduleMap)) {
  const dotRef = `ICONS.modules.${mapKey}.nerd`;
  const bracketRef = `ICONS.modules['${mapKey}'].nerd`;
  // Find the module block in the source
  const modBlockRe = new RegExp(`\\{ id:\\s*'${modId}'[^\\}]*\\}`);
  const modBlock = categoriesSrc.match(modBlockRe)?.[0] || '';
  assert(modBlock.includes(dotRef) || modBlock.includes(bracketRef),
    `categories.ts: module ${modId} icon should reference ${dotRef} or ${bracketRef}`);
}

// ── Task 4.1b: every category and module must have fallbackIcon ──
const catBlocks = categoriesSrc.match(/id:\s*'[a-z_]+'[\s\S]*?modules:\s*\[/g) || [];
for (const block of catBlocks) {
  const catId = block.match(/id:\s*'([a-z_]+)'/)?.[1];
  assert(block.includes('fallbackIcon'), `categories.ts: category ${catId} must have fallbackIcon`);
}

const modBlocks = categoriesSrc.match(/\{ id:\s*'[a-z_]+'[^}]*\}/g) || [];
for (const block of modBlocks) {
  const modId = block.match(/id:\s*'([a-z_]+)'/)?.[1];
  assert(block.includes('fallbackIcon'), `categories.ts: module ${modId} must have fallbackIcon`);
}

// Widget icons in TUIToggle.tsx
const togglePath = path.join(root, 'app/components/tui/TUIToggle.tsx');
const toggleSrc = fs.readFileSync(togglePath, 'utf8');
assert(toggleSrc.includes('ICONS.widgets.toggle_on.nerd'), 'TUIToggle.tsx: should use ICONS.widgets.toggle_on.nerd');
assert(toggleSrc.includes('ICONS.widgets.toggle_off.nerd'), 'TUIToggle.tsx: should use ICONS.widgets.toggle_off.nerd');

// Widget icons in TUIButton.tsx
const buttonPath = path.join(root, 'app/components/tui/TUIButton.tsx');
const buttonSrc = fs.readFileSync(buttonPath, 'utf8');
assert(buttonSrc.includes('ICONS.widgets.success.nerd'), 'TUIButton.tsx: should use ICONS.widgets.success.nerd');

// Widget icons in TUIInput.tsx
const inputPath = path.join(root, 'app/components/tui/TUIInput.tsx');
const inputSrc = fs.readFileSync(inputPath, 'utf8');
assert(inputSrc.includes('ICONS.widgets.error.nerd'), 'TUIInput.tsx: should use ICONS.widgets.error.nerd');
assert(inputSrc.includes('ICONS.widgets.success.nerd'), 'TUIInput.tsx: should use ICONS.widgets.success.nerd');
assert(inputSrc.includes('ICONS.widgets.confirm.nerd'), 'TUIInput.tsx: should use ICONS.widgets.confirm.nerd');

// Lock icon in App.tsx
const appPath = path.join(root, 'app/App.tsx');
const appSrc = fs.readFileSync(appPath, 'utf8');
assert(appSrc.includes('ICONS.widgets.lock.nerd'), 'App.tsx: should use ICONS.widgets.lock.nerd');
assert(appSrc.includes('ICONS.categories.system.nerd'), 'App.tsx: header should use ICONS.categories.system.nerd');

// ── Task 4.2: tokens.ts sync with colors.go ──
const tokensPath = path.join(root, 'app/components/tui/tokens.ts');
const tokensSrc = fs.readFileSync(tokensPath, 'utf8');

// Extract Dimmed from colors.go
const dimmedMatch = colorsGo.match(/Dimmed\s*=\s*"([^"]+)"/);
const dimmedValue = dimmedMatch ? dimmedMatch[1] : null;
assert(dimmedValue, `colors.go: should define Dimmed constant`);
assert(tokensSrc.includes(`dimmed:`), `tokens.ts: should define dimmed token`);
const dimmedTokenMatch = tokensSrc.match(/dimmed:\s*'([^']+)'/);
if (dimmedTokenMatch) {
  assert(dimmedTokenMatch[1] === dimmedValue, `tokens.ts: dimmed should be "${dimmedValue}", got "${dimmedTokenMatch[1]}"`);
} else {
  assert(false, `tokens.ts: dimmed token value not found`);
}

// Check comment referencing colors.go
assert(tokensSrc.includes('colors.go') || tokensSrc.includes('src/lambda-env/internal/tui/theme/colors.go'),
  `tokens.ts: should contain comment referencing Go source of truth colors.go`);

// ── Task 4.3: Makefile has sync-icons target ──
const makefilePath = path.join(__dirname, '../../../../Makefile');
const makefileSrc = fs.readFileSync(makefilePath, 'utf8');
assert(makefileSrc.includes('sync-icons'), `Makefile: should have sync-icons target`);
assert(makefileSrc.includes('generate-icon-map.js'), `Makefile: sync-icons should run generate-icon-map.js`);
assert(makefileSrc.includes('validate-sync.js'), `Makefile: sync-icons should run validate-sync.js`);

if (failures > 0) {
  console.error(`\n${failures} validation(s) FAILED`);
  process.exit(1);
} else {
  console.log('All validations passed ✓');
  process.exit(0);
}
