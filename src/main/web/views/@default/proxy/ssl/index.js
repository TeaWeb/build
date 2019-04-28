Tea.context(function () {
    /**
     * SSL管理
     */
    this.showSSLOptions = (this.server.ssl != null) ? this.server.ssl.on : false;
    this.sslCertFile = null;
    this.sslCertEditing = false;

    this.sslKeyFile = null;
    this.sslKeyEditing = false;

    this.switchSSLOn = function () {
        var message = (this.server.ssl != null && this.server.ssl.on) ? "确定要关闭HTTPS吗？" : "确定要开启HTTPS吗？";
        if (!window.confirm(message)) {
            return;
        }

        this.showSSLOptions = !this.showSSLOptions;
        if (this.server.ssl == null) {
            this.server.ssl = { "on": this.showSSLOptions };
        } else {
            this.server.ssl.on = !this.server.ssl.on;
        }

        if (this.server.ssl.on) {
            this.$post("/proxy/ssl/on")
                .params({
                    "serverId": this.server.id
                });
        } else {
            this.$post("/proxy/ssl/off")
                .params({
                    "serverId": this.server.id
                });
        }
    };

    this.changeSSLCertFile = function (event) {
        if (event.target.files.length > 0) {
            this.sslCertFile = event.target.files[0];
        }
    };

    this.uploadSSLCertFile = function () {
        if (this.sslCertFile == null) {
            alert("请先选择证书文件");
            return;
        }

        this.$post("/proxy/ssl/uploadCert")
            .params({
                "serverId": this.server.id,
                "certFile": this.sslCertFile
            });
    };

    this.changeSSLKeyFile = function (event) {
        if (event.target.files.length > 0) {
            this.sslKeyFile = event.target.files[0];
        }
    };

    this.uploadSSLKeyFile = function () {
        if (this.sslKeyFile == null) {
            alert("请先选择密钥文件");
            return;
        }

        this.$post("/proxy/ssl/uploadKey")
            .params({
                "serverId": this.server.id,
                "keyFile": this.sslKeyFile
            });
    };

	/**
	 * 启动
	 */
	this.startHttps = function () {
		if (!window.confirm("确定要启动此HTTPS服务吗？")) {
			return;
		}
		this.$post(".startHttps")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};

	this.shutdownHttps = function () {
		if (!window.confirm("确定要关闭此HTTPS服务吗？")) {
			return;
		}
		this.$post(".shutdownHttps")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};
});