FROM golang:1.16 AS base

ADD . /home

WORKDIR /home

CMD ["go","run","."]