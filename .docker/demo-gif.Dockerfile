FROM ghcr.io/charmbracelet/vhs:v0.10.0

RUN rm -rf /var/lib/apt/lists/* && \
	apt-get update --allow-releaseinfo-change && \
	apt-get -y install --no-install-recommends git && \
	git config --global --add safe.directory /vhs

ENTRYPOINT ["/usr/bin/vhs"]
