Tea.context(function () {
    this.$delay(function () {
        this.$find("form input[name='name']").focus();
    });

    this.submitSuccess = function (response) {
        alert("保存成功");
        window.location = "/agents/apps?agentId=" + this.agentId
    };

    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };
});