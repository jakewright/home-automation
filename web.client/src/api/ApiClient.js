import axios from 'axios';
import {toCamelCase} from './utils';
import {apiToDevice, apiToRoom, apiToRooms} from '../domain/marshalling';

export default class ApiClient {

    /**
     * @param {string} apiGateway
     */
    constructor(apiGateway) {
        this.axios = axios.create({
            baseURL: apiGateway + '/',
            validateStatus: (status) => {
                return status >= 200 && status < 300;
            },
            transformResponse: [function (data) {
                data = JSON.parse(data);

                // If the data field is empty, there was probably an error so just return
                // as-is and handle in the request method below.
                if (!data.data) return data;

                // The response will be enveloped but we no longer care about the
                // extra metadata. Also convert keys to camel case because JavaScript.
                return toCamelCase(data.data);
            }],
        });
    }

    async request(args) {
        try {
            const response = await this.axios.request.apply(this, arguments);
            console.log(`${args.method} ${args.url}`);
            if (args.data) {
                console.log(args.data);
            }
            console.log(response.data);
            return response.data;
        } catch (err) {
            console.error(err);

            // If the backend actually returned a response, extract useful information from it.
            if (err.response) {
                console.error(err.response);

                // The response should include an array of errors. Turn them into a comma-separated list.
                let errors = err.response.data.errors.reduce((s, e) => `${s}, ${e}`);

                // Add all of the extra messages from the response to the error message
                err.message += `: ${err.response.data.message}: errors: [${errors}]`;
            }

            throw err;
        }
    }

    /**
     * Fetch metadata for all devices known to the device registry
     */
    // async fetchDeviceHeaders() {
    //     return this.request({
    //         url: 'device-registry/devices',
    //         method: 'get',
    //     });
    // }

    /**
     * Fetch device information from the device's controller
     *
     * @param {DeviceHeader} deviceHeader Metadata about the device
     * @return {Device}
     */
    async fetchDevice(deviceHeader) {
        const rsp = await this.request({
            url: deviceHeader.controllerName + '/device/' + deviceHeader.identifier,
            method: 'get',
        });

        return apiToDevice(rsp);
    }

    /**
     * Update a single property on a device
     *
     * @param {Object} deviceHeader Object containing device identifier and controller name
     * @param {Object} properties A map of property names to their new values. Properties that are omitted will not be updated.
     *
     * @return {Device} The updated Device object
     */
    async updateDevice({ identifier, controllerName }, properties) {
        const rsp = await this.request({
            url: `${controllerName}/device/${identifier}`,
            method: 'patch',
            data: properties,
        });

        return apiToDevice(rsp);
    }

    /**
     * Fetch all rooms
     * @returns {Array.<Room>}
     */
    async fetchRooms() {
        const rsp = await this.request({
            url: 'service.registry.device/rooms',
            method: 'get',
        });

        return apiToRooms(rsp);
    }

    /**
     * Fetch a single room by ID
     *
     * @param identifier
     * @returns {Room}
     */
    async fetchRoom(identifier) {
        const rsp = await this.request({
            url: 'service.registry.device/room/' + identifier,
            method: 'get',
        });

        return apiToRoom(rsp);
    }
}



