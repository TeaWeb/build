Tea.context(function () {
    this.advancedOptionsVisible = false;

    this.$delay(function () {
        this.$find("form input[name='address']").focus();
    });

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    this.submitSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };
});