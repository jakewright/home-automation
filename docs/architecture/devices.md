# Devices

![](https://media.giphy.com/media/3o6Mb3dvYo1WSH1lw4/source.gif)

Devices are registered with the `device-registry` service. 
Various metadata about the devices are stored here, including which controller controls them. This allows the frontend to discover controllers.

 
Devices have _properties_ and _commands_. 
- A property has a value that can be changed, e.g. `power` can be `true` or `false`.
- A commend has no state but performs an action, e.g. `toggleState()`. A command can take arguments.

### Read a device

`GET http://[device-controller]/[device-type]?device_id=foo`

**Response**

```json
{
    "device": {
        ...
    }
}
```

The specific response received is dependent on the type of the device, and will be described in the controller's def file. 

### Updating a property

`PATCH http://[device-controller]/[device-type]`

**Request**

- `"device_id":string` the globally unique ID for this device
- See the controller's def file for the list of fields that can be set

**Response**

```json
{
    "device": {
        ...
    }
}
```

### Calling a command

`POST service.device-controller/device/cmd`

**Request**
- `"device_id":string` the globally unique ID for this device
- `"command":string` the name of the command to call
- `"args":object` a map of argument names to values

The response from reading the device will define the available commands and their arguments.

**Response**

```json
{}
```

## State providers

Devices do not have dependencies because of the complexity of implementing this in a generic way. This would solve the problem of one device needing another device to be in a particular state before a property can be set, e.g. a light needs the WiFi plug to be on before the brightness can be changed. These kinds of problems can be instead solved by _scenes_.

Devices can, however, have _state providers_. E.g. a WiFi plug might know whether the TV is on or off, and provide that state to the TV's controller which has no way to know on its own.

State providers are only implemented where needed. The state providers' controller names are listed as part of the device's metadata. The device's controller, when fetching state, will hit `/provide-state?device_id=<device-id>` on all of the state providers and merge the resulting state together.
