#!/usr/bin/env bash
# a simple install script for ctop

KERNEL=$(uname -s)

# extract github download url matching pattern
function extract_url() {
  match=$1; shift
  echo "$@" | while read line; do
    case $line in
      *browser_download_url*${match}*)
        url=$(echo $line | sed -e 's/^.*"browser_download_url":[ ]*"//' -e 's/".*//;s/\ //g')
        echo $url
        break
      ;;
    esac
  done
}

case $KERNEL in
  Linux) MATCH_BUILD="linux-amd64" ;;
  Darwin) MATCH_BUILD="darwin-amd64" ;;
  *)
    echo "platform not supported by this install script"
    exit 1
    ;;
esac

TMP=$(mktemp -d "${TMPDIR:-/tmp}/ctop.XXXXX")
cd ${TMP}

echo "fetching latest release info"
resp=$(curl -s https://api.github.com/repos/bcicen/ctop/releases/latest)

echo "fetching release checksums"
checksum_url=$(extract_url sha256sums.txt "$resp")
wget -q $checksum_url -O sha256sums.txt

# skip if latest already installed
cur_ctop=$(which ctop 2> /dev/null)
if [[ -n "$cur_ctop" ]]; then
  cur_sum=$(sha256sum $cur_ctop | sed 's/ .*//')
  (grep -q $cur_sum sha256sums.txt) && {
    echo "already up-to-date"
    exit 0
  }
fi

echo "fetching latest ctop"
url=$(extract_url $MATCH_BUILD "$resp")
wget -q --show-progress $url
(sha256sum -c --quiet --ignore-missing sha256sums.txt) || exit 1

echo "installing to /usr/local/bin"
chmod +x ctop-*
sudo mv ctop-* /usr/local/bin/ctop

echo "done"
