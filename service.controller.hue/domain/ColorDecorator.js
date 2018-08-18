const _ = require("lodash");

const HUE_MIN = 0;
const HUE_MAX = 65536;
const SAT_MIN = 0;
const SAT_MAX = 254;

const ColorDecorator = hueLight => {
  hueLight.color = { hue: 0, saturation: 0 };

  hueLight.setColor = ({ hue, saturation }) => {
    if (_.isEqual(hueLight.color, { hue, saturation })) return;

    if (hue < HUE_MIN || hue > HUE_MAX) throw new Error(`Invalid hue '${hue}'`);
    if (saturation < SAT_MIN || saturation > SAT_MAX)
      throw new Error(`Invalid saturation '${sat}'`);

    hueLight.color = { hue, saturation };
    hueLight.power = true;
  };

  const setState = hueLight.setState;
  hueLight.setState = state => {
    setState.call(hueLight, state);
    if ("color" in state) hueLight.setColor(state.color);
  };

  const prepareLight = hueLight.prepareLight;
  hueLight.prepareLight = () => {
    prepareLight.call(hueLight);

    // Don't try to set color if the light is off
    if (!hueLight.power) return;

    if (hueLight.cache.hue !== hueLight.color.hue) {
      hueLight.light.hue = hueLight.color.hue;
    }

    if (hueLight.cache.saturation !== hueLight.color.saturation) {
      hueLight.light.saturation = hueLight.color.saturation;
    }
  };

  const applyRemoteState = hueLight.applyRemoteState;
  hueLight.applyRemoteState = light => {
    hueLight.setColor({
      hue: light.hue,
      saturation: light.saturation
    });

    applyRemoteState.call(hueLight, light);
  };

  const toJSON = hueLight.toJSON;
  hueLight.toJSON = () => {
    let json = toJSON.call(hueLight);
    json["color"] = hueLight.color;
    return json;
  };

  const getProperties = hueLight.getProperties;
  hueLight.getProperties = () => {
    let properties = getProperties.call(hueLight);
    properties["color"] = { type: "color" };
    return properties;
  };
};

exports = module.exports = ColorDecorator;
