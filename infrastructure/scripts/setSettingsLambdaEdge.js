const https = require('https');
const fs = require('fs');
const path = require('path');


const getKeys = (region, cognitoPoolId) => new Promise((resolve, reject) => {
  https.get(`https://cognito-idp.${region}.amazonaws.com/${cognitoPoolId}/.well-known/jwks.json`, (resp) => {
    let data = '';
    resp.on('data', (chunk) => {
      data += chunk;
    });
    resp.on('end', () => {
      try {
        const keys = JSON.parse(data);
        resolve(keys);
      } catch (error) {
        reject(error)
      }

    });

  }).on("error", (err) => {
    reject(err)
  });
});

const template = ({ keys, region, cognitoPoolId, endpoint, cfdomain }) => `
const keys = ${JSON.stringify(keys)};
const region = "${region}";
const cognitoPoolId = "${cognitoPoolId}";
const endpoint = "${endpoint}";
const cfdomain = "${cfdomain}";
module.exports = {
  keys,
  region,
  cognitoPoolId,
  endpoint,
  cfdomain
}`;


(async () => {
  try {
    const args = process.argv.slice(2);
    const region = args[0];
    const cognitoPoolId = args[1];
    const endpoint = args[2];
    const cfdomain = args[3];
    const keys = await getKeys(region, cognitoPoolId);
    fs.writeFileSync(path.resolve(__dirname, '../lambdaEdge/settings.js'), template({ keys, region, cognitoPoolId, endpoint, cfdomain }));
    console.log(keys);
  } catch (error) {
    console.log(error)
  }
})();

