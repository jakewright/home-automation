FROM node:8.11

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install app dependencies
RUN npm install -g @vue/cli
COPY ./web.client/.npmrc /usr/src/app
COPY ./web.client/package.json /usr/src/app
RUN npm install

# Expose ports for web access and debugging
EXPOSE 8080 9229

CMD [ "npm", "run", "serve" ]
