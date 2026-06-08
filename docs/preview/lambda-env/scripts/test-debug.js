const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..');
console.log('root:', root);

const categoriesPath = path.join(root, 'app/data/categories.ts');
console.log('categoriesPath:', categoriesPath);
console.log('exists:', fs.existsSync(categoriesPath));

const categoriesSrc = fs.readFileSync(categoriesPath, 'utf8');
console.log('sistema section:', categoriesSrc.substring(categoriesSrc.indexOf("id: 'sistema'"), categoriesSrc.indexOf("id: 'sistema'") + 200));

const catId = 'sistema';
const expectedRef = 'ICONS.categories.system.nerd';
const re = new RegExp(`id:\\s*'${catId}'[\\s\\S]*?icon:\\s*${expectedRef.replace(/\\./g, '\\.')}`);
console.log('regex:', re);
console.log('test:', re.test(categoriesSrc));

// Also test the module regex
const modId = 'pantalla';
const modExpected = 'ICONS.modules.display.nerd';
const modRe = new RegExp(`id:\\s*'${modId}'[^}]*icon:\\s*${modExpected.replace(/\\./g, '\\.')}`);
console.log('mod regex:', modRe);
console.log('mod test:', modRe.test(categoriesSrc));
