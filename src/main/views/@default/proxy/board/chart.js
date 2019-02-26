Tea.context(function () {
	this.from = encodeURIComponent(window.location.toString());

	this.$delay(function () {
		this.loadJavascriptChart();
		this.loadClipboard();
	});

	/**
	 * Javascript chart
	 */
	this.loadJavascriptChart = function () {
		var editor = CodeMirror.fromTextArea(document.getElementById("javascript-code-editor"), {
			theme: "idea",
			lineNumbers: false,
			value: "",
			readOnly: true,
			showCursorWhenSelecting: true,
			height: "auto",
			//scrollbarStyle: null,
			viewportMargin: Infinity,
			lineWrapping: true,
			highlightFormatting: false,
			indentUnit: 4,
			indentWithTabs: true
		});

		editor.setValue(this.chart.options.code);
		editor.save();

		var info = CodeMirror.findModeByMIME("text/javascript");
		if (info != null) {
			editor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(editor, info.mode);
		}

		//editor.focus();
	};

	/**
	 * 测试
	 */
	this.charts = [];
	this.isLoaded = false;
	this.intervalSeconds = 5;
	this.events = [];

	this.test = function () {
		this.$post("/proxy/board/test")
			.params({
				"serverId": this.server.id,
				"name": this.chart.name,
				"javascriptCode": this.chart.options.code,
				"columns": this.chart.columns,
				"events": JSON.stringify(this.events)
			})
			.success(function (resp) {
				// output
				resp.data.output.$each(function (k, v) {
					console.log("[widget]" + v);
				});

				// charts
				this.charts = resp.data.charts;
				var that = this;
				new ChartRender(this.charts, function (events) {
					that.events.$pushAll(events);
					that.test();
				});
			})
			.fail(function (resp) {
				throw new Error("[widget]" + resp.message);
			})
			.done(function () {
				this.events = [];
				this.isLoaded = true;
			});
	};

	this.$delay(function () {
		this.test();
	});

	/**
	 * 删除这个图表
	 */
	this.deleteChart = function () {
		if (!window.confirm("确定要删除这个图表吗？")) {
			return;
		}
		this.$post("/proxy/board/deleteChart")
			.params({
				"serverId": this.server.id,
				"widgetId": this.widget.id,
				"chartId": this.chart.id
			})
			.success(function () {
				alert("删除成功");
				window.location = "/proxy/board/charts?serverId=" + this.server.id + "&boardType=" + this.boardType;
			});
	};

	/**
	 * 剪切板
	 */
	this.loadClipboard = function () {
		new ClipboardJS('[data-clipboard-target]');

		this.$find("[data-clipboard-target]").bind("click", function () {
			var em =  Tea.element(this).find("em")[0];
			em.style.cssText = "display:inline";
			setTimeout(function () {
				em.style.cssText = "display:none";
			}, 1000);
		});
	};
});