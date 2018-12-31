Tea.context(function () {
    var that = this;

    this.targetType = (this.rewrite.targetType == 2) ? "url" : "proxy";
    this.pattern = this.rewrite.pattern;
    this.redirectMode = this.rewrite.redirectMode;
    this.proxyId = this.rewrite.proxyId;
    this.on = this.rewrite.on;
    this.replace = this.rewrite.replace;

    this.$delay(function () {
        this.$find("form input[name='pattern']").focus();
    });

    this.updateSuccess = function () {
        alert("保存成功");
        window.location = this.from;
    };

    this.advancedOptionsVisible = false;
    this.showAdvancedOptions = function () {
        this.advancedOptionsVisible = !this.advancedOptionsVisible;
    };

    this.conds = this.rewrite.conds.$map(function (k, v) {
        return {
            "param": v.param,
            "value": v.value,
            "op": v.operator,
            "description": that.operators.$find(function (k1, v1) {
                return v.operator == v1.op;
            }).description
        };
    });
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