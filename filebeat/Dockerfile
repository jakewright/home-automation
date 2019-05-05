FROM docker.elastic.co/beats/filebeat:7.0.0
COPY ./filebeat/filebeat.yml /usr/share/filebeat/filebeat.yml

# The Elastic images use a different user that causes permission issues
# Set to root before changing the permissions of filebeat.yml below
USER root

# filebeat.yml can only be writable by the owner
# https://www.elastic.co/guide/en/beats/libbeat/current/config-file-permissions.html
RUN chown root:root /usr/share/filebeat/filebeat.yml
RUN chmod go-w /usr/share/filebeat/filebeat.yml

ENV strict.perms false
