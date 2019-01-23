Tea.context(function () {
	var that = this;
	this.chartType = "line";
	this.chartDescription = "";

	this.pieLimit = 100;

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

		editor.setValue("<div>\n    <!-- 这里写一些HTML内容 -->\n</div>");
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
	 * 线图
	 */
	this.lineLimit = 60;
	this.lineParams = [{
		"varName": "${0}"
	}];

	this.addLine = function () {
		this.lineParams.push({
			"varName": "${" + this.lineParams.length + "}"
		});
	};

	this.removeLine = function (index) {
		this.lineParams.$remove(index);
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

		editor.setValue("var chart = new charts.HTMLChart();\nchart.html = \"使用Javascript代码来构造图表\";\nchart.render();");
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
});