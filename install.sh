#!/bin/sh
# Installs FWSYNC release.

VERSION="v0.0.1-3"
OS=$(uname -s | tr -d '\n')
ARCH=$(uname -m | tr -d '\n')
RELEASE=https://github.com/jharshman/fwsync/releases/download/${VERSION}/fwsync_${OS}_${ARCH}.tar.gz

# install
mkdir -p $HOME/.local/bin
wget $RELEASE
tar -C $HOME/.local/bin/ --exclude README.md -zxvf fwsync_${OS}_${ARCH}.tar.gz
chmod +x $HOME/.local/bin/fwsync

# update PATH if required.
if ! grep -q '# ADDED BY FWSYNC' $HOME/.zshrc; then
  echo "export PATH=\$HOME/.local/bin:\$PATH # ADDED BY FWSYNC" >> $HOME/.zshrc
fi

cat <<EOM
/////////////
// FWSYNC has been installed at $HOME/.local/bin/fwsync
//
// Your PATH has been updated in .zshrc.
// Restart your Terminal for the changes to take effect.
///////////////////////////////////////////////////////////
EOM
