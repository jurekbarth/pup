const {
  getBase,
  stripBase,
  getRules,
  checkIfRules,
  getMatchedRules,
  getFirstMatchingRule,
  getGroupsForUri
} = require('./index');

const rules = require('./mock');

test('getBase from /a/b/c', () => {
  const uri = '/a/b/c'
  expect(getBase(uri)).toBe('/a/b');
});

test('getBase from /a/b/c/d', () => {
  const uri = '/a/b/c/d'
  expect(getBase(uri, { baseDepth: 3 })).toBe('/a/b/c');
});

test('getBase from /a', () => {
  const uri = '/a'
  expect(getBase(uri)).toBe('/a');
});

test('getBase from /a/index.html', () => {
  const uri = '/a/index.html'
  expect(getBase(uri)).toBe('/a');
});

test('getBase from /a/b/index.html', () => {
  const uri = '/a/b/index.html'
  expect(getBase(uri)).toBe('/a/b');
});

test('stripBase base: from /a/b/c', () => {
  const uri = '/a/b/c'
  const base = getBase(uri);
  expect(stripBase(base, uri)).toBe('/c');
});

test('stripBase from: /a/b/c/d', () => {
  const uri = '/a/b/c/d'
  const base = getBase(uri);
  expect(stripBase(base, uri)).toBe('/c/d');
});

test('stripBase from: /a/index.html', () => {
  const uri = '/a/index.html'
  const base = getBase(uri);
  expect(stripBase(base, uri)).toBe('/index.html');
});

test('checkIfRules for: /a/b/c/d', () => {
  const uri = '/a/b/c/d'
  const base = getBase(uri);
  expect(checkIfRules(rules, base)).toBe(true);
});

test('getRules for: /a/b/c/d', () => {
  const uri = '/a/b/c/d'
  const base = getBase(uri);
  expect(getRules(rules, base)).toBe(rules[base]);
});

test('getMatchedRules for: /a/b/c/d', () => {
  const uri = '/a/b/c/d'
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const shouldBe = [{ "allow": true, "rule": "/**/*", "triggers": { "groups": ["p1--client-view", "dev"] } }];
  expect(getMatchedRules(uriRules, strippedUri)).toEqual(shouldBe)
});

test('getMatchedRules for: /a/b/c/d/index.html', () => {
  const uri = '/a/b/c/d/index.html'
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri)
  const uriRules = getRules(rules, base); // returns an object of rules
  const shouldBe = [{ "allow": false, "rule": "!/**/index.html", "triggers": { "groups": ["p1--client-view-c-level"] } }, { "allow": true, "rule": "/**/*", "triggers": { "groups": ["p1--client-view", "dev"] } }];
  expect(getMatchedRules(uriRules, strippedUri)).toEqual(shouldBe)
});

test('getMatchedRules for: /login/b/index.html', () => {
  const uri = '/login/b/index.html'
  const base = getBase(uri, { baseDepth: 1 });
  const strippedUri = stripBase(base, uri)
  const uriRules = getRules(rules, base); // returns an object of rules
  const shouldBe = [{ "allow": true, "rule": "/**/*", "triggers": { "groups": ["public"] } }];
  expect(getMatchedRules(uriRules, strippedUri)).toEqual(shouldBe)
});

test('getFirstMatchingRule for: /a/b/c/d/index.html group: other-group', () => {
  const uri = '/a/b/c/d/index.html'
  const groups = ['other-group'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toBeUndefined()
});

test('getFirstMatchingRule for: /a/b/c/d/index.html group: p1--client-view', () => {
  const uri = '/a/b/c/d/index.html'
  const groups = ['p1--client-view'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toEqual({ "allow": true, "rule": "/**/*", "triggers": { "groups": ["p1--client-view", "dev"] } })
});

test('getFirstMatchingRule for: /a/b/master/d/index.html group: p1--client-view-c-level', () => {
  const uri = '/a/b/master/d/index.html'
  const groups = ['p1--client-view-c-level'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toEqual({ "allow": true, "rule": "/master/**", "triggers": { "groups": ["p1--client-view-c-level"] } })
});

test('getFirstMatchingRule for: /a/b/c/d/index.html group: p1--client-view-c-level, p1--client-view', () => {
  const uri = '/a/b/c/d/index.html'
  const groups = ['p1--client-view-c-level', 'p1--client-view'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toEqual({ "allow": false, "rule": "!/**/index.html", "triggers": { "groups": ["p1--client-view-c-level"] } })
});

test('getFirstMatchingRule for: /a/b/c/d/somepage.html group: p1--client-view-c-level, p1--client-view', () => {
  const uri = '/a/b/c/d/somepage.html'
  const groups = ['p1--client-view-c-level', 'p1--client-view'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toEqual({ "allow": true, "rule": "/**/*", "triggers": { "groups": ["p1--client-view", "dev"] } })
});

test('getFirstMatchingRule for: /a/b/master/d/index.html group: p1--client-view-c-level, p1--client-view', () => {
  const uri = '/a/b/master/d/index.html'
  const groups = ['p1--client-view-c-level', 'p1--client-view'];
  const base = getBase(uri); // /a/b
  const strippedUri = stripBase(base, uri);
  const uriRules = getRules(rules, base); // returns an object of rules
  const matchingRules = getMatchedRules(uriRules, strippedUri);
  expect(getFirstMatchingRule(matchingRules, groups)).toEqual({ "allow": true, "rule": "/master/**", "triggers": { "groups": ["p1--client-view-c-level"] } })
});

test('getGroupsForUri for: /a/b/master/d/index.html group: p1--client-view-c-level, p1--client-view', () => {
  const uri = '/a/b/master/d/index.html'
  const groups = ['p1--client-view-c-level', 'p1--client-view'];
  expect(getGroupsForUri(uri, rules, groups)).toEqual(["p1--client-view-c-level"])
});

test('getGroupsForUri for: /login/index.html group: public', () => {
  const uri = '/login/index.html'
  const groups = ['public'];
  expect(getGroupsForUri(uri, rules, groups)).toEqual(["public"])
});

test('getGroupsForUri for: /login/somethingElse.html group: somegroup, public', () => {
  const uri = '/login/somethingElse.html'
  const groups = ['somegroup', 'public'];
  expect(getGroupsForUri(uri, rules, groups)).toEqual(["public"])
});

test('getGroupsForUri for: /login/a/b.html group: public', () => {
  const uri = '/login/a/b.html'
  const groups = ['public'];
  // this should not work, because the base is '/login/a'
  // improvement would be to build an ast of rules, so we won't have a fixed base
  expect(getGroupsForUri(uri, rules, groups)).toEqual([])
});

test('getGroupsForUri for: /a/b/c/d/index.html group: p1--client-view-c-level, p1--client-view', () => {
  const uri = '/a/b/c/d/index.html'
  const groups = ['p1--client-view-c-level', 'p1--client-view'];
  expect(getGroupsForUri(uri, rules, groups)).toEqual([])
});
