FROM ubuntu:precise

MAINTAINER Drifty

# Install some utils
#wget -q -O - https://www.loggly.com/install/configure-syslog.py | sudo python - setup --auth e3cbb07f-a71b-4ccc-b391-6e2865a3d663 --account drifty


RUN apt-get update &&\
    apt-get install -y python-software-properties &&\
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 5862E31D &&\
    apt-get update &&\
    apt-get install -y build-essential wget rsyslog supervisor curl git &&\
    mkdir -p /var/log/supervisor

RUN curl -s https://go.googlecode.com/files/go1.2.1.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz

ENV PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH  /home/docker/code/
ENV GOROOT  /usr/local/go

ADD . /home/docker/code/
WORKDIR /home/docker/code

# Clone libsass including submodules


RUN ln -s /home/docker/code/docker/supervisord.conf /etc/supervisor/conf.d/

RUN go get -d &&\
    bash build.sh

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 8080

CMD ["/usr/bin/supervisord"]
