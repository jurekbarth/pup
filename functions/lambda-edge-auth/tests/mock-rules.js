const rules = {
  '/login': {
    "/**/*": {
      "groups": ["public"],
    },
  },
  '/aa/bb': {
    "/**/*": {
      "groups": ["superuser"],
    },
  }
};

module.exports = rules;
