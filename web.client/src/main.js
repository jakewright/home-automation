import ApiClient from './api/ApiClient';
import App from './App.vue';
import store from './store';
import Vue from 'vue';
import Vuex from 'vuex';
import router from './router';
import EventConsumer from './api/EventConsumer';

import { library } from '@fortawesome/fontawesome-svg-core';
import { faLightbulb, faSpinnerThird } from '@fortawesome/pro-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

library.add(faLightbulb, faSpinnerThird);
Vue.component('FontAwesomeIcon', FontAwesomeIcon);

Vue.use(Vuex);

Vue.config.productionTip = false;

const apiGateway = (process.env.NODE_ENV === 'production') ? 'http://192.168.1.210:5005' : 'http://localhost:5005';
export const apiClient = new ApiClient(apiGateway);

new Vue({
    render: h => h(App),
    store,
    router,
}).$mount('#app');

const eventBusUrl = (process.env.NODE_ENV === 'production') ? 'ws://192.168.1.210:5004' : 'ws://localhost:5004';
const eventConsumer = new EventConsumer(eventBusUrl, store);
eventConsumer.listen();
