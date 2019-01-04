Tea.context(function () {
    this.$delay(function () {
        this.$find("form input[name='pass']").focus();
    });

    this.addSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    this.params = [
        {
            "name": "DOCUMENT_ROOT",
            "value": "",
            "nameZh": "文档目录"
        },
        {
            "name": "SCRIPT_FILENAME",
            "value": "",
            "nameZh": "脚本文件"
        }
    ];

    this.addParam = function () {
        this.params.push({
            "name": "",
            "value": "",
            "nameZh": ""
        });
        this.$delay(function () {
            this.$find("form input[name='paramNames']").last().focus();
        });
    };

    this.removeParam = function (index) {
        this.params.$remove(index);
    };

    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;
    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };
});