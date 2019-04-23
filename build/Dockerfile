FROM ubuntu

RUN mkdir -p /opt/blog/swagger-ui
RUN mkdir -p /opt/blog/secure

COPY blogserver /usr/local/bin/
COPY start-service.sh /opt/blog/

ADD swagger-ui/ /opt/blog/swagger-ui

RUN chmod +x /usr/local/bin/blogserver
RUN chmod +x /opt/blog/start-service.sh
