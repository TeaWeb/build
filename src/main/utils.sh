#!/usr/bin/env bash

# shell utilities

function build() {
    VERSION_DATA=`cat ${GOPATH}/src/github.com/TeaWeb/code/teaconst/const.go`
    VERSION_DATA=${VERSION_DATA#*"Version = \""}
    VERSION=${VERSION_DATA%%[!0-9.]*}
    TARGET=${GOPATH}/dist/teaweb-v${VERSION}
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
    mkdir ${TARGET}/logs
    mkdir ${TARGET}/plugins
    mkdir ${TARGET}/web
    mkdir ${TARGET}/web/tmp
    mkdir ${TARGET}/web/upgrade
    mkdir ${TARGET}/configs

    echo "[build static file]"

    # remove plus
    if [ -f ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go ]
    then
        rm -f ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    # build main & plugin
    ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb${EXT} ${GOPATH}/src/github.com/TeaWeb/code/main/main.go

    if [ -d ${GOPATH}/src/github.com/TeaWeb/agent ]
    then
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/plugins/agent.tea${EXT} ${GOPATH}/src/github.com/TeaWeb/agent/main/main-plugin.go
    fi

    # restore plus
    if [ -f ${GOPATH}/drafts/src/plus.go ]
    then
        cp ${GOPATH}/drafts/src/plus.go ${GOPATH}/src/github.com/TeaWeb/code/teaweb/plus.go
    fi

    echo "[copy files]"
    cp -R ${GOPATH}/src/main/configs/admin.sample.conf ${TARGET}/configs/admin.conf
    cp -R ${GOPATH}/src/main/configs/server.sample.conf ${TARGET}/configs/server.conf
    cp -R ${GOPATH}/src/main/configs/db.sample.conf ${TARGET}/configs/db.conf
    cp -R ${GOPATH}/src/main/configs/mongo.sample.conf ${TARGET}/configs/mongo.conf
    cp -R ${GOPATH}/src/main/configs/mysql.sample.conf ${TARGET}/configs/mysql.conf
    cp -R ${GOPATH}/src/main/configs/postgres.sample.conf ${TARGET}/configs/postgres.conf
    cp -R ${GOPATH}/src/main/configs/server.sample.www.proxy.conf ${TARGET}/configs/server.www.proxy.conf
    cp -R ${GOPATH}/src/main/configs/widgets ${TARGET}/configs/
    cp -R ${GOPATH}/src/main/www ${TARGET}/

    cp -R ${GOPATH}/src/main/web/public ${TARGET}/web/
    cp -R ${GOPATH}/src/main/web/resources ${TARGET}/web/
    cp -R ${GOPATH}/src/main/web/views ${TARGET}/web/
    cp -R ${GOPATH}/src/main/web/libs ${TARGET}/web/
    cp -R ${GOPATH}/src/main/web/pages ${TARGET}/web/
    cp -R ${GOPATH}/src/main/configs/widgets ${TARGET}/web/libs/

    if [ -d ${GOPATH}/src/main/web/upgrade ]
    then
    	cp -R ${GOPATH}/src/main/web/upgrade ${TARGET}/web/
    fi

    cp -R ${GOPATH}/src/main/scripts ${TARGET}

	# windows
    if [ ${GOOS} = "windows" ]
    then
        cp ${GOPATH}/src/main/start.bat ${TARGET}
        cp ${GOPATH}/src/main/README_WINDOWS.txt ${TARGET}/README.txt

		# service manager
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ${GOPATH}/src/github.com/TeaWeb/code/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ${GOPATH}/src/github.com/TeaWeb/code/main/service_uninstall.go
    fi

    # linux
    if [ ${GOOS} = "linux" ]
    then
    	# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ${GOPATH}/src/github.com/TeaWeb/code/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ${GOPATH}/src/github.com/TeaWeb/code/main/service_uninstall.go
    fi

	# not windows
    if [ ${GOOS} != "windows" ]
    then
		cp ${GOPATH}/src/main/README_LINUX.md ${TARGET}/README.md
		cp -R ${GOPATH}/src/main/upgrade.sh ${TARGET}
    fi

    # installers
	if [ -d ${GOPATH}/src/main/web/installers ]
	then
		cp -R ${GOPATH}/src/main/web/installers ${TARGET}/web/
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
	VERSION_DATA=`cat ${GOPATH}/src/github.com/TeaWeb/agent/teaconst/const.go`
	VERSION_DATA=${VERSION_DATA#*"Version = \""}
	VERSION=${VERSION_DATA%%[!0-9.]*}
	TARGET=${GOPATH}/dist/teaweb-agent-v${VERSION}
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
    cp ${GOPATH}/src/main/configs/agent.sample.conf ${TARGET}/configs/agent.conf

    if [ ${GOOS} = "windows" ]
    then
		cp ${GOPATH}/src/main/start-agent.bat ${TARGET}/start.bat
		cp ${GOPATH}/src/main/README_AGENT_WINDOWS.txt ${TARGET}/README.txt

		# service manager
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ${GOPATH}/src/github.com/TeaWeb/agent/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ${GOPATH}/src/github.com/TeaWeb/agent/main/service_uninstall.go
    fi

    # linux
    if [ ${GOOS} = "linux" ]
    then
    	mkdir ${TARGET}/scripts
    	cp ${GOPATH}/src/main/scripts/teaweb-agent ${TARGET}/scripts/

    	# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ${GOPATH}/src/github.com/TeaWeb/agent/main/service_install.go
        ${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ${GOPATH}/src/github.com/TeaWeb/agent/main/service_uninstall.go
    fi

    if [ ${GOOS} != "windows" ]
    then
		cp ${GOPATH}/src/main/README_AGENT_LINUX.md ${TARGET}/README.md
    fi

	${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb-agent${EXT} ${GOPATH}/src/github.com/TeaWeb/agent/main/main-agent.go

	if [ ! -d "${GOPATH}/src/main/web/upgrade/${VERSION}/${GOOS}/${GOARCH}" ]
	then
		mkdir -p "${GOPATH}/src/main/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"
	fi
	rm -f "${GOPATH}/src/main/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/*
	cp ${TARGET}/bin/teaweb-agent${EXT} "${GOPATH}/src/main/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/teaweb-agent${EXT}

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


	if [ ! -d ${GOPATH}/src/main/web/installers ]
	then
		mkdir ${GOPATH}/src/main/web/installers
	fi

	if [ ! -d ${GOPATH}/src/main/web/installers ]
	then
		rm -f ${GOPATH}/src/main/web/installers/*
	fi

	${GO_CMD} build -ldflags="-s -w" -o ${GOPATH}/src/main/web/installers/agentinstaller_${GOOS}_${GOARCH}${EXT} ${GOPATH}/src/github.com/TeaWeb/agentinstaller/main/main.go

	echo "[done]"
}

