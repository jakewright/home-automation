class Config {
  constructor(config) {
    this.config = config;
  }

  has(key) {
    const value = this.get(key);
    return value !== undefined
  }

  get(path, def, config = this.config) {
    const reduce = (parts, config) => {
      // If this is the last part of the key
      if (parts.length === 0) {
        // Return the value
        return config;
      }

      // If config is not an object then we can't continue
      if (config == null || typeof config !== 'object') {
        // Return the default
        return def;
      }

      // Take the first part of the path
      const key = parts.shift();

      // If the key we are searching for is not defined
      if (!(key in config)) {
        // Return the default
        return def;
      }

      // Recurse
      return reduce(parts, config[key]);
    };

    return reduce(path.split("."), this.config);
  }
}

exports = module.exports = Config;
