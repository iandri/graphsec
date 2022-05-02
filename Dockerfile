FROM golang:alpine AS build-env

# ARG APPNAME
# ENV SRCPATH $GOPATH/src/github.com/iandri/graphsec
ENV SRCPATH /src
WORKDIR ${SRCPATH}
# RUN go version
ENV CGO_ENABLED=0

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build

FROM alpine

RUN apk add --no-cache dumb-init ca-certificates 

WORKDIR /app
COPY --from=build-env /src/graphsec  /app/
COPY --from=build-env /src/testdata /app/testdata


RUN addgroup -S appuser && adduser -S -G appuser appuser
RUN chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["/app/graphsec"]
