const edgeAuth = require('../index');

const rules = require('./mock-rules');
const cfmocks = require('./mock-cf-request');
const settings = require('./mock-settings');
const expiredToken = require('./mock-expired-token');
const validToken = require('./mock-token');

test('create response header', () => {
  const uri = 'https://google.com'
  expect(edgeAuth.redirectResponse(uri)).toEqual({
    status: '302',
    statusDescription: 'Found',
    headers: {
      location: [{
        key: 'Location',
        value: 'https://google.com',
      }],
    },
  });
});


// test('check empty jwt token', async () => {
//   expect.assertions(1);
//   const data = await edgeAuth.verifyToken();
//   expect(data).toEqual({
//     auth: false,
//     message: 'no token',
//     groups: [],
//   });
// });

// test('check empty string jwt token', async () => {
//   expect.assertions(1);
//   const data = await edgeAuth.verifyToken('');
//   expect(data).toEqual({
//     auth: false,
//     message: 'no token',
//     groups: [],
//   });
// });

test('key pem conversion lazy', () => {
  const { keys } = settings;
  const pems = edgeAuth.setupPem(keys);
  const tof = typeof pems
  expect(tof).toBe('object');
});


test('checkAccessGroups with loginRequest', () => {
  const groups = ['public'];
  const data = edgeAuth.checkAccessGroups(cfmocks.loginRequest, rules, groups);
  expect(data).toEqual(['public']);
});

test('checkAccessGroups with loginRequest', () => {
  const groups = ['public'];
  const data = edgeAuth.checkAccessGroups(cfmocks.loginRequest, rules, groups);
  expect(data).toEqual(['public']);
});

test('checkAccessGroups with userRequest', () => {
  const groups = ['public', 'superuser'];
  const data = edgeAuth.checkAccessGroups(cfmocks.userRequest, rules, groups);
  expect(data).toEqual(['superuser']);
});

test('getCookie with default name', () => {
  const data = edgeAuth.getAuthCookie(cfmocks.userRequest);
  expect(data).toBe('authtoken');
});


test('verify token', async () => {
  expect.assertions(1);
  const pems = edgeAuth.setupPem(settings.keys);
  const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
  const token = validToken;
  const data = await edgeAuth.verifyToken(token, pems, iss);
  expect(data).toEqual({ "auth": true, "groups": ["superuser"], "message": "user authenticated" });
});

test('verify expired token', async () => {
  expect.assertions(1);
  const pems = edgeAuth.setupPem(settings.keys);
  const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
  const token = expiredToken;
  const data = await edgeAuth.verifyToken(token, pems, iss);
  expect(data.auth).toBe(false);
});




test('verify token: private site', async () => {
  expect.assertions(1);
  const pems = edgeAuth.setupPem(settings.keys);
  const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
  const token = validToken;
  const cfEvent = cfmocks.userRequest;
  const tokenGroups = await edgeAuth.verifyToken(token, pems, iss);
  const groups = tokenGroups.groups;
  groups.push('public');
  const data = edgeAuth.checkAccessGroups(cfEvent, rules, groups);
  expect(data).toEqual(["superuser"]);
});

test('verify token: private site expired', async () => {
  expect.assertions(1);
  const pems = edgeAuth.setupPem(settings.keys);
  const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
  const token = expiredToken;
  const cfEvent = cfmocks.userRequest;
  const tokenGroups = await edgeAuth.verifyToken(token, pems, iss);
  const groups = tokenGroups.groups;
  groups.push('public');
  const data = edgeAuth.checkAccessGroups(cfEvent, rules, groups);
  expect(data).toEqual([]);
});


test('verify token: public site', async () => {
  expect.assertions(1);
  const pems = edgeAuth.setupPem(settings.keys);
  const iss = `https://cognito-idp.${settings.region}.amazonaws.com/${settings.cognitoPoolId}`;
  const token = expiredToken;
  const cfEvent = cfmocks.loginRequest;
  const tokenGroups = await edgeAuth.verifyToken(token, pems, iss);
  const groups = tokenGroups.groups;
  groups.push('public');
  const data = edgeAuth.checkAccessGroups(cfEvent, rules, groups);
  expect(data).toEqual(["public"]);
});


