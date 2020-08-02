var widgets = {};

widgets.Widget = function (options) {
	this.name = "";
	this.code = "";
	this.author = "";
	this.version = "";
	this.description = "";
	this.params = {};
	this.options = {};
	this.requirements = [];

	for (var key in options) {
		if (!options.hasOwnProperty(key)) {
			continue;
		}
		var value = options[key];
		if (typeof (this[key]) == typeof (value)) {
			this[key] = value;
		}
	}

	this.run = function () {
		// STUB
	};

	this.callRun = function () {
		if (this.requirements.length > 0) {
			var that = this;
			this.requirements.$each(function (k, v) {
				if (!context.features.$contains(v)) {
					throw new Error("'" + that.name + "' need feature '" + v + "'");
				}
			});
		}

		this.run();
	};
};