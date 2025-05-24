FROM alpine:3.16.2
RUN apk add --no-cache \
	curl \
	git \
	musl-dev

ENV VERSION_PROJECT=/project-vf
RUN mkdir $VERSION_PROJECT && \
	git config --global --add safe.directory $VERSION_PROJECT

ENV TAG_PROJECT=/project-tag
RUN mkdir $TAG_PROJECT && \
	git config --global --add safe.directory $TAG_PROJECT

COPY .scripts/init-project-vf.sh ${VERSION_PROJECT}/init-project-vf.sh
COPY .scripts/init-project-tag.sh ${TAG_PROJECT}/init-project-one.sh

RUN ${VERSION_PROJECT}/init-project-vf.sh && \
	${TAG_PROJECT}/init-project-tag.sh

COPY vrsn /usr/bin/vrsn
ENTRYPOINT [ "sh" ]
