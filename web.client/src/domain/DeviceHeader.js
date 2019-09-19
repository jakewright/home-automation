export default class DeviceHeader {
  constructor(identifier, name, type, kind, controllerName) {
    this.identifier = identifier;
    this.name = name;
    this.type = type;
    this.kind = kind;
    this.controllerName = controllerName;
  }
}
