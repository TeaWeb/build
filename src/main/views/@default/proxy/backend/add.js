Tea.context(function () {
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    this.addSuccess = function () {
        alert("保存成功");
        window.location = "/proxy/backend?server=" + this.proxy.filename;
    };
});