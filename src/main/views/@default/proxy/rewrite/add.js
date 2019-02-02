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

	/**
	 * 重写条件
	 */
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

	/**
	 * 匹配测试
	 */
	this.testingVisible = false;
	this.testingFinished = false;
	this.testingSuccess = false;
	this.testingMapping = null;
	this.testingReplace = "";
	this.testingError = "";

	this.showTesting = function () {
		this.testingVisible = !this.testingVisible;
		if (this.testingVisible) {
			this.$delay(function () {
				this.$find("form input[name='testingPath']").focus();
			});
		}
	};

	this.testSubmit = function () {
		this.testingFinished = false;
		this.testingError = "";
		this.testingMapping = null;
		this.testingReplace = "";

		var form = this.$find("#rewrite-form")[0];
		var formData = new FormData(form);
		this.$post("/proxy/rewrite/test")
			.params(formData)
			.success(function (resp) {
				this.testingMapping = resp.data.mapping;
				this.testingReplace = resp.data.replace;
				this.testingFinished = true;
				this.testingSuccess = true;
			})
			.fail(function (resp) {
				if (resp.message != null && resp.message.length > 0) {
					this.testingError = resp.message;
				}

				this.testingFinished = true;
				this.testingSuccess = false;
			});
	};
});