#!/bin/sh 
this_dir=`pwd`
mkdir Peer/

chmod 777 Peer
cp -r peer/src/peer.yaml Peer/

export GOPATH=$this_dir:$this_dir/peer/:$this_dir/vendor/
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin


go build -o Peer/peer peer/src/main.go
