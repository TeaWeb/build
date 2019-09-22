var values = {};

function NewQuery() {
	return new values.Query();
}

values.Query = function () {
	var query = {
		"action": "",
		"cond": {
			"error": {
				"eq": ""
			}
		},
		"duration": "",
		"offset": -1,
		"size": -1,
		"sorts": [],
		"cache": 0,
		"timePast": null,
		"aggregationFields": null
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

	this.past = function (number, unit) {
		query.timePast = {
			"number": number,
			"unit": unit
		};
		return this;
	};

	this.timeLabel = function (one) {
		if (context.timeUnit != null && context.timeUnit.length > 0) {
			query.timePast = {
				"unit": context.timeUnit
			};
		}
		if (query.timePast == null) {
			var minute = one.timeFormat.minute.substring(8);
			return minute.substr(0, 2) + ":" + minute.substr(2, 2);
		}
		switch (query.timePast.unit) {
			case time.SECOND:
				var second = one.timeFormat.second.substring(8);
				return second.substr(0, 2) + ":" + second.substr(2, 2) + ":" + second.substr(4, 2);
			case time.MINUTE:
				var minute = one.timeFormat.minute.substring(8);
				return minute.substr(0, 2) + ":" + minute.substr(2, 2);
			case time.HOUR:
				return one.timeFormat.hour.substr(8);
			case time.DAY:
				return one.timeFormat.day.substr(4, 2) + "-" + one.timeFormat.day.substr(6, 2);
			case time.MONTH:
				return one.timeFormat.month.substr(0, 4) + "-" + one.timeFormat.month.substr(4, 2);
			case time.YEAR:
				return one.timeFormat.year;
		}
		var minute = one.timeFormat.minute.substring(8);
		return minute.substr(0, 2) + ":" + minute.substr(2, 2);
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

	this.action = function (action) {
		query["action"] = action;
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

	this.avg = function () {
		var fields = [];
		for (var i = 0; i < arguments.length; i++) {
			fields.push(arguments[i]);
		}
		query.aggregationFields = fields;
		var ones = this.action("avgValues")
			.execute();
		for (var i = 0; i < ones.length; i++) {
			ones[i].label = this.timeLabel(ones[i]);
		}
		return ones;
	};

	this.findAll = function () {
		return this.action("findAll")
			.execute();
	};

	this.find = function () {
		return this.action("find")
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
	if (value == null) {
		return "";
	}
	var v = param.replace(/(\${[\w\\.\s]+})/, function (match) {
		var varName = match.substring(2, match.length - 1)
			.replace(/\s+/g, "");
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
						if (lastObject != null && typeof (lastObject) == "object" && typeof (lastObject[piece]) != "undefined") {
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