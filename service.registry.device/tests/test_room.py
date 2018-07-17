import unittest
import os
from device_registry import app
import json
from tests import add_device, add_room, decode_response


class TestRoom(unittest.TestCase):
    def setUp(self):
        # Create an instance of the test client
        self.app = app.test_client()

        # Use a temporary file as the database file
        app.config['DATABASE'] = '/tmp/test'

    def tearDown(self):
        # Delete the temporary file
        os.unlink(app.config['DATABASE'])

    def test_get_invalid_room(self):
        """Test that getting an invalid room's details will return a 404."""

        response = self.app.get('/room/unknown')
        self.assertEqual(404, response.status_code)

    def test_get_room(self):
        """Test that getting a room with devices works."""

        device1 = {
            'identifier': 'device1',
            'name': 'Test Device 1',
            'device_type': 'switch',
            'controller_name': 'controller-1',
        }

        device2 = {
            'identifier': 'device2',
            'name': 'Test Device 2',
            'device_type': 'switch',
            'controller_name': 'controller-1',
        }

        room = {
            'identifier': 'bedroom',
            'name': "Jake's Bedroom",
            'devices': [device1, device2],
        }

        add_room(room['identifier'], room['name'])
        add_device(
            device1['identifier'],
            device1['name'],
            device1['device_type'],
            device1['controller_name'],
            room['identifier'],
        )

        add_device(
            device2['identifier'],
            device2['name'],
            device2['device_type'],
            device2['controller_name'],
            room['identifier'],
        )

        response = self.app.get('/room/bedroom')
        decoded_response = decode_response(response)

        self.assertEqual(room, decoded_response)

    def test_delete_room(self):
        """Test that a room is not longer available after deleting it."""

        # Add a room
        add_room('test-del', 'name')

        # Try to get the room
        response = self.app.get('/room/test-del')
        self.assertEqual(200, response.status_code)

        # Delete the room
        response = self.app.delete('/room/test-del')
        self.assertEqual(204, response.status_code)

        # Try to get the room again
        response = self.app.get('/room/test-del')
        self.assertEqual(404, response.status_code)

        # Try to delete the room again
        response = self.app.delete('/room/test-del')
        self.assertEqual(404, response.status_code)
