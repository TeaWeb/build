Tea.context(function () {
    this.$delay(function () {
        this.$find("form input[name='name']").focus();
    });

    this.addSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };
});