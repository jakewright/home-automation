import requests_mock
import json
from device_registry import app


def add_device(identifier, name, device_type, controller_name,
               room_identifier):
    """Helper function to post a new device to the registry and return the response."""

    # Create a test client
    with app.test_client() as c:

        return c.post(
            '/devices',
            data={
                'identifier': identifier,
                'name': name,
                'device_type': device_type,
                'controller_name': controller_name,
                'room_identifier': room_identifier,
            })


def add_room(identifier, name):
    """Helper function to post a new room to the registry and return the response."""

    # Create a test client
    with app.test_client() as c:

        return c.post(
            '/rooms', data={
                'identifier': identifier,
                'name': name,
            })


def decode_response(response):
    return json.loads(response.data.decode('utf-8'))['data']
