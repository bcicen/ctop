FROM quay.io/vektorcloud/glibc:latest

RUN ctop_url=$(wget -q -O - https://api.github.com/repos/bcicen/ctop/releases/latest | grep 'browser_' | cut -d\" -f4 |grep 'linux-amd64') && \
    wget -q $ctop_url -O /ctop && \
    chmod +x /ctop

ENTRYPOINT ["/ctop"]
