Tea.context(function () {
    this.allowAllIP = true;
	this.from = encodeURIComponent(window.location.toString());

    this.$delay(function () {
        this.$find("form input[name='name']").focus();
    });

    this.submitSuccess = function (response) {
        alert("保存成功");
        window.location = "/agents/board?agentId=" + response.data.agentId;
    };

    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };
    
    /**
     * 允许的IP
     */
    this.ipAdding = false;
    this.addingIPName = "";
    this.ips = [];

    this.addIP = function () {
        this.ipAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingIPName']").focus();
        });
    };

    this.confirmAddIP = function () {
        this.addingIPName = this.addingIPName.trim();
        if (this.addingIPName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingIPName']").focus();
            return;
        }
        this.ips.push(this.addingIPName);
        this.cancelIPAdding();
    };

    this.cancelIPAdding = function () {
        this.ipAdding = !this.ipAdding;
        this.addingIPName = "";
    };

    this.removeIP = function (index) {
        this.ips.$remove(index);
    };
});