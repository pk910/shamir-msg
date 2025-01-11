
# go build env
FROM golang:1.23 AS go-env
COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG release=
RUN make build RELEASE=$release GOOS=$TARGETOS GOARCH=$TARGETARCH

# final stage
FROM debian:stable-slim
WORKDIR /app
COPY --from=go-env /src/bin /app
ENTRYPOINT ["./shamir-msg"]
CMD []
