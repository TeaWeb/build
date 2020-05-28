Tea.context(function () {
	var scriptEditor = null;
	var isLoaded = false;

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
		isLoaded = true;
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/notices/medias";
	};

	/**
	 * 名称
	 */
	this.name = "";
	this.rateNoticeVisible = false;

	this.changeName = function (name) {
		if (name.indexOf("短信") > -1 || name.indexOf("钉钉") > -1 || name.indexOf("微信") > -1) {
			this.rateNoticeVisible = true;
		} else {
			this.rateNoticeVisible = false;
		}
	};

	/**
	 * 类型
	 */
	this.mediaType = "email";
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

	/**
	 * 邮箱
	 */
	this.emailUsername = "";
	this.emailUsernameHelp = "";

	this.changeEmailUsername = function () {
		this.emailUsernameHelp = "";
		if (this.emailUsername.indexOf("qq.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://service.mail.qq.com/cgi-bin/help?id=28\" target='_blank'>QQ邮箱相关设置帮助</a>";
		} else if (this.emailUsername.indexOf("163.com") > 0) {
			this.emailUsernameHelp = "，<a href=\"https://help.mail.163.com/faqDetail.do?code=d7a5dc8471cd0c0e8b4b8f4f8e49998b374173cfe9171305fa1ce630d7f67ac22dc0e9af8168582a\" target='_blank'>网易邮箱相关设置帮助</a>";
		}
	};

	this.changeMediaType("email");

	/**
	 * 脚本
	 */
	this.scriptTab = "path";
	this.scriptLang = "shell";
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
				scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
				var info = CodeMirror.findModeByMIME("text/x-sh");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "bat":
				scriptEditor.setValue("");
				break;
			case "php":
				scriptEditor.setValue("#!/usr/bin/env php\n\n<?php\n// your PHP codes here");
				var info = CodeMirror.findModeByMIME("text/x-php");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "python":
				scriptEditor.setValue("#!/usr/bin/env python\n\n''' your Python codes here '''");
				var info = CodeMirror.findModeByMIME("text/x-python");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "ruby":
				scriptEditor.setValue("#!/usr/bin/env ruby\n\n# your Ruby codes here");
				var info = CodeMirror.findModeByMIME("text/x-ruby");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "nodejs":
				scriptEditor.setValue("#!/usr/bin/env node\n\n// your javascript codes here");
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
		scriptEditor.setValue("#!/usr/bin/env bash\n\n# your commands here\n");
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

	/**
	 * 阿里云短信模板
	 */
	this.aliyunSmsTemplateVars = [];
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
	 * 企业微信
	 */
	this.qyWeixinTextFormat = "text";

	/**
	 * 企业微信群机器人
	 */
	this.qyWeixinRobotTextFormat = "text";

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};
});