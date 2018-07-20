import _ from 'lodash';

export default class Device {
  constructor(identifier, name, deviceType, controllerName, availableProperties, properties) {
    this.identifier = identifier;
    this.name = name;
    this.deviceType = deviceType;
    this.controllerName = controllerName;
    this.availableProperties = availableProperties;
    this._properties = properties;
  }

  /**
     * Combine properties and available properties into a single object
     *      brightness: {
     *          value: 60,
     *          min: 0,
     *          max: 100,
     *          interpolation: "continuous",
     *      }
     */
  get properties() {
    return _.mapValues(this.availableProperties, (property, propertyName) => {
      property.value = this._properties[propertyName];
      return property;
    });
  }
}
