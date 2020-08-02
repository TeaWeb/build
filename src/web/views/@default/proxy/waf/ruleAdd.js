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
		this.operatorCase = this.rules[index].case;
		this.checkpoint.options = JSON.parse(this.rules[index].options);
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
				"value": this.value,
				"case": this.operatorCase,
				"options": JSON.stringify(this.checkpoint.options)
			});
		} else {
			this.rules[this.ruleIndex] = {
				"checkpoint": this.checkpoint,
				"param": this.checkpointParam,
				"operator": this.operator,
				"value": this.value,
				"case": this.operatorCase,
				"options": JSON.stringify(this.checkpoint.options)
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
	this.operatorCase = "";

	this.changeOperator = function () {
		var that = this;
		this.operator = this.operators.$find(function (k, v) {
			return v.code == that.operatorCode;
		});
		this.operatorCase = (this.operator.case == "yes");
	};
	this.changeOperator();

	/**
	 * 对比值
	 */
	this.value = "";

	/**
	 * 关系
	 */
	this.connectorValue = "or";
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
	this.action = "block";

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
		this.$post("/proxy/waf/group/rule/add")
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

	/**
	 * action options
	 */
	this.actionGroupId = "";
	this.actionSetId = "";

	if (this.config.inbound.length > 0) {
		this.actionGroupId = this.config.inbound[0].id;
		if (this.config.inbound[0].ruleSets.length > 0) {
			this.actionSetId = this.config.inbound[0].ruleSets[0].id;
		}
	}

	var that = this;
	this.groupSets = function (groupId) {
		var group = that.config.inbound.$find(function (k, v) {
			return v.id == groupId;
		});
		if (group == null) {
			return [];
		}
		return group.ruleSets;
	};
});