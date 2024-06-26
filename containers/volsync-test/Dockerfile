FROM debian:stable-slim
WORKDIR /
COPY volsync-test /volsync-test
RUN chmod +x /volsync-test
RUN ls -la
ENTRYPOINT ["/volsync-test"]