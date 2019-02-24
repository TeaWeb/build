Tea.context(function () {
    this.location = null;

    this.deleteLocation = function (locationId) {
        if (!window.confirm("确定要删除此路径配置吗？")) {
            return;
        }
        this.$post("/proxy/locations/delete")
            .params({
                "serverId": this.server.id,
                "locationId": locationId
            })
            .refresh();
    };

    this.moveUp = function (index) {
        this.$post("/proxy/locations/moveUp")
            .params({
                "serverId": this.server.id,
                "index": index
            });
    };

    this.moveDown = function (index) {
        this.$post("/proxy/locations/moveDown")
            .params({
                "serverId": this.server.id,
                "index": index
            });
    };
});