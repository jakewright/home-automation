const req = require("../http");

/**
 * @param {Object} state
 * @param {Object[]}dependencies
 * @param {string} dependencies[].localProperty
 * @param {string} dependencies[].localValue
 * @param {string} dependencies[].remoteDeviceIdentifier
 * @param {string} dependencies[].remoteProperty
 * @param {string} dependencies[].remoteValue
 *
 * @return {Promise}
 */
const updateDependencies = (state, dependencies) => {
  // Filter the dependencies such that only ones that need updating are left
  dependencies = dependencies.filter(
    d => state[d.localProperty] === d.localValue
  );

  // Execute all requests in parallel
  const promises = dependencies.map(updateDependency);
  return Promise.all(promises);
};

const updateDependency = async dependency => {
  const id = dependency.remoteDeviceIdentifier;
  return req.get(`service.registry.device/device/${id}`).then(header => {
    const data = { [dependency.remoteProperty]: dependency.remoteValue };
    return req.patch(`${header.controllerName}/device/${id}`, data);
  });
};

exports = module.exports = updateDependencies;
