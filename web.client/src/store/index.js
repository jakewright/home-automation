import { createLogger, createStore } from "vuex";
import devices from "./modules/devices";
import errors from "./modules/errors";
import rooms from "./modules/rooms";

const debug = true;

const store = createStore({
  modules: {
    devices,
    errors,
    rooms
  },
  strict: debug,
  plugins: debug ? [createLogger()]: []
});

export default store
