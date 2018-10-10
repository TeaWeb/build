Tea.context(function () {
    this.switchHttpOn = function () {
        var message = this.proxy.http ? "确定要关闭HTTP吗？" : "确定要开启HTTP吗？";
        if (!window.confirm(message)) {
            return false;
        }

        this.proxy.http = !this.proxy.http;
        if (this.proxy.http) {
            this.$get(".httpOn").params({
                "filename": this.filename
            });
        } else {
            this.$get(".httpOff").params({
                "filename": this.filename
            });
        }
    };

    // 代理说明
    this.proxyDescriptionEditing = false;
    this.editDescription = function () {
        this.proxyDescriptionEditing = !this.proxyDescriptionEditing;
    };

    this.editDescriptionSave = function () {
        this.$post(".updateDescription")
            .params({
                "filename": this.filename,
                "description": this.proxy.description
            });
    };

    /**
     * 域名管理
     */
    this.newName = "";
    this.nameAdding = false;
    this.addName = function () {
        this.nameAdding = true;
    };

    this.addNameSave = function () {
        this.$post(".addName")
            .params({
                "filename": this.filename,
                "name": this.newName
            });
    };

    this.editNameIndex = -1;
    this.editName = function (index, name) {
        this.editNameIndex = index;
    };

    this.editNameSave = function (index, name) {
        this.$post(".updateName").params({
                "filename": this.filename,
                "index": index,
                "name": name
            });
    };

    this.editNameCancel = function () {
        this.editNameIndex = -1;
    };

    this.deleteName = function (index) {
        if (!window.confirm("确定要删除此域名吗？")) {
            return;
        }

        this.$post(".deleteName").params({
            "filename": this.filename,
            "index": index
        });
    };

    /**
     * 监听地址管理
     */
    this.newListen = "";
    this.listenAdding = false;
    this.addListen = function () {
        this.listenAdding = true;
    };

    this.addListenSave = function () {
        this.$post(".addListen")
            .params({
                "filename": this.filename,
                "listen": this.newListen
            });
    };

    this.editListenIndex = -1;
    this.editListen = function (index, listen) {
        this.editListenIndex = index;
    };

    this.editListenSave = function (index, listen) {
        this.$post(".updateListen").params({
            "filename": this.filename,
            "index": index,
            "listen": listen
        });
    };

    this.deleteListen = function (index) {
        if (!window.confirm("确定要删除此域名吗？")) {
            return;
        }

        this.$post(".deleteListen").params({
            "filename": this.filename,
            "index": index
        });
    };

    /**
     * 根目录
     */
    this.rootEditing = false;
    var root = "";
    this.editRoot = function () {
        this.rootEditing = !this.rootEditing;
        if (this.rootEditing) {
            root = this.proxy.root;
        } else {
            this.proxy.root = root;
        }
    };

    this.editRootSave = function () {
        this.$post("/proxy/updateRoot")
            .params({
                "filename": this.filename,
                "root": this.proxy.root
            })
            .success(function () {
                this.rootEditing = false;
            });
    };
});