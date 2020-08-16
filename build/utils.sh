#!/usr/bin/env bash

# shell utilities

function build() {
	ROOT=$CWD/..
	VERSION_DATA=$(cat ${ROOT}/internal/teaconst/const.go)
	VERSION_DATA=${VERSION_DATA#*"Version = \""}
	VERSION=${VERSION_DATA%%[!0-9.]*}
	TARGET=${ROOT}/dist/teaweb-v${VERSION}
	GO_CMD="go"
	GOROOT=""

	EXT=""
	if [ ${GOOS} = "windows" ]; then
		EXT=".exe"
		GO_CMD="go"
	fi

	echo "[================ building ${GOOS}/${GOARCH}/v${VERSION}] ================]"

	echo "[goversion]using" $(${GO_CMD} version)
	echo "[create target directory]"

	if [ ! -d ${ROOT}/dist ]; then
		mkdir ${ROOT}/dist
	fi

	if [ -d ${TARGET} ]; then
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

	# build main & plugin
	${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb${EXT} ${ROOT}/cmd/teaweb/main.go
	${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/plugins/agent.tea${EXT} ${ROOT}/cmd/agent-plugin/main.go

	echo "[copy files]"
	cp -R ${ROOT}/build/configs/admin.sample.conf ${TARGET}/configs/admin.conf
	cp -R ${ROOT}/build/configs/server.sample.conf ${TARGET}/configs/server.conf
	cp -R ${ROOT}/build/configs/db.sample.conf ${TARGET}/configs/db.conf
	cp -R ${ROOT}/build/configs/mongo.sample.conf ${TARGET}/configs/mongo.conf
	cp -R ${ROOT}/build/configs/mysql.sample.conf ${TARGET}/configs/mysql.conf
	cp -R ${ROOT}/build/configs/postgres.sample.conf ${TARGET}/configs/postgres.conf
	cp -R ${ROOT}/build/configs/server.sample.www.proxy.conf ${TARGET}/configs/server.www.proxy.conf
	cp -R ${ROOT}/build/configs/widgets ${TARGET}/configs/
	cp -R ${ROOT}/build/www ${TARGET}/
	cp -R ${ROOT}/build/scripts ${TARGET}

	cp -R ${ROOT}/web/certs ${TARGET}/web/
	cp -R ${ROOT}/web/public ${TARGET}/web/
	cp -R ${ROOT}/web/resources ${TARGET}/web/
	cp -R ${ROOT}/web/views ${TARGET}/web/
	cp -R ${ROOT}/web/libs ${TARGET}/web/
	cp -R ${ROOT}/web/pages ${TARGET}/web/
	cp -R ${ROOT}/build/configs/widgets ${TARGET}/web/libs/

	if [ -d ${ROOT}/web/upgrade ]; then
		cp -R ${ROOT}/web/upgrade ${TARGET}/web/
	fi

	# windows
	if [ ${GOOS} = "windows" ]; then
		cp ${ROOT}/build/start.bat ${TARGET}
		cp ${ROOT}/build/README_WINDOWS.txt ${TARGET}/README.txt

		# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ${ROOT}/cmd/service-install/main.go
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ${ROOT}/cmd/service-uninstall/main.go
	fi

	# linux
	if [ ${GOOS} = "linux" ]; then
		# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ${ROOT}/cmd/service-install/main.go
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ${ROOT}/cmd/service-uninstall/main.go
	fi

	# not windows
	if [ ${GOOS} != "windows" ]; then
		cp ${ROOT}/build/README_LINUX.md ${TARGET}/README.md
		cp -R ${ROOT}/build/upgrade.sh ${TARGET}
	fi

	# installers
	if [ -d ${ROOT}/web/installers ]; then
		cp -R ${ROOT}/web/installers ${TARGET}/web/
	fi

	echo "[zip files]"
	cd ${TARGET}/../
	if [ -f teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip ]; then
		rm -f teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip
	fi
	zip -r -X -q teaweb-${GOOS}-${GOARCH}-v${VERSION}.zip teaweb-v${VERSION}/
	cd -

	rm -rf ${TARGET}

	echo "[done]"
}

function buildAgent() {
	ROOT=$CWD/..
	VERSION_DATA=$(cat ${ROOT}/internal/teaagent/agentconst/const.go)
	VERSION_DATA=${VERSION_DATA#*"Version = \""}
	VERSION=${VERSION_DATA%%[!0-9.]*}
	TARGET=${ROOT}/dist/teaweb-agent-v${VERSION}
	GO_CMD="go"
	GOROOT=""

	EXT=""
	if [ ${GOOS} = "windows" ]; then
		EXT=".exe"
		GO_CMD="go"
	fi

	echo "[================ building agent ${GOOS}/${GOARCH}/v${VERSION}] ================]"

	if [ -d ${TARGET} ]; then
		rm -rf ${TARGET}
	fi

	mkdir ${TARGET}
	mkdir ${TARGET}/bin
	mkdir ${TARGET}/configs
	mkdir ${TARGET}/configs/agents/
	mkdir ${TARGET}/logs
	mkdir ${TARGET}/plugins

	# config
	cp ${ROOT}/build/configs/agent.sample.conf ${TARGET}/configs/agent.conf

	if [ ${GOOS} = "windows" ]; then
		cp ${ROOT}/build/start-agent.bat ${TARGET}/start.bat
		cp ${ROOT}/build/README_AGENT_WINDOWS.txt ${TARGET}/README.txt

		# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install.exe ${ROOT}/cmd/agent-service-install/main.go
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall.exe ${ROOT}/cmd/agent-service-uninstall/main.go
	fi

	# linux
	if [ ${GOOS} = "linux" ]; then
		mkdir ${TARGET}/scripts
		cp ${ROOT}/build/scripts/teaweb-agent ${TARGET}/scripts/

		# service manager
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-install ${ROOT}/cmd/agent-service-install/main.go
		${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/service-uninstall ${ROOT}/cmd/agent-service-uninstall/main.go
	fi

	if [ ${GOOS} != "windows" ]; then
		cp ${ROOT}/build/README_AGENT_LINUX.md ${TARGET}/README.md
	fi

	${GO_CMD} build -ldflags="-s -w" -o ${TARGET}/bin/teaweb-agent${EXT} ${ROOT}/cmd/agent/main.go

	if [ ! -d "${ROOT}/web/upgrade/${VERSION}/${GOOS}/${GOARCH}" ]; then
		mkdir -p "${ROOT}/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"
	fi
	rm -f "${ROOT}/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/*
	cp ${TARGET}/bin/teaweb-agent${EXT} "${ROOT}/web/upgrade/${VERSION}/${GOOS}/${GOARCH}"/teaweb-agent${EXT}

	echo "[zip files]"
	cd ${TARGET}/../
	if [ -f teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip ]; then
		rm -f teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip
	fi
	zip -r -X -q teaweb-agent-${GOOS}-${GOARCH}-v${VERSION}.zip teaweb-agent-v${VERSION}/
	cd -

	echo "[clean files]"
	rm -rf ${TARGET}

	echo "[done]"
}

function buildAgentInstaller() {
	ROOT=$CWD/..
	GO_CMD="go"
	GOROOT=""

	EXT=""
	if [ ${GOOS} = "windows" ]; then
		EXT=".exe"
		GO_CMD="go"
	fi

	echo "[================ building agent installer ${GOOS}/${GOARCH}/v${VERSION}] ================]"

	if [ ! -d ${ROOT}/web/installers ]; then
		mkdir ${ROOT}/web/installers
	fi

	if [ ! -d ${ROOT}/web/installers ]; then
		rm -f ${ROOT}/web/installers/agentinstaller_*
		echo "remove"
	fi

	${GO_CMD} build -ldflags="-s -w" -o ${ROOT}/web/installers/agentinstaller_${GOOS}_${GOARCH}${EXT} ${ROOT}/cmd/agent-installer/main.go

	echo "[done]"
}
