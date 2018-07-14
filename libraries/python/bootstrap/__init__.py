import markdown
import os
import requests
import sys
from flask import Flask, g, Markup, render_template
from flask_restful import Resource, Api, reqparse


class Service:
    def __init__(self):
        # Create the application with the instance config option on
        self.app = Flask(__name__, instance_relative_config=True)

        # Load the default configuration
        self.app.config.from_object('config.default')

        # Load the file specified by the APP_CONFIG_FILE environment variable
        # Variables defined here will override those in the default configuration
        self.app.config.from_envvar('APP_CONFIG_FILE')

        self.app.add_url_rule('/', 'index', self.present_documentation)

        # Create the API
        self.api = Api(self.app)

    def present_documentation(self):
        """Present some documentation."""

        # Get the path of the running script (not this script)
        path = os.path.abspath(os.path.dirname(sys.argv[0]))

        # Open the README file
        with open(path + '/README.md', 'r') as markdown_file:
            # Read the markdown contents
            content = markdown_file.read()

            # Convert the markdown to HTML and then treat it as actual HTML so it's not escaped
            html = Markup(markdown.markdown(content, extensions=['markdown.extensions.fenced_code']))

        return render_template('index.html', content=html)


class ApiClient:
    def __init__(self, service):
        self.service = service
        self.api_gateway = service.app.config['API_GATEWAY']

    def register_room(self, identifier, name):
        r = requests.post(self.api_gateway + '/service.registry.device/rooms', data = {
            'identifier': identifier,
            'name': name,
        })

        if r.status_code != 201:
            raise Exception('Status code was not 201')

    def register_device(self, identifier, name, device_type, controller_name, room_identifier):
        r = requests.post(self.api_gateway + '/service.registry.device/devices', data = {
            'identifier': identifier,
            'name': name,
            'device_type': device_type,
            'controller_name': controller_name,
            'room_identifier': room_identifier,
        })

        if r.status_code != 201:
            raise Exception('Status code was not 201')
