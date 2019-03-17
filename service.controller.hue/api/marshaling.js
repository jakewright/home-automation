const conversions = require("./conversions");

const fromDomain = domain => {
  const api = {};

  if ("power" in domain) {
    api.on = Boolean(domain.power);
  }

  if ("brightness" in domain) {
    api.bri = domain.brightness;
  }

  if ("color" in domain) {
    api.hue = domain.color.hue;
    api.sat = domain.color.saturation;
  }

  if ("colorTemp" in domain) {
    // Convert from Kelvin to Mirek (Huejay wants a value between 153 and 500)
    api.ct = Math.floor(1000000 / domain.colorTemp);
  }

  if ("rgb" in domain) {
    api.xy = conversions.rgbHexToXy(domain.rgb);
  }

  return api;
};

const toDomain = api => {
  const domain = {};

  if (api.state.on !== undefined) {
    domain.power = api.state.on;
  }

  if (api.state.bri !== undefined) {
    domain.brightness = api.state.bri;
  }

  if (api.state.hue !== undefined && api.state.saturation !== undefined) {
    domain.color = { hue: api.state.hue, saturation: api.state.saturation };
  }

  if (api.state.ct !== undefined) {
    domain.colorTemp = Math.ceil(1000000 / api.state.ct);
  }

  if (api.state.xy !== undefined) {
    domain.rgb = conversions.xyToRgbHex(api.state.xy[0], api.state.xy[1]);
  }

  return domain;
};

exports = module.exports = { fromDomain, toDomain };
