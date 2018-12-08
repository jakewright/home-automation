export default class Room {
  /**
   * @param {string} identifier
   * @param {string} name
   * @param {Array.<DeviceHeader>} deviceHeaders
   */
  constructor(identifier, name, deviceHeaders) {
    this.identifier = identifier;
    this.name = name;
    this.deviceHeaders = deviceHeaders;
  }
}
