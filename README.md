# Home Automation

Distributed home automation system. Largely a learning opportunity rather than a production-ready system.

The home automation system is made up of separate microservices which run on various devices that are distributed physically around the home.
Most of the core services run within a Kubernetes cluster, however there are some periphery services that sit outside of the cluster.
The hardware used is a combination of Raspberry Pis and a Synology NAS.

A series of YouTube videos accompanies this project. They can be found in [this YouTube playlist](https://www.youtube.com/playlist?list=PLlj9BrHKq9WI4R30l_M_tdRMPF4AZ6dcs).

## Will this work for me?

This is not designed as a general-purpose home automation system. It is pretty specific to my use cases. If you’re looking for something generic, check out [Home Assistant](https://www.home-assistant.io) or [openHAB](https://www.openhab.org). If, however, it does work for you, feel free to use it but don’t expect any support. I also don’t plan to take feature requests. If you would like to make any changes then I suggest [forking the repository](https://help.github.com/en/github/getting-started-with-github/fork-a-repo).


## Getting started

There are various tools that can be installed to aid development.

```shell
./tools/install
``` 

## Project structure
- `docs/`
  - Documentation about the system
- `libraries/`
  - Library code shared between all services
- `private/`
  - A git submodule containing mostly private configuration
- `services/`
  - The services that make up the system
- `tools/`
  - Useful tools for working with the system
- `web.x`
  - A web-based application
