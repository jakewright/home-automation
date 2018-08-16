import array
import re
import sys

from bootstrap import ApiClient, Service
from flask_restful import Resource, reqparse
from ola.ClientWrapper import ClientWrapper


def create_app():
    service = Service()
    controller_name = service.app.config['CONTROLLER_NAME']
    api_client = ApiClient(service)

    devices = api_client.get_devices(controller_name)
    if len(devices) == 0:
        raise Exception('Did not find any devices')

    for device in devices:
        service.api.add_resource(
            Device,
            '/device/' + device['identifier'],
            resource_class_kwargs={
                'service': service,
                'device_identifier': device['identifier'],
                'device_name': device['name'],
                'controller_name': controller_name,
            })

    return service.app


# This is awful but store state as global variables
# @todo: make this work with multiple devices
rgb = '#000000'
brightness = 0
strobe = 0


def dmx_sent(state):
    wrapper.Stop()


class Device(Resource):
    def __init__(self, service, device_identifier, device_name,
                 controller_name):
        self.service = service
        self.identifier = device_identifier
        self.name = device_name
        self.controller_name = controller_name

    def to_json(self):
        return {
            'message': 'Device',
            'data': {
                'identifier': self.identifier,
                'name': self.name,
                'type': 'dmx',
                'controller_name': self.controller_name,
                'available_properties': {
                    'rgb': {
                        'type': 'rgb'
                    },
                    'brightness': {
                        'type': 'int',
                        'min': 0,
                        'max': 255,
                        'interpolation': 'continuous'
                    },
                    'strobe': {
                        'type': 'int',
                        'min': 0,
                        'max': 240,
                        'interpolation': 'continuous'
                    },
                },
                'rgb': rgb,
                'brightness': brightness,
                'strobe': strobe,
            }
        }

    def hash(self):
        global rgb, brightness, strobe
        return rgb + str(brightness) + str(strobe)

    def get(self):
        return self.to_json(), 200

    def patch(self):
        global rgb, brightness, strobe
        cache = self.hash()

        parser = reqparse.RequestParser()
        parser.add_argument('rgb', type=unicode, location='json')
        parser.add_argument('strobe', type=int, location='json')
        parser.add_argument('brightness', type=int, location='json')

        # Parse the arguments into an object
        args = parser.parse_args()

        if args['rgb'] is not None:
            if not re.search(r'^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$',
                             args['rgb']):
                return {'error': 'Invalid hex string'}, 400
            rgb = args['rgb']

        if args['brightness'] is not None:
            if args['brightness'] < 0 or args['brightness'] > 255:
                return {'error': 'Invalid brightness value'}, 400
            brightness = args['brightness']

        if args['strobe'] is not None:
            if args['strobe'] < 0 or args['strobe'] > 240:
                return {'error': 'Invalid strobe value'}, 400
            strobe = args['strobe']

        hex = rgb.lstrip('#')
        rgb_parts = tuple(int(hex[i:i + 2], 16) for i in (0, 2, 4))

        data = array.array(
            'B',
            [
                rgb_parts[0],  # red
                rgb_parts[1],  # green
                rgb_parts[2],  # blue
                0,  # color macros
                strobe + 15,  # strobing/program speed
                0,  # programs
                brightness  # master dimmer
            ])

        global wrapper
        wrapper = ClientWrapper()

        client = wrapper.Client()
        client.SendDmx(1, data, dmx_sent)
        wrapper.Run()

        # If the state has changed, publish a state-change event.
        if cache != self.hash():
            self.service.publish('device-state-changed.{}'.format(self.identifier),
                                 self.to_json())

        return self.to_json(), 200
