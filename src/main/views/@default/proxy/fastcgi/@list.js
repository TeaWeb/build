Tea.context(function () {
    this.loaded = false;
    this.from = encodeURIComponent(window.location.toString());

    this.$delay(function () {
        this.loadData();
    });

    this.fastcgiList = [];
    this.query = Tea.serialize(this.queryParams);

    this.loadData = function () {
        this.$get("/proxy/fastcgi/data")
            .params(this.queryParams)
            .success(function (resp) {
                this.fastcgiList = resp.data.fastcgiList;
                this.loaded = true;
            });
    };

    this.deleteFastcgi = function (fastcgiId) {
        if (!window.confirm("确定要删除此Fastcgi设置吗？")) {
            return;
        }
        var query = this.queryParams;
        query["fastcgiId"] = fastcgiId
        this.$post("/proxy/fastcgi/delete")
            .params(query)
            .refresh();
    };
});