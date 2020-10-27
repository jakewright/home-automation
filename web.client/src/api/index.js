import http from "../../../libraries/javascript/http";
import { apiToDevice, apiToRoom, apiToRooms } from "../domain/marshalling";

/**
 * Fetch device information from the device's controller
 *
 * @param {DeviceHeader} deviceHeader Metadata about the device
 * @return {Device}
 */
const fetchDevice = async deviceHeader => {
  const url = `${deviceHeader.controllerName}/device/${
    deviceHeader.identifier
  }`;
  const rsp = await http.get(url);
  return apiToDevice(rsp);
};

/**
 * Update a single property on a device
 *
 * @param {Object} deviceHeader Object containing device identifier and controller name
 * @param {Object} properties A map of property names to their new values. Properties that are omitted will not be updated.
 *
 * @return {Device} The updated Device object
 */
const updateDevice = async ({ identifier, controllerName }, properties) => {
  const url = `${controllerName}/device/${identifier}`;
  const rsp = await http.patch(url, properties);
  return apiToDevice(rsp);
};

/**
 * Fetch all rooms
 * @returns {Array.<Room>}
 */
const fetchRooms = async () => {
  const rsp = await http.get("device-registry/rooms");
  return apiToRooms(rsp.rooms);
};

/**
 * Fetch a single room by ID
 *
 * @param identifier
 * @returns {Room}
 */
const fetchRoom = async identifier => {
  const rsp = await http.get(`device-registry/room/${identifier}`);
  return apiToRoom(rsp);
};

export default { fetchDevice, updateDevice, fetchRooms, fetchRoom };
