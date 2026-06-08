FROM docker.m.daocloud.io/library/ubuntu:24.04

RUN echo 'hosts: files dns' >> /etc/nsswitch.conf

RUN set -ex && \
    mkdir -p /usr/bin /etc/categraf 

COPY docker/categraf  /usr/bin/categraf

COPY conf /etc/categraf/conf

COPY docker/entrypoint.sh /entrypoint.sh

RUN chmod +x /usr/bin/categraf /entrypoint.sh

CMD ["/entrypoint.sh"]
