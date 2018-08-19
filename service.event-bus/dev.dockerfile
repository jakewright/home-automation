FROM node:8.4
ENV NODE_ENV=development

# Create app directory
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

# Install nodemon
RUN npm install nodemon -g

# Install app dependencies
COPY ./service.event-bus/package.json /usr/src/app/
RUN npm install

# Bundle app source
COPY ./service.event-bus /usr/src/app

# Expose ports for web access and debugging
EXPOSE 80 9229
CMD [ "npm", "run", "debug" ]
