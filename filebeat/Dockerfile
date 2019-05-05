FROM docker.elastic.co/beats/filebeat:7.0.0
COPY ./filebeat/filebeat.yml /usr/share/filebeat/filebeat.yml

# filebeat.yml can only be writable by the owner
# https://www.elastic.co/guide/en/beats/libbeat/current/config-file-permissions.html
RUN chmod go-w /usr/share/filebeat/filebeat.yml

ENV strict.perms false
USER root