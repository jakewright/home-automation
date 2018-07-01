import markdown
import os
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

    def serve(self):
        self.app.run(host='0.0.0.0', port=int(self.app.config['PORT']), debug=True)

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
