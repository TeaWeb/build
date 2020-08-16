Tea.context(function () {
	var scriptEditor = null;
	var that = this;

	this.from = encodeURIComponent(window.location.toString());
	if (this.sourceOptions != null) {
		this.selectedSource = this.sourceOptions.summary;
	}

	this.$delay(function () {
		// 数据源
		if (this.sourcePresentation != null && this.sourcePresentation.javascript != null && this.sourcePresentation.javascript.length > 0) {
			eval(this.sourcePresentation.javascript);
		}
	});

	if (this.item.thresholds != null) {
		this.item.thresholds.$each(function (k, v) {
			v.levelName = that.noticeLevels.$find(function (k, v1) {
				return v.noticeLevel == v1.code;
			}).name;
		});
	}

	this.loadCodeEditor = function (scriptLang, scriptCode) {
		if (scriptEditor == null) {
			scriptEditor = CodeMirror(document.getElementById("script-code-editor"), {
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
		}
		scriptEditor.setValue(scriptCode);

		var lang = "sh";
		if (scriptLang != null && scriptLang.length > 0) {
			lang = scriptLang;
		}
		var mimeType = "text/x-" + lang;
		if (lang == "nodejs") {
			mimeType = "text/javascript";
		} else if (lang == "shell") {
			mimeType = "text/x-sh";
		}
		var info = CodeMirror.findModeByMIME(mimeType);
		if (info != null) {
			scriptEditor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(scriptEditor, info.mode);
		}
	};

	/**
	 * 外部监控API
	 */
	this.agentItemAPIVisible = false;

	this.showAgentItemAPI = function () {
		this.agentItemAPIVisible = !this.agentItemAPIVisible;
		if (this.agentItemAPIVisible) {
			this.$delay(function () {
				window.scroll(0, 10000);
			});
		}
	};

	/**
	 * 显示数据源变量
	 */
	this.sourceVariablesVisible = false;

	this.showSourceVariables = function () {
		this.sourceVariablesVisible = !this.sourceVariablesVisible;
	};

	/**
	 * 状态切换
	 */
	this.putOn = function () {
		this.$post(".itemOn")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id
			})
			.refresh();
	};

	this.putOff = function () {
		this.$post(".itemOff")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id
			})
			.refresh();
	};

	/**
	 * 执行一次数据源
	 */
	this.execSource = function (itemId) {
		if (!window.confirm("确定要立即从数据源采集数据吗？")) {
			return;
		}
		this.$post(".execItemSource")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id
			})
			.success(function () {
				alert("执行成功，请在数值记录里查看获取的数据");
			});
	};
});