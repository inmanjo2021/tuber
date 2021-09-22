FROM bitnami/kubectl:1.15-ol-7 AS deps

FROM golang:1.16.5-alpine3.13 AS build

RUN mkdir /app
WORKDIR /app

COPY go.mod   ./go.mod
COPY go.sum   ./go.sum
RUN go mod download

COPY pkg      ./pkg
COPY cmd      ./cmd
COPY main.go  ./main.go
COPY data     ./data
COPY graph    ./graph
COPY .tuber   /.tuber

ENV CGO_ENABLED=0
RUN go build

FROM scratch AS run

COPY --from=deps /opt/bitnami/kubectl/bin/kubectl /usr/bin/kubectl
COPY --from=build /app/tuber /app/tuber

CMD ["/app/tuber", "start", "-y"]