charts.Menu = function (chart) {
	this.chart = null;
	this.items = [];


	this.addItem = function (name, code, isActive) {
		var item = new charts.MenuItem();
		item.name = name;
		item.code = code;
		item.isActive = isActive;
		this.items.push(item);
	};

	this.onChange = function (callback) {
		if (chart == null) {
			return;
		}
		if (chart.options.events == null) {
			return;
		}
		var result = [];
		for (var i = 0; i < chart.options.events.length; i++) {
			var info = chart.options.events[i];
			if (info["event"] == "click" && info["target"] == "menu.item") {
				var data = info["data"];
				var found = false;
				this.items.$each(function (k, v) {
					if (v.name == data["name"] && v.code == data["code"]) {
						found = true;
						if (typeof(callback) == "function") {
							callback(data);
						}
					}
				});
				if (found) {
					continue;
				}
			}
			result.push(info);
		}
		chart.options.events = result;
	};
};

charts.MenuItem = function () {
	this.name = "";
	this.code = "";
	this.isActive = false;
};