var agents = {};

agents.Agent = function (options) {
	this.isOn = false;
	this.id = "";
	this.isLocal = false;
	this.name = "";
	this.host = "";
	this.apps = [];

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];

			// apps
			if (key == "apps") {
				for (var i = 0; i < value.length; i++) {
					var app = new agents.App(value[i]);
					this.apps.push(app);
				}
			}
			// 其他
			else if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

agents.App = function (options) {
	this.id = "";
	this.isOn = false;
	this.name = "";
	this.tasks = [];

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];

			// apps
			if (key == "tasks") {
				for (var i = 0; i < value.length; i++) {
					var task = new agents.Task(value[i]);
					this.tasks.push(task);
				}
			}
			// 其他
			else if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

agents.Task = function (options) {
	this.id = "";
	this.isOn = false;
	this.name = "";
	this.isBooting = false;
	this.isManual = false;
	this.isScheduling = false;

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];

			if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

agents.Item = function (options) {
	this.id = "";
	this.isOn = false;
	this.name = "";
	this.interval = 0;

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];

			if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};
