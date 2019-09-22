Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
		this.sortable();
	});

	/**
	 * 提交成功
	 */
	this.submitSuccess = function () {
		alert("保存成功");
		window.location = this.from;
	};

	/**
	 * 数据源分类
	 */
	this.selectedCategory = "basic";

	this.selectCategory = function (category) {
		this.selectedCategory = category;
		this.sourceCode = this.sources.$findAll(function (k, v) {
			return v.category == category;
		})[0].code;
		this.changeSource();
	};

	/**
	 * 数据源
	 */
	this.isLoaded = false;
	this.sourceCode = this.sources[0].code;
	this.sourceDescription = "";
	this.defaultThresholds = [];
	this.sourcePlatforms = [];
	this.selectedSource = null;

	this.changeSource = function () {
		var that = this;
		var source = this.sources.$find(function (k, v) {
			return v.code == that.sourceCode;
		});
		this.selectedSource = source;
		this.sourceDescription = source.description;
		if (source.thresholds != null) {
			this.defaultThresholds = source.thresholds;
		}
		if (source.platforms != null) {
			this.sourcePlatforms = source.platforms;
		}

		if (!this.isLoaded) {
			this.isLoaded = true;
			return;
		}

		this.$delay(function () {
			this.$find("input[name^='" + this.sourceCode + "_']:not([type='hidden'])").first().focus();
		});
	};

	this.changeSource();

	this.sourceVariablesVisible1 = false;
	this.sourceVariablesVisible2 = false;

	this.showSourceVariables = function (n) {
		if (n == 1) {
			this.sourceVariablesVisible1 = !this.sourceVariablesVisible1;
		}
		if (n == 2) {
			this.sourceVariablesVisible2 = !this.sourceVariablesVisible2;
		}
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;
	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * 数据格式
	 */
	this.dataFormat = this.dataFormats[0].code;
	this.dataFormatDescription = "";
	this.$delay(function () {
		this.changeDataFormat();
	});

	this.changeDataFormat = function () {
		var that = this;
		this.dataFormatDescription = this.dataFormats.$find(function (k, v) {
			return v.code == that.dataFormat;
		}).description;
	};

	/**
	 * 阈值
	 */
	this.conds = [];
	this.addingCond = null;
	this.condIndex = 0;
	this.condActions = [];

	this.addCond = function () {
		this.addingCond = {
			"id": this.condIndex++,
			"param": "${0}",
			"op": "eq",
			"value": "",
			"description": "",
			"noticeLevel": 2,
			"noticeLevelName": "",
			"noticeMessage": "",
			"isAdding": true,
			"actions": [],
			"maxFails": 1,
			"moreOptions": false
		};
		this.changeCondOp(this.addingCond);
		this.$delay(function () {
			this.$find("form input[name='addingParam']").focus();
		});
	};

	this.changeCondOp = function (cond) {
		cond.description = this.operators.$find(function (k, v) {
			return cond.op == v.op;
		}).description;
	};

	this.removeCond = function (index) {
		if (!window.confirm("确定要删除该阈值设置吗？")) {
			return;
		}
		this.conds.$remove(index);
		this.addingCond = null;
	};

	this.confirmAddingCond = function () {
		if (this.addingCond.param.length == 0) {
			alert("请输入参数");
			this.$find("form input[name='addingParam']").focus();
			return;
		}
		var that = this;
		this.addingCond.noticeLevelName = this.noticeLevels.$find(function (k, v) {
			return v.code == that.addingCond.noticeLevel;
		}).name;
		this.addingCond.isAdding = false;
		this.conds.push(this.addingCond);
		this.addingCond = null;
	};

	this.cancelAddingCond = function () {
		this.addingCond = null;
	};

	this.editCond = function (cond) {
		this.addingCond = {
			"id": cond.id,
			"param": cond.param,
			"op": cond.op,
			"value": cond.value,
			"description": cond.description,
			"noticeLevel": cond.noticeLevel,
			"noticeLevelName": cond.noticeLevelName,
			"noticeMessage": cond.noticeMessage,
			"actions": cond.actions,
			"isAdding": false,
			"maxFails": cond.maxFails,
			"moreOptions": false
		};
	};

	this.saveCond = function () {
		if (this.addingCond.isAdding) {
			this.confirmAddingCond();
			return;
		}

		var index = -1;
		var that = this;
		this.addingCond.noticeLevelName = this.noticeLevels.$find(function (k, v) {
			return v.code == that.addingCond.noticeLevel;
		}).name;

		this.conds.$each(function (k, v) {
			if (v.id == that.addingCond.id) {
				v.param = that.addingCond.param;
				v.op = that.addingCond.op;
				v.value = that.addingCond.value;
				v.description = that.addingCond.description;
				v.noticeLevel = that.addingCond.noticeLevel;
				v.noticeLevelName = that.addingCond.noticeLevelName;
				v.noticeMessage = that.addingCond.noticeMessage;
				v.actions = that.addingCond.actions;
				v.maxFails = that.addingCond.maxFails;
			}
		});
		this.addingCond = null;
	};

	this.addCondAction = function () {
		this.addingCond.actions.push({
			"code": "script",
			"options": {
				"scriptType": "path"
			}
		});
	};

	this.removeCondAction = function (index) {
		this.addingCond.actions.$remove(index);
	};

	this.showCondMoreOptions = function (cond) {
		cond.moreOptions = !cond.moreOptions;
	};

	/**
	 * 默认阈值
	 */
	this.addThreshold = function (threshold) {
		this.conds.push({
			"id": this.condIndex++,
			"param": threshold.param,
			"op": threshold.operator,
			"value": threshold.value,
			"description": "",
			"noticeLevel": threshold.noticeLevel,
			"noticeLevelName": this.noticeLevels.$find(function (k1, v1) {
				return v1.code == threshold.noticeLevel
			}).name,
			"noticeMessage": threshold.noticeMessage,
			"actions": (threshold.actions == null) ? [] : threshold.actions,
			"maxFails": (threshold.maxFails == 0) ? 1 : threshold.maxFails,
			"isAdding": false
		});
	};

	this.addDefaultThresholds = function () {
		var that = this;
		this.defaultThresholds.$each(function (k, v) {
			that.addThreshold(v);
		});
	};

	/**
	 * 阈值拖动
	 */
	this.sortable = function () {
		var box = this.$find(".threshold-box")[0];
		Sortable.create(box, {
			draggable: "span.label",
			handle: ".handle.icon"
		});
	};
});