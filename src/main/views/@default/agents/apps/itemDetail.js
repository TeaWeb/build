Tea.context(function () {
	var scriptEditor = null;
	var that = this;

	this.from = encodeURIComponent(window.location.toString());

	if (this.item.thresholds != null) {
		this.item.thresholds.$each(function (k, v) {
			v.levelName = that.noticeLevels.$find(function (k, v1) {
				return v.noticeLevel == v1.code;
			}).name;
		});
	}

	if (this.item.sourceCode == "script" && this.item.sourceOptions.scriptType == "code") {
		this.$delay(function () {
			this.loadEditor();
		});
	}

	this.loadEditor = function () {
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
		scriptEditor.setValue(this.item.sourceOptions.script);

		var lang = "sh";
		if (this.item.sourceOptions.scriptLang != null && this.item.sourceOptions.scriptLang.length > 0) {
			lang = this.item.sourceOptions.scriptLang;
		}
		var info = CodeMirror.findModeByMIME("text/x-" + lang);
		if (info != null) {
			scriptEditor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(scriptEditor, info.mode);
		}
	};
});