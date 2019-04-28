Tea.context(function () {
    this.isStarting = false;
    this.startMongo = function () {
        this.isStarting = true;
        this.$post(".install")
            .success(function () {
                window.location.reload();
            })
            .fail(function () {
                this.isStarting = false;
            });
    };
});