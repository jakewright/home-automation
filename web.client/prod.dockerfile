FROM nginx
COPY ./web.client/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./web.client/dist /usr/share/nginx/html
