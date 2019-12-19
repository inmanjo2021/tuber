FROM bitnami/kubectl:1.15-ol-7

FROM golang:1.13.5-alpine3.10

COPY --from=0 /opt/bitnami/kubectl/bin/kubectl /usr/bin/kubectl

RUN mkdir /app
WORKDIR /app

COPY go.mod   ./go.mod
COPY go.sum   ./go.sum
COPY justfile ./justfile
COPY pkg      ./pkg
COPY cmd      ./cmd
COPY main.go  ./main.go

ENV GO111MODULE on

RUN go build

CMD ["/app/tuber", "start"]

# tuber
COPY .tuber /.tuber
