FROM ghcr.io/charmbracelet/vhs:v0.10.0

RUN apt-get update && \
	apt-get -y install --no-install-recommends git && \
	git config --global --add safe.directory /vhs

ENTRYPOINT ["/usr/bin/vhs"]
