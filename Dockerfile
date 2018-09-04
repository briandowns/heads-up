FROM scratch
COPY bin/heads-up /heads-up
ENTRYPOINT ["/heads-up"]
