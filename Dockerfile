ARG PKG_VER_GO=1.11
ARG PKG_VER_ALPINE=3.8
ARG VERSION

FROM golang:${PKG_VER_GO}-alpine${PKG_VER_ALPINE} as cached_base

# Deps + useful debugging tools
RUN apk add --no-cache ca-certificates tzdata

FROM cached_base as build_base

WORKDIR /go/src/github.com/mindoktor/fillpdf

COPY . .

RUN go build -ldflags="-X main.version=$VERSION"

# Alpine base image
FROM alpine:$PKG_VER_ALPINE

RUN apk add --no-cache ca-certificates tzdata pdftk

# Setup an unprivileged user
RUN addgroup -S -g 101 docly && adduser -u 100 -S -G docly docly

COPY --from=build_base /go/src/github.com/mindoktor/fillpdf/fillpdf .
COPY --from=build_base /go/src/github.com/mindoktor/fillpdf/certificate /certificate

EXPOSE 8082

# Run as unprivileged user docly:docly
USER 100:101

CMD PORT=8082 ./fillpdf
