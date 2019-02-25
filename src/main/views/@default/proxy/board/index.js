Tea.context(function () {
    var that = this;

    this.charts = [];

    this.$delay(function () {
        this.loadCharts();
    });

    this.widgetIsLoaded = false;
    this.widgetError = "";
    this.events = [];

    this.loadCharts = function () {
        this.$post("/proxy/board")
            .params({
                "serverId": this.server.id,
				"type": this.boardType,
				"events": JSON.stringify(this.events)
            })
            .timeout(10)
			.success(function (resp) {
				// output
				if (resp.data.output != null) {
					resp.data.output.$each(function (k, v) {
						console.log("[widget]" + v);
					});
				}

				// charts
				this.charts = resp.data.charts;
				var that = this;
				new ChartRender(this.charts, function (events) {
					that.events = events;
					that.loadCharts();
				});
			})
			.fail(function (resp) {
				throw new Error("[widget]" + resp.message);
			})
			.done(function () {
				this.$delay(function () {
					this.widgetIsLoaded = true;
				}, 500);

				if (this.boardType == "realtime") {
					this.$delay(function () {
						this.loadCharts();
					}, 5000);
				}

				//this.events = [];
			});
    };
});