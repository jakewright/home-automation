FROM node:15

# Add the libraries
RUN mkdir -p /usr/src/libraries/javascript
WORKDIR /usr/src/libraries/javascript
COPY ./libraries/javascript .
RUN npm install

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install app dependencies
RUN npm install -g @vue/cli@4.5.8
COPY ./web.client/package.json .
COPY ./web.client/package-lock.json .
RUN npm install

# Move one level up so node_modules is not overwritten by a mounted directory
RUN mv node_modules /usr/src/node_modules

# Expose ports for web access and debugging
EXPOSE 8080 9229

CMD [ "npm", "run", "serve" ]
