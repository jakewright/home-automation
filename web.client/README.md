# Web client

The web client is a [Vue.js](https://vuejs.org) project. This is the main interface to the home automation system.
 
## Setup
 
While it's possible to do all of this using Docker, it's not worth the hassle. A useful addition could be a Docker container that is able to execute all of the npm and Vue-related commands, but the current recommendation is to install the tools locally.

First, install node using `brew`.

```sh
brew update && brew install node
```

If node is already installed, make sure it is up-to-date.

```sh
brew update && brew upgrade node
```

Install or update npm

```
npm install -g npm
```
 
Manage the project using the [Vue CLI](https://github.com/vuejs/vue-cli). Version 4.5 of the CLI is required, and this should be installed globally on your local machine.

```sh
npm install -g @vue/cli
```
