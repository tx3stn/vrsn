FROM alpine:3.21.3
RUN apk add --no-cache \
  curl \
  git \
  musl-dev

ENV WORKDIR=/repo
RUN mkdir $WORKDIR && \
  git config --global --add safe.directory $WORKDIR

WORKDIR ${WORKDIR}

COPY vrsn /usr/bin/vrsn
ENTRYPOINT ["vrsn"]
