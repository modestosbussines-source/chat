const fs = require('fs');
const data = JSON.parse(fs.readFileSync('C:\\\\Users\\\\Igor\\\\Documents\\\\chatwo\\\\config.json', 'utf8'));

console.log("Root keys:", Object.keys(data));
if (Array.isArray(data.data)) {
  console.log("Data array length:", data.data.length);
  data.data.forEach(item => {
    console.log(`- name: ${item.name}, keys: ${Object.keys(item).join(', ')}, row_count: ${item.rows?.length || 0}`);
  });
} else {
  console.log("Data is not an array:", typeof data.data);
}
