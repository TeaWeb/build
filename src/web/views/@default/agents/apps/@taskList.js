Tea.context(function () {
    this.from = encodeURIComponent(window.location.toString());

    this.$delay(function () {
        this.loadEditor();
        this.testMongo();
    });

    this.loadEditor = function () {
        var that = this;
        this.$find(".task-script-box").each(function (k, v) {
            var code = that.tasks[k].script;
            var editor = CodeMirror(v, {
                theme: "idea",
                lineNumbers: false,
                value: code,
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
        });
    };

    this.deleteTask = function (taskId) {
        if (!window.confirm("确定要删除这个任务吗？")) {
            return;
        }
        this.$post("/agents/apps/deleteTask")
            .params({
                "agentId": this.agentId,
                "appId": this.app.id,
                "taskId": taskId
            })
            .refresh();
    };
});