FROM node:11
ENV NODE_ENV=development

# Install nodemon
RUN npm install -g nodemon

# Add the libraries
RUN mkdir -p /usr/src/libraries/javascript
COPY ./libraries/javascript /usr/src/libraries/javascript
WORKDIR /usr/src/libraries/javascript
RUN npm install

# Move one level up so node_modules is not overwritten by a mounted directory
RUN mv node_modules /usr/src/libraries/node_modules

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install app dependencies
COPY ./service.event-bus/package.json .
RUN npm install

# Move one level up so node_modules is not overwritten by a mounted directory
RUN mv node_modules /usr/src/node_modules

# Bundle app source
COPY ./service.event-bus .

# Expose ports for web access and debugging
EXPOSE 80 9229

CMD nodemon --inspect=0.0.0.0:9229 --watch . --watch /usr/src/libraries/javascript index.js
