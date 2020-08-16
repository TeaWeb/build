package forms

import (
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
)

type CodeEditor struct {
	Element

	Lang     string
	Readonly bool
}

func NewCodeEditor(title string, subTitle string) *CodeEditor {
	return &CodeEditor{
		Element: Element{
			Title:    title,
			Subtitle: subTitle,
		},
	}
}

func (this *CodeEditor) Super() *Element {
	return &this.Element
}

func (this *CodeEditor) Compose() string {
	code := this.Value

	this.Javascript = `
	var codeEditor` + this.Code + `Code = ` + stringutil.JSONEncode(code) + `;
	if (codeEditor` + this.Code + `Code  == null) {
		codeEditor` + this.Code + `Code = "SELECT 1";
	}
	var codeEditor` + this.Code + `Lang = ` + stringutil.JSONEncode(this.Lang) + `;
	var codeEditor` + this.Code + ` = null;
	this.$delay(function () {
		Tea.Vue.$watch("selectedSource", function (v) {
			Tea.delay(function () {
				this.loadCodeEditor` + this.Code + `();
			}, 1000);
		});
	});

	this.loadCodeEditor` + this.Code + ` = function () {
		//if (codeEditor` + this.Code + `  == null) {
			codeEditor` + this.Code + ` = CodeMirror.fromTextArea(document.getElementById("code-editor-` + this.Code + `"), {
				//theme: "idea",
				lineNumbers: true,
				value: "",
				readOnly: ` + stringutil.JSONEncode(this.Readonly) + `,
				showCursorWhenSelecting: true,
				height: "auto",
				//scrollbarStyle: null,
				viewportMargin: Infinity,
				lineWrapping: true,
				highlightFormatting: false,
				indentUnit: 4,
				indentWithTabs: true
			});
		//}
		var that = this;
		codeEditor` + this.Code + `.setValue(codeEditor` + this.Code + `Code);
		codeEditor` + this.Code + `.save();

		var info = CodeMirror.findModeByMIME(codeEditor` + this.Code + `Lang);
		if (info != null) {
			codeEditor` + this.Code + `.setOption("mode", info.mode);
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
			CodeMirror.autoLoadMode(codeEditor` + this.Code + `, info.mode);
		}

		codeEditor` + this.Code + `.on("change", function () {
			codeEditor` + this.Code + `.save();
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

	return `<textarea name="` + this.Code + `_code" id="code-editor-` + this.Code + `" rows="1"></textarea>`
}

func (this *CodeEditor) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	return map[string]interface{}{
		"code": req.Form.Get(this.Namespace + "_code"),
	}, false, nil
}
