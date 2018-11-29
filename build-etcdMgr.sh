#!/bin/sh 
this_dir=`pwd`
mkdir EtcdMgr
chmod 777 EtcdMgr
cp -r etcdMgr/src/etcdMgr.yaml EtcdMgr/
export GOPATH=$this_dir:$this_dir/etcdMgr:$this_dir/vendor/
go build -o EtcdMgr/etcdMgr etcdMgr/src/main.go
