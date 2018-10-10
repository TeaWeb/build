Tea.context(function () {
    this.widgetGroups = [];

    this.$delay(function () {
        this.loadWidgets();
    });

    this.loadWidgets = function () {
        this.$get(".widgets")
            .success(function (response) {
                this.CHART.updateWidgetGroups(response.data.widgetGroups);
            })
            .done(function () {
                this.$delay(function () {
                    this.loadWidgets();
                }, 1000);
            });
    };
});