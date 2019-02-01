#!/usr/local/bin/python3

import argparse
import configparser
import os
import subprocess

# Parse arguments
parser = argparse.ArgumentParser(description='Deploy services')
parser.add_argument('service', help='The name of the service to deploy')
parser.add_argument('--status', help='Get the currently deployed revision', action='store_true')
args = parser.parse_args()

# Read config
dir_path = os.path.dirname(os.path.realpath(__file__))
config = configparser.ConfigParser(interpolation=configparser.ExtendedInterpolation())
config.read(dir_path + '/config.ini')

if args.service not in config:
    print("Service not found")
    exit(1)

if 'System' not in config[args.service]:
    print("Deployment system not defined for service")
    exit(1)

if config[args.service]['System'] != 'docker':
    print("Deployment system not supported")
    exit(1)

if args.status:
    subprocess.run(dir_path + "/docker/status.sh", shell=True, env={
        "SERVICE": args.service,
        "DEPLOYMENT_TARGET": config[args.service]['DeploymentTarget'],
        "TARGET_USERNAME": config[args.service]['TargetUsername'],
        "TARGET_DIRECTORY": config[args.service]['TargetDirectory'],
    })
    exit(0)

subprocess.run(dir_path + "/docker/deploy.sh", shell=True, env={
    "SERVICE": args.service,
    "DEPLOYMENT_TARGET": config[args.service]['DeploymentTarget'],
    "TARGET_USERNAME": config[args.service]['TargetUsername'],
    "TARGET_DIRECTORY": config[args.service]['TargetDirectory'],
})
exit(0)