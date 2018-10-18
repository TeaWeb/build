#!/usr/bin/env bash

# initialize project

export GOPATH=`pwd`/../../

# go_get function
function go_get() {
    echo "go get ${1} ..."
    if [ ! -d ${GOPATH}/${1} ]
    then
        go get ${1}
    fi
}

# TeaGO
go_get "github.com/iwind/TeaGo"
go_get "github.com/pquerna/ffjson"

# fastcgi
go_get "github.com/iwind/gofcgi"

# geo ip
go_get "github.com/oschwald/maxminddb-golang"
go_get "github.com/oschwald/geoip2-golang"

# cache memory
go_get "github.com/pbnjay/memory"

# mongodb
go_get "github.com/mongodb/mongo-go-driver"
echo "   Don't worry, you can ignore 'no Go files' warning in mongodb"

# usa parser
go_get "github.com/ua-parser/uap-go"

# system cpu, memory, disk ...
go_get "github.com/shirou/gopsutil"

# javascript
go_get "github.com/robertkrimen/otto"

# TeaWeb
go_get "github.com/TeaWeb/code"
echo "   Don't worry, you can ignore 'no Go files' warning in TeaWeb code"

echo "[done]"