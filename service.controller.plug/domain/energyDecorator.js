const energyDecorator = {
  getProperties(properties) {
    properties.watts = { type: "int", immutable: true };
    return properties;
  }
};

exports = module.exports = energyDecorator;
