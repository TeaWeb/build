Tea.context(function () {
	var scriptEditor = null;
	var isLoaded = false;

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
		isLoaded = true;
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/notices/mediaDetail?mediaId=" + this.media.id;
	};

	/**
	 * 类型
	 */
	this.mediaType = this.media.type;
	this.typeDescription = "";

	this.changeMediaType = function () {
		var that = this;
		this.typeDescription = this.mediaTypes.$find(function (k, v) {
			return v.code == that.mediaType;
		}).description;

		if (!isLoaded) {
			return;
		}
		if (this.mediaType == "email") {
			this.$delay(function () {
				this.$find("form input[name='emailSmtp']").focus();
			});
		} else if (this.mediaType == "webhook") {
			this.$delay(function () {
				this.$find("form input[name='webhookURL']").focus();
			});
		} else if (this.mediaType == "script") {
			this.$delay(function () {
				this.$find("form input[name='scriptPath']").focus();
			});
		}
	};

	/**
	 * webhook
	 */
	this.webhookMethod = "GET";
	if (this.media.type == "webhook") {
		this.webhookMethod = this.media.options.method;
	}

	/**
	 * 脚本
	 */
	this.scriptTab = "path";
	this.scriptLang = "shell";

	if (this.media.type == "script") {
		if (this.media.options.scriptType == "path") {
			this.scriptTab = "path";
		} else {
			this.scriptTab = "code";
			this.scriptLang = this.media.options.scriptLang;
			this.$delay(function () {
				this.loadEditor();
			});
		}
	}

	this.scriptLangs = [
		{
			"name": "Shell",
			"code": "shell"
		},
		{
			"name": "批处理(bat)",
			"code": "bat"
		},
		{
			"name": "PHP",
			"code": "php"
		},
		{
			"name": "Python",
			"code": "python"
		},
		{
			"name": "Ruby",
			"code": "ruby"
		},
		{
			"name": "NodeJS",
			"code": "nodejs"
		}
	];

	this.selectScriptTab = function (tab) {
		this.scriptTab = tab;

		if (tab == "path") {
			this.$delay(function () {
				this.$find("form input[name='scriptPath']").focus();
			});
		} else if (tab == "code") {
			this.$delay(function () {
				this.loadEditor();
			});
		}
	};

	this.selectScriptLang = function (lang) {
		this.scriptLang = lang;
		switch (lang) {
			case "shell":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "shell") {
					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
				}
				var info = CodeMirror.findModeByMIME("text/x-sh");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "bat":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "bat") {
					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("");
				}
				break;
			case "php":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "php") {
					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env php\n\n<?php\n// your PHP codes here");
				}
				var info = CodeMirror.findModeByMIME("text/x-php");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "python":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "python") {

					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env python\n\n''' your Python codes here '''");
				}
				var info = CodeMirror.findModeByMIME("text/x-python");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "ruby":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "ruby") {
					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env ruby\n\n# your Ruby codes here");
				}
				var info = CodeMirror.findModeByMIME("text/x-ruby");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "nodejs":
				if (this.media.type == "script" && this.media.options.scriptType == "code" && this.media.options.scriptLang == "nodejs") {
					scriptEditor.setValue(this.media.options.script);
				} else {
					scriptEditor.setValue("#!/usr/bin/env node\n\n// your javascript codes here");
				}
				var info = CodeMirror.findModeByMIME("text/javascript");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
		}

		scriptEditor.save();
		scriptEditor.focus();
	};

	this.loadEditor = function () {
		if (scriptEditor == null) {
			scriptEditor = CodeMirror.fromTextArea(document.getElementById("script-code-editor"), {
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
		}
		if (this.media.options.script != null && this.media.options.script.length > 0) {
			scriptEditor.setValue(this.media.options.script);
		} else {
			scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
		}
		scriptEditor.save();
		scriptEditor.focus();

		var info = CodeMirror.findModeByMIME("text/x-sh");
		if (info != null) {
			scriptEditor.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(scriptEditor, info.mode);
		}

		scriptEditor.on("change", function () {
			scriptEditor.save();
		});
	};

	/**
	 * 环境变量
	 */
	this.env = [];
	if (this.media.type == "script" && this.media.options.env != null) {
		this.env = this.media.options.env;
	}

	this.envAdding = false;
	this.envAddingName = "";
	this.envAddingValue = "";

	this.addEnv = function () {
		this.envAdding = !this.envAdding;
		this.$delay(function () {
			this.$find("form input[name='envAddingName']").focus();
		});
	};

	this.confirmAddEnv = function () {
		if (this.envAddingName.length == 0) {
			alert("请输入变量名");
			this.$find("form input[name='envAddingName']").focus();
		}
		this.env.push({
			"name": this.envAddingName,
			"value": this.envAddingValue
		});
		this.envAdding = false;
		this.envAddingName = "";
		this.envAddingValue = "";
	};

	this.removeEnv = function (index) {
		this.env.$remove(index);
	};

	this.cancelEnv = function () {
		this.envAdding = false;
	};

	/**
	 * 时间设置
	 */
	this.timeFromHour = 0;
	this.timeFromMinute = 0;
	this.timeFromSecond = 0;
	this.timeToHour = 23;
	this.timeToMinute = 59;
	this.timeToSecond = 59;
	if (this.media.timeFrom.length > 0) {
		var pieces = this.media.timeFrom.split(":");
		this.timeFromHour = parseInt(pieces[0], 10);
		this.timeFromMinute = parseInt(pieces[1], 10);
		this.timeFromSecond = parseInt(pieces[2], 10);
	}
	if (this.media.timeTo.length > 0) {
		var pieces = this.media.timeTo.split(":");
		this.timeToHour = parseInt(pieces[0], 10);
		this.timeToMinute = parseInt(pieces[1], 10);
		this.timeToSecond = parseInt(pieces[2], 10);
	}

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});