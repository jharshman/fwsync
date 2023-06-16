#!/bin/sh
# Installs FWSYNC release.

VERSION="v0.0.1-7"
OS=$(uname -s | tr -d '\n')
ARCH=$(uname -m | tr -d '\n')
RELEASE=https://github.com/jharshman/fwsync/releases/download/${VERSION}/fwsync_${OS}_${ARCH}.tar.gz

which wget > /dev/null 2>&1
if [[ $? != 0 ]]; then
  echo "FATAL missing wget"
  exit 1
fi

# install
mkdir -p $HOME/.local/bin
wget -q $RELEASE
tar -C $HOME/.local/bin/ --exclude README.md -zxvf fwsync_${OS}_${ARCH}.tar.gz
chmod +x $HOME/.local/bin/fwsync

rcfile="$HOME/.zshrc"
if [[ $SHELL == "/bin/bash" ]]; then
  rcfile="$HOME/.bashrc"
fi

# update PATH if required.
if ! grep -q '# ADDED BY FWSYNC' $rcfile; then
  echo "export PATH=\$HOME/.local/bin:\$PATH # ADDED BY FWSYNC" >> $rcfile
fi

cat <<EOM
********************************************************
* FWSYNC has been installed at $HOME/.local/bin/fwsync
*
* Your PATH has been updated in $rcfile
* Restart your Terminal for the changes to take effect.
********************************************************
EOM
