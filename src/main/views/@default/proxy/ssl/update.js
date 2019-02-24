Tea.context(function () {
    this.httpsOn = (this.server.ssl != null && this.server.ssl.on);

    if (this.server.ssl == null) {
       this.server.ssl = {
           "certificate": "",
           "certificateKey": "",
           "listen": []
       };
    }
    if (this.server.ssl.listen == null) {
        this.server.ssl.listen = [];
    }
    if (this.server.ssl.listen.length == 0) {
        this.server.ssl.listen = [ "0.0.0.0:443" ];
    }

    this.submitSuccess = function () {
        alert("修改成功");

        window.location = "/proxy/ssl?serverId=" + this.server.id;
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
        this.server.ssl.listen.push(this.addingListenName);
        this.cancelListenAdding();
    };

    this.cancelListenAdding = function () {
        this.listenAdding = !this.listenAdding;
        this.addingListenName = "";
    };

    this.removeListen = function (index) {
        this.server.ssl.listen.$remove(index);
    };
});