FROM ubuntu:22.04

RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-docker

RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-docker

RUN DEBIAN_FRONTEND=noninteractive \
  && apt-get update \
  && apt-get upgrade -y \
  && apt-get install -y wget curl ca-certificates vim \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir robot

COPY anote-robot robot/anote-robot

WORKDIR /robot

CMD [ "./anote-robot" ]