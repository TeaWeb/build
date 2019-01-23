Tea.context(function () {
    this.deleteApp = function (appId) {
        if (!window.confirm("确定要删除此App吗？")) {
            return;
        }
        this.$post("/agents/apps/delete")
            .params({
                "agentId": this.agentId,
                "appId": appId
            })
            .refresh();
    };
});