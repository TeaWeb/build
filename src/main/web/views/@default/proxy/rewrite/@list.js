Tea.context(function () {
    this.isLoaded = false;
    this.from = encodeURIComponent(window.location.toString());
    this.query = Tea.serialize(this.queryParams);
    this.rewriteList = [];

    this.$delay(function () {
        this.loadData();
    });

    this.loadData = function () {
        this.$get("/proxy/rewrite/data")
            .params(this.queryParams)
            .success(function (resp) {
                this.rewriteList = resp.data.rewriteList;
            })
            .done(function () {
                this.isLoaded = true;
            });
    };

    this.deleteRewrite = function (rewrite) {
        if (!window.confirm("确定要删除此重写规则吗？")) {
            return;
        }
        var params = this.queryParams;
        params["rewriteId"] = rewrite.id;
        this.$post("/proxy/rewrite/delete")
            .params(params)
            .refresh();
    };
});