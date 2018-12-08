import createLogger from "vuex/dist/logger";
import Vue from "vue";
import Vuex from "vuex";
import devices from "./modules/devices";
import errors from "./modules/errors";
import rooms from "./modules/rooms";

Vue.use(Vuex);

const debug = true;

export default new Vuex.Store({
  modules: {
    devices,
    errors,
    rooms
  },
  strict: debug,
  plugins: debug ? [createLogger()] : []
});
