FROM centos:7
LABEL maintainer="root@teaos.cn"
ENV TZ "Asia/Shanghai"
ENV TEAWEB_VERSION "0.1.3"
ENV TEAWEB_URL "http://dl.teaos.cn/v${TEAWEB_VERSION}/teaweb-linux-amd64-v${TEAWEB_VERSION}.zip"
ENV MONGO_VERSION "4.0.6"
ENV MONGO_URL "https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-${MONGO_VERSION}.tgz"

RUN yum -y install unzip; \
	cd /opt; \
	echo "downloading ${TEAWEB_URL}"; \
	curl ${TEAWEB_URL} -o ./teaweb-linux-amd64-v${TEAWEB_VERSION}.zip; \
	unzip teaweb-linux-amd64-v${TEAWEB_VERSION}.zip; \
	mv teaweb-v${TEAWEB_VERSION} teaweb; \
	echo "downloading ${MONGO_URL}"; \
	curl ${MONGO_URL} -o ./mongodb-linux-x86_64-${MONGO_VERSION}.tgz; \
	tar -zxvf mongodb-linux-x86_64-${MONGO_VERSION}.tgz; \
	mv mongodb-linux-x86_64-${MONGO_VERSION} mongodb; \
	cd mongodb; \
	mkdir data;
COPY teaweb.sh /opt/teaweb.sh
EXPOSE 7777
ENTRYPOINT [ "/opt/teaweb.sh" ]
