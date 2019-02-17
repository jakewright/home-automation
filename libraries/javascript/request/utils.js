const _ = require("lodash");

/** Recursively convert an object's keys to camel case */
const toCamelCase = input => {
  if (_.isArray(input)) {
    return input.map(toCamelCase);
  }

  if (!_.isPlainObject(input)) {
    return input;
  }

  const result = {};

  _.forEach(input, (value, key) => {
    const newKey = _.camelCase(key);
    result[newKey] = toCamelCase(value);
  });

  return result;
};

/** Recursively convert an object's keys to snake case */
const toSnakeCase = input => {
  if (_.isArray(input)) {
    return input.map(toSnakeCase);
  }

  if (!_.isPlainObject(input)) {
    return input;
  }

  const result = {};

  _.forEach(input, (value, key) => {
    const newKey = _.snakeCase(key);
    result[newKey] = toSnakeCase(value);
  });

  return result;
};

exports = module.exports = { toCamelCase, toSnakeCase };
