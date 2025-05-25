FROM bats/bats:1.12.0

RUN apk add --no-cache \
	curl \
	git \
	musl-dev

COPY vrsn /usr/bin/vrsn

ENTRYPOINT [ "bash" ]
