FROM alpine:3.16.2
RUN apk add --no-cache \
	curl \
	git \
	musl-dev

ENV VERSION_PROJECT=/project-vf
ENV TAG_PROJECT=/project-tag

RUN mkdir $VERSION_PROJECT && \
	mkdir $TAG_PROJECT

COPY .scripts/init-git.sh ${VERSION_PROJECT}/init-git.sh
COPY .scripts/init-project-vf.sh ${VERSION_PROJECT}/init-project-vf.sh
COPY .scripts/init-project-tag.sh ${TAG_PROJECT}/init-project-tag.sh

RUN chmod +x ${VERSION_PROJECT}/init-git.sh && \ 
	chmod +x ${VERSION_PROJECT}/init-project-vf.sh && \
	chmod +x ${TAG_PROJECT}/init-project-tag.sh && \
	sh ${VERSION_PROJECT}/init-git.sh && \
	sh ${VERSION_PROJECT}/init-project-vf.sh && \
	sh ${TAG_PROJECT}/init-project-tag.sh

COPY vrsn /usr/bin/vrsn
ENTRYPOINT [ "sh" ]
