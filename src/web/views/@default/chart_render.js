function ChartRender(charts, eventCallback) {
	var that = this;

	var chartObjects = [];

	charts.$map(function (k, v) {
		if (v.options.id != null) {
			v.id = v.options.id;
		} else {
			v.id = Math.random().toString().replace(".", "_");
		}
		Tea.delay(function () {
			that.render(v);
		});
		return v;
	});

	this.render = function (v) {
		// 处理菜单
		if (v.menus != null) {
			var items = [];
			v.menus.$each(function (k, menu) {
				if (k > 0) {
					items.push({
						"name": "|"
					});
				}
				menu.items.$each(function (k2, item) {
					item.menuIndex = k;
				});
				items.$pushAll(menu.items);
			});
			v.menus = [
				{
					"items": items
				}
			];

			var that = this;
			Tea.delay(function () {
				var itemElements = Tea.element("#chart-menu-box-" + v.id + " .item");
				itemElements.each(function (k, item) {
					if (typeof (item.clickAttached) == "undefined") {
						item.clickAttached = true;
					} else {
						return;
					}
					item = Tea.element(item)
						.bind("click", function () {
							var itemName = item.attr("data-name");
							var itemCode = item.attr("data-code");
							var menuIndex = item.attr("data-menuindex");
							if (itemName == "|") { // 分隔符
								return;
							}

							// 传递选中项
							var events = [];
							events.push({
								"chart": (v.options != null) ? v.options.id : "",
								"event": "click",
								"target": "menu.item",
								"data": {
									"code": itemCode,
									"name": itemName
								}
							});

							// 先前选中的
							var activeItem = Tea.element("#chart-menu-box-" + v.id + " .item[data-active='1']:not([data-menuindex='" + menuIndex + "'])");
							if (activeItem.length > 0) {
								var itemName = activeItem.attr("data-name");
								var itemCode = activeItem.attr("data-code");
								events.push({
									"chart": (v.options != null) ? v.options.id : "",
									"event": "click",
									"target": "menu.item",
									"data": {
										"code": itemCode,
										"name": itemName
									}
								});
							}


							if (typeof (eventCallback) == "function") {
								eventCallback(events);
							}
						});
				});
			});
		}

		// 开始绘制
		var chartType = v.type;
		var f = "render" + chartType[0].toUpperCase() + chartType.substring(1) + "Chart";
		if (typeof (this[f]) != "function") {
			throw new Error("can not find render function '" + f + "(chart)'");
		}
		if (v.html != null && v.html.length > 0) {
			Tea.element("#chart-box-" + v.id).html(v.html);
			return;
		}
		var html = this[f](v);
		if (html != null && html.length > 0) {
			Tea.element("#chart-box-" + v.id).html(html);
		}
	};

	this.renderHtmlChart = function (chart) {
		return chart.html;
	};

	this.renderUrlChart = function (chart) {
		return '<iframe src="' + chart.url + '" frameborder="0" scrolling="yes"></iframe>';
	};

	this.renderLineChart = function (chart) {
		var chartBox = Tea.element("#chart-box-" + chart.id)[0];
		if (chartBox == null) {
			return "";
		}
		var c = echarts.init(chartBox);
		chartObjects.push(c);

		var bottomHeight = (chart.labels == null || chart.labels.length == 0) ? 16 : 20;
		if (chart.lines.$exist(function (k, v) {
			return v.name.length > 0;
		})) {
			bottomHeight += 20;
		}
		var maxValue = 0;
		chart.lines.$each(function (k, v) {
			var m = v.values.$max();
			if (m > maxValue) {
				maxValue = m;
			}
		});
		if (maxValue == 0) {
			maxValue = 1;
		}

		var option = {
			textStyle: {
				fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
			},
			title: {
				text: "",
				subText: "",
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
				},
				axisLabel: {
					rotate: 0
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
				max: (chart.max != 0) ? chart.max : null,
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
				left: this.charactersWidth([(Math.ceil(maxValue * ((maxValue <= 1) ? 1000 : 10))).toString()]) + 10,
				right: 10,
				bottom: bottomHeight,
				top: 16
			},
			tooltip: {
				formatter: 'X:{b0} Y:{c0}',
				show: true,
				trigger: "axis",
				axisPointer: {
					type: "cross"
				}
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
		chartObjects.push(c);

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
			tooltip: {
				formatter: 'X:{b0} Y:{c0}',
				show: true
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

	var clockTimer = null;
	var clockSeconds = 0;

	this.showClockTime = function (chart, timestamp) {
		clockSeconds++;

		var date = new Date((timestamp + clockSeconds) * 1000);
		var hour = date.getHours().toString();
		var minute = date.getMinutes().toString();
		var second = date.getSeconds().toString();
		var time = ((hour.length == 1) ? "0" + hour : hour) + ":" + ((minute.length == 1) ? "0" + minute : minute) + ":" + ((second.length == 1) ? "0" + second : second);
		Tea.element("#time-" + chart.id).html(time);
	};

	this.renderClockChart = function (chart) {
		var canvasId = "canvas-" + chart.id;
		setTimeout(function () {
			var timestamp = chart.timestamp;
			var diff = new Date().getTime() / 1000 - timestamp;
			var options = {
				rimColour: "#ccc",
				colour: "rgba(255, 0, 0, 0.2)",
				rim: 2,
				markerType: "dot",
				markerDisplay: true,
				addHours: parseInt(diff / 3600, 10),
				addMinutes: parseInt(diff % 3600 / 60, 10),
				addSeconds: parseInt(diff % 60, 10)
			};
			var myClock = new clock(canvasId, options);
			try {
				myClock.start();
			} catch (e) {
			}

			if (clockTimer != null) {
				clearInterval(clockTimer);
				clockSeconds = 0;
			}
			clockTimer = setInterval(function () {
				that.showClockTime(chart, timestamp);
			}, 1000);

			that.showClockTime(chart, timestamp);
		});
		return "<div style=\"position: relative\"> \
				<canvas id=\"" + canvasId + "\" style=\"width:20em;display: block; margin: 0 auto\"></canvas> \
				<div style='text-align:center;margin-top:1em' id='" + "time-" + chart.id + "'></div></div>";
	};

	this.renderStackbarChart = function (chart) {
		var chartBox = Tea.element("#chart-box-" + chart.id)[0];
		if (chartBox == null) {
			return "";
		}

		if (chart.options.height != null) {
			chartBox.style.cssText += ";height:" + chart.options.height + "em";
		}

		var c = echarts.init(chartBox);
		chartObjects.push(c);
		var seriesIndexes = [0, 1];
		if (chart.values.length > 0) {
			seriesIndexes = Array.$range(0, chart.values[0].length - 1);
		}
		var option = {
			textStyle: {
				fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
			},
			xAxis: {
				type: "value"
			},
			yAxis: {
				type: "category",
				data: chart.labels
			},
			series: seriesIndexes.$map(function (index, _) {
				return {
					name: "",
					type: 'bar',
					stack: "总量",
					data: chart.values.$map(function (k, v) {
						return v[index];
					}),
					barWidth: 10,
					itemStyle: {
						opacity: 0.8
					}
				};
			}),

			grid: {
				left: this.charactersWidth(chart.labels),
				right: 10,
				bottom: 0,
				top: -10
			},
			axisPointer: {
				show: false
			},

			tooltip: {
				formatter: 'X:{b0} Y:{c0}',
				show: true,
				trigger: "item",
				axisPointer: {
					type: "cross"
				}
			},
			animation: false,
			color: chart.colors
		};

		c.setOption(option);
	};

	this.renderTableChart = function (chart) {
		return "";
	};

	this.charactersWidth = function (labels) {
		var width = 0;
		labels.$each(function (k, v) {
			var span = document.createElement("span");
			span.innerHTML = v;
			span.style.cssText = "visibility:hidden";
			document.body.appendChild(span);
			var w = span.offsetWidth;
			if (w > width) {
				width = w;
			}
			span.parentNode.removeChild(span);
		});
		return width;
	};

	window.addEventListener("resize", function () {
		chartObjects.$each(function (k, v) {
			v.resize();
		});
	});
}