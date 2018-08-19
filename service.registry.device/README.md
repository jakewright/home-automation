# Home Automation Device Registry

## Usage
All responses will have the form:

```json
{
    "message": "Description of what happened",
    "data": "Mixed type holding the content of the response"
}
```

Subsequent response definitions will only detail the expected value of the `data` field.

### List all devices
**Definition**

`GET /devices`

**Response**

- 200: success

```json
[
    {
        "identifier": "id1",
        "name": "Device 1",
        "device_type": "switch",
        "controller_name": "controller-1",
        "room": {
            "identifier": "bedroom",
            "name": "Jake's Bedroom"
        }
    },
    {
        "identifier": "id2",
        "name": "Device 2",
        "device_type": "bulb",
        "controller_name": "controller-2",
        "room": {
            "identifier": "kitchen",
            "name": "Kitchen"
        }
    }
]
```

### Register a new device
**Definition**

`POST /devices`

**Arguments**

- `"identifier":string` a globally unique identifier for this device
- `"name":string` a friendly name for the device
- `"device_type":string` the type of the device as understood by the client
- `"room_identifier":string` the globally unique identifier of the room
- `"controller_name":string` the name of the device's controller
- `"attributes":object` controller specific information about the device
- `"depends_on":array` an array of dependencies for this device
    - `"local_property":string` the name of the local property that has a dependency
    - `"local_value":string` the value of the local property that has a dependency
    - `"remote_device":string` the identifier of the device this property depends on
    - `"remote_property":string` the remote property that must be set
    - `"remote_value":string` the value of the remote property

If the identifier already exists, the existing device will be overwritten.

**Response**

- 400: unknown room
- 201: created successfully

Returns the new device if successful.

```json
{
    "identifier": "id1",
    "name": "Device 1",
    "device_type": "switch",
    "controller_name": "controller-2",
    "room": {
        "identifier": "bedroom",
        "name": "Jake's Bedroom"
    }
}
```

### Lookup device details
**Definition**

`GET /device/<identifier>`

**Response**

- 404: device not found
- 200: success

```json
{
    "identifier": "id1",
    "name": "Device 1",
    "device_type": "switch",
    "controller_name": "controller-1",
    "room": {
        "identifier": "bedroom",
        "name": "Jake's Bedroom"
    }
}
```

### Delete a device
**Definition**

`DELETE /device/<identifier>`

**Response**

- 404: device not found
- 204: success

### List rooms
**Definition**

`GET /rooms`

**Response**

- 200: success

```json
[
    {
        "identifier": "bedroom",
        "name": "Jake's Bedroom",
        "devices": [
            {
                "identifier": "lamp1",
                "name": "Lamp",
                "device_type": "bulb",
                "controller_name": "controller-1"
            }
        ]
    },
    {
        "identifier": "kitchen",
        "name": "Kitchen",
        "devices": [
            {
                "identifier": "tv2",
                "name": "TV",
                "device_type": "switch",
                "controller_name": "controller-2"
            }
        ]
    }
]

```

### Register new room
**Definition**

`POST /rooms`

**Arguments**

- `"identifier":string` a globally unique identifier for the room
- `"name":string` a friendly name for the room

If the identifier already exists, the existing room will be overwritten.
Devices belonging to an existing room will not be modified.

**Response**

- 201: created successfully

Returns the new room is created successfully.

```json
{
    "identifier": "bedroom",
    "name": "Jake's Bedroom",
    "devices": []
}
```

### Lookup room details
**Definition**
`GET /room/<identifier>`

**Response**

- 404: room not found
- 200: success

```json
{
    "identifier": "bedroom",
    "name": "Jake's Bedroom",
    "devices": [
        {
            "identifier": "id1",
            "name": "Device 1",
            "device_type": "switch",
            "controller_name": "controller-1"
        }
    ]
}
```

### Delete a room
**Definition**

`DELETE /rooms/<identifier>`

**Response**

- 404: room not found
- 204: success

