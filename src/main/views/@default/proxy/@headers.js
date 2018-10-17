Tea.context(function () {
    /**
     * 添加Header
     */
    this.headerAdding = false;
    this.headerAddingName = "";
    this.headerAddingValue = "";

    var prefix = "";
    if (this.locationIndex != null) {
        prefix = ".headers";
    }

    this.addHeader = function () {
        this.headerAdding = !this.headerAdding;
    };

    this.addHeaderSave = function () {
        this.$post(prefix + ".set")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "name": this.headerAddingName,
                "value": this.headerAddingValue
            });
    };

    /**
     * 删除Header
     */
    this.deleteHeader = function (index) {
        if (!window.confirm("确定要删除此Header吗？")) {
            return;
        }
        this.$post(prefix + ".delete")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "index": index
            })
            .success(function () {
                window.location.reload();
            });
    };

    /**
     * 修改Header
     */
    this.editHeader = function (index, header) {
        this.headers.$each(function (_, v) {
            v.editing = false;
        });

        this.headerEditingName = header.name;
        this.headerEditingValue = header.value;

        header.editing = true;
        this.$set(this.headers, index, header);
    };

    this.editHeaderCancel = function (index, header) {
        header.editing = false;
        this.$set(this.headers, index, header);
    };

    this.editHeaderSave = function (index) {
        this.$post(prefix + ".update")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "index": index,
                "name": this.headerEditingName,
                "value": this.headerEditingValue
            })
            .success(function () {
                window.location.reload();
            });
    };

    /**
     * 屏蔽Header
     */
    this.ignoreHeaderAdding = false;
    this.ignoreHeaderAddingName = "";
    this.addIgnoreHeader = function () {
        this.ignoreHeaderAdding = !this.ignoreHeaderAdding;
    };

    this.addIgnoreHeaderSave = function () {
        this.$post(prefix + ".addIgnore")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "name": this.ignoreHeaderAddingName
            });
    };

    this.editIgnoreHeader = function (index, header) {
        header.isEditing = true;
        this.ignoreHeaderEditingName = header.name;
        this.$set(this.ignoreHeaders, index, header);
    };

    this.editIgnoreHeaderSave = function (index, header) {
        this.$post(prefix + ".updateIgnore")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "index": index,
                "name": this.ignoreHeaderEditingName
            })
            .success(function () {
                window.location.reload();
            });
    };

    this.deleteIgnoreHeader = function (index) {
        if (!window.confirm("确定要删除此Header吗？")) {
            return;
        }

        this.$post(prefix + ".deleteIgnore")
            .params({
                "filename": this.filename,
                "locationIndex": this.locationIndex,
                "index": index
            })
            .success(function () {
                window.location.reload();
            });
    };
});