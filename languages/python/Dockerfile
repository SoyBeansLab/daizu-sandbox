FROM python:3.7.3-slim-stretch

RUN mkdir /tmp/daizu/ && chmod 700 /tmp/daizu/

WORKDIR /tmp/daizu

COPY run /tmp/daizu/

RUN chmod +x /tmp/daizu/run
