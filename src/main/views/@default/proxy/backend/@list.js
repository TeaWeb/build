Tea.context(function () {
    this.from = encodeURIComponent(window.location.toString());
    this.query = Tea.serialize(this.queryParams);

    this.padZero = function (s) {
        s = s.toString();
        if (s.length == 1) {
            return "0" + s;
        }
        return s;
    };

    var that = this;
    this.normalBackends = [];
    this.backupBackends = [];
    this.scheduling = null;
    this.isLoaded = false;

    this.$delay(function () {
        this.loadData();
    });

    this.$delay(function () {
        // scroll to bottom
        if (window.location.hash == "#scheduling") {
            window.scrollTo(0, 10000);
        }
    }, 300);

    this.loadData = function () {
        this.$get("/proxy/backend/data")
            .params(this.queryParams)
            .success(function (resp) {
                this.normalBackends = resp.data.normalBackends.$map(function (k, v) {
                    if (v.isDown) {
                        var date = new Date(v.downTime);
                        v.downTime = that.padZero(date.getMonth() + 1) + "-" + that.padZero(date.getDate()) + " " + that.padZero(date.getHours()) + ":" + that.padZero(date.getMinutes()) + ":" + that.padZero(date.getSeconds());
                    }
                    return v;
                });
                this.backupBackends = resp.data.backupBackends.$map(function (k, v) {
                    if (v.isDown) {
                        var date = new Date(v.downTime);
                        v.downTime = that.padZero(date.getMonth() + 1) + "-" + that.padZero(date.getDate()) + " " + that.padZero(date.getHours()) + ":" + that.padZero(date.getMinutes()) + ":" + that.padZero(date.getSeconds());
                    }
                    return v;
                });
                this.scheduling = resp.data.scheduling;
            })
            .done(function () {
                this.isLoaded = true;
                this.$delay(function () {
                    this.loadData();
                }, 5000);
            })
            .fail(function (resp) {
                console.log(resp.message);
            });
    };

    this.deleteBackend = function (backendId) {
        if (!window.confirm("确定要删除此服务器吗？")) {
            return;
        }
        var query = this.queryParams;
        query["backendId"] = backendId;
        this.$post("/proxy/backend/delete")
            .params(query);
    };

    this.putOnline = function (backend) {
        if (!window.confirm("确定要上线此服务器吗？")) {
            return;
        }
        var query = this.queryParams;
        query["backendId"] = backend.id;
        this.$post("/proxy/backend/online")
            .params(query)
            .success(function () {
                backend.isDown = false;
                backend.currentFails = 0;
            });
    };

    this.clearFails = function (backend) {
        if (!window.confirm("确定要清除此服务器的失败次数吗？此操作不会改变上线状态")) {
            return;
        }
        var query = this.queryParams;
        query["backendId"] = backend.id;
        this.$post("/proxy/backend/clearFails")
            .params(query)
            .success(function () {
                backend.currentFails = 0;
            });
    };
});