/**
 * @param {string} rgb hexadecimal RGB string e.g. #FBEE13
 *
 * @return {array} x and y values
 */
const rgbHexToXy = rgb => {
  return rgbToXy(
    parseInt(rgb.substring(1, 3), 16),
    parseInt(rgb.substring(3, 5), 16),
    parseInt(rgb.substring(5, 7), 16)
  );
};

/**
 * https://developers.meethue.com/documentation/color-conversions-rgb-xy
 *
 * @param {number} r Red 0-255 (inclusive)
 * @param {number} g Green 0-255 (inclusive)
 * @param {number} b Blue 0-255 (inclusive)
 *
 * @return {array} x and y values
 */
const rgbToXy = (r, g, b) => {
  // The algorithm below can't handle #000000 so just set it to white instead
  if (r === 0 && g === 0 && b === 0) {
    r = 255;
    g = 255;
    b = 255;
  }

  r = r > 0.04045 ? Math.pow((r + 0.055) / 1.055, 2.4) : r / 12.92;
  g = g > 0.04045 ? Math.pow((g + 0.055) / 1.055, 2.4) : g / 12.92;
  b = b > 0.04045 ? Math.pow((b + 0.055) / 1.055, 2.4) : b / 12.92;

  const X = r * 0.664511 + g * 0.154324 + b * 0.162028;
  const Y = r * 0.283881 + g * 0.668433 + b * 0.047685;
  const Z = r * 0.000088 + g * 0.07231 + b * 0.986039;

  let x = X / (X + Y + Z);
  let y = Y / (X + Y + Z);

  if (isNaN(x)) x = 0;
  if (isNaN(y)) y = 0;

  return [x, y];
};

/**
 * @param {number} 0-1 (inclusive)
 * @param {number} 0-1 (inclusive)
 *
 * @return {string} Hexadecimal RGB string
 */
const xyToRgbHex = (x, y) => {
  const rgb = xyToRgb(x, y);
  return `#${rgb.map(x => x.toString(16).padStart(2, "0")).join("")}`;
};

/**
 * https://developers.meethue.com/documentation/color-conversions-rgb-xy
 *
 * @param {number} x 0-1 (inclusive)
 * @param {number} y 0-1 (inclusive)
 *
 * @return {array} red, green and blue values
 */
const xyToRgb = (x, y) => {
  const z = 1.0 - x - y;
  const Y = 1; // The given brightness value
  const X = (Y / y) * x;
  const Z = (Y / y) * z;

  let r = X * 1.656492 - Y * 0.354851 - Z * 0.255038;
  let g = -X * 0.707196 + Y * 1.655397 + Z * 0.036152;
  let b = X * 0.051713 - Y * 0.121364 + Z * 1.01153;

  if (r > b && r > g && r > 1) {
    // red is too big
    g /= r;
    b /= r;
    r = 1;
  } else if (g > b && g > r && g > 1) {
    // green is too big
    r /= g;
    b /= g;
    g = 1;
  } else if (b > r && b > g && b > 1) {
    // blue is too big
    r /= b;
    g /= b;
    b = 1;
  }

  // Apply gamma correction
  r = r <= 0.0031308 ? 12.92 * r : 1.055 * Math.pow(r, 1.0 / 2.4) - 0.055;
  g = g <= 0.0031308 ? 12.92 * g : 1.055 * Math.pow(g, 1.0 / 2.4) - 0.055;
  b = b <= 0.0031308 ? 12.92 * b : 1.055 * Math.pow(b, 1.0 / 2.4) - 0.055;

  if (r > b && r > g) {
    // red is biggest
    if (r > 1) {
      g /= r;
      b /= r;
      r = 1;
    }
  } else if (g > b && g > r) {
    // green is biggest
    if (g > 1.0) {
      r /= g;
      b /= g;
      g = 1.0;
    }
  } else if (b > r && b > g) {
    // blue is biggest
    if (b > 1.0) {
      r /= b;
      g /= b;
      b = 1.0;
    }
  }

  r = Math.min(r, 1);
  r = Math.max(r, 0);
  g = Math.min(g, 1);
  g = Math.max(g, 0);
  b = Math.min(b, 1);
  b = Math.max(b, 0);

  return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
};

exports = module.exports = { rgbHexToXy, xyToRgbHex };
