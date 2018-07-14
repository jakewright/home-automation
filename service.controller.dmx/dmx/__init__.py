import array
import re
import sys

from bootstrap import ApiClient, Service
from flask_restful import Resource, reqparse
from ola.ClientWrapper import ClientWrapper

service = Service()
app = service.app
api_client = ApiClient(service)

device_identifier = service.app.config['DEVICE_IDENTIFIER']
device_name = service.app.config['DEVICE_NAME']
controller_name = service.app.config['CONTROLLER_NAME']
api_client.register_room('living-room', 'Living Room')

api_client.register_device(
    device_identifier,
    device_name,
    'dmx',
    controller_name,
    'living-room',
)

# This is awful but store state as global variables
rgb = '#000000'
brightness = 0
strobe = 0

def dmx_sent(state):
    wrapper.Stop()

class Device(Resource):
    def __init__(self, device_identifier, device_name, controller_name):
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
                    'rgb': {'type': 'rgb'},
                    'brightness': {'type': 'int', 'min': 0, 'max': 255},
                    'strobe': {'type': 'int', 'min': 0, 'max': 240},
                },

                'rgb': rgb,
                'brightness': brightness,
                'strobe': strobe,
            }
        }

    def get(self):
        return self.to_json(), 200

    def patch(self):
        global rgb
        global brightness
        global strobe

        parser = reqparse.RequestParser()
        parser.add_argument('rgb', type=unicode)
        parser.add_argument('strobe', type=int)
        parser.add_argument('brightness', type=int)

        # Parse the arguments into an object
        args = parser.parse_args()

        if args['rgb'] is not None:
            if not re.search(r'^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$', args['rgb']):
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
        rgb_parts = tuple(int(hex[i:i+2], 16) for i in (0, 2 ,4))

        data = array.array('B', [
            rgb_parts[0], # red
            rgb_parts[1], # green
            rgb_parts[2], # blue
            0, # color macros
            strobe + 15, # strobing/program speed
            0, # programs
            brightness # master dimmer
        ])

        global wrapper
        wrapper = ClientWrapper()

        client = wrapper.Client()
        client.SendDmx(1, data, dmx_sent)
        wrapper.Run()

        return self.to_json(), 200

service.api.add_resource(Device, '/device/' + device_identifier, resource_class_kwargs={
    'device_identifier': device_identifier,
    'device_name': device_name,
    'controller_name': controller_name,
})
