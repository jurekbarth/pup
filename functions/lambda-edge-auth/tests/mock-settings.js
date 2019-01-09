
const keys = {};
const region = "eu-central-1";
const cognitoPoolId = "eu-central-1_xxxxxxx";
const endpoint = "testing";
const cfdomain = "www.testing.com";
const loginUrl = "https://testing.auth.eu-central-1.amazoncognito.com/login?response_type=token&client_id=SPACLIENTf&redirect_uri=https://www.testing.com/login/index.html";

module.exports = {
  keys,
  region,
  cognitoPoolId,
  endpoint,
  cfdomain,
  loginUrl
}
