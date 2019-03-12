#!/usr/bin/env bash

# shell utilities

function build() {
    VERSION_DATA=`cat ${GOPATH}/src/github.com/TeaWeb/code/teaconst/const.go`
    VERSION_DATA=${VERSION_DATA#*"Version = \""}
    VERSION=${VERSION_DATA%%[!0-9.]*}
    TARGET=${GOPATH}/dist/teaweb-v${VERSION}
    EXT=""
    if [ ${GOOS} = "windows" ]
    then
        EXT=".exe"
    fi

    echo "[================ building ${GOOS}/${GOARCH}/v${VERSION}] ================]"

    echo "[goversion]using" `go version`
    echo "[create target directory]"

    if [ ! -d ${GOPATH}/dist ]
    then
		mkdir ${GOPATH}/dist
    fi

    if [ -d ${TARGET} ]
    then
        rm -rf ${TARGET}
    fi

    mkdir ${TARGET}
    mkdir ${TARGET}/bin
    mkdir ${TARGET}/plugins
    mkdir ${TARGET}/tmp
    mkdir ${TARGET}/configs

    echo "[build static file]"

    # remove plus
    if [ -f ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go ]
    then
        rm -f ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    # build main & plugin
    go build -o ${TARGET}/bin/teaweb${EXT} ${GOPATH}/src/github.com/TeaWeb/code/main/main.go

    if [ -d ${GOPATH}/src/github.com/TeaWeb/agent ]
    then
        go build -o ${TARGET}/plugins/agent.tea${EXT} ${GOPATH}/src/github.com/TeaWeb/agent/main/main-plugin.go
    fi

    # restore plus
    if [ -f ${GOPATH}/drafts/src/plus.go ]
    then
        cp ${GOPATH}/drafts/src/plus.go ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    echo "[copy files]"
    cp -R ${GOPATH}/src/main/configs/admin.sample.conf ${TARGET}/configs/admin.conf
    cp -R ${GOPATH}/src/main/configs/server.sample.conf ${TARGET}/configs/server.conf
    cp -R ${GOPATH}/src/main/configs/mongo.sample.conf ${TARGET}/configs/mongo.conf
    cp -R ${GOPATH}/src/main/configs/server.sample.www.proxy.conf ${TARGET}/configs/server.www.proxy.conf
    cp -R ${GOPATH}/src/main/configs/widgets ${TARGET}/configs/
    cp -R ${GOPATH}/src/main/www ${TARGET}/

    cp -R ${GOPATH}/src/main/public ${TARGET}/
    cp -R ${GOPATH}/src/main/resources ${TARGET}/
    cp -R ${GOPATH}/src/main/views ${TARGET}/
    cp -R ${GOPATH}/src/main/libs ${TARGET}
    cp -R ${GOPATH}/src/main/upgrade ${TARGET}
    if [ -d ${TARGET}/libs/.idea ]
    then
		rm -rf ${TARGET}/libs/.idea
    fi

    if [ ${GOOS} = "windows" ]
    then
        cp ${GOPATH}/src/main/start.bat ${TARGET}
        cp ${GOPATH}/src/main/README_WINDOWS.txt ${TARGET}/README.txt
    fi

    if [ ${GOOS} != "windows" ]
    then
		cp ${GOPATH}/src/main/README_LINUX.md ${TARGET}/README.md
		cp -R ${GOPATH}/src/main/upgrade.sh ${TARGET}
    fi

    # remove plus files
    rm -rf ${TARGET}/views/@default/plus

    echo "[zip files]"
    cd ${TARGET}/../
    if [ -f teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip ]
    then
        rm -f teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip
    fi
    zip -r -X -q teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip  teaweb-v${VERSION}/
    cd -

    echo "[clean files]"
    rm -rf ${TARGET}

    echo "[done]"
}

function buildAgent() {
	VERSION_DATA=`cat ${GOPATH}/src/github.com/TeaWeb/code/teaconst/const.go`
	VERSION_DATA=${VERSION_DATA#*"Version = \""}
	VERSION=${VERSION_DATA%%[!0-9.]*}
	TARGET=${GOPATH}/dist/teaweb-agent-v${VERSION}
    EXT=""
	if [ ${GOOS} = "windows" ]
	then
    	EXT=".exe"
	fi

	echo "[================ building agent ${GOOS}/${GOARCH}/v${VERSION}] ================]"

	if [ -d ${TARGET} ]
    then
        rm -rf ${TARGET}
    fi

    mkdir ${TARGET}
    mkdir ${TARGET}/bin
    mkdir ${TARGET}/configs
    mkdir ${TARGET}/configs/agents/
    mkdir ${TARGET}/logs
    mkdir ${TARGET}/plugins

    cp ${GOPATH}/src/main/configs/agent.sample.conf ${TARGET}/configs/agent.conf

    if [ ${GOOS} = "windows" ]
    then
		cp ${GOPATH}/src/main/start-agent.bat ${TARGET}/start.bat
		cp ${GOPATH}/src/main/README_AGENT_WINDOWS.txt ${TARGET}/README.txt
    fi

    if [ ${GOOS} != "windows" ]
    then
		cp ${GOPATH}/src/main/README_AGENT_LINUX.md ${TARGET}/README.md
    fi

	go build -o ${TARGET}/bin/teaweb-agent${EXT} ${GOPATH}/src/github.com/TeaWeb/agent/main/main-agent.go

	if [ ! -d "${GOPATH}/src/main/upgrade/${VERSION}/${GOOS}/${GOARCH}" ]
	then
		mkdir -p "${GOPATH}/src/main/upgrade/${VERSION}/${GOOS}/${GOARCH}"
	fi
	rm -f "${GOPATH}/src/main/upgrade/${VERSION}/${GOOS}/${GOARCH}"/*
	cp ${TARGET}/bin/teaweb-agent${EXT} "${GOPATH}/src/main/upgrade/${VERSION}/${GOOS}/${GOARCH}"/teaweb-agent${EXT}

	echo "[zip files]"
    cd ${TARGET}/../
    if [ -f teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip ]
    then
        rm -f teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip
    fi
    zip -r -X -q teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip  teaweb-agent-v${VERSION}/
    cd -

    echo "[clean files]"
    rm -rf ${TARGET}

	echo "[done]"
}