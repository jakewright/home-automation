FROM node:8.11

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install nodemon
RUN npm install -g nodemon

# Install app dependencies
COPY ./service.controller.hue/package.json .
RUN npm install

# Bundle app source
COPY ./service.controller.hue .

# Expose ports for web access and debugging
EXPOSE 80 9229
CMD [ "npm", "run", "debug" ]
