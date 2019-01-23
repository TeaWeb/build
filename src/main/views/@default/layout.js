Tea.context(function () {
    this.mongoFailed = false;

    this.testMongo = function () {
        this.$get("/mongo/test")
            .fail(function () {
                this.mongoFailed = true;
            });
    };
});