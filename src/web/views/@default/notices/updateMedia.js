Tea.context(function () {
	var scriptEditor = null;
	var isLoaded = false;

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
		isLoaded = true;

		if (this.media.type == "email") {
			this.changeEmailUsername();
		}
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/notices/mediaDetail?mediaId=" + this.media.id;
	};

	/**
	 * 名称
	 */
	this.rateNoticeVisible = false;

	this.changeName = function (name) {
		if (name.indexOf("短信") > -1 || name.indexOf("钉钉") > -1 || name.indexOf("微信") > -1) {
			this.rateNoticeVisible = true;
		} else {
			this.rateNoticeVisible = false;
		}
	};
	this.changeName(this.media.name);

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
		} else if (this.mediaType == "dingTalk") {
			this.$delay(function () {
				this.$find("form textarea[name='dingTalkWebhookURL']").focus();
			});
		} else if (this.mediaType == "qyWeixin") {
			this.$delay(function () {
				this.$find("form input[name='qyWeixinCorporateId']").focus();
			});
		} else if (this.mediaType == "qyWeixinRobot") {
			this.$delay(function () {
				this.$find("form textarea[name='qyWeixinRobotWebhookURL']").focus();
			});
		}
	};
	this.changeMediaType();

	/**
	 * 邮箱
	 */
	this.emailUsernameHelp = "";

	this.changeEmailUsername = function () {
		this.emailUsernameHelp = "";
		if (this.media.options.username.indexOf("qq.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://service.mail.qq.com/cgi-bin/help?id=28\" target='_blank'>QQ邮箱相关设置帮助</a>";
		} else if (this.media.options.username.indexOf("163.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://help.mail.163.com/faqDetail.do?code=d7a5dc8471cd0c0e8b4b8f4f8e49998b374173cfe9171305fa1ce630d7f67ac22dc0e9af8168582a\" target='_blank'>网易邮箱相关设置帮助</a>";
		}
	};

	/**
	 * webhook
	 */
	this.webhookMethod = "GET";
	this.webhookHeadersAdding = false;
	this.webhookHeaders = [];
	this.webhookHeadersAddingName = "";
	this.webhookHeadersAddingValue = "";

	this.addWebhookHeader = function () {
		this.webhookHeadersAdding = true;
		this.$delay(function () {
			this.$find("form input[name='webhookHeaderName']").focus();
		});
	};

	this.cancelWebhookHeadersAdding = function () {
		this.webhookHeadersAdding = false;
	};

	this.confirmWebhookHeadersAdding = function () {
		this.webhookHeaders.push({
			"name": this.webhookHeadersAddingName,
			"value": this.webhookHeadersAddingValue
		});
		this.webhookHeadersAddingName = "";
		this.webhookHeadersAddingValue = "";
		this.webhookHeadersAdding = false;
	};

	this.removeWebhookHeader = function (index) {
		if (!window.confirm("确定要删除此Header吗？")) {
			return;
		}
		this.webhookHeaders.$remove(index);
	};

	this.webhookContentType = "params";

	this.selectWebhookContentType = function (contentType) {
		this.webhookContentType = contentType;
		this.$delay(function () {
			if (contentType == "params") {

			} else if (contentType == "body") {
				this.$find("form textarea[name='webhookBody']").focus();
			}
		});
	};

	this.webhookParamsAdding = false;
	this.webhookParams = [];
	this.webhookParamsAddingName = "";
	this.webhookParamsAddingValue = "";

	this.addWebhookParam = function () {
		this.webhookParamsAdding = true;
		this.$delay(function () {
			this.$find("form input[name='webhookParamName']").focus();
		});
	};

	this.cancelWebhookParamsAdding = function () {
		this.webhookParamsAdding = false;
	};

	this.confirmWebhookParamsAdding = function () {
		this.webhookParams.push({
			"name": this.webhookParamsAddingName,
			"value": this.webhookParamsAddingValue
		});
		this.webhookParamsAddingName = "";
		this.webhookParamsAddingValue = "";
		this.webhookParamsAdding = false;
	};

	this.removeWebhookParam = function (index) {
		if (!window.confirm("确定要删除此参数吗？")) {
			return;
		}
		this.webhookParams.$remove(index);
	};

	this.webhookBody = "";

	if (this.media.type == "webhook") {
		this.webhookMethod = this.media.options.method;
		if (this.media.options.headers != null) {
			this.webhookHeaders = this.media.options.headers;
		}

		if (this.media.options.contentType == "params") {
			this.webhookContentType = "params";
			if (this.media.options.params != null) {
				this.webhookParams = this.media.options.params;
			}
		}

		if (this.media.options.contentType == "body") {
			this.webhookContentType = "body";
			this.webhookBody = this.media.options.body;
		}
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
			return;
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
	 * 阿里云短信模板
	 */
	this.aliyunSmsTemplateVars = [];
	if (this.media.options.variables != null) {
		this.aliyunSmsTemplateVars = this.media.options.variables;
	}
	this.aliyunSmsTemplateVarAdding = false;
	this.aliyunSmsTemplateVarAddingName = "";
	this.aliyunSmsTemplateVarAddingValue = "";

	this.addAliyunSmsTemplateVar = function () {
		this.aliyunSmsTemplateVarAdding = !this.aliyunSmsTemplateVarAdding;
		this.$delay(function () {
			this.$find("form input[name='aliyunSmsTemplateVarAddingName']").focus();
		});
	};

	this.confirmAddAliyunSmsTemplateVar = function () {
		if (this.aliyunSmsTemplateVarAddingName.length == 0) {
			alert("请输入变量名");
			this.$find("form input[name='aliyunSmsTemplateVarAddingName']").focus();
			return;
		}
		this.aliyunSmsTemplateVars.push({
			"name": this.aliyunSmsTemplateVarAddingName,
			"value": this.aliyunSmsTemplateVarAddingValue
		});
		this.aliyunSmsTemplateVarAdding = false;
		this.aliyunSmsTemplateVarAddingName = "";
		this.aliyunSmsTemplateVarAddingValue = "";
	};

	this.removeAliyunSmsTemplateVar = function (index) {
		this.aliyunSmsTemplateVars.$remove(index);
	};

	this.cancelAliyunSmsTemplateVar = function () {
		this.aliyunSmsTemplateVarAdding = false;
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});