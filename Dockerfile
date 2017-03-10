FROM debian:jessie

RUN BUILD_PACKAGES="curl wget" && \
    apt-get update && \
    apt-get install -y $BUILD_PACKAGES && \
    wget $(curl -s https://api.github.com/repos/bcicen/ctop/releases/latest | \
           grep 'browser_' | cut -d\" -f4 |grep 'linux-amd64') \
      -O /usr/local/bin/ctop && \
    chmod +x /usr/local/bin/ctop &&
    AUTO_ADDED_PACKAGES=`apt-mark showauto` && \
    apt-get remove --purge -y $BUILD_PACKAGES $AUTO_ADDED_PACKAGES

CMD ["ctop"]
