import array
import re
import sys

from bootstrap import ApiClient, Service
from flask_restful import Resource, reqparse
from ola.ClientWrapper import ClientWrapper

service = Service()
api_client = ApiClient(service)

api_client.register_room('living-room', 'Living Room')

api_client.register_device(
    service.app.config['DEVICE_IDENTIFIER'],
    service.app.config['DEVICE_NAME'],
    'dmx',
    service.app.config['CONTROLLER_NAME'],
    'living-room',
)

def dmx_sent(state):
    wrapper.Stop()

class Device(Resource):
    def __init__(self):
        self.rgb = '#000000'
        self.brightness = 0
        self.strobe = 0

    def to_json(self):
        return {
            'message': 'Device',
            'data': {
                'identifier': 'dmx',
                'name': 'Sofa lights',
                'type': 'dmx',
                'available_properties': {
                    'rgb': {'type': 'rgb'},
                    'brightness': {'type': 'int', 'min': 0, 'max': 255},
                    'strobe': {'type': 'int', 'min': 0, 'max': 240},
                },

                'rgb': self.rgb,
                'brightness': self.brightness,
                'strobe': self.strobe,
            }
        }

    def get(self):
        return self.to_json(), 200

    def patch(self):
        parser = reqparse.RequestParser()
        parser.add_argument('rgb', type=unicode)
        parser.add_argument('strobe', type=int)
        parser.add_argument('brightness', type=int)

        # Parse the arguments into an object
        args = parser.parse_args()

        if args['rgb']:
            if not re.search(r'^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$', args['rgb']):
                return {'error': 'Invalid hex string'}, 400
            self.rgb = args['rgb']

        if args['brightness']:
            if args['brightness'] < 0 or args['brightness'] > 255:
                return {'error': 'Invalid brightness value'}, 400
            self.brightness = args['brightness']

        if args['strobe']:
            if args['strobe'] < 0 or args['strobe'] > 240:
                return {'error': 'Invalid strobe value'}, 400
            self.strobe = args['strobe']


        hex = self.rgb.lstrip('#')
        rgb = tuple(int(hex[i:i+2], 16) for i in (0, 2 ,4))

        data = array.array('B', [
            rgb[0], # red
            rgb[1], # green
            rgb[2], # blue
            0, # color macros
            self.strobe + 15, # strobing/program speed
            0, # programs
            self.brightness # master dimmer
        ])

        global wrapper
        wrapper = ClientWrapper()

        client = wrapper.Client()
        client.SendDmx(1, data, dmx_sent)
        wrapper.Run()

        return self.to_json(), 200

service.api.add_resource(Device, '/device/dmx')
