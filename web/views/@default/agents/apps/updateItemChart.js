Tea.context(function () {
	var that = this;
	this.chartType = this.chart.type;
	this.chartDescription = "";

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = this.from;
	};

	/**
	 * 图表类型
	 */
	this.changeChartType = function () {
		this.chartDescription = this.chartTypes.$find(function (k, v) {
			return v.code == that.chartType;
		}).description;

		this.$delay(function () {
			switch (this.chartType) {
				case "html":
					this.loadHTMLChart();
					break;
				case "javascript":
					this.loadJavascriptChart();
					break;
				case "url":
					this.$find("form input[name='urlURL']").focus();
					break;
			}
		});
	};

	this.changeChartType();

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;
	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * HTML
	 */
	this.loadHTMLChart = function () {
		var editor = CodeMirror.fromTextArea(document.getElementById("html-code-editor"), {
			theme: "idea",
			lineNumbers: true,
			value: "",
			readOnly: false,
			showCursorWhenSelecting: true,
			height: "auto",
			//scrollbarStyle: null,
			viewportMargin: Infinity,
			lineWrapping: true,
			highlightFormatting: false,
			indentUnit: 4,
			indentWithTabs: true
		});

		if (this.chart.type == "html") {
			editor.setValue(this.chart.options.html);
		} else {
			editor.setValue("<div>\n    <!-- 这里写一些HTML内容 -->\n</div>");
		}
		editor.save();

		var info = CodeMirror.findModeByMIME("text/html");
		if (info != null) {
			editor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(editor, info.mode);
		}

		var that = this;
		editor.on("change", function () {
			editor.save();
			that.chartCode = editor.getValue();
		});

		editor.focus();
	};


	/**
	 * 饼图
	 */
	this.pieParam = {
		"varName": "${0}",
		"key": ""
	};

	if (this.chart.options.param != null) {
		this.pieParam.varName = this.chart.options.param;
	} else if (this.chart.options.Param != null) { // 兼容老版本 v0.1.8
		this.pieParam.varName = this.chart.options.Param;
	}

	this.pieLimit = 100;
	if (this.chart.options.limit != null) {
		this.pieLimit = this.chart.options.limit;
	} else if (this.chart.options.Limit != null) { // 兼容老版本 v0.1.8
		this.pieLimit = thsi.chart.options.Limit;
	}


	/**
	 * 线图
	 */
	this.colors = [
		{
			"name": "白色",
			"value": "WHITE"
		},
		{
			"name": "红色",
			"value": "RED"
		},
		{
			"name": "蓝色",
			"value": "BLUE"
		},
		{
			"name": "绿色",
			"value": "GREEN"
		},
		{
			"name": "黄色",
			"value": "YELLOW"
		},
		{
			"name": "棕色",
			"value": "BROWN"
		},
		{
			"name": "粉红",
			"value": "PINK"
		}
	];
	this.lineParams = [{
		"varName": "${0}",
		"isFilled": 0,
		"color": "",
		"key": "",
		"name": "",
		"moreVisible": false
	}];
	this.lineMax = 0;
	if (this.chart.options.max != null) {
		this.lineMax = this.chart.options.max;
	}
	if (this.chart.type == "line") {
		if (this.chart.options.lines != null && this.chart.options.lines.length > 0) {
			this.lineParams = this.chart.options.lines.$map(function (k, v) {
				return {
					"varName": v.param,
					"isFilled": v.isFilled ? 1 : 0,
					"color": v.color,
					"key": "",
					"name": v.name,
					"moreVisible": false
				};
			});
		} else if (this.chart.options.Params != null) {
			// 老的版本
			this.lineParams = this.chart.options.Params.$map(function (k, v) {
				return {
					"varName": v,
					"isFilled": 0,
					"color": "",
					"key": "",
					"name": "",
					"moreVisible": false
				};
			});
		}
	}

	this.addLine = function () {
		this.lineParams.push({
			"varName": "${" + this.lineParams.length + "}",
			"isFilled": 0,
			"color": "",
			"key": "",
			"name": "",
			"moreVisible": false
		});
	};

	this.removeLine = function (index) {
		this.lineParams.$remove(index);
	};

	this.changeValueKey = function (param) {
		if (param.key.length > 0) {
			param.varName = "${" + param.key + "}";
		}
	};

	this.showMoreParamOptions = function (param) {
		param.moreVisible = !param.moreVisible;
	};

	/**
	 * Javascript chart
	 */
	this.loadJavascriptChart = function () {
		var editor = CodeMirror.fromTextArea(document.getElementById("javascript-code-editor"), {
			theme: "idea",
			lineNumbers: true,
			value: "",
			readOnly: false,
			showCursorWhenSelecting: true,
			height: "auto",
			//scrollbarStyle: null,
			viewportMargin: Infinity,
			lineWrapping: true,
			highlightFormatting: false,
			indentUnit: 4,
			indentWithTabs: true
		});

		if (this.chart.type == "javascript") {
			editor.setValue(this.chart.options.code);
		} else {
			editor.setValue("var chart = new charts.HTMLChart();\nchart.html = \"使用Javascript代码来构造图表\";\nchart.render();");
		}
		editor.save();

		var info = CodeMirror.findModeByMIME("text/javascript");
		if (info != null) {
			editor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(editor, info.mode);
		}

		var that = this;
		editor.on("change", function () {
			editor.save();
			that.chartCode = editor.getValue();
		});

		editor.focus();
	};

	/**
	 * 预览
	 */
	this.charts = [];

	this.preview = function () {
		this.charts = [];
		var chartForm = document.getElementById("chart-form");
		var form = new FormData(chartForm);
		this.$post(".previewItemChart")
			.params(form)
			.success(function (resp) {
				// output
				resp.data.output.$each(function (k, v) {
					console.log("[widget]" + v);
				});

				// charts
				this.charts = resp.data.charts;
				new ChartRender(this.charts);
			});
	};

	this.$delay(function () {
		this.preview();
	});
});