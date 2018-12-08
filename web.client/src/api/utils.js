import _ from "lodash";

/**
 * Recursively convert an object's property names to camel case using Lodash functions.
 */
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

export { toCamelCase };
