const CT_MIN = 2000;
const CT_MAX = 6536;

const colorTempDecorator = {
  validate(state) {
    if ("colorTemp" in state) {
      if (state.colorTemp < CT_MIN || state.colorTemp > CT_MAX)
        return `Invalid colour temperature '${state.colorTemp}'`;
    }
  },

  transform(state, t) {
    if ("colorTemp" in state) {
      t.colorTemp = state.colorTemp;
      t.power = true;
    }

    return t;
  },

  state: {
    colorTemp: {
      prettyName: "colour temperature",
      type: "int",
      min: CT_MIN,
      max: CT_MAX,
      interpolation: "continuous"
    }
  }
};

exports = module.exports = colorTempDecorator;
