Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='pattern']").focus();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = this.from;
	};

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
	this.changePatternType(this.location.type);

	/**
	 * 请求条件
	 */
	if (this.location.conds == null) {
		this.location.conds = [];
	}
	var that = this;
	this.conds = this.location.conds.$map(function (k, v) {
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
		this.location.index.push(this.addingIndexName);
		this.cancelIndexAdding();
	};

	this.cancelIndexAdding = function () {
		this.indexAdding = !this.indexAdding;
		this.addingIndexName = "";
	};

	this.removeLocationIndex = function (index) {
		this.location.index.$remove(index);
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
	this.maxBodySize = 0;
	if (this.location.maxBodySize.length > 0) {
		this.maxBodyUnit = this.location.maxBodySize[this.location.maxBodySize.length - 1];
		this.maxBodySize = this.location.maxBodySize.substring(0, this.location.maxBodySize.length - 1);
	}

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
	if (this.location.gzipMinLength.length > 0) {
		this.gzipMinUnit = this.location.gzipMinLength[this.location.gzipMinLength.length - 1];
		this.location.gzipMinLength = this.location.gzipMinLength.substring(0, this.location.gzipMinLength.length - 1);
	}
});