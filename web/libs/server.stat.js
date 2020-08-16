var stat = {};

stat.Query = function () {
	var query = {
		"action": "",
		"item": "",
		"period": "",
		"cond": {},
		"offset": -1,
		"size": -1,
		"sorts": []
	};

	this.attr = function (field, value) {
		if (value != null && value instanceof Array) {
			return this.op("in", field, value);
		}
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

	this.param = function (name, value) {
		this.attr("params." + name, value);
		return this;
	};

	this.period = function (period) {
		this.attr("period", period);
		query.period = period;
		return this;
	};

	this.time = function (period, timeString) {
		this.attr("timeFormat." + period, timeString);
		return this;
	};

	this.second = function (secondString) {
		if (secondString == null || secondString.length == 0) {
			secondString = new times.Time().format("YmdHis");
		}
		return this.time("second", secondString);
	};

	this.minute = function (minuteString) {
		if (minuteString == null || minuteString.length == 0) {
			minuteString = new times.Time().format("YmdHi");
		}
		return this.time("minute", minuteString);
	};

	this.hour = function (hourString) {
		if (hourString == null || hourString.length == 0) {
			hourString = new times.Time().format("YmdH");
		}
		return this.time("hour", hourString);
	};

	this.day = function (dayString) {
		if (dayString == null || dayString.length == 0) {
			dayString = new times.Time().format("Ymd");
		}
		return this.time("day", dayString);
	};

	this.week = function (weekString) {
		if (weekString == null || weekString.length == 0) {
			weekString = new times.Time().format("YW");
		}
		return this.time("week", weekString);
	};

	this.month = function (monthString) {
		if (monthString == null || monthString.length == 0) {
			monthString = new times.Time().format("Ym");
		}
		return this.time("month", monthString);
	};

	this.year = function (yearString) {
		if (yearString == null || yearString.length == 0) {
			yearString = new times.Time().format("Y");
		}
		return this.time("year", yearString);
	};

	this.seconds = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("YmdHis"));
			date = date.addTime(0, 0, 0, 0, 0, -1);
		}
		this.attr("timeFormat.second", timeStrings);
		return this;
	};

	this.minutes = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("YmdHi"));
			date = date.addTime(0, 0, 0, 0, -1);
		}
		this.attr("timeFormat.minute", timeStrings);
		return this;
	};

	this.hours = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("YmdH"));
			date = date.addTime(0, 0, 0, -1);
		}
		this.attr("timeFormat.hour", timeStrings);
		return this;
	};

	this.days = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("Ymd"));
			date = date.addTime(0, 0, -1);
		}
		this.attr("timeFormat.day", timeStrings);
		return this;
	};

	this.weeks = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("YW"));
			date = date.addTime(0, 0, -7);
		}
		this.attr("timeFormat.week", timeStrings);
		return this;
	};

	this.months = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("Ym"));
			date = date.addTime(0, -1);
		}
		this.attr("timeFormat.month", timeStrings);
		return this;
	};

	this.years = function (count, date) {
		if (count < 1) {
			return;
		}
		if (date == null) {
			date = new times.Time();
		}
		var timeStrings = [];
		for (var i = 0; i < count; i++) {
			timeStrings.push(date.format("Y"));
			date = date.addTime(-1);
		}
		this.attr("timeFormat.year", timeStrings);
		return this;
	};

	this.item = function (item) {
		this.attr("item", item);
		query.item = item;
		return this;
	};

	this.offset = function (offset) {
		query.offset = offset;
		return this;
	};

	this.limit = function (size) {
		query.size = size;
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
		query.action = action;
		return this;
	};

	this.execute = function () {
		return callStatExecuteQuery(query);
	};

	this.find = function () {
		return this.action("find")
			.execute();
	};

	this.findAll = function () {
		return this.action("findAll")
			.execute();
	};

	this.latest = function (size, defValue, date) {
		if (size == null) {
			size = 10;
		}
		if (defValue == null) {
			defValue = {};
		}

		var period = query.period;
		if (period.length == 0 && query.item.length > 0) {
			var dotIndex = query.item.lastIndexOf(".");
			if (dotIndex > 0) {
				var last = query.item.substr(dotIndex + 1);
				if (["second", "minute", "hour", "day", "week", "month", "year"].$contains(last)) {
					period = last;
				}
			}
		}
		if (period.length == 0) {
			throw new Error("'item' or 'period' should be specified");
		}

		var timeStrings = [];
		if (size <= 0) {
			return [];
		}
		if (date == null) {
			date = new times.Time();
		}
		switch (period) {
			case "second":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("YmdHis"));
					date = date.addTime(0, 0, 0, 0, 0, -1);
				}
				break;
			case "minute":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("YmdHi"));
					date = date.addTime(0, 0, 0, 0, -1);
				}
				break;
			case "hour":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("YmdH"));
					date = date.addTime(0, 0, 0, -1);
				}
				break;
			case "day":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("Ymd"));
					date = date.addTime(0, 0, -1);
				}
				break;
			case "week":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("YW"));
					date = date.addTime(0, 0, -7);
				}
				break;
			case "month":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("Ym"));
					date = date.addTime(0, -1);
				}
				break;
			case "year":
				for (var i = 0; i < size; i++) {
					timeStrings.push(date.format("Y"));
					date = date.addTime(-1);
				}
				break
		}

		if (timeStrings.length == 0) {
			return [];
		}
		timeStrings.reverse();
		this.attr("timeFormat." + period, timeStrings);
		var ones = this.desc().findAll();

		// 填充
		var m = {};
		ones.$each(function (k, one) {
			m[one.timeFormat[period]] = one;
		});
		return timeStrings.$map(function (k, v) {
			if (typeof m[v] == "undefined") {
				return {
					"time": v,
					"value": defValue,
					"params": {}
				};
			} else {
				return {
					"time": v,
					"value": m[v].value,
					"params": m[v].params
				};
			}
		});
	};

	this.group = function (param) {
		var ones = this.findAll();
		var mapping = {};
		ones.$each(function (k, one) {
			var key = one.params[param];
			if (key == null) {
				return;
			}
			if (typeof (mapping[key]) == "undefined") {
				mapping[key] = one.value;
			} else {
				for (var k in one.value) {
					mapping[key][k] += one.value[k];
				}
			}
		});
		var result = [];
		for (var key in mapping) {
			result.push({
				"param": key,
				"value": mapping[key]
			});
		}
		return result;
	};

	this.inspect = function () {
		return query;
	};
};
