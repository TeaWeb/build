Tea.context(function () {

    /**
     * 后端地址
     */
    this.backendAdding = false;
    this.newBackendAddress = "";
    this.backendEditing = false;

    this.addBackend = function () {
        this.backendAdding = !this.backendAdding;
    };

    this.addBackendSave = function () {
        this.$post("/proxy/backend/add")
            .params({
                "filename": this.filename,
                "address": this.newBackendAddress
            });
    };

    this.editBackendCancel = function (index, backend) {
        backend.isEditing = !backend.isEditing;
        this.$set(this.proxy.backends, index, backend);
    };

    this.editBackend = function (index, backend) {
        backend.isEditing = true;
        this.backendEditing = !this.backendEditing;

        this.$set(this.proxy.backends, index, backend);
    };

    this.editBackendSave = function (index, backend) {
        this.$post("/proxy/backend/update")
            .params({
                "filename": this.filename,
                "index": index,
                "address": backend.address
            });
    };

    this.deleteBackend = function (index) {
        if (!window.confirm("确定要删除此服务器吗？")) {
            return;
        }
        this.$post("/proxy/backend/delete")
            .params({
                "filename": this.filename,
                "index": index
            });
    };

});