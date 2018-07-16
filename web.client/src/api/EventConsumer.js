import {apiToDevice} from '../domain/marshalling';

export default class EventConsumer {
    constructor(url, store) {
        this.url = url;
        this.store = store;
    }

    listen() {
        console.log("Connecting to web socket...");
        this.socket = new WebSocket(this.url);
        this.socket.onmessage = event => {

            try {
                const data = JSON.parse(event.data);
                const [ eventType ] = data.channel.split('.');

                switch (eventType) {
                    case 'state-change':
                        this.handleStateChangeEvent(data.message);
                }
            } catch (err) {
                // Ignore events that are not JSON encoded
            }
        }
    }

    handleStateChangeEvent(msg) {
        this.store.commit('setDevice', apiToDevice(msg));
    }
}
