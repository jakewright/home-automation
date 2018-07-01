from device_registry import app
import os
import unittest
from tests import add_device, add_room, decode_response

class TestRoomList(unittest.TestCase):

    def setUp(self):
        # Create an instance of the test client
        self.app = app.test_client();

        # Use a temporary file as the database file
        app.config['DATABASE'] = '/tmp/test'

    def tearDown(self):
        # Delete the temporary file
        os.unlink(app.config['DATABASE'] + '.db')

    def test_no_rooms(self):
        """Test that an empty array is returned when there are no rooms known by the registry."""

        response = self.app.get('/rooms')
        decoded_response = decode_response(response)

        # Assert that the status code is correct
        self.assertEqual(200, response.status_code)


        # Assert that no rooms were returned
        self.assertEqual([], decoded_response)

    def test_add_room(self):
        """Test that a new room can be added to the registry."""

        room = {
            'identifier': 'bedroom',
            'name': "Jake's Bedroom",
            'devices': [],
        }

        # Add a new room and get the response
        response = add_room(room['identifier'], room['name'])
        decoded_response = decode_response(response)

        self.assertEqual(201, response.status_code)
        self.assertEqual(room, decoded_response)

    def test_add_duplicate_room(self):
        """Test that adding a room with an identifier that is already used replaces the room."""

        # Add the same room twice
        add_room('room1', "room")
        response = add_room('room1', "room new")

        # Assert that we get a 201 response code the second time
        self.assertEqual(201, response.status_code)

        response = self.app.get('/room/room1')
        decoded_response = decode_response(response)

        self.assertEqual('room new', decoded_response['name'])
