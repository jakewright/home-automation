# Devices

Devices are registered with `service.registry.device`. 
Various metadata about the devices are stored here, including which controller controls them. This allows the frontend to discover controllers.

 
Devices have _properties_ and _commands_. 
- A property has a value that can be changed, e.g. `power` can be `true` or `false`.
- A commend has no state but performs an action, e.g. `toggleState()`. A command can take arguments.

Devices do not have dependencies because of the complexity of implementing this in a generic way. This would solve the problem of one device needing another device to be in a particular state before a property can be set, e.g. a light needs the WiFi plug to be on before the brightness can be changed. These kinds of problems can be instead solved by _scenes_.

Devices can, however, have _state providers_. E.g. a WiFi plug might know whether the TV is on or off, and provide that state to the TV's controller which has no way to know on its own.

State providers are only implemented where needed. The state providers controller names are listed as part of the device's metadata. The device's controller, when fetching state, will hit `/provide-state/<device-id>` on all of the state providers and merge the resulting state together.
