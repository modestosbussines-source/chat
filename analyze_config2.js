const fs = require('fs');
const data = JSON.parse(fs.readFileSync('C:\\\\Users\\\\Igor\\\\Documents\\\\chatwo\\\\config.json', 'utf8'));

data.data.forEach(item => {
  if (['groups', 'items', 'scripts'].includes(item.name)) {
    console.log(`\n--- Example of ${item.name} ---`);
    console.log(JSON.stringify(item.rows[0], null, 2));
  }
});
