Tea.context(function () {
    this.username = "";
    this.password = "";

    if (this.teaDemoEnabled) {
        this.username = "admin";
        this.password = "123456";
    }

    this.$delay(function () {
        this.$find("form input[name='username']").focus();
    });
});