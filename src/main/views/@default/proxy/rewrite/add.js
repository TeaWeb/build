Tea.context(function () {
    this.targetType = "url";
    this.pattern = "";
    this.redirectMode = "p";
    this.proxyId = "";

    this.$delay(function () {
        this.$find("form input[name='pattern']").focus();
    });

    this.addSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    this.advancedOptionsVisible = false;
    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    this.conds = [];
    this.addCond = function () {
        this.conds.push({
            "param": "",
            "op": "eq",
            "value": "",
            "description": ""
        });
        this.changeCondOp(this.conds.$last());
        this.$delay(function () {
            this.$find("form input[name='condParams']").last().focus();
            window.scroll(0, 10000);
        });
    };

    this.changeCondOp = function (cond) {
        cond.description = this.operators.$find(function (k, v) {
            return cond.op == v.op;
        }).description;
    };

    this.removeCond = function (index) {
        this.conds.$remove(index);
    };
});