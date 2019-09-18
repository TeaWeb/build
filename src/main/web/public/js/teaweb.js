window.teaweb = {
	set: function (key, value) {
		localStorage.setItem(key, JSON.stringify(value));
	},
	get: function (key) {
		var item = localStorage.getItem(key);
		if (item == null || item.length == 0) {
			return null;
		}

		return JSON.parse(item);
	},
	getString: function (key) {
		var value = this.get(key);
		if (typeof (value) == "string") {
			return value;
		}
		return "";
	},
	getBool: function (key) {
		return Boolean(this.get(key));
	},
	remove: function (key) {
		localStorage.removeItem(key)
	},
	match: function (source, keyword) {
		if (source == null) {
			return false;
		}
		if (keyword == null) {
			return true;
		}
		source = source.trim();
		keyword = keyword.trim();
		if (keyword.length == 0) {
			return true;
		}
		if (source.length == 0) {
			return false;
		}
		var pieces = keyword.split(/\s+/);
		for (var i = 0; i < pieces.length; i++) {
			var pattern = pieces[i];
			pattern = pattern.replace(/(\+|\*|\?|[|]|{|}|\||\\|\(|\)|\.)/g, "\\$1");
			var reg = new RegExp(pattern, "i");
			if (!reg.test(source)) {
				return false;
			}
		}
		return true;
	},

	datepicker: function (element, callback) {
		if (typeof (element) == "string") {
			element = document.getElementById(element);
		}
		var year = new Date().getFullYear();
		var picker = new Pikaday({
			field: element,
			firstDay: 1,
			minDate: new Date(year - 1, 0, 1),
			maxDate: new Date(year + 10, 11, 31),
			yearRange: [year - 1, year + 10],
			format: "YYYY-MM-DD",
			i18n: {
				previousMonth: '上月',
				nextMonth: '下月',
				months: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月'],
				weekdays: ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
				weekdaysShort: ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
			},
			theme: 'triangle-theme',
			onSelect: function () {
				if (typeof (callback) == "function") {
					callback(picker.toString());
				}
			}
		});
	},

	formatBytes: function (bytes) {
		bytes = Math.ceil(bytes);
		if (bytes < 1024) {
			return bytes + " bytes";
		}
		if (bytes < 1024 * 1024) {
			return (Math.ceil(bytes * 100 / 1024) / 100) + " k";
		}
		return (Math.ceil(bytes * 100 / 1024 / 1024) / 100) + " m";
	}
};
