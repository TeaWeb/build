Tea.context(function () {
    this.addProbe = function (probeId) {
        this.$post(".copyProbe")
            .params({
                "probeId": probeId
            })
            .success(function () {
                alert("添加成功");
                window.location.reload();
            });
    };
});