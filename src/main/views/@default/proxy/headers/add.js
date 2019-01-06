Tea.context(function () {
    this.$delay(function () {
        this.$find("form input[name='name']").focus();
    });

    this.addSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    /**
     * 状态码相关
     */
    this.supportAllStatus = true;
    this.statusList = [ 200, 201, 204, 206, 301, 302, 303, 304, 307, 308 ];
    this.statusAdding = false;
    this.addingStatus = "";

    this.addStatus = function () {
        this.statusAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingStatus']").focus();
        });
    };

    this.cancelAdding = function () {
        this.statusAdding = false;
    };

    this.addStatusConfirm = function (e) {
        if (this.addingStatus.length != 3) {
            alert("状态码必须是3位数字");
            this.$find("form input[name='addingStatus']").focus();
            return;
        }
        if (this.statusList.$contains(this.addingStatus)) {
            alert("状态码已存在");
            this.$find("form input[name='addingStatus']").focus();
            return;
        }
        this.statusList.push(this.addingStatus);
        this.statusAdding = false;
        this.addingStatus = "";

        return false;
    };

    this.deleteStatus = function (index) {
        this.statusList.$remove(index);
    };

    /**
     * 高级选项
     */
    this.advancedOptionsVisible = false;

    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };
});