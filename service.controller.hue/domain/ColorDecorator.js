const HUE_MIN = 0;
const HUE_MAX = 65536;
const SAT_MIN = 0;
const SAT_MAX = 254;

const colorDecorator = {
  transform(state, t) {
    if ("color" in state) {
      ({ hue, saturation } = state.color);

      if (hue < HUE_MIN || hue > HUE_MAX)
        throw new Error(`Invalid hue '${hue}'`);

      if (saturation < SAT_MIN || saturation > SAT_MAX)
        throw new Error(`Invalid saturation '${saturation}'`);

      t.hue = hue;
      t.saturation = saturation;
      t.on = true;
    }

    return t;
  },

  applyRemoteState({ hue, saturation }) {
    this.color = { hue, saturation };
  },

  getProperties(properties) {
    properties.color = { type: "color" };
    return properties;
  }
};

exports = module.exports = colorDecorator;
