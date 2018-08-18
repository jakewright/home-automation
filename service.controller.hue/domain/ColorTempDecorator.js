const CT_MIN = 2000;
const CT_MAX = 6536;

const ColorTempDecorator = hueLight => {
  hueLight.colorTemp = 2700;

  hueLight.setColorTemp = temp => {
    if (hueLight.colorTemp === temp) return;

    // Despite what the docs say, huejay sometimes returns numbers higher than 6500
    if (temp < CT_MIN || temp > CT_MAX) {
      throw new Error(`Invalid colour temperature "${temp}"`);
    }

    hueLight.colorTemp = temp;
    hueLight.power = true;
  };

  const setState = hueLight.setState;
  hueLight.setState = state => {
    setState.call(hueLight, state);
    if ("colorTemp" in state) hueLight.setColorTemp(state.colorTemp);
  };

  const prepareLight = hueLight.prepareLight;
  hueLight.prepareLight = () => {
    prepareLight.call(hueLight);

    // Don't try to set colorTemp if the light is off
    if (!hueLight.power) return;

    // Don't try to set colorTemp if the value hasn't changed
    if (hueLight.cache.colorTemp === hueLight.colorTemp) return;

    // Convert from Kelvin to Mirek (Huejay wants a valye between 153 and 500)
    hueLight.light.colorTemp = Math.floor(1000000 / hueLight.colorTemp);
  };

  const applyRemoteState = hueLight.applyRemoteState;
  hueLight.applyRemoteState = light => {
    hueLight.setColorTemp(Math.ceil(1000000 / light.colorTemp));
    applyRemoteState.call(hueLight, light);
  };

  const getProperties = hueLight.getProperties;
  hueLight.getProperties = () => {
    let properties = getProperties.call(hueLight);
    properties["colorTemp"] = {
      type: "int",
      min: CT_MIN,
      max: CT_MAX,
      interpolation: "continuous"
    };
    return properties;
  };

  const toJSON = hueLight.toJSON;
  hueLight.toJSON = () => {
    let json = toJSON.call(hueLight);
    json["colorTemp"] = hueLight.colorTemp;
    return json;
  };
};

exports = module.exports = ColorTempDecorator;
