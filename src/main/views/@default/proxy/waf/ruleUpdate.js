Tea.context(function () {
	this.$delay(function () {
		this.$find("input[name='name']").focus();
		this.sortable();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/proxy/waf/group?wafId=" + this.config.id + "&groupId=" + this.group.id;
	};

	/**
	 * rules
	 */
	this.isAddingRule = false;
	this.ruleIndex = -1;
	this.rules = [];

	if (this.oldRules != null) {
		var that = this;
		this.rules = this.oldRules.$map(function (k, v) {
			return {
				"checkpoint": that.checkpoints.$find(function (k1, v1) {
					return v1.prefix == v.prefix;
				}),
				"param": v.param,
				"operator": that.operators.$find(function (k1, v1) {
					return v1.code == v.operator;
				}),
				"value": v.value
			}
		});
		this.$delay(function () {
			this.refreshTestParams();
		});
	}

	this.testParams = [];
	this.testVisible = false;

	this.addRule = function () {
		this.isAddingRule = true;
		this.ruleIndex = -1;
		this.checkpointPrefix = "";
		this.checkpoint = null;
		this.operatorCode = "match";
		this.changeOperator();
		this.value = "";
	};

	this.cancelAdding = function () {
		this.isAddingRule = false;
		this.ruleIndex = -1;
	};

	this.editRule = function (index) {
		this.ruleIndex = index;
		this.isAddingRule = true;
		this.checkpointPrefix = this.rules[index].checkpoint.prefix;
		this.checkpoint = this.rules[index].checkpoint;
		this.operatorCode = this.rules[index].operator.code;
		this.operator = this.rules[index].operator;
		this.value = this.rules[index].value;
		this.checkpointParam = this.rules[index].param;
	};

	this.confirmAdding = function () {
		if (this.checkpoint == null) {
			alert("请选择一个参数");
			return;
		}
		this.isAddingRule = false;
		if (this.ruleIndex < 0) {
			this.rules.push({
				"checkpoint": this.checkpoint,
				"param": this.checkpointParam,
				"operator": this.operator,
				"value": this.value
			});
		} else {
			this.rules[this.ruleIndex] = {
				"checkpoint": this.checkpoint,
				"param": this.checkpointParam,
				"operator": this.operator,
				"value": this.value
			};
			this.ruleIndex = -1;
		}
		this.refreshTestParams();
		this.isTested = false;
	};

	this.removeRule = function (index) {
		this.rules.$remove(index);
		this.ruleIndex = -1;
		this.isAddingRule = false;
		this.refreshTestParams();
	};

	this.refreshTestParams = function () {
		var testParams = [];
		this.rules.$each(function (k, rule) {
			if (testParams.$find(function (k1, v1) {
				return rule.checkpoint.prefix == v1.prefix && rule.param == v1.param;
			}) != null) {
				return;
			}
			testParams.push({
				"name": rule.checkpoint.name,
				"prefix": rule.checkpoint.prefix,
				"param": rule.param,
				"value": rule.value,
				"description": rule.checkpoint.description
			});
		});
		this.testParams = testParams;
	};

	/**
	 * checkpoint
	 */
	this.checkpointPrefix = "";
	this.checkpoint = null;
	this.checkpointParam = "";
	this.changeCheckpoint = function () {
		if (this.checkpointPrefix.length == 0) {
			this.checkpoint = null;
			return;
		}
		var that = this;
		this.checkpoint = this.checkpoints.$find(function (k, v) {
			return v.prefix == that.checkpointPrefix;
		});
		this.checkpointParam = "";
	};

	/**
	 * operator
	 */
	this.operatorCode = "match";

	this.changeOperator = function () {
		var that = this;
		this.operator = this.operators.$find(function (k, v) {
			return v.code == that.operatorCode;
		});
	};
	this.changeOperator();

	/**
	 * 对比值
	 */
	this.value = "";

	/**
	 * 关系
	 */
	this.connectorValue = this.set.connector;
	this.connectorDescription = "";

	this.changeConnector = function () {
		var that = this;
		this.connectorDescription = this.connectors.$find(function (k, v) {
			return v.value == that.connectorValue;
		}).description;
	};
	this.changeConnector();

	/**
	 * action
	 */
	this.action = this.set.action;

	/**
	 * Test
	 */
	this.matchedIndex = -1;
	this.breakIndex = -1;
	this.isTested = false;

	this.showTestForm = function () {
		this.testVisible = !this.testVisible;
	};

	this.test = function () {
		this.isTested = true;
		this.matchedIndex = -1;
		var form = new FormData(this.$find("#rule-form")[0]);
		form.append("test", "1");
		this.$post("/proxy/waf/group/rule/update")
			.params(form)
			.success(function (resp) {
				this.matchedIndex = resp.data.matchedIndex;
				this.breakIndex = resp.data.breakIndex;
				resp.data.matchLogs.$each(function (k, v) {
					console.log(v);
				});
			});
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		var box = this.$find(".rules-box")[0];
		var that = this;
		Sortable.create(box, {
			draggable: ".label",
			handle: ".label",
			onStart: function () {

			},
			onUpdate: function (event) {
			}
		});
	};
});