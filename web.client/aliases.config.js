const path = require("path");

function resolveSrc(_path) {
  return path.join(__dirname, _path);
}

const aliases = {
  "@design": "src/design/index.scss",
  "@variables": "src/design/variables/index.scss"
};

module.exports = {
  webpack: {},
  jest: {}
};

for (const alias in aliases) {
  module.exports.webpack[alias] = resolveSrc(aliases[alias]);
  module.exports.jest[`^${alias}/(.*)$`] = `<rootDir>/${aliases[alias]}/$1`;
}
