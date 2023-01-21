#!/bin/bash -eu

# Download the nebud binary with the same version as mainnet and unpack it

# USAGE: ./download-mainnet-nebud.sh

is_macos() {
  [[ "$OSTYPE" == "darwin"* ]]
}

architecture=$(uname -m)

CWD="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
MAINNET_VERSION=${MAINNET_VERSION:-"v3.0.2"}

download_mainnet_binary(){
  # Checks for the nebud v1 file
  if [ ! -f "$NEBUD_BIN_MAINNET" ]; then
    echo "$NEBUD_BIN_MAINNET doesn't exist"

    if [ -z $NEBUD_BIN_MAINNET_URL_TARBALL ]; then
      echo You need to set the NEBUD_BIN_MAINNET_URL_TARBALL variable
      exit 1
    fi

    NEBUD_RELEASES_PATH=$CWD/nebud-releases
    mkdir -p $NEBUD_RELEASES_PATH
    wget -c $NEBUD_BIN_MAINNET_URL_TARBALL -O - | tar -xz -C $NEBUD_RELEASES_PATH

    NEBUD_BIN_MAINNET_BASENAME=$(basename $NEBUD_BIN_MAINNET_URL_TARBALL .tar.gz)
    NEBUD_BIN_MAINNET=$NEBUD_RELEASES_PATH/$NEBUD_BIN_MAINNET_BASENAME/nebud
  fi
}

mac_mainnet() {
  if [[ "$architecture" == "arm64" ]];then
    NEBUD_BIN_MAINNET_URL_TARBALL=${NEBUD_BIN_MAINNET_URL_TARBALL:-"https://github.com/tessornetwork/nebula/releases/download/${MAINNET_VERSION}/nebud-${MAINNET_VERSION}-darwin-arm64.tar.gz"}
    NEBUD_BIN_MAINNET=${NEBUD_BIN_MAINNET:-"$CWD/nebud-releases/nebud-${MAINNET_VERSION}-darwin-arm64/nebud"}
  else
    NEBUD_BIN_MAINNET_URL_TARBALL=${NEBUD_BIN_MAINNET_URL_TARBALL:-"https://github.com/tessornetwork/nebula/releases/download/${MAINNET_VERSION}/nebud-${MAINNET_VERSION}-darwin-amd64.tar.gz"}
    NEBUD_BIN_MAINNET=${NEBUD_BIN_MAINNET:-"$CWD/nebud-releases/nebud-${MAINNET_VERSION}-darwin-amd64/nebud"}
  fi
}

linux_mainnet(){
  if [[ "$architecture" == "arm64" ]];then
    NEBUD_BIN_MAINNET_URL_TARBALL=${NEBUD_BIN_MAINNET_URL_TARBALL:-"https://github.com/tessornetwork/nebula/releases/download/${MAINNET_VERSION}/nebud-${MAINNET_VERSION}-linux-arm64.tar.gz"}
    NEBUD_BIN_MAINNET=${NEBUD_BIN_MAINNET:-"$CWD/nebud-releases/nebud-${MAINNET_VERSION}-linux-arm64/nebud"}
  else
    NEBUD_BIN_MAINNET_URL_TARBALL=${NEBUD_BIN_MAINNET_URL_TARBALL:-"https://github.com/tessornetwork/nebula/releases/download/${MAINNET_VERSION}/nebud-${MAINNET_VERSION}-linux-amd64.tar.gz"}
    NEBUD_BIN_MAINNET=${NEBUD_BIN_MAINNET:-"$CWD/nebud-releases/nebud-${MAINNET_VERSION}-linux-amd64/nebud"}
  fi
}

if is_macos ; then
  mac_mainnet
  download_mainnet_binary
else
  linux_mainnet
  download_mainnet_binary
fi
