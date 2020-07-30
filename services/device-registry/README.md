# service.device-registry

## Usage

### List all devices
**Definition**

`GET /devices`

**Response**

- 200: success

```json
[
    {
        "id": "id1",
        "name": "Device 1",
        "type": "hs100",
        "kind": "switch",
        "controller_name": "controller-1",
        "room": {
            "id": "bedroom",
            "name": "Jake's Bedroom"
        }
    },
    {
        "id": "id2",
        "name": "Device 2",
        "type": "huelight",
        "kind": "lamp",
        "controller_name": "controller-2",
        "room": {
            "id": "kitchen",
            "name": "Kitchen"
        }
    }
]
```

### Register a new device
**Definition**

`POST /devices`

**Arguments**

- `"id":string` a globally unique ID for this device
- `"name":string` a friendly name for the device
- `"type":string` the type of the device as understood by the client e.g. hs100
- `"kind":string` the kind of device e.g. lamp
- `"room_id":string` the globally unique ID of the room
- `"controller_name":string` the name of the device's controller
- `"attributes":object` arbitrary controller-specific information about the device
- `"state_providers":array` names of external services that provide state

If the ID already exists, the existing device will be overwritten.

**Response**

- 400: unknown room
- 201: created successfully

Returns the new device if successful.

```json
{
    "id": "id1",
    "name": "Device 1",
    "type": "hs100",
    "kind": "switch",
    "controller_name": "controller-2",
    "room": {
        "id": "bedroom",
        "name": "Jake's Bedroom"
    }
}
```

### Lookup device details
**Definition**

`GET /device/<id>`

**Response**

- 404: device not found
- 200: success

```json
{
    "id": "id1",
    "name": "Device 1",
    "type": "hs100",
    "kind": "switch",
    "controller_name": "controller-1",
    "room": {
        "id": "bedroom",
        "name": "Jake's Bedroom"
    }
}
```

### Delete a device
**Definition**

`DELETE /device/<id>`

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
        "id": "bedroom",
        "name": "Jake's Bedroom",
        "devices": [
            {
                "id": "lamp1",
                "name": "Lamp",
                "type": "huelight",
                "kind": "lamp",
                "controller_name": "controller-1"
            }
        ]
    },
    {
        "id": "kitchen",
        "name": "Kitchen",
        "devices": [
            {
                "id": "tv2",
                "name": "TV",
                "type": "philips48",
                "kind": "tv",
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

- `"id":string` a globally unique id for the room
- `"name":string` a friendly name for the room

If the id already exists, the existing room will be overwritten.
Devices belonging to an existing room will not be modified.

**Response**

- 201: created successfully

Returns the new room is created successfully.

```json
{
    "id": "bedroom",
    "name": "Jake's Bedroom",
    "devices": []
}
```

### Lookup room details
**Definition**
`GET /room/<id>`

**Response**

- 404: room not found
- 200: success

```json
{
    "id": "bedroom",
    "name": "Jake's Bedroom",
    "devices": [
        {
            "id": "id1",
            "name": "Device 1",
            "type": "hs100",
            "kind": "switch",
            "controller_name": "controller-1"
        }
    ]
}
```

### Delete a room
**Definition**

`DELETE /rooms/<id>`

**Response**

- 404: room not found
- 204: success
