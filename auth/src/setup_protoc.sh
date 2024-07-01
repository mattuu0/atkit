#!/bin/bash
apt update
PROTOC_VERSION=$(curl -s "https://api.github.com/repos/protocolbuffers/protobuf/releases/latest" | grep -Po '"tag_name": "v\K[0-9.]+')
curl -Lo protoc.zip "https://github.com/protocolbuffers/protobuf/releases/latest/download/protoc-${PROTOC_VERSION}-linux-x86_64.zip"
apt install unzip
unzip -q protoc.zip bin/protoc 'include/*' -d /usr/local
rm protoc.zip

chmod +x ./proto/compile.sh
chmod +x ./proto/bin/protoc