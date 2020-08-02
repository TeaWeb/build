Tea.context(function () {
    var that = this;
    this.modeDescription = "";

    this.changeMode = function () {
        this.modeDescription = this.modes.$find(function (k, v) {
            return v.mode == that.websocket.forwardMode;
        }).description;
    };
    this.changeMode();

    this.submitSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    /**
     * 域名
     */
    this.originAdding = false;
    this.addingOriginName = "";

    this.addOrigin = function () {
        this.originAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingOriginName']").focus();
        });
    };

    this.confirmAddOrigin = function () {
        this.addingOriginName = this.addingOriginName.trim();
        if (this.addingOriginName.length == 0) {
            alert("域名不能为空");
            this.$find("form input[name='addingOriginName']").focus();
            return;
        }
        this.origins.push(this.addingOriginName);
        this.cancelOriginAdding();
    };

    this.cancelOriginAdding = function () {
        this.originAdding = !this.originAdding;
        this.addingOriginName = "";
    };

    this.removeOrigin = function (index) {
        this.origins.$remove(index);
    };


    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;
    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };
});