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

const template = ({ keys, region, cognitoPoolId, endpoint, cfdomain, loginUrl }) => `
const keys = ${JSON.stringify(keys)};
const region = "${region}";
const cognitoPoolId = "${cognitoPoolId}";
const endpoint = "${endpoint}";
const cfdomain = "${cfdomain}";
const loginUrl = "${loginUrl}";
module.exports = {
  keys,
  region,
  cognitoPoolId,
  endpoint,
  cfdomain,
  loginUrl
}`;

// ⚠️ Change here to get something useful

const exchangeThis = {
  region: 'eu-central-1',
  cognitoPoolId: "eu-central-1_xxxx",
  endpoint: "testing",
  cfdomain: "www.testing.com",
  clientId: "SPACLIENTID",
};

(async () => {
  try {
    const { region, cognitoPoolId, endpoint, cfdomain, clientId } = exchangeThis;
    const keys = await getKeys(region, cognitoPoolId);
    const cognitoUrl = `https://${endpoint}.auth.${region}.amazoncognito.com`
    const loginUrl = `${cognitoUrl}/login?response_type=token&client_id=${clientId}&redirect_uri=https://${cfdomain}/login/index.html`
    fs.writeFileSync(path.resolve(__dirname, './mock-settings.js'), template({ keys, region, cognitoPoolId, endpoint, cfdomain, loginUrl }));
  } catch (error) {
    console.log(error)
  }
})();
