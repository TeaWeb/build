#!/usr/bin/env bash

# shell utilities

function build() {
    VERSION_DATA=`cat ../src/github.com/TeaWeb/code/teaconst/const.go`
    VERSION_DATA=${VERSION_DATA#*"Version = \""}
    VERSION=${VERSION_DATA%%[!0-9.]*}
    TARGET=../dist/teaweb-v${VERSION}
    GO_CMD="go"
    GOROOT=""

    EXT=""
    if [ ${GOOS} = "windows" ]
    then
        EXT=".exe"

		# we use go 1.11 to build 386 binary
        if [ ${GOARCH} = "386" ]
		then
			echo "check go 1.11 for old windows"
			result=`go.1.11 version|wc -c`
			if [ ${result} -gt 0 ]
			then
				GO_CMD="go.1.11"
				GOROOT=""
			else
				GO_CMD="go"
			fi
        fi
    fi

    echo "[================ building ${GOOS}/${GOARCH}/v${VERSION}] ================]"

    echo "[goversion]using" `${GO_CMD} version`
    echo "[create target directory]"

    if [ ! -d ../dist ]
    then
		mkdir ../dist
    fi

    if [ -d ${TARGET} ]
    then
        rm -rf ${TARGET}
    fi

    mkdir ${TARGET}
    mkdir ${TARGET}/bin
    mkdir ${TARGET}/logs
    mkdir ${TARGET}/plugins
    mkdir ${TARGET}/web
    mkdir ${TARGET}/web/tmp
    mkdir ${TARGET}/web/upgrade
    mkdir ${TARGET}/configs

    echo "[build static file]"

    # remove plus
    if [ -f ../src/github.com/TeaWeb/code/teaweb/plus.go ]
    then
        rm -f ../src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    # build main & plugin
    ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb${EXT} ../src/github.com/TeaWeb/code/main/main.go

    if [ -d ../src/github.com/TeaWeb/agent ]
    then
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/plugins/agent.tea${EXT} ../src/github.com/TeaWeb/agent/main/main-plugin.go
    fi

    # restore plus
    if [ -f ../drafts/src/plus.go ]
    then
        cp ../drafts/src/plus.go ../src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    echo "[copy files]"
    cp -R ../src/configs/admin.sample.conf ${TARGET}/configs/admin.conf
    cp -R ../src/configs/server.sample.conf ${TARGET}/configs/server.conf
    cp -R ../src/configs/db.sample.conf ${TARGET}/configs/db.conf
    cp -R ../src/configs/mongo.sample.conf ${TARGET}/configs/mongo.conf
    cp -R ../src/configs/mysql.sample.conf ${TARGET}/configs/mysql.conf
    cp -R ../src/configs/postgres.sample.conf ${TARGET}/configs/postgres.conf
    cp -R ../src/configs/server.sample.www.proxy.conf ${TARGET}/configs/server.www.proxy.conf
    cp -R ../src/configs/widgets ${TARGET}/configs/
    cp -R ../src/www ${TARGET}/

	cp -R ../src/web/certs ${TARGET}/web/
    cp -R ../src/web/public ${TARGET}/web/
    cp -R ../src/web/resources ${TARGET}/web/
    cp -R ../src/web/views ${TARGET}/web/
    cp -R ../src/web/libs ${TARGET}/web/
    cp -R ../src/web/pages ${TARGET}/web/
    cp -R ../src/configs/widgets ${TARGET}/web/libs/

    if [ -d ../src/web/upgrade ]
    then
    	cp -R ../src/web/upgrade ${TARGET}/web/
    fi

    cp -R ../src/scripts ${TARGET}

	# windows
    if [ ${GOOS} = "windows" ]
    then
        cp ../src/start.bat ${TARGET}
        cp ../src/README_WINDOWS.txt ${TARGET}/README.txt

		# service manager
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ../src/github.com/TeaWeb/code/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ../src/github.com/TeaWeb/code/main/service_uninstall.go
    fi

    # linux
    if [ ${GOOS} = "linux" ]
    then
    	# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ../src/github.com/TeaWeb/code/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ../src/github.com/TeaWeb/code/main/service_uninstall.go
    fi

	# not windows
    if [ ${GOOS} != "windows" ]
    then
		cp ../src/README_LINUX.md ${TARGET}/README.md
		cp -R ../src/upgrade.sh ${TARGET}
    fi

    # installers
	if [ -d ../src/web/installers ]
	then
		cp -R ../src/web/installers ${TARGET}/web/
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
	VERSION_DATA=`cat ../src/github.com/TeaWeb/agent/teaconst/const.go`
	VERSION_DATA=${VERSION_DATA#*"Version = \""}
	VERSION=${VERSION_DATA%%[!0-9.]*}
	TARGET=../dist/teaweb-agent-v${VERSION}
    GO_CMD="go"
    GOROOT=""

    EXT=""
    if [ ${GOOS} = "windows" ]
    then
        EXT=".exe"

		# we use go 1.11 to build 386 binary
        if [ ${GOARCH} = "386" ]
		then
			echo "check go 1.11 for old windows"
			result=`go.1.11 version|wc -c`
			if [ ${result} -gt 0 ]
			then
				GO_CMD="go.1.11"
				GOROOT=""
			else
				GO_CMD="go"
			fi
        fi
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

	# config
    cp ../src/configs/agent.sample.conf ${TARGET}/configs/agent.conf

    if [ ${GOOS} = "windows" ]
    then
		cp ../src/start-agent.bat ${TARGET}/start.bat
		cp ../src/README_AGENT_WINDOWS.txt ${TARGET}/README.txt

		# service manager
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ../src/github.com/TeaWeb/agent/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ../src/github.com/TeaWeb/agent/main/service_uninstall.go
    fi

    # linux
    if [ ${GOOS} = "linux" ]
    then
    	mkdir ${TARGET}/scripts
    	cp ../src/scripts/teaweb-agent ${TARGET}/scripts/

    	# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ../src/github.com/TeaWeb/agent/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ../src/github.com/TeaWeb/agent/main/service_uninstall.go
    fi

    if [ ${GOOS} != "windows" ]
    then
		cp ../src/README_AGENT_LINUX.md ${TARGET}/README.md
    fi

	${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb-agent${EXT} ../src/github.com/TeaWeb/agent/main/main-agent.go

	if [ ! -d "../src/web/upgrade/${VERSION}/${GOOS}/${GOARCH}" ]
	then
		mkdir -p "../src/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"
	fi
	rm -f "../src/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/*
	cp ${TARGET}/bin/teaweb-agent${EXT} "../src/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/teaweb-agent${EXT}

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

function buildAgentInstaller() {
    GO_CMD="go"
    GOROOT=""

    EXT=""
    if [ ${GOOS} = "windows" ]
    then
        EXT=".exe"

			# we use go 1.11 to build 386 binary
    	if [ ${GOARCH} = "386" ]
			then
				echo "check go 1.11 for old windows"
				result=`go.1.11 version|wc -c`
				if [ ${result} -gt 0 ]
				then
					GO_CMD="go.1.11"
					GOROOT=""
				else
					GO_CMD="go"
				fi
      fi
    fi

	echo "[================ building agent installer ${GOOS}/${GOARCH}/v${VERSION}] ================]"


	if [ ! -d ../src/web/installers ]
	then
		mkdir ../src/web/installers
	fi

	if [ ! -d ../src/web/installers ]
	then
		rm -f ../src/web/installers/*
	fi

	${GO_CMD} build -ldflags="-s -w" -o ../src/web/installers/agentinstaller_${GOOS}_${GOARCH}${EXT} ../src/github.com/TeaWeb/agentinstaller/main/main.go

	echo "[done]"
}

