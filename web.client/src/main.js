import Vue from "vue";
import Vuex from "vuex";

import { library } from "@fortawesome/fontawesome-svg-core";
import { faLightbulb, faSpinnerThird } from "@fortawesome/pro-solid-svg-icons";
import { faHome as farHome, faArrowLeft as farArrowLeft } from "@fortawesome/pro-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

import httpClient from "../../libraries/javascript/http";

import App from "./App.vue";
import store from "./store";
import router from "./router";
import EventConsumer from "./api/EventConsumer";

httpClient.setApiGateway(process.env.VUE_APP_API_GATEWAY);

library.add(faLightbulb, faSpinnerThird, farArrowLeft, farHome);
Vue.component("FontAwesomeIcon", FontAwesomeIcon);

Vue.use(Vuex);
Vue.config.productionTip = false;

new Vue({
  render: h => h(App),
  store,
  router
}).$mount("#app");

const eventBusUrl =
  process.env.NODE_ENV === "production"
    ? "ws://192.168.1.100:7004"
    : "ws://localhost:7004";
const eventConsumer = new EventConsumer(eventBusUrl, store);
eventConsumer.listen();
