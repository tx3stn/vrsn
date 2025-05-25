FROM bats/bats:1.12.0

RUN apk add --no-cache \
	curl \
	git \
	musl-dev

ENV VF_PROJECT=/project-vf
ENV TAG_PROJECT=/project-tag
ENV SCRIPTS_DIR=/scripts

RUN mkdir $SCRIPTS_DIR && \
	mkdir $VF_PROJECT && \
	mkdir $TAG_PROJECT

COPY .scripts/init-git.sh ${SCRIPTS_DIR}/init-git.sh
COPY .scripts/init-project-vf.sh ${SCRIPTS_DIR}/init-project-vf.sh
COPY .scripts/init-project-tag.sh ${SCRIPTS_DIR}/init-project-tag.sh

RUN sh ${SCRIPTS_DIR}/init-git.sh && \
	sh ${SCRIPTS_DIR}/init-project-vf.sh ${VF_PROJECT} && \
	sh ${SCRIPTS_DIR}/init-project-tag.sh ${TAG_PROJECT}

COPY vrsn /usr/bin/vrsn

ENTRYPOINT [ "sh" ]
