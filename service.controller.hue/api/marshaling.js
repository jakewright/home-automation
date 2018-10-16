const conversions = require("./conversions");

const fromDomain = domain => {
  const api = {};

  if ("power" in domain) {
    api.on = Boolean(domain.power);
  }

  if ("brightness" in domain) {
    api.brightness = domain.brightness;
  }

  if ("color" in domain) {
    api.hue = domain.color.hue;
    api.saturation = domain.color.saturation;
  }

  if ("colorTemp" in domain) {
    // Convert from Kelvin to Mirek (Huejay wants a value between 153 and 500)
    api.colorTemp = Math.floor(1000000 / domain.colorTemp);
  }

  if ("rgb" in domain) {
    api.xy = conversions.rgbHexToXy(domain.rgb);
  }

  return api;
};

const toDomain = api => {
  const domain = {};

  if (api.on !== undefined) {
    domain.power = api.on;
  }

  if (api.brightness !== undefined) {
    domain.brightness = api.brightness;
  }

  if (api.hue !== undefined && api.saturation !== undefined) {
    domain.color = { hue: api.hue, saturation: api.saturation };
  }

  if (api.colorTemp !== undefined) {
    domain.colorTemp = Math.ceil(1000000 / api.colorTemp);
  }

  if (api.xy !== undefined) {
    domain.rgb = conversions.xyToRgbHex(api.xy[0], api.xy[1]);
  }

  return domain;
};

exports = module.exports = { fromDomain, toDomain };
