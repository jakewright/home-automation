from device_registry import app
import os
import unittest
from tests import add_device, add_room, decode_response


class TestDeviceList(unittest.TestCase):
    def setUp(self):
        # Create an instance of the test client
        self.app = app.test_client()

        # Use a temporary file as the database file
        app.config['DATABASE'] = '/tmp/test'

        add_room('room', 'Room')

    def tearDown(self):
        # Delete the temporary file
        os.unlink(app.config['DATABASE'])

    def test_no_devices(self):
        """Test that an empty array is returned when there are no devices known by the registry."""

        # Make a get request to /devices
        response = self.app.get('/devices')
        decoded_response = decode_response(response)

        # Assert that the status code is correct
        self.assertEqual(200, response.status_code)

        # Assert that no devices were returned
        self.assertEqual([], decoded_response)

    def test_add_device_invalid_room(self):
        """Test that adding a device with an unknown room returns a 400"""

        response = add_device('invalid-room-test', 'name', 'type', 'kind',
                              'controller', 'unknown-room')
        self.assertEqual(400, response.status_code)

    def test_add_device(self):
        """Test that a new device can be added to the registry."""

        device = {
            'identifier': 'device1',
            'name': 'Device 1',
            'type': 'hs100',
            'kind': 'switch',
            'depends_on': None,
            'controller_name': 'controller',
            'room': {
                'identifier': 'room',
                'name': 'Room',
            },
            'state_providers': None,
        }

        # Add a new device and get the response
        response = add_device(
            device['identifier'],
            device['name'],
            device['type'],
            device['kind'],
            device['controller_name'],
            device['room']['identifier'],
        )
        decoded_response = decode_response(response)

        self.assertEqual(201, response.status_code)
        self.assertEqual(device, decoded_response)

    def test_add_duplicate_device(self):
        """Test that adding a device with an identifier that is already used replaces the device."""

        # Add the same device twice
        add_device('device1', 'device 1', 'hs100', 'switch', 'controller', 'room')
        response = add_device('device1', 'device new', 'hs100', 'switch', 'controller',
                              'room')

        # Assert that we get a 201 response code the second time
        self.assertEqual(201, response.status_code)

        response = self.app.get('/device/device1')
        decoded_response = decode_response(response)

        self.assertEqual('device new', decoded_response['name'])
