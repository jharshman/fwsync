#!/bin/sh
# Installs FWSYNC release.

VERSION="v0.0.1"
OS=$(uname -s | tr -d '\n')
ARCH=$(uname -m | tr -d '\n')
RELEASE=https://github.com/jharshman/fwsync/releases/download/${VERSION}/fwsync_${OS}_${ARCH}.tar.gz

# delete any existing archives otherwise noop
rm -f fwsync_${OS}_${ARCH}.tar.gz* || :

which wget > /dev/null 2>&1
if [[ $? != 0 ]]; then
  echo "FATAL missing wget"
  exit 1
fi

# install
mkdir -p $HOME/.local/bin
wget -q $RELEASE
tar -C $HOME/.local/bin/ --exclude README.md --exclude LICENSE -zxvf fwsync_${OS}_${ARCH}.tar.gz
chmod +x $HOME/.local/bin/fwsync

# migrate existing bitly users' transaction file from .bitly_firewall to .fwsync
if [[ -e $HOME/.bitly_firewall ]]; then
  echo "project: bitly-devvm" > $HOME/.fwsync
  cat $HOME/.bitly_firewall >> $HOME/.fwsync
  rm -f $HOME/.bitly_firewall
fi

rcfile="$HOME/.zshrc"
if [[ $SHELL == "/bin/bash" ]]; then
  rcfile="$HOME/.bashrc"
fi

# update PATH if required.
if ! grep -q '# ADDED BY FWSYNC' $rcfile; then
  echo ""
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

rm -f fwsync_${OS}_${ARCH}.tar.gz
