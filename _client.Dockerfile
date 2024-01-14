########## BUILD CLIENT ##########
FROM golang:1.21-alpine AS build-client

RUN apk add --no-cache make

WORKDIR /app
COPY . /app

COPY cmd/client /cmd/client/
COPY internal /internal/

RUN make client-build

CMD ["./bin/client"]