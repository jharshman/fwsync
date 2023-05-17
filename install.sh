#!/bin/sh

VERSION="v0.0.1-1"
OS=$(uname -s | tr -d '\n')
ARCH=$(uname -m | tr -d '\n')
RELEASE=https://github.com/jharshman/fwsync/releases/download/${VERSION}/fwsync_${OS}_${ARCH}.tar.gz

# install
mkdir -p $HOME/.local/bin
wget $RELEASE
tar -zxvf -C $HOME/.local/bin/ fwsync_${OS}_${ARCH}
chmod +x $HOME/.local/bin/fwsync

# update PATH
echo "export PATH=$HOME/.local/bin:$PATH" >> $HOME/.bashrc

cat <<EOM
/////////////
// FWSYNC has been installed at $HOME/.local/bin/fwsync
//
// Your PATH has been updated in .bashrc.
// Restart your Terminal for the changes to take effect.
///////////////////////////////////////////////////////////
EOM
