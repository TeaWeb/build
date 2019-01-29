Tea.context(function () {
	var scriptEditor = null;
	var that = this;

	this.from = encodeURIComponent(window.location.toString());

	if (this.media.type == "script" && this.media.options.scriptType == "code") {
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
		scriptEditor.setValue(this.media.options.script);

		var lang = "shell";
		if (this.media.options.scriptLang != null && this.media.options.scriptLang.length > 0) {
			lang = this.media.options.scriptLang;
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
});