#!/usr/bin/env bash

# initialize project
export GOPATH=`pwd`/../../
export GO111MODULE=off

# mkdir
if [ ! -d ${GOPATH}/src/main/web/tmp ]
then
	echo "mkdir web/tmp/"
	mkdir ${GOPATH}/src/main/web/tmp
fi

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
go_get "github.com/mailru/easyjson"

# fastcgi
go_get "github.com/iwind/gofcgi"

# geo ip
go_get "github.com/oschwald/maxminddb-golang"
go_get "github.com/oschwald/geoip2-golang"

# system cpu, memory, disk ...
go_get "github.com/shirou/gopsutil"
go_get "github.com/shirou/w32"
go_get "github.com/StackExchange/wmi"
go_get "github.com/go-ole/go-ole"

# javascript
go_get "github.com/robertkrimen/otto"

# msg pack
go_get "github.com/vmihailenco/msgpack"

# redis
go_get "github.com/go-redis/redis"

# markdown
go_get "github.com/russross/blackfriday"

# fsnotify
go_get "github.com/fsnotify/fsnotify"

# websocket
go_get "github.com/gorilla/websocket"

# leveldb
go_get "github.com/syndtr/goleveldb/leveldb"

# go winio
go_get "github.com/Microsoft/go-winio"

# ping
go_get "github.com/tatsushid/go-fastping"

# ssh
go_get "github.com/pkg/sftp"

# mysql
go_get "github.com/go-sql-driver/mysql"

# pqsql
go_get "github.com/lib/pq"

# siphash
go_get "github.com/dchest/siphash"

# captcha
go_get "github.com/dchest/captcha"

# aliyun
go_get "github.com/aliyun/alibaba-cloud-sdk-go"

# ftp
go_get "github.com/jlaffaye/ftp"

# TeaWeb
go_get "github.com/TeaWeb/code"
echo "   [TeaWeb]Don't worry, you can ignore 'no Go files' warning in TeaWeb code"

go_get "github.com/TeaWeb/plugin"
go_get "github.com/TeaWeb/uaparser"
go_get "github.com/TeaWeb/agent"
go_get "github.com/TeaWeb/agentinstaller"

echo "[done]"