const CT_MIN = 2000;
const CT_MAX = 6536;

const colorTempDecorator = {
  transform(state, t) {
    if ("colorTemp" in state) {
      if (state.colorTemp < CT_MIN || state.colorTemp > CT_MAX)
        throw new Error(`Invalid colour temperature '${state.colorTemp}'`);

      // Convert from Kelvin to Mirek (Huejay wants a valye between 153 and 500)
      t.colorTemp = Math.floor(1000000 / state.colorTemp);
      t.on = true;
    }

    return t;
  },

  applyRemoteState({ colorTemp }) {
    this.colorTemp = Math.ceil(1000000 / colorTemp);
  },

  getProperties(properties) {
    properties.colorTemp = {
      type: "int",
      min: CT_MIN,
      max: CT_MAX,
      interpolation: "continuous"
    };
    return properties;
  }
};

exports = module.exports = colorTempDecorator;
