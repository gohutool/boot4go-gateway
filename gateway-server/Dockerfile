FROM golang:1.18.0-alpine as build
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn,direct"
RUN go mod tidy
RUN go mod download
RUN go build -o server .

FROM alpine
WORKDIR /app
COPY --from=build /app/server /app
COPY ./log4go.xml /app

## fix the bug that dyncimal library is loss in alpine image.
RUN  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


EXPOSE 9000 9443

ENTRYPOINT  ["./server"]
