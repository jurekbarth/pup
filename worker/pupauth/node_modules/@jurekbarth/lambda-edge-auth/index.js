const jwt = require('jsonwebtoken');
const jwkToPem = require('jwk-to-pem');
const pathmatcher = require('@jurekbarth/pathmatcher');


const setupPem = (keys) => {
  const pems = {};
  for (const key of keys) {
    //Convert each key to PEM
    const { kid, n, e, kty } = key;
    const jwk = { kty, n, e };
    const pem = jwkToPem(jwk);
    pems[kid] = pem;
  }
  return pems;
}

const verifyToken = async (jwtToken, pems, iss) => {
  const decodedJwt = jwt.decode(jwtToken, { complete: true });
  if (!decodedJwt) {
    return {
      auth: false,
      message: 'no token',
      groups: [],
    };
  }
  if (decodedJwt.payload.iss != iss) {
    console.log('issuer wrong')
    return {
      auth: false,
      message: 'issuer wrong',
      groups: [],
    }
  }
  const kid = decodedJwt.header.kid;
  const pem = pems[kid];
  if (!pem) {
    return {
      auth: false,
      message: 'wrong key',
      groups: [],
    }
  }

  const v = new Promise((resolve, _) => {
    jwt.verify(jwtToken, pem, { issuer: iss }, function (err, _) {
      if (err) {
        resolve({
          auth: false,
          message: err,
          groups: [],
        });
      } else if (!decodedJwt.payload.hasOwnProperty('cognito:groups')) {
        resolve({
          auth: false,
          message: 'user not part of any group',
          groups: [],
        });
      } else {
        const groups = decodedJwt.payload['cognito:groups'];
        resolve({
          auth: true,
          message: 'user authenticated',
          groups,
        });
      }
    });
  })
  const data = await v;
  return data;
}

const getAuthCookie = (event, cookieName = 'pupauthcookie') => {
  const cfrequest = event.Records[0].cf.request;
  const headers = cfrequest.headers;
  let cookie = '';
  if (headers.cookie) {
    for (const c of headers.cookie) {
      if (c.value.indexOf(`${cookieName}=`) > -1) {
        let str = c.value;
        let idx = str.indexOf(`${cookieName}=`);
        str = str.slice(idx);
        idx = str.indexOf(";");
        if (idx !== -1) {
          str = str.slice(0, idx);
        }
        cookie = str.replace(`${cookieName}=`, "");
        break;
      }
    }
  }
  return cookie;
}


const checkAccessGroups = (event, rules, groups) => {
  const cfrequest = event.Records[0].cf.request;
  const uri = cfrequest.uri;
  return pathmatcher.getGroupsForUri(uri, rules, groups);
}

const redirectResponse = (uri) => ({
  status: '302',
  statusDescription: 'Found',
  headers: {
    location: [{
      key: 'Location',
      value: uri,
    }],
  },
})

module.exports = {
  setupPem,
  verifyToken,
  getAuthCookie,
  checkAccessGroups,
  redirectResponse
}
