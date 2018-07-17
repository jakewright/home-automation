import unittest
import os
from device_registry import app
from tests import add_device, add_room, decode_response


class TestDevice(unittest.TestCase):
    def setUp(self):
        # Create an instance of the test client
        self.app = app.test_client()

        # Use a temporary file as the database file
        app.config['DATABASE'] = '/tmp/test'

    def tearDown(self):
        # Delete the temporary file
        os.unlink(app.config['DATABASE'])

    def test_get_invalid_device(self):
        """Test that getting an invalid device's details will return a 404."""

        response = self.app.get('/device/1')
        self.assertEqual(404, response.status_code)

    def test_get_device(self):
        """Test that adding and then getting a device works."""

        add_room('bedroom', "Jake's Bedroom")

        device = {
            'identifier': 'test',
            'name': 'Test',
            'device_type': 'switch',
            'controller_name': 'controller-1',
            'room': {
                'identifier': 'bedroom',
                'name': "Jake's Bedroom",
            },
        }

        # Add the device to the registry
        add_device(device['identifier'], device['name'], device['device_type'],
                   device['controller_name'], device['room']['identifier'])

        # Ask the registry for the device's details
        response = self.app.get('/device/test')

        # Assert that we got an OK response code
        self.assertEqual(200, response.status_code)

        # Decode the json
        decoded_response = decode_response(response)

        self.assertEqual(device, decoded_response)

    def test_delete_device(self):
        """Test that a device is no longer available after deleting it"""

        # Add a device
        add_room('bedroom', "Jake's Bedroom")
        add_device('test-del', 'name', 'type', 'controller', 'bedroom')

        # Try to get the device
        response = self.app.get('/device/test-del')
        self.assertEqual(200, response.status_code)

        # Delete the device
        response = self.app.delete('/device/test-del')
        self.assertEqual(204, response.status_code)

        # Try to get the device again
        response = self.app.get('/device/test-del')
        self.assertEqual(404, response.status_code)

        # Try to delete the device again
        response = self.app.delete('/device/test-del')
        self.assertEqual(404, response.status_code)
