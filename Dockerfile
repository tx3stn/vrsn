FROM alpine:3.24.1

RUN apk upgrade --no-cache && \
	apk add --no-cache \
	curl=8.21.0-r0 \
	git=2.54.0-r0 \
	musl-dev=1.2.6-r2 && \
	rm -rf /var/cache/apk/*

ENV WORKDIR=/repo
RUN mkdir $WORKDIR && \
	git config --global --add safe.directory $WORKDIR

WORKDIR ${WORKDIR}

COPY --chmod=755 vrsn /usr/bin/vrsn

ENTRYPOINT ["vrsn"]
