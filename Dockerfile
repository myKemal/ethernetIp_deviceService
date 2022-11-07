
## build app
FROM golang:1.18-alpine3.15 AS builder
RUN apk update && apk add --no-cache make git zeromq-dev gcc pkgconfig musl-dev

WORKDIR /device-ethernetip-go

# Install our build time packages.
COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."
COPY . .
RUN make build



# Next image - Copy built Go binary into new workspace
FROM alpine:3.15

# dumb-init needed for injected secure bootstrapping entrypoint script when run in secure mode.
RUN apk add --update --no-cache zeromq dumb-init


# expose command data port
ENV APP_PORT=59996
#expose meta data port
EXPOSE $APP_PORT

COPY --from=builder /device-ethernetip-go/cmd/ /
COPY --from=builder /device-ethernetip-go/cmd/res /res

ENTRYPOINT ["/device-ethernetip-go"]
CMD ["--cp=consul://edgex-core-consul:8500", "--registry", "--confdir=/res"] 