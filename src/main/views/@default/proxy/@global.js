Tea.context(function () {
   this.statusChanged = false;

    this.refreshStatus = function () {
        this.$get("/proxy/status")
            .success(function (response) {
                this.statusChanged = response.data.changed;
            })
            .done(function () {
                this.$delay(function () {
                    this.refreshStatus();
                }, 3000);
            });
    };

    this.$delay(function () {
        this.refreshStatus();
    }, 100);

    this.restart = function () {
        this.$get("/proxy/restart")
            .success(function () {
                this.statusChanged = false;
            });
    };
});