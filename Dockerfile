FROM golang:1.8

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...

# 安装glide
# RUN curl https://glide.sh/get | sh
# RUN make dev
# glide install
# RUN go install -v ./...