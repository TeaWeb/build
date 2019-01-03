Tea.context(function () {
    this.httpsOn = (this.proxy.ssl != null && this.proxy.ssl.on);

    if (this.proxy.ssl == null) {
       this.proxy.ssl = {
           "certificate": "",
           "certificateKey": "",
           "listen": []
       };
    }
    if (this.proxy.ssl.listen == null) {
        this.proxy.ssl.listen = [];
    }
    if (this.proxy.ssl.listen.length == 0) {
        this.proxy.ssl.listen = [ "0.0.0.0:443" ];
    }

    this.submitSuccess = function () {
        alert("修改成功");

        window.location = "/proxy/ssl?server=" + this.proxy.filename;
    };

    /**
     * 绑定的网络地址
     */
    this.listenAdding = false;
    this.addingListenName = "";

    this.addListen = function () {
        this.listenAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingListenName']").focus();
        });
    };

    this.confirmAddListen = function () {
        this.addingListenName = this.addingListenName.trim();
        if (this.addingListenName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingListenName']").focus();
            return;
        }
        this.proxy.ssl.listen.push(this.addingListenName);
        this.cancelListenAdding();
    };

    this.cancelListenAdding = function () {
        this.listenAdding = !this.listenAdding;
        this.addingListenName = "";
    };

    this.removeListen = function (index) {
        this.proxy.ssl.listen.$remove(index);
    };
});