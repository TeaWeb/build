Tea.context(function () {
    this.$delay(function () {
        this.$find("form input[name='pass']").focus();
    });

    this.addSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    /**
     * 更多选项
     */
    this.advancedOptionsVisible = false;
    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    /**
     * 参数
     */
    var params = [];
    var nameZhMap = {
        "SCRIPT_FILENAME": "脚本文件",
        "DOCUMENT_ROOT": "文档目录",
        "HTTP_HOST": "主机名",
        "REQUEST_METHOD": "请求方法",
        "CONTENT_TYPE": "文档类型",
        "SERVER_SOFTWARE": "软件版本"
    };
    for (var key in this.fastcgi.params) {
        if (this.fastcgi.params.hasOwnProperty(key)) {
            var nameZh = (typeof(nameZhMap[key]) == "string") ? nameZhMap[key] : "";
            if (nameZh.length > 0) {
                params.push({
                    "name": key,
                    "value": this.fastcgi.params[key],
                    "nameZh": nameZh
                });
            }
        }
    }

    // 让自定义的放到下面去
    for (var key in this.fastcgi.params) {
        if (this.fastcgi.params.hasOwnProperty(key)) {
            var nameZh = (typeof(nameZhMap[key]) == "string") ? nameZhMap[key] : "";
            if (nameZh.length == 0) {
                params.push({
                    "name": key,
                    "value": this.fastcgi.params[key],
                    "nameZh": nameZh
                });
            }
        }
    }
    this.params = params;

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
});