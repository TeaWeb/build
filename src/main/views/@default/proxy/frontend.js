Tea.context(function () {
    this.$delay(function () {
        this.$find("pre code.hljs").each(function (k, v) {
            hljs.highlightBlock(v);
        });
    }, 500);
});