var values = {};

values.Query = function () {
	var query = {
		"action": "",
		"group": null,
		"cond": {
			"error": {
				"eq": ""
			}
		},
		"duration": "",
		"for": null,
		"offset": -1,
		"size": -1,
		"sorts": [],
		"cache": 0
	};

	this.group = function (field) {
		query["group"] = field;
		return this;
	};

	this.attr = function (field, value) {
		return this.op("eq", field, value);
	};

	this.op = function (op, field, value) {
		if (typeof (query.cond[field]) == "undefined") {
			query.cond[field] = {};
		}
		query.cond[field][op] = value;
		return this;
	};

	this.not = function (field, value) {
		return this.op("not", field, value);
	};

	this.lt = function (field, value) {
		return this.op("lt", field, value);
	};

	this.lte = function (field, value) {
		return this.op("lte", field, value);
	};

	this.gt = function (field, value) {
		return this.op("gt", field, value);
	};

	this.gte = function (field, value) {
		return this.op("gte", field, value);
	};

	this.offset = function (offset) {
		query.offset = offset;
		return this;
	};

	this.limit = function (size) {
		query.size = size;
		return this;
	};

	this.cache = function (seconds) {
		query.cache = seconds;
		return this;
	};

	this.asc = function (field) {
		if (field == null) {
			field = "";
		}
		var m = {};
		m[field] = 1;
		query.sorts.push(m);
		return this;
	};

	this.desc = function (field) {
		if (field == null) {
			field = "";
		}
		var m = {};
		m[field] = -1;
		query.sorts.push(m);
		return this;
	};

	this.action = function (action, forField) {
		query["action"] = action;
		query["for"] = forField;
		return this;
	};

	this.execute = function () {
		var cacheKey = null;
		if (query.cache > 0) {
			cacheKey = JSON.stringify({
				"query": query,
				"agentId": (context.agent == null) ? "" : context.agent.id,
				"itemId": (context.item == null) ? "" : context.item.id,
				"appId": (context.app == null) ? "" : context.app.id
			});
			var result = caches.get(cacheKey);
			if (result != null) {
				return result;
			}
		}

		var result = callExecuteQuery(query);
		if (query.cache > 0 && result != null) {
			caches.set(cacheKey, result, query.cache);
		}
		return result;
	};

	this.count = function () {
		return this.action("count")
			.execute();
	};

	this.sum = function (field) {
		return this.action("sum", field)
			.execute();
	};

	this.avg = function (field) {
		return this.action("avg", field)
			.execute();
	};

	this.min = function (field) {
		return this.action("min", field)
			.execute();
	};

	this.max = function (field) {
		return this.action("max", field)
			.execute();
	};

	this.findAll = function () {
		return this.action("findAll")
			.execute();
	};

	this.latest = function (size) {
		if (typeof (size) == "undefined") {
			size = 10;
		}
		return this.action("findAll")
			.desc()
			.limit(size)
			.findAll();
	};

	this.latestValues = function (size) {
		return this.latest(size).$map(function (k, v) {
			return v.value;
		});
	};
};

/**
 * 获取参数值
 */
values.valueOf = function (value, param) {
	var v = param.replace(/(\${[\w\\.]+})/, function (match) {
		var varName = match.substring(2, match.length - 1);
		if (value instanceof Array) {
			var index = parseInt(varName, 10);
			if (index < 0 || index >= value.length) {
				return "";
			}
			return value[index];
		} else if (typeof (value) == "object" && value != null) {
			if (typeof (value[varName]) != "undefined") {
				return value[varName];
			} else {
				var pieces = varName.split(".");
				if (pieces.length > 1) {
					var lastObject = value;
					for (var i = 0; i < pieces.length; i++) {
						var piece = pieces[i];
						if (lastObject != null && typeof (lastObject[piece]) != "undefined") {
							lastObject = lastObject[piece];

							if (i == pieces.length - 1) {
								return lastObject;
							}
						} else {
							break;
						}
					}
				}

				return "";
			}
		}
		if (varName == "0") {
			return value;
		}
		return ""
	});

	// 是否含有+-*/%运算符
	if (v.indexOf(v, "+") > -1 || v.indexOf(v, "-") > -1 || v.indexOf(v, "*") > -1 || v.indexOf(v, "/") > -1 || v.indexOf(v, "%") > -1) {
		try {
			var result = v;
			eval("result = " + v);
			return result;
		} catch (e) {
			console.log(v, e);
		}
	}
	return v;
};