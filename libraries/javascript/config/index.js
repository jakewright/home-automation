class Config {
  setContents(config) {
    this.config = config;
  }

  has(path) {
    const value = this.get(path);
    return value !== undefined;
  }

  get(path, def) {
    // Throw an error if the config hasn't been loaded
    if (this.config === undefined) {
      throw new Error("Config not loaded");
    }

    const reduce = (parts, config) => {
      // If this is the last part of the key
      if (parts.length === 0) {
        // Return the value
        return config;
      }

      // If config is not an object then we can't continue
      if (config == null || typeof config !== "object") {
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

const config = new Config();
exports = module.exports = config;
