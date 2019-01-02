function ChartRender(charts) {
    var that = this;

    charts.$map(function (k, v) {
        v.id = Math.random().toString().replace(".", "_");
        Tea.delay(function () {
            that.render(v);
        });
        return v;
    });

    this.render = function (v) {
        var chartType = v.type;
        var f = "render" + chartType[0].toUpperCase() + chartType.substring(1) + "Chart";
        var html = this[f](v);
        if (html != null && html.length > 0) {
            Tea.element("#chart-box-" + v.id).html(html);
        }
    };

    this.renderHtmlChart = function (chart) {
        return chart.html;
    };

    this.renderLineChart = function (chart) {
        var chartBox = Tea.element("#chart-box-" + chart.id)[0];
        if (chartBox == null) {
            return "";
        }
        var c = echarts.init(chartBox);

        var option = {
            textStyle: {
                fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
            },
            title: {
                text: "",
                subText : "",
                top: 0,
                x: "left",
                textStyle: {
                    fontSize: 12,
                    fontWeight: "bold",
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                }
            },
            legend: {
                data: chart.lines.$map(function (_, line) {
                    if (line != null) {
                        return line.name;
                    }
                    return "";
                }),
                bottom: -10,
                y: "bottom",
                textStyle: {
                    fontSize: 10
                }
            },
            xAxis: {
                data: chart.labels,
                axisTick: {
                    show: chart.xShowTick
                }
            },
            axisLabel: {
                formatter: function (v) {
                    return v;
                },
                textStyle: {
                    fontSize: 10
                }
            },
            yAxis: {
                max: (chart.max != 0 ) ? chart.max : null,
                splitNumber: (chart.yTickCount >= 1) ? chart.yTickCount : null,
                axisTick: {
                    show: chart.yShowTick
                }
            },
            series: chart.lines.$map(function (_, line) {
                return {
                    name: line.name,
                    type: 'line',
                    data: line.values,
                    lineStyle: {
                        width: 1.2,
                        color: line.color,
                        opacity: 0.5
                    },
                    itemStyle: {
                        color: line.color,
                        opacity: line.showItems ? 1 : 0
                    },
                    areaStyle: {
                        color: line.color,
                        opacity: line.isFilled ? 0.5 : 0
                    },
                    smooth: line.smooth
                };
            }),
            grid: {
                left: 40,
                right: 40,
                bottom: 16,
                top: 16
            },
            axisPointer: {
                show: false
            },
            tooltip: {
                formatter: 'X:{b0} Y:{c0}',
                show: false
            },
            animation: false
        };

        c.setOption(option);
        return "";
    };

    this.renderPieChart = function (chart) {
        var chartBox = Tea.element("#chart-box-" + chart.id)[0];
        if (chartBox == null) {
            return "";
        }
        var c = echarts.init(chartBox);

        var option = {
            textStyle: {
                fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
            },
            title: {
                text: "",
                top: 10,
                bottom: 10,
                x: "center",
                textStyle: {
                    fontSize: 12,
                    fontWeight: "bold",
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                }
            },
            legend: {
                orient: 'vertical',
                x: 'right',
                y: 'center',
                data: chart.labels,
                itemWidth: 6,
                itemHeight: 6,
                textStyle: {
                    fontSize: 12
                }
            },
            xAxis: {
                data: []
            },
            yAxis: {},
            series: [{
                name: '',
                type: 'pie',
                data: chart.values.$map(function (k, v) {
                    return {
                        name: chart.labels.$get(k),
                        value: v
                    };
                }),
                radius: ['0%', '70%'],
                center: ['50%', '54%']/**,
                label: {
                    normal: {
                        show: false,
                        position: 'center'
                    },
                    emphasis: {
                        show: false,
                        textStyle: {
                            fontSize: '30',
                            fontWeight: 'bold'
                        }
                    }
                }**/
            }],

            grid: {
                left: -10,
                right: 0,
                bottom: 0,
                top: -10
            },
            axisPointer: {
                show: false
            },

            tooltip: {
                formatter: 'X:{b0} Y:{c0}',
                show: false
            },
            animation: false,
            color: chart.colors
        };

        c.setOption(option);
    };

    this.renderProgressChart = function () {

    };

    this.renderGaugeChart = function () {

    };
}