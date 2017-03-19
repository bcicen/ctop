FROM scratch
COPY ./ctop /ctop
ENTRYPOINT ["/ctop"]
