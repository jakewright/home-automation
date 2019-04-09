const HUE_MIN = 0;
const HUE_MAX = 65536;
const SAT_MIN = 0;
const SAT_MAX = 254;

const colorDecorator = {
  validate(state) {
    if ("color" in state) {
      ({ hue, saturation } = state.color);

      if (hue < HUE_MIN || hue > HUE_MAX) return `Invalid hue '${hue}'`;

      if (saturation < SAT_MIN || saturation > SAT_MAX)
        return `Invalid saturation '${saturation}'`;
    }
  },

  transform(state, t) {
    if ("color" in state) {
      t.color = state.color;
      t.power = true;
    }

    return t;
  },

  state: {
    color: { type: "color" }
  }
};

exports = module.exports = colorDecorator;
