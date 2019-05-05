FROM docker.elastic.co/logstash/logstash:7.0.0

COPY ./logstash/logstash.yml /usr/share/logstash/config/logstash.yml
COPY ./logstash/logstash.conf /usr/share/logstash/pipeline/logstash.conf

# Stop Filebeat from harvesting logs from a container running this image
LABEL "co.elastic.logs/disable"="true"

# Set the user to root to avoid permissions errors on the mounted log volume
USER root

# Ports for Filebeat and to receive Syslog messages
EXPOSE 5044 7514