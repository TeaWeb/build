#!/usr/bin/env bash

# usage: ./upgrade.sh teaweb-vx.x.x

FROM=${1}

if [ ${FROM} ]
then
	yes|cp -R ${FROM}/libs .
	yes|cp -R ${FROM}/public .
	yes|cp -R ${FROM}/resources .
	yes|cp -R ${FROM}/views .
	yes|cp -R ${FROM}/www

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
