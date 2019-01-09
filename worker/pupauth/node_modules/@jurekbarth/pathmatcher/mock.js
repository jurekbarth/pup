const rules = {
  '/login': {
    "/**/*": {
      "groups": ["public"],
    },
  },
  '/a/b': {
    "/master/**": {
      "groups": [
        "p1--client-view-c-level"
      ]
    },
    "!/**/index.html": {
      "groups": [
        "p1--client-view-c-level"
      ]
    },
    "/*/resources/js/**/*.private.js": {
      "groups": [
        "dev"
      ]
    },
    "/*/resources/js/**": {
      "groups": [
        "public"
      ]
    },
    "/**/*": {
      "groups": [
        "p1--client-view",
        "dev"
      ]
    },
  },
};

module.exports = rules;
