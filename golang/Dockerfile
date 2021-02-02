# builder image
FROM golang:1.15-alpine as builder
RUN mkdir /build

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o noobaa-bucket-example .

# generate clean, runtime image
FROM registry.redhat.io/ubi8-minimal
COPY --from=builder /build/noobaa-bucket-example /usr/bin/

# executable
ENTRYPOINT ["/usr/bin/noobaa-bucket-example"]
