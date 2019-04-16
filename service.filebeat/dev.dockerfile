FROM docker.elastic.co/beats/filebeat:7.0.0
COPY ./service.filebeat/filebeat.yml /usr/share/filebeat/filebeat.yml
ENV strict.perms false
USER ROOT