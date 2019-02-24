Tea.context(function () {
    this.deleteProxy = function () {
        if (!window.confirm("确定要删除此代理服务吗？")) {
            return;
        }
        this.$post("/proxy/delete")
            .params({
                "serverId": this.server.id
            })
            .success(function () {
                alert("删除成功");
                window.location = "/proxy";
            });
    };
});