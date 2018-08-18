const conversions = require("./conversions");

const RgbDecorator = hueLight => {
  hueLight.rgb = "#000000";

  const setState = hueLight.setState;
  hueLight.setState = state => {
    setState.call(hueLight, state);
    if ("rgb" in state) hueLight.setRgb(state.rgb);
  };

  const prepareLight = hueLight.prepareLight;
  hueLight.prepareLight = () => {
    prepareLight.call(hueLight);

    // Don't try to set xy if the light is off
    if (!hueLight.power) return;

    // Don't try to set xy if the RGB value hasn't changed
    if (hueLight.cache.rgb === hueLight.rgb) return;

    hueLight.light.xy = conversions.rgbHexToXy(hueLight.rgb);

  };

  const applyRemoteState = hueLight.applyRemoteState;
  hueLight.applyRemoteState = light => {
    const rgb = conversions.xyToRgbHex(light.xy[0], light.xy[1]);
    hueLight.setRgb(rgb);

    applyRemoteState.call(hueLight, light);
  };

  hueLight.setRgb = rgb => {
    if (hueLight.rgb === rgb) return;

    const ok = /^#[0-9A-F]{6}$/i.test(rgb);
    if (!ok) throw new Error(`Invalid hex color '${rgb}'`);

    hueLight.rgb = rgb;
    hueLight.power = true;
  };

  const toJSON = hueLight.toJSON;
  hueLight.toJSON = () => {
    let json = toJSON.call(hueLight);
    json["rgb"] = hueLight.rgb;
    return json;
  };

  const getProperties = hueLight.getProperties;
  hueLight.getProperties = () => {
    let properties = getProperties.call(hueLight);
    properties["rgb"] = { type: "rgb" };
    return properties;
  };
};

exports = module.exports = RgbDecorator;
