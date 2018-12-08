module.exports = {
  configureWebpack: {
    resolve: {
      alias: require("./aliases.config").webpack
    }
  }
};
