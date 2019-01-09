const mm = require('minimatch');

const getBase = (uri, options = { baseDepth: 2 }) => {
  let depth = options.baseDepth + 1;
  const uriParts = uri.split('/');
  const length = uriParts.length - 1;
  if (depth >= length && uriParts[length].indexOf('.') != -1) {
    depth = uriParts.length - 1
  }
  return uriParts.slice(0, depth).join("/");
}

const stripBase = (base, uri) => uri.replace(base, '');

const checkIfRules = (rules, base) => rules.hasOwnProperty(base);

const getRules = (rules, base) => {
  return rules[base];
}

const getMatchedRules = (rules, uri) => {
  const arr = [];
  Object.keys(rules).forEach(rule => {
    if (rule.startsWith('!') && mm(uri, rule.substring(1))) {
      arr.push({
        allow: false,
        rule,
        triggers: rules[rule],
      })
    }
    if (!rule.startsWith('!') && mm(uri, rule)) {
      arr.push({
        allow: true,
        rule,
        triggers: rules[rule],
      })
    }
  });
  return arr;
};

const getFirstMatchingRule = (matchedRules, groups) => matchedRules.find(matchedRule => {
  let isGroupMember = false;
  for (let group of groups) {
    if (matchedRule.triggers.groups.indexOf(group) != -1) {
      isGroupMember = true;
      break;
    };
  }
  return isGroupMember
});

const getGroupsForUri = (uri = '', rules = {}, groups = []) => {
  const base = getBase(uri);
  if (!checkIfRules(rules, base)) {
    return [];
  }
  const uriRules = getRules(rules, base);
  const strippedUri = stripBase(base, uri);
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  const firstRule = getFirstMatchingRule(matchingRules, groups);
  if (firstRule === undefined) {
    return [];
  }
  if (firstRule.allow) {
    return firstRule.triggers.groups;
  }
  return [];
}


module.exports = {
  getBase,
  stripBase,
  checkIfRules,
  getRules,
  getMatchedRules,
  getFirstMatchingRule,
  getGroupsForUri
}
