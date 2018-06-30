import array
import re

from bootstrap import Service
from flask_restful import Resource, reqparse
from ola.ClientWrapper import ClientWrapper

service = Service()

wrapper = ClientWrapper()
client = wrapper.Client()

@service.app.teardown_appcontext
def close_wrapper(exception):
    wrapper.Stop()

class Device(Resource):
    def __init__(self):
        self.rgb = '#000000'
        self.brightness = 0
        self.strobe = 0

    def get(self):
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
        }, 200

    def patch(self):
        parser = reqparse.RequestParser()
        parser.add_argument('rgb', type=unicode)
        parser.add_argument('strobe', type=int)
        parser.add_argument('brightness', type=int)

        # Parse the arguments into an object
        args = parser.parse_args()

        print args['rgb']

        if not re.search(r'^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$', args['rgb']):
            return {'error': 'Invalid hex string'}, 400

        if args['brightness'] < 0 or args['brightness'] > 255:
            return {'error': 'Invalid brightness value'}, 400

        if args['strobe'] < 0 or args['strobe'] > 240:
            return {'error': 'Invalid strobe value'}, 400

        self.rgb = args['rgb']
        self.brightness = args['brightness']
        self.strobe = args['strobe']

        hex = args['rgb'].lstrip('#')
        rgb = tuple(int(hex[i:i+2], 16) for i in (0, 2 ,4))

        data = array.array('B', [
            rgb[0], # red
            rgb[1], # green
            rgb[2], # blue
            0, # color macros
            args['strobe'] + 15, # strobing/program speed
            0, # programs
            args['brightness'] # master dimmer
        ])

        client.SendDmx(1, data)
        wrapper.Run()

        return self.get()

service.api.add_resource(Device, '/device/dmx')