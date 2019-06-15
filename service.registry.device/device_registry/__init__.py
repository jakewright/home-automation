import markdown
import os
import shelve
import sys

# Import the framework
from flask import Flask, g, Markup, render_template, request
from flask_restful import Resource, Api, reqparse

from bootstrap import Service

service = Service()
app = service.app


def get_db():
    db = getattr(g, '_database', None)
    if db is None:
        db = g._database = shelve.open(service.app.config['DATABASE'])
    return db


def get_device_repo():
    repo = getattr(g, '_device_repository', None)
    if repo is None:
        repo = g._device_repository = Repository(get_db(), 'device')
    return repo


def get_room_repo():
    repo = getattr(g, '_room_repository', None)
    if repo is None:
        repo = g._room_repository = Repository(get_db(), 'room')
    return repo


@service.app.teardown_appcontext
def close_connection(exception):
    db = getattr(g, '_database', None)
    if db is not None:
        db.close()


class DeviceList(Resource):
    def get(self):
        repo = get_device_repo()
        controller_name = request.args.get('controller_name')

        if controller_name is not None:
            devices = repo.find_by('controller_name', controller_name)
        else:
            devices = repo.find_all()

        for device in devices:
            decorate_with_room(get_room_repo(), device)

        return {'message': 'Retrieved all devices', 'data': devices}, 200

    def post(self):
        device_repo = get_device_repo()
        room_repo = get_room_repo()

        parser = reqparse.RequestParser()

        parser.add_argument('identifier', required=True, location='json')
        parser.add_argument('name', required=True, location='json')
        parser.add_argument('type', required=True, location='json')
        parser.add_argument('controller_name', required=True, location='json')
        parser.add_argument('room_identifier', required=True, location='json')

        parser.add_argument(
            'attributes', type=dict, required=False, location='json')
        parser.add_argument(
            'depends_on', type=list, required=False, location='json')
        parser.add_argument(
            'state_providers', type=list, required=False, location='json')

        # Parse the arguments into an object
        args = parser.parse_args()

        # Make sure the room exists
        room = room_repo.find(args['room_identifier'])
        if room is None:
            return {'message': 'Room does not exist'}, 400

        # Save the device to the data store
        device_repo.save(args)

        # Get the newly saved device
        device = device_repo.find(args['identifier'])

        # Decorate with the room
        decorate_with_room(room_repo, device)

        print('Registered device: ' + args['identifier'], file=sys.stderr)
        return {'message': 'Device registered', 'data': device}, 201


class Device(Resource):
    def get(self, identifier: str):
        device = get_device_repo().find(identifier)

        if device is None:
            return {'message': 'Device not found'}, 404

        decorate_with_room(get_room_repo(), device)

        return {'message': 'Device found', 'data': device}, 200

    def delete(self, identifier: str):
        device_repo = get_device_repo()

        # If the key does not exist in the database
        device = device_repo.find(identifier)

        if device is None:
            return {'message': 'Device not found'}, 404

        device_repo.delete(identifier)

        return '', 204


class RoomList(Resource):
    def get(self):
        rooms = get_room_repo().find_all()

        for room in rooms:
            decorate_with_devices(get_device_repo(), room)

        return {'message': 'Retrieved all rooms', 'data': rooms}, 200

    def post(self):
        device_repo = get_device_repo()
        room_repo = get_room_repo()

        parser = reqparse.RequestParser()

        parser.add_argument('identifier', required=True)
        parser.add_argument('name', required=True)

        # Parse the arguments into an object
        args = parser.parse_args()

        # Save the room to the data store
        room_repo.save(args)

        # Get the newly saved room
        room = room_repo.find(args['identifier'])

        # Decorate with devices
        decorate_with_devices(device_repo, room)

        print('Registered room: ' + args['identifier'], file=sys.stderr)
        return {'message': 'Room registered', 'data': room}, 201


class Room(Resource):
    def get(self, identifier: str):
        room = get_room_repo().find(identifier)

        if room is None:
            return {'message': 'Room not found'}, 404

        decorate_with_devices(get_device_repo(), room)

        return {'message': 'Room found', 'data': room}, 200

    def delete(self, identifier: str):
        device_repo = get_device_repo()
        room_repo = get_room_repo()

        # If the room does not exist in the database
        room = room_repo.find(identifier)
        if room is None:
            return {'message': 'Room not found'}, 404

        # If the room has devices
        devices = device_repo.find_by('room_identifier', identifier)
        if devices:
            return {'message': 'Cannot delete a room that has devices'}, 400

        room_repo.delete(identifier)

        return '', 204


service.api.add_resource(DeviceList, '/devices')
service.api.add_resource(Device, '/device/<string:identifier>')
service.api.add_resource(RoomList, '/rooms')
service.api.add_resource(Room, '/room/<string:identifier>')


# Database access
class Repository:
    def __init__(self, shelf, prefix: str):
        self.shelf = shelf
        self.prefix = prefix + ':'

    def find_all(self):
        """Return an array of all objects in the data store."""

        keys = list(self.shelf.keys())
        objects = []

        for key in keys:
            if key.startswith(self.prefix):
                objects.append(self.shelf[key])

        return objects

    def find(self, identifier: str):
        """Return the object with the given identifier or None if it cannot be found."""

        identifier = self.prefix + identifier

        if not (identifier in self.shelf):
            return None

        return self.shelf[identifier]

    def find_by(self, key: str, value: str):
        """Return an array of objects that match the given condition."""

        objects = []

        for obj in self.find_all():
            if obj[key] == value:
                objects.append(obj)

        return objects

    def save(self, obj):
        self.shelf[self.prefix + obj['identifier']] = obj

    def delete(self, identifier: str):
        del self.shelf[self.prefix + identifier]


# Decorators
def decorate_with_room(room_repository, device):
    # Decorate the device with the room
    device['room'] = room_repository.find(device['room_identifier'])

    # Delete the room identifier from the device
    del device['room_identifier']


def decorate_with_devices(device_repository, room):
    # Find the devices that belong to this room
    devices = device_repository.find_by('room_identifier', room['identifier'])

    # Remove the room_identifier key from each device
    for device in devices:
        del device['room_identifier']

    # Decorate the room with its devices
    room['devices'] = devices
