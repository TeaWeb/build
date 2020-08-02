Tea.context(function () {
    this.loading = false;
    this.success = false;
    this.hasNew = false;
    this.latest = "";

    this.errorMessage = "";

    this.checkVersion = function () {
        this.loading = true;
        this.success = false;
        this.errorMessage = "";

        this.$post("/settings/update")
            .success(function (resp) {
                this.success = true;
                this.hasNew = resp.data.hasNew;
                this.latest = resp.data.latest;
            })
            .fail(function (resp) {
                this.errorMessage = resp.message;
            })
            .done(function () {
                this.loading = false;
            });
    };
});