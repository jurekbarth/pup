const auth = require('@jurekbarth/lambda-edge-auth');

const rules = require('./rules.js');
const clientId = require('./clientId');
const settings = require('./settings');

const pems = auth.setupPem(settings.keys.keys);
const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
const cognitoUrl = `https://${settings.endpoint}.auth.${settings.region}.amazoncognito.com`;
const loginUrl = `${cognitoUrl}/login?response_type=token&client_id=${clientId}&redirect_uri=https://${settings.cfdomain}/login/index.html`

const unauthorizedUri = `https://${settings.cfdomain}/login/unauthorized.html`;



exports.handler = async (event, context, callback) => {
  const cfrequest = event.Records[0].cf.request;
  const uri = cfrequest.uri;
  const token = auth.getAuthCookie(event);
  const authInformation = await auth.verifyToken(token, pems, iss);
  console.log('authinformation');
  console.log(authInformation);
  // check if public
  const accessPublic = auth.checkAccessGroups(event, rules, ['public']);
  if (accessPublic.length > 0) {
    callback(null, cfrequest);
  }
  // user not authenticated
  if (!authInformation.auth) {
    const response = auth.redirectResponse(`${loginUrl}&state=${uri}`);
    callback(null, response);
    return;
  }
  const groups = authInformation.groups;
  console.log('authInformation groups')
  console.log(groups);
  groups.push('public');
  const accessGroups = auth.checkAccessGroups(event, rules, groups);
  console.log('accessgroups');
  console.log(accessGroups);
  // user unauthorized
  if (accessGroups.length === 0) {
    const response = auth.redirectResponse(unauthorizedUri);
    callback(null, response);
    return;
  }
  // all good
  callback(null, cfrequest);
  return;
};
