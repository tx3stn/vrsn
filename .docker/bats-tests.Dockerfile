FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/tx3stn/vrsn/cmd.Version=e2e-test" -o vrsn

FROM bats/bats:1.12.0

RUN apk add --no-cache \
	curl \
	git \
	musl-dev \
	expect

COPY --from=builder /app/vrsn /usr/bin/vrsn

ENTRYPOINT [ "bash" ]
