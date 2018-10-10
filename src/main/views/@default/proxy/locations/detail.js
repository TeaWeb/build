Tea.context(function () {
    this.switchOn = function () {
        this.location.on = !this.location.on;

        if (this.location.on) {
            this.$post(".on")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".off")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.patternEditing = false;
    this.editPattern = function () {
        this.patternEditing = !this.patternEditing;
    };

    this.editPatternSave = function () {
        this.$post(".updatePattern")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "pattern": this.location.pattern
            });
    };

    this.typeEditing = false;
    this.editType = function () {
        this.typeEditing = !this.typeEditing;
    };

    this.editTypeSave = function () {
        this.$post(".updateType")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "typeId": this.location.type
            });
    };

    this.reverse = function () {
        this.location.reverse = !this.location.reverse;

        if (this.location.reverse) {
            this.$post(".updateReverse")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "reverse": 1
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".updateReverse")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "reverse": 0
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.switchCaseInsensitive = function () {
        this.location.caseInsensitive = !this.location.caseInsensitive;

        if (this.location.caseInsensitive) {
            this.$post(".updateCaseInsensitive")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "caseInsensitive": 1
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post(".updateCaseInsensitive")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "caseInsensitive": 0
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    /**
     * 重写规则
     */
    this.rewriteAdding = false;
    this.addingPattern = "";
    this.addingReplace = "";
    this.targetType = "url";
    this.proxyId = "";

    this.location.rewrite = this.location.rewrite.$map(function (k, rewrite) {
        if (/^proxy:\/\//.test(rewrite.replace)) {
            var index = rewrite.replace.indexOf("/", "proxy://".length);
            rewrite.proxy = rewrite.replace.substring(0, index);
            rewrite.replace = rewrite.replace.substring(index);
            rewrite.type = "proxy";
            rewrite.proxyId = rewrite.proxy.substr("proxy://".length);
        } else {
            rewrite.proxy = "";
            rewrite.type = "url";
            rewrite.proxyId = "";
        }
        return rewrite;
    });

    this.addRewrite = function () {
        this.rewriteAdding = !this.rewriteAdding;
    };

    this.cancelRewrite = function () {
        this.rewriteAdding = false;
    };

    this.saveRewrite = function () {
        this.$post("/proxy/rewrite/add")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "pattern": this.addingPattern,
                "replace": this.addingReplace,
                "targetType": this.targetType,
                "proxyId": this.proxyId
            });
    };

    this.deleteRewrite = function (index) {
        if (!window.confirm("确定要删除此重写规则吗？")) {
            return;
        }
        this.$post("/proxy/rewrite/delete")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "rewriteIndex": index
            });
    };

    /**
     * 修改重写规则
     */
    this.editRewrite = function (rewrite, index) {
        var that = this;
        this.location.rewrite.$each(function (k, v) {
            if (k == index) {
                if (typeof(rewrite.isEditing) == "undefined") {
                    rewrite.isEditing = true;
                } else {
                    rewrite.isEditing = !rewrite.isEditing;
                }
                that.$set(that.location.rewrite, index, rewrite);
            } else {
                v.isEditing = false;
                that.$set(that.location.rewrite, k, v);
            }
        });
    };

    this.switchRewriteIndex = function (rewrite, index) {
        rewrite.on = !rewrite.on;
        if (rewrite.on) {
            this.$post("/proxy/rewrite/on")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "rewriteIndex": index
                })
                .fail(function () {
                    window.location.reload();
                });
        } else {
            this.$post("/proxy/rewrite/off")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex,
                    "rewriteIndex": index
                })
                .fail(function () {
                    window.location.reload();
                });
        }
    };

    this.cancelEditRewrite = function (rewrite, index) {
        rewrite.isEditing = false;
        this.$set(this.location.rewrite, index, rewrite);
    };

    this.updateRewrite = function (rewrite, index) {
        this.$post("/proxy/rewrite/update")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "rewriteIndex": index,
                "pattern": rewrite.pattern,
                "replace": rewrite.replace,
                "targetType": rewrite.type,
                "proxyId": rewrite.proxyId
            });
    };

    /**
     * fastcgi设置
     */
    this.fastcgiAdding = false;
    this.newFastcgiOn = true;
    this.newFastcgiPass = "";
    this.newFastcgiTimeout = "";
    this.newFastcgiParams = [
        {
            "name": "SCRIPT_FILENAME",
            "value": "",
            "nameZh": "脚本文件"
        },
        {
            "name": "DOCUMENT_ROOT",
            "value": "",
            "nameZh": "文档目录"
        }
    ];

    this.addFastcgi = function () {
        this.fastcgiAdding = !this.fastcgiAdding;
    };

    this.switchNewFastcgiOn = function () {
        this.newFastcgiOn = !this.newFastcgiOn;
    };

    this.addNewFastcgiParam = function () {
        this.newFastcgiParams.push({
            "name": "",
            "value": "",
            "nameZh": ""
        });
    };

    this.removeNewFastcgiParam = function (index) {
        this.newFastcgiParams.$remove(index);
    };

    this.addFastcgiSave = function () {
        var m = {};
        for (var i = 0; i < this.newFastcgiParams.length; i ++) {
            m[this.newFastcgiParams[i]["name"]] = this.newFastcgiParams[i]["value"];
        }

        this.$post("/proxy/fastcgi/add")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "on": this.newFastcgiOn ? 1 : 0,
                "pass": this.newFastcgiPass,
                "readTimeout": this.newFastcgiTimeout,
                "params": JSON.stringify(m)
            });
    };

    this.switchFastcgiOn = function () {
        this.location.fastcgi.on = !this.location.fastcgi.on;
        if (this.location.fastcgi.on) {
            this.$post("/proxy/fastcgi/on")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                });
        }
        else {
            this.$post("/proxy/fastcgi/off")
                .params({
                    "filename": this.filename,
                    "index": this.locationIndex
                });
        }
    };

    this.deleteFastcgi = function () {
        if (!window.confirm("确定要删除Fastcgi配置吗？")) {
            return;
        }
        this.$post("/proxy/fastcgi/delete")
            .params({
                "filename": this.filename,
                "index": this.locationIndex
            });
    };

    this.fastcgiPassEditing = false;

    this.editFastcgiPass = function () {
        this.fastcgiPassEditing = !this.fastcgiPassEditing;
    };

    this.editFastcgiPassSave = function () {
        this.$post("/proxy/fastcgi/updatePass")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "pass": this.location.fastcgi.pass
            });
    };

    this.deleteFastcgiParam = function (name) {
        if (!window.confirm("确定要删除此参数吗？")) {
            return;
        }
        this.$post("/proxy/fastcgi/deleteParam")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "name": name
            })
            .success(function () {
                this.$delete(this.location.fastcgi.params, name);
            });
    };

    this.fastcgiParamAdding = false;
    this.fastcgiNewParamName = "";
    this.fastcgiNewParamValue = "";

    this.addFastcgiParam = function () {
        this.fastcgiParamAdding = !this.fastcgiParamAdding;
    };

    this.addFastcgiParamSave = function () {
        this.$post("/proxy/fastcgi/addParam")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "name": this.fastcgiNewParamName,
                "value": this.fastcgiNewParamValue
            })
            .success(function () {
                this.$set(this.location.fastcgi.params, this.fastcgiNewParamName, this.fastcgiNewParamValue);
                this.fastcgiParamAdding = false;
                this.fastcgiNewParamName = "";
                this.fastcgiNewParamValue = "";
            });
    };

    this.fastcgiParamEditingName = "";

    this.editFastcgiParam = function (name, value) {
        this.fastcgiParamEditingName = name;
        this.fastcgiNewParamName = name;
        this.fastcgiNewParamValue = value;
    };

    this.cancelFastcgiParamEdit = function () {
        this.fastcgiParamEditingName = "";
    };

    this.editFastcgiParamSave = function (name) {
        this.$post("/proxy/fastcgi/updateParam")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "oldName": name,
                "name": this.fastcgiNewParamName,
                "value": this.fastcgiNewParamValue
            })
            .success(function () {
                this.$delete(this.location.fastcgi.params, name);
                this.$set(this.location.fastcgi.params, this.fastcgiNewParamName, this.fastcgiNewParamValue);
                this.fastcgiParamEditingName = "";
            });
    };

    /**
     * 修改超时时间
     */
    this.fastcgiTimeoutEditing = false;

    this.editFastcgiTimeout = function () {
        this.newFastcgiTimeout = parseInt(this.location.fastcgi.readTimeout.replace("s", ""));
        this.fastcgiTimeoutEditing = !this.fastcgiTimeoutEditing;
    };

    this.editFastcgiTimeoutSave = function () {
        this.$post("/proxy/fastcgi/updateTimeout")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "timeout": this.newFastcgiTimeout
            })
            .success(function () {
                this.location.fastcgi.readTimeout = this.newFastcgiTimeout + "s";
                this.fastcgiTimeoutEditing = false;
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
            root = this.location.root;
        } else {
            this.location.root = root;
        }
    };

    this.editRootSave = function () {
        this.$post("/proxy/locations/updateRoot")
            .params({
                "filename": this.filename,
                "index": this.locationIndex,
                "root": this.location.root
            })
            .success(function () {
                this.rootEditing = false;
            });
    };
});