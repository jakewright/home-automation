import Device from "./Device";
import _ from "lodash";
import Room from "./Room";
import DeviceHeader from "./DeviceHeader";

/**
 * Converts a JSON API response into a domain object
 *
 * @param {Object} rsp
 * @returns {Device}
 */
const apiToDevice = rsp => {
  return new Device(
    rsp.identifier,
    rsp.name,
    rsp.type,
    rsp.controllerName,
    rsp.state
  );
};

/**
 * Converts a JSON API response into a domain object
 * @param {Object} rsp
 * @returns {Room}
 */
const apiToRoom = rsp => {
  const deviceHeaders = rsp.devices.map(
    d => new DeviceHeader(d.identifier, d.name, d.deviceType, d.controllerName)
  );
  return new Room(rsp.identifier, rsp.name, deviceHeaders);
};

/**
 * Converts a JSON API response into an array of domain objects
 * @param {Array} rsp
 * @returns {Array.<Room>}
 */
const apiToRooms = rsp => rsp.map(r => apiToRoom(r));

export { apiToDevice, apiToRoom, apiToRooms };
