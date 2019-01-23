Tea.context(function () {
    this.deleteAgent = function () {
        if (!window.confirm("确定要删除此主机吗？")) {
            return;
        }
        this.$post("/agents/delete")
            .params({
                "agentId": this.agentId
            })
            .success(function () {
                alert("删除成功");
                window.location = "/agents";
            });
    };
});