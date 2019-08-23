var times = {};

var time = {
	"YEAR": "YEAR",
	"MONTH": "MONTH",
	"DAY": "DAY",
	"HOUR": "HOUR",
	"MINUTE": "MINUTE",
	"SECOND": "SECOND"
};


times.Time = function (date) {
	if (date == null) {
		date = new Date();
	}

	this.raw = function () {
		return date;
	};

	this.year = function () {
		return date.getFullYear();
	};

	this.setYear = function (year) {
		date.setFullYear(year);
	};

	this.month = function () {
		return date.getMonth() + 1;
	};

	this.setMonth = function (month) {
		date.setMonth(month - 1);
	};

	this.day = function () {
		return date.getDate();
	};

	this.setDay = function (day) {
		date.setDate(day);
	};

	this.hour = function () {
		return date.getHours();
	};

	this.setHour = function (hour) {
		date.setHours(hour);
	};

	this.minute = function () {
		return date.getMinutes();
	};

	this.setMinute = function (minute) {
		date.setMinutes(minute);
	};

	this.second = function () {
		date.getSeconds();
	};

	this.setSecond = function (second) {
		date.setSeconds(second);
	};

	this.unix = function () {
		return parseInt(date.getTime() / 1000);
	};

	this.weekday = function () {
		return date.getDay();
	};

	this.addTime = function (years, months, days, hours, minutes, seconds) {
		var newDate = new Date(date.getTime());
		if (!isNaN(years) && years != 0) {
			newDate.setFullYear(newDate.getFullYear() + years);
		}
		if (!isNaN(months) && months != 0) {
			newDate.setMonth(newDate.getMonth() + months);
		}
		if (!isNaN(days) && days != 0) {
			newDate.setDate(newDate.getDate() + days);
		}
		if (!isNaN(hours) && hours != 0) {
			newDate.setHours(newDate.getHours() + hours);
		}
		if (!isNaN(minutes) && minutes != 0) {
			newDate.setMinutes(newDate.getMinutes() + minutes);
		}
		if (!isNaN(seconds) && seconds != 0) {
			newDate.setSeconds(newDate.getSeconds() + seconds);
		}
		return new times.Time(newDate);
	};

	this.format = function (format) {
		var result = "";
		if (format.length > 0) {
			for (var i = 0; i < format.length; i++) {
				var chr = format.charAt(i);
				result += this.formatChar(chr);
			}
		}
		return result;
	};

	//timezone
	this._parse_O = function () {
		var hours = (Math.abs(date.getTimezoneOffset() / 60)).toString();
		if (hours.length == 1) {
			hours = "0" + hours;
		}
		return "+" + hours + "00";
	};

	this._parse_r = function () {
		return this.format("D, d M Y H:i:s O");
	};

	//parse year
	this._parse_Y = function () {
		return date.getFullYear().toString();
	};

	this._parse_y = function () {
		var y = this._parse_Y();
		return y.substr(2);
	};

	this._parse_L = function () {
		var y = parseInt(this.formatChar("Y"));
		if (y % 4 == 0 && (y % 100 > 0 || y % 400 == 0)) {
			return "1";
		}
		return "0";
	};

	//month
	this._parse_m = function () {
		var n = this._parse_n();
		if (n.length < 2) {
			n = "0" + n;
		}
		return n;
	};

	this._parse_n = function () {
		return (date.getMonth() + 1).toString();
	};

	this._parse_t = function () {
		var t = 32 - new Date(this.formatChar("Y"), this.formatChar("m") - 1, 32).getDate();
		return t;
	};

	this._parse_F = function () {
		var n = parseInt(this.formatChar("n"));
		var months = ["", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
		return months[n];
	};

	this._parse_M = function () {
		var n = parseInt(this.formatChar("n"));
		var months = ["", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
		return months[n];
	};

	//week
	this._parse_w = function () {
		return date.getDay().toString();
	};

	this._parse_D = function () {
		var w = parseInt(this._parse_w());
		var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
		return days[w];
	};

	this._parse_l = function () {
		var w = parseInt(this._parse_w());
		var days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
		return days[w];
	};

	//day
	this._parse_d = function () {
		var j = this._parse_j();
		if (j.length < 2) {
			j = "0" + j;
		}
		return j;
	};

	this._parse_j = function () {
		return date.getDate().toString();
	};

	this._parse_W = function () {
		var _date = new times.Time();
		_date.setMonth(1);
		_date.setDay(1);
		var w = parseInt(_date.formatChar("w"));
		var m = parseInt(this.formatChar("m"), 10);
		var total = 0;
		for (var i = 1; i < m; i++) {
			var date2 = new times.Time();
			date2.setMonth(i);
			var t = parseInt(date2.formatChar("t"));
			total += t;
		}
		total += parseInt(this.formatChar("d"), 10);
		var w2 = parseInt(this.formatChar("w"));
		total = total - w2 + (w - 1);
		var weeks = 0;
		if (w2 != 0) {
			weeks = parseInt(total / 7 + 1);
		} else {
			weeks = parseInt(total / 7);
		}
		if (weeks.toString().length == 1) {
			weeks = "0" + weeks;
		}
		return weeks;
	};

	this._parse_z = function () {
		var m = parseInt(this.formatChar("m"), 10);
		var total = 0;
		for (var i = 1; i < m; i++) {
			var date2 = new times.Time();
			date2.set("m", i);
			var t = parseInt(date2.formatChar("t"));
			total += t;
		}
		total += parseInt(this.formatChar("d"), 10) - 1;
		return total;
	};

	//minute
	this._parse_i = function () {
		var i = date.getMinutes().toString();
		if (i.length < 2) {
			i = "0" + i;
		}
		return i;
	};

	//second
	this._parse_s = function () {
		var s = date.getSeconds().toString();
		if (s.length < 2) {
			s = "0" + s;
		}
		return s;
	};

	//hour
	this._parse_H = function () {
		var H = this._parse_G();
		if (H.length < 2) {
			H = "0" + H;
		}
		return H;
	};

	this._parse_G = function () {
		return date.getHours().toString();
	};

	this._parse_h = function () {
		var h = this._parse_g();
		if (h.length < 2) {
			h = "0" + h;
		}
		return h;
	};

	this._parse_g = function () {
		var g = parseInt(this._parse_G(), 10);
		if (g > 12) {
			g = g - 12;
		}
		return g.toString();
	};

	//time
	this._parse_U = function () {
		return this.time().toString();
	};

	//am/pm
	this._parse_a = function () {
		var hour = this.formatChar("H");
		return (hour < 12) ? "am" : "pm";
	};

	this._parse_A = function () {
		return this.formatChar("a").toUpperCase();
	};

	this.formatChar = function (chr) {
		if ((chr >= "a" && chr <= "z") || (chr >= "A" && chr <= "Z")) {
			var func = "_parse_" + chr;
			if (this[func]) {
				return this[func]();
			}
		}
		return chr;
	};

	this.toString = function () {
		return this.format("r");
	};
};

times.now = function () {
	return new times.Time();
};

times.new = function (year, month, day, hour, minute, second) {
	if (isNaN(year)) {
		year = 0;
	}
	if (isNaN(month)) {
		month = 1;
	}
	if (isNaN(day)) {
		day = 1;
	}
	if (isNaN(hour)) {
		hour = 0;
	}
	if (isNaN(minute)) {
		minute = 0;
	}
	if (isNaN(second)) {
		second = 0;
	}
	return new times.Time(new Date(year, month - 1, day, hour, minute, second));
};

times.unix = function (timestamp) {
	return new times.Time(new Date(timestamp * 1000));
};