Tea.context(function () {
    this.$delay(function () {
        this.$find(".code").each(function (k, v) {
            var code = Tea.element(v).text();
            Tea.element(v).text("");
            CodeMirror(v, {
                theme: "idea",
                lineNumbers: true,
                styleActiveLine: false,
                matchBrackets: true,
                value: code,
                readOnly: true,
                height: "auto",
                viewportMargin: Infinity,
                lineWrapping: true,
                highlightFormatting: true,
                mode: "nginx",
                indentUnit: 4,
                indentWithTabs: true
            })
        });
    }, 500);
});