FROM alpine:3.20
COPY dockguard /usr/local/bin/dockguard
ENTRYPOINT ["dockguard"]
