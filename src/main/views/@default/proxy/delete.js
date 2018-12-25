Tea.context(function () {
    this.deleteProxy = function () {
        this.$post("/proxy/delete")
            .params({
                "server": this.server.filename
            })
            .success(function () {
                alert("删除成功");
                window.location = "/proxy";
            });
    };
});