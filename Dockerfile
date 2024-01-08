FROM golang:bullseye AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-test-api

FROM scratch

WORKDIR /

EXPOSE 8117

COPY --from=build-stage /go-test-api /go-test-api

ENTRYPOINT [ "/go-test-api" ]
