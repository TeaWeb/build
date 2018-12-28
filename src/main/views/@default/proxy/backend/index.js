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
});