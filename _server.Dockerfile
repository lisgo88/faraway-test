########## BUILD SERVER ##########
FROM golang:1.21-alpine AS build-server

RUN apk add --no-cache make

WORKDIR /app
COPY ../_faraway-test /app

COPY cmd/server /cmd/server/
COPY internal /internal/

RUN make server-build

EXPOSE 8080

CMD ["./bin/server"]
