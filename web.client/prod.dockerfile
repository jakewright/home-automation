FROM node:11

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

# Copy source code
COPY ./web.client/ .

# Build the app
RUN npm run build

FROM nginx
COPY --from=0 /usr/src/app/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=0 /usr/src/app/dist /usr/share/nginx/html
