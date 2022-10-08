FROM ghcr.io/charmbracelet/vhs:v0.2.0

RUN apt-get -y install --no-install-recommends git

ENTRYPOINT ["/usr/bin/vhs"]
