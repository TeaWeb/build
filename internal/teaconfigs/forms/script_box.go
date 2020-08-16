package forms

import (
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
)

type ScriptBox struct {
	Element
}

func NewScriptBox(title string, subTitle string) *ScriptBox {
	return &ScriptBox{
		Element: Element{
			Title:    title,
			Subtitle: subTitle,
		},
	}
}

func (this *ScriptBox) Super() *Element {
	return &this.Element
}

func (this *ScriptBox) Compose() string {
	value := maps.NewMap(this.Value)
	scriptType := value.GetString("scriptType")
	if len(scriptType) == 0 {
		scriptType = "path"
	}

	scriptLang := value.GetString("scriptLang")
	if len(scriptLang) == 0 {
		scriptLang = "shell"
	}

	scriptPath := value.GetString("scriptPath")
	scriptCode := value.GetString("scriptCode")

	this.Javascript = `/**
	 * 脚本
	 */
	var scriptEditor = null;
	this.scriptTab = ` + stringutil.JSONEncode(scriptType) + `;
	this.scriptLang = ` + stringutil.JSONEncode(scriptLang) + `;
	this.scriptPath = ` + stringutil.JSONEncode(scriptPath) + `;
	this.scriptCode = ` + stringutil.JSONEncode(scriptCode) + `;
	if (this.scriptTab == "code") {
		this.$delay(function () {
			this.loadEditor();
		});
	}
	this.scriptLangs = [
		{
			"name": "Shell",
			"code": "shell",
			"value":  (this.scriptLang == "shell" && this.scriptCode.length > 0) ? this.scriptCode : "#!/usr/bin/env bash\n\n# your commands here\n"
		},
		{
			"name": "批处理(bat)",
			"code": "bat",
			"value": (this.scriptLang == "bat") ? this.scriptCode : ""
		},
		{
			"name": "PHP",
			"code": "php",
			"value": (this.scriptLang == "php") ? this.scriptCode : "#!/usr/bin/env php\n\n<?php\n// your PHP codes here"
		},
		{
			"name": "Python",
			"code": "python",
			"value": (this.scriptLang == "python") ? this.scriptCode : "#!/usr/bin/env python\n\n''' your Python codes here '''"
		},
		{
			"name": "Ruby",
			"code": "ruby",
			"value": (this.scriptLang == "ruby") ? this.scriptCode : "#!/usr/bin/env ruby\n\n# your Ruby codes here"
		},
		{
			"name": "NodeJS",
			"code": "nodejs",
			"value": (this.scriptLang == "nodejs") ? this.scriptCode : "#!/usr/bin/env node\n\n// your javascript codes here"
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
		var value = this.scriptLangs.$find(function (k, v) {
			return v.code == lang;
		}).value; 
		switch (lang) {
			case "shell":
				scriptEditor.setValue(value);
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
				scriptEditor.setValue(value);
				var info = CodeMirror.findModeByMIME("text/x-php");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "python":
				scriptEditor.setValue(value);
				var info = CodeMirror.findModeByMIME("text/x-python");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "ruby":
				scriptEditor.setValue(value);
				var info = CodeMirror.findModeByMIME("text/x-ruby");
				if (info != null) {
					scriptEditor.setOption("mode", info.mode);
					CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
					CodeMirror.autoLoadMode(scriptEditor, info.mode);
				}
				break;
			case "nodejs":
				scriptEditor.setValue(value);
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
		var that = this;
		scriptEditor.setValue(this.scriptLangs.$find(function (k, v) {
				return v.code == that.scriptLang;
			}).value);
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
			that.scriptLangs.$find(function (k, v) {
				return v.code == that.scriptLang;
			}).value = scriptEditor.getValue();
		});
	};`

	this.CSS = `/** codemirror **/
/** codemirror **/
.CodeMirror {
    border: 1px solid #eee;
    height: auto!important;
}

.CodeMirror-vscrollbar {
    width: 6px;
    border-radius: 3px!important;
}

.CodeMirror-vscrollbar::-webkit-scrollbar-thumb {
    border-radius: 2px;
}
`

	return `<input type="hidden" name="` + this.Namespace + `_scriptType" :value="scriptTab"/>
					<input type="hidden" name="` + this.Namespace + `_scriptLang" :value="scriptLang"/>
					<div class="ui tabular menu attached small">
						<a class="item" :class="{active:scriptTab == 'path'}" @click.prevent="selectScriptTab('path')">脚本文件</a>
						<a class="item" :class="{active:scriptTab == 'code'}" @click.prevent="selectScriptTab('code')" v-if="!teaDemoEnabled">脚本代码</a>
					</div>
					<div class="ui bottom segment attached" v-show="scriptTab == 'path'">
                    	<input type="text" name="` + this.Namespace + `_scriptPath" v-model="scriptPath"/>
						<p class="comment">如果是Shell脚本，请不要忘记在头部添加 <em>#!脚本解释工具</em>，比如 <em>#!/bin/bash</em></p>
					</div>
					<div class="ui bottom segment attached" v-show="scriptTab == 'code'" style="padding-top:0">
						<div class="ui menu text small">
							<a class="item" v-for="lang in scriptLangs" :class="{active:lang.code == scriptLang}" @click.prevent="selectScriptLang(lang.code)">{{lang.name}}</a>
						</div>
						<textarea name="` + this.Namespace + `_scriptCode" id="script-code-editor" rows="1"></textarea>
						<p class="comment">如果是Shell脚本，请不要忘记在头部添加 <em>#!脚本解释工具</em>，比如 <em>#!/bin/bash</em></p>
					</div>`
}

func (this *ScriptBox) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	return map[string]interface{}{
		"scriptType": req.Form.Get(this.Namespace + "_scriptType"),
		"scriptLang": req.Form.Get(this.Namespace + "_scriptLang"),
		"scriptPath": req.Form.Get(this.Namespace + "_scriptPath"),
		"scriptCode": req.Form.Get(this.Namespace + "_scriptCode"),
	}, false, nil
}
