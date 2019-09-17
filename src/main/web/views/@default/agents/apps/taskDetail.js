Tea.context(function () {
	this.from = encodeURIComponent(window.location.toString());

	this.$delay(function () {
		this.loadEditor();
	});

	/**
	 * 编辑器
	 */
	this.loadEditor = function () {
		var editor = CodeMirror(document.getElementById("code-editor"), {
			theme: "idea",
			lineNumbers: false,
			value: this.task.script,
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

		var info = CodeMirror.findModeByMIME("text/x-sh");
		if (info != null) {
			editor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(editor, info.mode);
		}
	};

	/**
	 * 执行任务
	 */
	this.runTask = function () {
		if (!window.confirm("确定要手动执行一次此任务吗？")) {
			return;
		}

		this.$get("/agents/apps/runTask")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"taskId": this.task.id
			})
			.success(function () {
				alert("已向Agent发送执行请求，稍后请在\"日志\"界面查看执行结果");
			});
	};

	/**
	 * 启用和停止某个任务
	 */
	this.enableTask = function () {
		this.$post(".taskOn")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"taskId": this.task.id
			})
			.refresh();
	};

	this.disableTask = function () {
		this.$post(".taskOff")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"taskId": this.task.id
			})
			.refresh();
	};
});