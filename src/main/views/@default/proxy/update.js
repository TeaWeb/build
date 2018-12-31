Tea.context(function () {
    this.submitSuccess = function () {
        alert("保存成功");
        window.location = "/proxy/detail?server=" + this.proxy.filename;
    };

    /**
     * 域名
     */
    this.nameAdding = false;
    this.addingNameName = "";

    this.addName = function () {
        this.nameAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingNameName']").focus();
        });
    };

    this.confirmAddName = function () {
        this.addingNameName = this.addingNameName.trim();
        if (this.addingNameName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingNameName']").focus();
            return;
        }
        this.proxy.name.push(this.addingNameName);
        this.cancelNameAdding();
    };

    this.cancelNameAdding = function () {
        this.nameAdding = !this.nameAdding;
        this.addingNameName = "";
    };

    this.removeName = function (index) {
        this.proxy.name.$remove(index);
    };

    /**
     * 监听地址
     */
    this.listenAdding = false;
    this.addingListenName = "";

    this.addListen = function () {
        this.listenAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingListenName']").focus();
        });
    };

    this.confirmAddListen = function () {
        this.addingListenName = this.addingListenName.trim();
        if (this.addingListenName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingListenName']").focus();
            return;
        }
        this.proxy.listen.push(this.addingListenName);
        this.cancelListenAdding();
    };

    this.cancelListenAdding = function () {
        this.listenAdding = !this.listenAdding;
        this.addingListenName = "";
    };

    this.removeListen = function (index) {
        this.proxy.listen.$remove(index);
    };

    /**
     * index
     */
    this.indexAdding = false;
    this.addingIndexName = "";

    this.addIndex = function () {
        this.indexAdding = true;
        this.$delay(function () {
            this.$find("form input[name='addingIndexName']").focus();
        });
    };

    this.confirmAddIndex = function () {
        this.addingIndexName = this.addingIndexName.trim();
        if (this.addingIndexName.length == 0) {
            alert("文件名不能为空");
            this.$find("form input[name='addingIndexName']").focus();
            return;
        }
        this.proxy.index.push(this.addingIndexName);
        this.cancelIndexAdding();
    };

    this.cancelIndexAdding = function () {
        this.indexAdding = !this.indexAdding;
        this.addingIndexName = "";
    };

    this.removeIndex = function (index) {
        this.proxy.index.$remove(index);
    };
});