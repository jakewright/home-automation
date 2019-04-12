export default class Device {
  constructor(identifier, name, deviceType, controllerName, state) {
    this.identifier = identifier;
    this.name = name;
    this.deviceType = deviceType;
    this.controllerName = controllerName;
    this.state = state;
  }
}
