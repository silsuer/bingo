FROM golang:1.8

MAINTAINER silsuer <silsuer.liu@gmail.com>

WORKDIR /go/src/github.com/silsuer/bingo

# COPY . .

#RUN curl https://glide.sh/get | sh && \

#    glide install

#    go get -d -v ./...


ENTRYPOINT ["make","dev"]


