FROM docker.elastic.co/beats/filebeat-oss:7.15.0

COPY ./filebeat/filebeat.yml /usr/share/filebeat/filebeat.yml
USER root
RUN chown -R root /usr/share/filebeat/
RUN chmod -R 777 /usr/share/filebeat