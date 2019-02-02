#!/usr/local/bin/python3

import argparse
import configparser
import os
import subprocess

def docker(service, status, config):
    script = dir_path + '/docker/status.sh' if status else dir_path + '/docker/deploy.sh'
    return subprocess.run(script, shell=True, env={
        'SERVICE': service,
        'DEPLOYMENT_TARGET': config['DeploymentTarget'],
        'TARGET_USERNAME': config['TargetUsername'],
        'TARGET_DIRECTORY': config['TargetDirectory'],
    }).returncode


def systemd(service, status, config):
    script = dir_path + '/systemd/status.sh' if status else dir_path + '/systemd/deploy.sh'

    if config['Language'] != 'python':
        print('Unsupported language')
        return 1

    return subprocess.run(script, shell=True, env={
        "SERVICE": service,
        "DEPLOYMENT_TARGET": config['DeploymentTarget'],
        "TARGET_USERNAME": config['TargetUsername'],
        "TARGET_DIRECTORY": config['TargetDirectory'],
        "SYSTEMD_SERVICE": get_python_systemd(service, config),
        "LANG": config['Language']
    }).returncode


def get_python_systemd(service, config):
    project_root = config['TargetDirectory'] + '/src'
    service_root =  project_root + '/' + service
    return '''\
[Unit]
Description={service}

[Service]
Environment=APP_CONFIG_FILE={service_root}/config/production.py
Environment=PYTHONPATH=$PYTHONPATH:{project_root}/libraries/python:/usr/lib/python2.7/dist-packages
Environment=FLASK_APP={service_root}/{flask_app_name}
Environment=FLASK_ENV=production
Environment=FLASK_RUN_HOST=0.0.0.0
Environment=FLASK_RUN_PORT={port}
Type=idle
ExecStart={service_root}/env/bin/flask run

[Install]
WantedBy=multi-user.target
'''.format(service=service, project_root=project_root, service_root=service_root, flask_app_name=config['FlaskAppName'], port=config['Port'])


def run():
    # Parse arguments
    parser = argparse.ArgumentParser(description='Deploy services')
    parser.add_argument('service', help='The name of the service to deploy')
    parser.add_argument('--status', help='Get the currently deployed revision', action='store_true')
    args = parser.parse_args()

    # Read config
    config = configparser.ConfigParser(interpolation=configparser.ExtendedInterpolation())
    config.read(dir_path + '/config.ini')

    if args.service not in config:
        print("Service not found")
        return 1

    if 'System' not in config[args.service]:
        print("Deployment system not defined for service")
        return 1

    system = config[args.service]['System']

    if system == 'docker':
        return docker(args.service, args.status, config[args.service])

    if system == 'systemd':
        return systemd(args.service, args.status, config[args.service])

    print("Deployment system not supported")
    return 1


dir_path = os.path.dirname(os.path.realpath(__file__))
exit(run())