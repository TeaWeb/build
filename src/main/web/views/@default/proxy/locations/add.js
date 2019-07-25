Tea.context(function () {
	this.location = null;

	this.showSpeical = false;
	this.pattern = "";
	this.type = 1;
	this.root = "";
	this.charset = "";
	this.indexes = [];
	this.on = true;
	this.isCaseInsensitive = false;
	this.isReverse = false;

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/proxy/locations?serverId=" + this.server.id;
	};

	this.$delay(function () {
		this.$find("form input[name='pattern']").focus();
		this.sortable();
	});

	/**
	 * advanced settings
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * special settings
	 */
	this.specialSettingsVisible = this.showSpecial; // 从参数中获取

	this.showSpecialSettings = function () {
		this.specialSettingsVisible = !this.specialSettingsVisible;
	};

	/**
	 * pattern
	 */
	this.patternDescription = "";
	this.changePatternType = function (patternType) {
		this.patternDescription = this.patternTypes.$find(function (k, v) {
			return v.type == patternType;
		}).description;
	};
	this.changePatternType(this.type);

	/**
	 * 匹配条件
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

	this.showCondVariables = function (index, cond) {
		cond.showVariables = !cond.showVariables;
		Vue.set(this.conds, index, cond);
	};

	this.selectCondVariable = function (cond, variable) {
		cond.param = variable.code;
		cond.showVariables = false;
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
			alert("首页文件名不能为空");
			this.$find("form input[name='addingIndexName']").focus();
			return;
		}
		this.indexes.push(this.addingIndexName);
		this.cancelIndexAdding();
	};

	this.cancelIndexAdding = function () {
		this.indexAdding = !this.indexAdding;
		this.addingIndexName = "";
	};

	this.removeLocationIndex = function (index) {
		this.indexes.$remove(index);
	};

	/**
	 * 匹配测试
	 */
	this.testingVisible = false;
	this.testingFinished = false;
	this.testingSuccess = false;
	this.testingMapping = null;
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

		var form = this.$find("#location-form")[0];
		var formData = new FormData(form);
		this.$post("/proxy/locations/test")
			.params(formData)
			.success(function (resp) {
				this.testingMapping = resp.data.mapping;
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

	/**
	 * 单位
	 */
	this.maxBodyUnits = [{
		"code": "k",
		"name": "K"
	}, {
		"code": "m",
		"name": "M"
	}, {
		"code": "g",
		"name": "G"
	}];
	this.maxBodyUnit = "m";

	/**
	 * 压缩级别
	 */
	this.gzipLevels = Array.$range(1, 9);
	this.gzipMinUnits = [
		{
			"code": "b",
			"name": "B"
		},
		{
			"code": "k",
			"name": "K"
		}, {
			"code": "m",
			"name": "M"
		}];
	this.gzipMinUnit = "k";

	/**
	 * 访问日志相关
	 */
	this.enableAccessLog = true;

	/**
	 * 状态页
	 */
	this.pageAdding = false;
	this.addingPage = {
		"status": "",
		"url": ""
	};
	this.editingPageIndex = -1;
	this.shutdownPageOn = false;
	this.shutdownPage = "";

	if (this.pages == null) {
		this.pages = [];
	} else {
		this.pages = this.pages.$map(function (k, v) {
			return {
				"status": v.status[0],
				"url": v.url
			};
		});
	}

	this.addPage = function () {
		this.pageAdding = true;
		this.addingPage = {
			"status": "",
			"url": ""
		};
		this.editingPageIndex = -1;
		this.$delay(function () {
			this.$find("form input[name='addingPageStatus']").focus();
		});
	};

	this.editPage = function (index) {
		this.pageAdding = true;
		this.editingPageIndex = index;
		this.$delay(function () {
			this.$find("form input[name='addingPageName']").focus();
		});
		var page = this.pages[index];
		this.addingPage = {
			"status": page.status,
			"url": page.url
		};
	};

	this.confirmAddPage = function () {
		if (this.addingPage.status.length == 0) {
			alert("请输入状态码");
			this.$find("form input[name='addingPageStatus']").focus();
			return;
		}
		if (this.addingPage.status.length != 3) {
			alert("状态码必须是3位");
			this.$find("form input[name='addingPageStatus']").focus();
			return;
		}
		if (!this.addingPage.status.match(/^[0-9x]+$/)) {
			alert("状态码中只能包含数字或者小写字母x");
			this.$find("form input[name='addingPageStatus']").focus();
			return;
		}
		if (this.addingPage.url.length == 0) {
			alert("请输入URL地址");
			this.$find("form input[name='addingPageURL']").focus();
			return;
		}
		if (this.editingPageIndex > -1) {
			this.pages[this.editingPageIndex] = {
				"status": this.addingPage.status,
				"url": this.addingPage.url
			};
		} else {
			this.pages.push(this.addingPage);
		}
		this.cancelPageAdding();
	};

	this.cancelPageAdding = function () {
		this.pageAdding = false;
		this.addingPageName = "";
		this.editingPageIndex = -1;
	};

	this.removePage = function (index) {
		this.pages.$remove(index);
		this.cancelPageAdding();
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		var that = this;
		[".indexes-box"].$each(function (k, box) {
			var box = that.$find(box)[0];
			if (!box) {
				return;
			}
			Sortable.create(box, {
				draggable: ".label",
				handle: ".handle",
				onStart: function () {

				},
				onUpdate: function (event) {
				}
			});
		});
	};
});