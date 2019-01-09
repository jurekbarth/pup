
const fs = require('fs');
const path = require('path');

const template = ({ clientId }) => `
const clientId = "${clientId}";
module.exports = clientId;
`;


(async () => {
  try {
    const args = process.argv.slice(2);
    const clientId = args[0];
    fs.writeFileSync(path.resolve(__dirname, '../lambdaEdge/clientId.js'), template({ clientId }));
  } catch (error) {
    console.log(error)
  }
})();

