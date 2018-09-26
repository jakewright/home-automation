const conversions = require("./conversions");

const rgbDecorator = {
  transform(state, t) {
    if ("rgb" in state) {
      const ok = /^#[0-9A-F]{6}$/i.test(state.rgb);
      if (!ok) throw new Error(`Invalid hex color '${state.rgb}'`);

      t.xy = conversions.rgbHexToXy(state.rgb);
      t.on = true;
    }

    return t;
  },

  applyRemoteState({ xy }) {
    this.rgb = conversions.xyToRgbHex(xy[0], xy[1]);
  },

  getProperties(properties) {
    properties.rgb = { type: "rgb" };
    return properties;
  }
};

exports = module.exports = rgbDecorator;
