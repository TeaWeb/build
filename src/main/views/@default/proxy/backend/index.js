Tea.context(function () {
    this.$delay(function () {
        // scroll to bottom
        if (window.location.hash == "#scheduling") {
            window.scrollTo(0, 10000);
        }
    }, 300);

    this.deleteBackend = function (backendId) {
        if (!window.confirm("确定要删除此服务器吗？")) {
            return;
        }
        this.$post("/proxy/backend/delete")
            .params({
                "server": this.proxy.filename,
                "backendId": backendId
            });
    };

    this.putOnline = function (backend) {
        if (!window.confirm("确定要上线此服务器吗？")) {
            return;
        }
        this.$post("/proxy/backend/online")
            .params({
                "server": this.proxy.filename,
                "backendId": backend.id
            })
            .success(function () {
                backend.isDown = false;
                backend.currentFails = 0;
            });
    };

    this.clearFails = function (backend) {
        if (!window.confirm("确定要清除此服务器的失败次数吗？")) {
            return;
        }
        this.$post("/proxy/backend/clearFails")
            .params({
                "server": this.proxy.filename,
                "backendId": backend.id
            })
            .success(function () {
                backend.isDown = false;
                backend.currentFails = 0;
            });
    };
});