var logs = {};

logs.Query = function () {
	var query = {
		"action": "",
		"timeFrom": null,
		"timeTo": null,
		"group": null,
		"cond": {},
		"duration": "",
		"for": null,
		"offset": -1,
		"size": -1,
		"sorts": [],
		"cache": 0,
		"result": []
	};

	this.from = function (time) {
		if (time != null) {
			query["timeFrom"] = time.unix();
		}
		return this;
	};

	this.to = function (time) {
		if (time != null) {
			query["timeTo"] = time.unix();
		}
		return this;
	};

	this.group = function (field) {
		query["group"] = field;
		return this;
	};

	this.monthly = function () {
		query.duration = "monthly";
		return;
	};

	this.daily = function () {
		query.duration = "daily";
		return this;
	};

	this.hourly = function () {
		query.duration = "hourly";
		return this;
	};

	this.minutely = function () {
		query.duration = "minutely";
	};

	this.secondly = function () {
		query.duration = "secondly";
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

	this.result = function (field1) {
		for (var i = 0; i < arguments.length; i++) {
			query.result.push(arguments[i]);
		}
		return this;
	};

	this.execute = function () {
		var cacheKey = null;
		if (query.cache > 0) {
			var cacheQuery = query;
			delete (cacheQuery["timeFrom"]);
			delete (cacheQuery["timeTo"]);
			cacheKey = JSON.stringify({
				"query": cacheQuery,
				"serverId": (context.server == null) ? "" : context.server.id
			});
			var result = caches.get(cacheKey);
			if (result != null) {
				return result;
			}
		}

		var result = callLogExecuteQuery(query);
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
		return this.action("findAll")
			.desc()
			.limit(size)
			.findAll();
	};
};