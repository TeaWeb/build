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
           value: this.func,
           height: "auto",
           viewportMargin: Infinity,
           scrollbarStyle: null,
           lineWrapping: true,
           indentUnit: 4,
           indentWithTabs: true
       });
   };

    this.results = [];
    this.loading = false;
    this.test = function () {
        this.loading = true;

        this.results = [];
        this.$post(".probeScript")
            .params({
                "isTesting": 1,
                "script": editor.getValue()
            })
            .success(function (resp) {
                this.results = resp.data.apps;
                this.$delay(function () {
                    window.scroll(0, 10000);
                });
            })
            .done(function () {
                this.loading = false;
            });
    };

    this.save = function () {
        this.results = [];
        this.$post(".probeScript")
            .params({
                "isTesting": 0,
                "probeId": this.probeId,
                "script": editor.getValue()
            })
            .success(function () {
                alert("保存成功");
                window.location = "/apps/probes";
            });
    };
});