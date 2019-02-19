FROM node:8.11

# Install nodemon
RUN npm install -g nodemon

# Add the libraries
RUN mkdir -p /usr/src/libraries/javascript
COPY ./libraries/javascript /usr/src/libraries/javascript
WORKDIR /usr/src/libraries/javascript
RUN npm install

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install app dependencies
COPY ./service.controller.plug/package.json .
RUN npm install

# Bundle app source
COPY ./service.controller.plug .

# Expose ports for web access and debugging
EXPOSE 80 9229
CMD [ "npm", "run", "debug" ]
