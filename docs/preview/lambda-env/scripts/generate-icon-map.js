const fs = require('fs');

const json = JSON.parse(fs.readFileSync('/home/lambda/Projects/LambdaOS/src/lambda-env/internal/tui/icons/icon-map.json', 'utf8'));

function escapeStr(s) {
  return s.replace(/\\/g, '\\\\').replace(/'/g, "\\'").replace(/\n/g, '\\n');
}

let out = '// Generated from src/lambda-env/internal/tui/icons/icon-map.json\n';
out += '// Run: make sync-icons  (do NOT edit manually)\n\n';
out += 'export const ICONS = {\n';

for (const [section, entries] of Object.entries(json)) {
  out += '  ' + section + ': {\n';
  for (const [key, vals] of Object.entries(entries)) {
    out += '    ' + key + ': { nerd: \'' + escapeStr(vals.nerd) + '\', fallback: \'' + escapeStr(vals.fallback) + '\' },\n';
  }
  out += '  },\n';
}
out += '};\n';

fs.writeFileSync('/home/lambda/Projects/LambdaOS/docs/preview/lambda-env/app/data/icon-map.ts', out);
console.log('Generated icon-map.ts');
