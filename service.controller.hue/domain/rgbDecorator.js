const rgbDecorator = {
  validate(state) {
    if ("rgb" in state) {
      const ok = /^#[0-9A-F]{6}$/i.test(state.rgb);
      if (!ok) return `Invalid hex color '${state.rgb}'`;
    }
  },

  transform(state, t) {
    if ("rgb" in state) {
      t.rgb = state.rgb;
      t.power = true;
    }
    return t;
  },

  getProperties(properties) {
    properties.rgb = { type: "rgb" };
    return properties;
  }
};

exports = module.exports = rgbDecorator;
