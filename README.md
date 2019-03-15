# Home Automation

[![CircleCI](https://circleci.com/gh/jakewright/home-automation/tree/master.svg?style=svg)](https://circleci.com/gh/jakewright/home-automation/tree/master)

Distributed home automation system written in a variety of languages. Largely a learning opportunity rather than a production-ready system.


## API Specification

All responses will have the form:

```json
{
    "data": "Mixed type holding the content of the response"
}
```

Individual service's READMEs will only detail the expected value of the `data` field.

## Errors

An error will be indicated by a non-2xx status code. The response will include a message.

```json
{
    "message": "Description of what went wrong"
}
```

### Controllers

Controllers must implement a standardised interface for fetching and updating device state.

`GET service.controller.x/device/<device-identifier>`

- 200: success

```json
{
    "identifier": "table-lamp",
    "name": "Table Lamp",
    "type": "light",
    "controller_name": "service.controller.hue",
    "availableProperties": {
        "brightness": {
            "type": "int",
            "min": 0,
            "max": 254,
            "interpolation": "continuous",
        }
    },
    "brightness": 100
}
```
