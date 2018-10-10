Tea.context(function () {
    /**
     * SSL管理
     */
    this.showSSLOptions = (this.proxy.ssl != null) ? this.proxy.ssl.on : false;
    this.sslCertFile = null;
    this.sslCertEditing = false;

    this.sslKeyFile = null;
    this.sslKeyEditing = false;

    this.switchSSLOn = function () {
        var message = (this.proxy.ssl != null && this.proxy.ssl.on) ? "确定要关闭HTTPS吗？" : "确定要开启HTTPS吗？";
        if (!window.confirm(message)) {
            return;
        }

        this.showSSLOptions = !this.showSSLOptions;
        if (this.proxy.ssl == null) {
            this.proxy.ssl = { "on": this.showSSLOptions };
        } else {
            this.proxy.ssl.on = !this.proxy.ssl.on;
        }

        if (this.proxy.ssl.on) {
            this.$post("/proxy/ssl/on")
                .params({
                    "filename": this.filename
                });
        } else {
            this.$post("/proxy/ssl/off")
                .params({
                    "filename": this.filename
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
                "filename": this.filename,
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
                "filename": this.filename,
                "keyFile": this.sslKeyFile
            });
    };
});