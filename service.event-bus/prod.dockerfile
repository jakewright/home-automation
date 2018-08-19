FROM node:8.4
ENV NODE_ENV=production

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install app dependencies
COPY ./service.event-bus/package.json /usr/src/app/
RUN npm install

# Bundle app source
COPY ./service.event-bus /usr/src/app

CMD [ "npm", "run", "debug" ]
