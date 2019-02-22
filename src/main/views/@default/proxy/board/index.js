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
                "server": this.server.filename
            })
            .timeout(10)
			.success(function (resp) {
				// output
				resp.data.output.$each(function (k, v) {
					console.log("[widget]" + v);
				});

				// charts
				this.charts = resp.data.charts;
				new ChartRender(this.charts);
			})
			.fail(function (resp) {
				throw new Error("[widget]" + resp.message);
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