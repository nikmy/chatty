FROM golang:1.21.4-alpine

WORKDIR /app

ADD . .

RUN go build ./cmd/rpc_server/
RUN rm go.mod

EXPOSE 80
ENTRYPOINT [ "./rpc_server", "-port=80" ]
