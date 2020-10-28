import Vue from "vue";
import Vuex from "vuex";
import { createApp } from "vue";

//import { library } from "@fortawesome/fontawesome-svg-core";
// import { faLightbulb, faSpinnerThird } from "@fortawesome/pro-solid-svg-icons";
// import { faHome as farHome, faArrowLeft as farArrowLeft } from "@fortawesome/pro-regular-svg-icons";
// import { faLightbulb as falLightbulb } from "@fortawesome/pro-light-svg-icons";
//import { fal } from "@fortawesome/pro-light-svg-icons";
//import { far } from "@fortawesome/pro-regular-svg-icons";
//import { fas } from "@fortawesome/pro-solid-svg-icons";
//import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

import httpClient from "../../libraries/javascript/http";

import App from "./App.vue";
import store from "./store";
import router from "./router";
import EventConsumer from "./api/EventConsumer";

httpClient.setApiGateway(process.env.VUE_APP_API_GATEWAY);

// library.add(faLightbulb, falLightbulb, faSpinnerThird, farArrowLeft, farHome);
// library.add(fas, far, fal);
// Vue.component("FontAwesomeIcon", FontAwesomeIcon);

const app = createApp(App);
app.use(router)
app.use(store)
app.mount("#app")


const eventBusUrl =
  process.env.NODE_ENV === "production"
    ? "ws://192.168.1.100:7004"
    : "ws://localhost:7004";
const eventConsumer = new EventConsumer(eventBusUrl, store);
eventConsumer.listen();
