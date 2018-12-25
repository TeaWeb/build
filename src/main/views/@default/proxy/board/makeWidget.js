Tea.context(function () {
    var editor = null;

    this.$delay(function () {
        this.loadEditor();
    });

    this.loadEditor = function () {
        editor = CodeMirror(document.getElementById("editor"), {
            theme: "idea",
            lineNumbers: true,
            styleActiveLine: true,
            matchBrackets: true,
            value: this.code,
            height: "auto",
            viewportMargin: Infinity,
            scrollbarStyle: null,
            lineWrapping: true,
            indentUnit: 4,
            indentWithTabs: true
        });
    };
});