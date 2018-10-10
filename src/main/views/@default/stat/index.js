Tea.context(function () {
    var that = this;

    this.dataType = "req";
    this.dataRange = "hourly";
    this.chartTitle = "";

    // 处理数据
    [this.topRequests, this.topOS, this.topBrowsers, this.topRegions, this.topStates ].$each(function (_, ranks) {
        ranks.$each(function (k, v) {
            var max = ranks[0].percent;
            var current = v.percent;
            if (max == 0) {
                v.compareMax = 0;
                return;
            }

            v.compareMax = current * 100;
        });
    });

    [this.topRequests, this.topCostRequests].$each(function (_, ranks) {
        ranks.$each(function (k, v) {
            var shortURL = v.url.substring(v.url.indexOf("//") + 2);
            v.uri = shortURL.substring(shortURL.indexOf("/"));
        });
    });

    this.topCostRequests.$each(function (k, v) {
        var max =  that.topCostRequests[0].cost;
        var current = v.cost;
        if (max == 0) {
            v.compareMax = 0;
            return;
        }
        v.compareMax = current * 50 / max;
    });

    this.$delay(function () {
        this.loadChart();

        setInterval(function () {
            that.loadChart();
        }, 60000);

        var that = this;
        window.addEventListener("resize", function () {
            var chart = echarts.init(that.$find(".main-box .chart")[0]);
            chart.resize();
        });
    });

    this.loadChart = function () {
        this.$get("/stat/data?type=" + this.dataType + "&range=" + this.dataRange)
            .success(function (response) {
                this.chartTitle = response.data.title;

                var chart = echarts.init(this.$find(".main-box .chart")[0]);

                // 指定图表的配置项和数据
                var option = {
                    textStyle: {
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    },
                    title: {
                        text: "",
                        top: 0,
                        x: "center"
                    },
                    tooltip: {},
                    legend: {
                        data: []
                    },
                    xAxis: {
                        data: response.data.labels
                    },
                    yAxis: {
                        axisLabel: {
                            formatter: function (value) {
                                if (value < 10000) {
                                    return value;
                                }
                                return (Math.round(value * 100 / 10000) / 100) + "万"
                            }
                        }
                    },
                    series: [{
                        name: '',
                        type: 'line',
                        data: response.data.data,
                        areaStyle: {
                            color: "#2185d0",
                            opacity: 0.8
                        },
                        lineStyle: {
                            color: "rgba(0, 0, 0, 0)"
                        },
                        itemStyle: {
                            color: "#2185d0"
                        }
                    }],
                    grid: {
                        left: 50,
                        right: 50,
                        bottom: 50,
                        top: 10
                    },
                    axisPointer: {
                        show: true
                    },
                    tooltip: {
                        formatter: 'X:{b0} Y:{c0}'
                    }
                };

                chart.setOption(option);
            });
    };

    this.changeType = function (dataType) {
        this.dataType = dataType;
        this.loadChart();
    };

    this.changeRange = function (dataRange) {
        this.dataRange = dataRange;
        this.loadChart();
    };

    this.formatPercent = function (num) {
        return Math.ceil(num * 10000) / 100;
    };

    this.formatMS = function (seconds) {
        return Math.ceil(seconds * 10000) / 10;
    };
});