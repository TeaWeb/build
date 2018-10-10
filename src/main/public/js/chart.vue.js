Tea.context(function () {
    var that = this;

    this.CHART = {
        "id": ""
    };

    this.CHART.id = function (options) {
        if (typeof(options["id"]) == "string") {
            return "chart-" + options["id"];
        }
        return "chart-" + Math.random().toString().replace(".", "");
    };

    this.CHART.progressBar = function (options) {
        var chartId = this.id(options);
        var chartBox = Tea.element("#" + chartId);
        var inner =  '   <div class="ui progress ' + options.color + ' tiny">' +
            '       <div class="bar" style="width:' + options.value.toString() + '%"></div>' +
            '       <div class="label">' + options.name + " &nbsp; <em>" + ((options.detail.length > 0) ? "(" + options.detail + ")" : "" ) + '</em></div>' +
            '   </div>';
        if (chartBox.length == 0) {
            return '<div class="chart-box progress" id="' + chartId + '">' + inner + '</div>';
        } else {
            chartBox.html(inner);
        }
        return "";
    };

    this.CHART.line = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chartBox = document.getElementById(chartId);
            if (chartBox == null) {
                return;
            }
            var chart = echarts.init(chartBox);
            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: -4,
                    x: "center",
                    textStyle: {
                        fontSize: 12,
                        fontWeight: "bold",
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    }
                },
                legend: {
                    data: options.lines.$map(function (_, line) {
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
                    data: options.labels,
                    axisTick: {
                        show: options.xShowTick
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
                    max: (options.max != 0 ) ? options.max : null,
                    splitNumber: (options.yTickCount >= 1) ? options.yTickCount : null,
                    axisTick: {
                        show: options.yShowTick
                    }
                },
                series: options.lines.$map(function (_, line) {
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
                            opacity: line.filled ? 0.5 : 0
                        }
                    };
                }),
                grid: {
                    left: 30,
                    right: 0,
                    bottom: 50,
                    top: 20
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

            chart.setOption(option);
        });
        return '<div class="chart-box line" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.pie = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chartBox = document.getElementById(chartId);
            if (chartBox == null) {
                return;
            }
            var chart = echarts.init(chartBox);

            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: 1,
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
                    data: options.labels,
                    itemWidth: 6,
                    itemHeight: 6,
                    textStyle: {
                        fontSize: 10
                    }
                },
                xAxis: {
                    data: []
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'pie',
                    data: options.values.$map(function (k, v) {
                        return {
                            name: options.labels.$get(k),
                            value: v
                        };
                    }),
                    radius: ['68%', '75%'],
                    center: ['50%', '56%'],
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
                    }
                }],

                grid: {
                    left: -2,
                    right: 0,
                    bottom: 0,
                    top: 0
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

            chart.setOption(option);
        });
        return '<div class="chart-box pie" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.gauge = function (options) {
        var chartId = this.id(options);

        setTimeout(function () {
            var chartBox = document.getElementById(chartId);
            if (chartBox == null) {
                return;
            }
            var chart = echarts.init(chartBox);

            var option = {
                textStyle: {
                    fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                },
                title: {
                    text: options.name,
                    top: 1,
                    bottom: 0,
                    x: "center",
                    textStyle: {
                        fontSize: 12,
                        fontWeight: "bold",
                        fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
                    }
                },
                legend: {
                    data: [""]
                },
                xAxis: {
                    data: []
                },
                yAxis: {},
                series: [{
                    name: '',
                    type: 'gauge',
                    min: options.min,
                    max: options.max,

                    data: [
                        {
                            "name": options.detail,
                            "value": Math.round(options.value * 100) / 100
                        }
                    ],
                    radius: "80%",
                    center: ["50%", "60%"],

                    splitNumber: 5,
                    splitLine: {
                        length: 6
                    },

                    axisLine: {
                        lineStyle: {
                            width: 8
                        }
                    },
                    axisTick: {
                        show: true
                    },
                    axisLabel: {
                        formatter: function (v) {
                            return v + options.unit
                        },
                        textStyle: {
                            fontSize: 8
                        }
                    },
                    detail: {
                        formatter: function (v) {
                            return v + options.unit
                        },
                        textStyle: {
                            fontSize: 12
                        }
                    },

                    pointer: {
                        width: 2
                    }
                }],

                grid: {
                    left: -2,
                    right: 0,
                    bottom: 0,
                    top: 0
                },
                axisPointer: {
                    show: false
                },
                tooltip: {
                    formatter: 'X:{b0} Y:{c0}',
                    show: false
                },
                animation: true
            };

            chart.setOption(option);
        });
        return '<div class="chart-box gauge" id="' + chartId + '">&nbsp;</div>'
    };

    this.CHART.table = function (options) {
        var chartId = this.id(options);
        var chartBox = Tea.element("#" + chartId);
        var s = '<table class="ui table selectable">';
        // s += '<thead><tr><th colspan="' + ((options.rows.length > 0) ? options.rows[0].columns.length : 1) + '">' + options.name  + '</th></tr></thead>';
        if (options.rows.length == 0) {
            s += '<tr><td>还没有数据</td></tr>';
        } else {
            options.rows.$each(function (_, row) {
                s += "<tr>";
                row.columns.$each(function (_, column) {
                    if (column.width > 0) {
                        s += "<td width=\"" + column.width + "%\">" + column.text + "</td>";
                    } else {
                        s += "<td>" + column.text + "</td>";
                    }
                });
                s += "</tr>";
            });
        }
        s += '</table>';

        if (chartBox.length > 0) {
            chartBox.html(s);
        } else {
            s = '<div class="chart-box table" id="' + chartId + '">' + s + '</div>';
        }

        return s;
    };

    this.CHART.updateWidgetGroups = function (newGroups) {
        if (that.widgetGroups == null || that.widgetGroups.length == 0) {
            that.widgetGroups = newGroups;
            return;
        }

        newGroups.$each(function (_, group) {
            group.widgets.$each(function (_, widget) {
                widget.charts.$each(function (_, chart) {
                    switch (chart.type) {
                        case "line":
                            that.CHART.line(chart);
                            break;
                        case "progressBar":
                            that.CHART.progressBar(chart);
                            break;
                        case "pie":
                            that.CHART.pie(chart);
                            break;
                        case "gauge":
                            that.CHART.gauge(chart);
                            break;
                        case "table":
                            that.CHART.table(chart);
                            break;
                    }
                });
            });
        });
    };
});