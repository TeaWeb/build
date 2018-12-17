Tea.context(function () {
    this.deleteProbe = function (probeId) {
        if (!window.confirm("确定要删除此探针吗？")) {
            return;
        }
        this.$post(".deleteProbe")
            .params({
                "probeId": probeId
            })
            .success(function () {
                window.location.reload();
            });
    };
});