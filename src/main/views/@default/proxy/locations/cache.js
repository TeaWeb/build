Tea.context(function () {
    /**
     * 缓存策略
     */
    this.cacheEditing = false;
    this.selectedCachePolicy = this.cachePolicyFile;

    this.editCache = function () {
        this.cacheEditing = !this.cacheEditing;
        if (this.cacheEditing) {
            this.$delay(function () {
                window.scroll(0, 10000);
            });
        }
    };

    this.cancelCacheEditing = function () {
        this.cacheEditing = false;
    };

    this.saveCacheEditing = function () {
        this.$post(".updateCache")
            .params({
                "serverId": this.server.id,
                "locationId": this.location.id,
                "policy": this.selectedCachePolicy
            })
            .success(function (resp) {
                window.location.reload();
            });
    };
});