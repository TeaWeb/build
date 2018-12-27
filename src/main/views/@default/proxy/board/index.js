Tea.context(function () {
    var that = this;

    this.charts = [];

    this.$delay(function () {
        this.loadWidgets();
    });

    this.widgetIsLoaded = false;
    this.widgetError = ""

    this.loadWidgets = function () {
        this.$post("/proxy/board")
            .params({
                "server": this.server.filename,
                "config": this.config
            })
            .timeout(10)
            .success(function (resp) {
                //console.log(JSON.stringify(resp.data.charts, 1, "  "));
                this.widgetError = resp.data.widgetError;
                if (this.widgetError != null && this.widgetError.length > 0) {
                    console.log("WIDGET ERROR:", this.widgetError);
                }

                this.charts = resp.data["charts"];
                new ChartRender(this.charts);
            })
            .done(function () {
                this.$delay(function () {
                    this.widgetIsLoaded = true;
                }, 500);

                this.$delay(function () {
                    this.loadWidgets();
                }, 3000);
            });
    };
});