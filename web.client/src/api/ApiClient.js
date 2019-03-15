import Client from "../../../libraries/javascript/request/index";
import { apiToDevice, apiToRoom, apiToRooms } from "../domain/marshalling";

export default class ApiClient {
  constructor(apiGateway) {
    this.client = new Client(apiGateway);
  }

  /**
   * Fetch device information from the device's controller
   *
   * @param {DeviceHeader} deviceHeader Metadata about the device
   * @return {Device}
   */
  async fetchDevice(deviceHeader) {
    const url = `${deviceHeader.controllerName}/device/${
      deviceHeader.identifier
    }`;
    const rsp = await this.client.get(url);
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
    const url = `${controllerName}/device/${identifier}`;
    const rsp = await this.client.patch(url, properties);
    return apiToDevice(rsp);
  }

  /**
   * Fetch all rooms
   * @returns {Array.<Room>}
   */
  async fetchRooms() {
    const rsp = await this.client.get("service.registry.device/rooms");
    return apiToRooms(rsp);
  }

  /**
   * Fetch a single room by ID
   *
   * @param identifier
   * @returns {Room}
   */
  async fetchRoom(identifier) {
    const rsp = await this.client.get(
      `service.registry.device/room/${identifier}`
    );
    return apiToRoom(rsp);
  }
}
