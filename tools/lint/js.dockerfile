FROM node:8.14
WORKDIR /usr/src

RUN npm i -g prettier babel-eslint eslint eslint-plugin-vue eslint-plugin-import eslint-config-airbnb-base

COPY tools/lint/js_fmt.sh ./
RUN chmod +x js_fmt.sh

COPY ./tools/lint/.eslintrc ./

# WORKDIR /usr/src/home-automation/web.client
# COPY ./web.client/package*.json ./
# RUN npm install

COPY . ./home-automation

CMD ["/usr/src/js_fmt.sh"]
