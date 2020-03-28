#!/usr/bin/env bash

# usage: ./upgrade.sh teaweb-vx.x.x

FROM=${1}

if [ -d ${FROM} ]
then
	if [ ! -d ./web ]
	then
		mkdir ./web
	fi

	yes|cp -R ${FROM}/web/certs ./web/
	yes|cp -R ${FROM}/web/installers ./web/
	yes|cp -R ${FROM}/web/libs ./web/
	yes|cp -R ${FROM}/web/public ./web/
	yes|cp -R ${FROM}/web/resources ./web/
	yes|cp -R ${FROM}/scripts .
	yes|cp -R ${FROM}/web/upgrade ./web/
	yes|cp -R ${FROM}/web/views ./web/
	yes|cp -R ${FROM}/www .

	# bin & plugins
	bin/teaweb stop
	sleep 1
	yes|cp -R ${FROM}/bin .
	yes|cp -R ${FROM}/plugins .

	bin/teaweb start

	echo "[done]"
else
	echo "usage: ./upgrade.sh teaweb-vx.x.x"
fi
