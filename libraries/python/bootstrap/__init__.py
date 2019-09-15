import markdown
import os
import requests
from flask import Flask, g, Markup, render_template
from flask_restful import Api, reqparse, Resource


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

        # Connect to Redis
        if 'REDIS' in self.app.config:
            import redis
            self.r = redis.StrictRedis(
                host=self.app.config['REDIS']['host'],
                port=self.app.config['REDIS']['port'])

        # Create the API
        self.api = Api(self.app)

    def present_documentation(self):
        """Present some documentation."""

        # Get the path of the running script (not this script)
        path = os.path.abspath(os.getcwd())

        # Open the README file
        with open(path + '/README.md', 'r') as markdown_file:
            # Read the markdown contents
            content = markdown_file.read()

            # Convert the markdown to HTML and then treat it as actual HTML so it's not escaped
            html = Markup(
                markdown.markdown(
                    content, extensions=['markdown.extensions.fenced_code']))

        return render_template('index.html', content=html)

    def publish(self, channel, data):
        self.r.publish(channel, data)


class ApiClient:
    def __init__(self, service):
        self.api_gateway = service.app.config['API_GATEWAY']

    def get_devices(self, controller_name):
        r = requests.get(
            self.api_gateway + '/service.device-registry/devices',
            params={
                'controller_name': controller_name,
            })

        if r.status_code != 200:
            raise Exception(
                'service.device-registry returned status code of {}, expected 200'.
                format(r.status_code))

        return r.json()['data']
