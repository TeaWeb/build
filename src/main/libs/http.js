var http = {};

http.Server = function (options) {
	this.isOn = true;
	this.id = "";
	this.filename = "";
	this.name = [];
	this.description = "";
	this.listen = [];
	this.backends = [];
	this.locations = [];
	this.http = true;
	this.ssl = {};

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];

			// backends
			if (key == "backends") {
				for (var i = 0; i < value.length; i++) {
					var backend = new http.Backend(value[i]);
					this.backends.push(backend);
				}
			}
			// locations
			else if (key == "locations") {
				for (var i = 0; i < value.length; i++) {
					var location = new http.Location(value[i]);
					this.locations.push(location);
				}
			}
			// 其他
			else if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

http.Backend = function (options) {
	this.isOn = true;
	this.id = "";
	this.address = "";
	this.weight = 0;
	this.isDown = false;
	this.isBackup = false;
	this.name = [];

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];
			if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

http.Location = function (options) {
	this.isOn = true;
	this.pattern = "";
	this.cachePolicy = "";
	this.fastcgi = [];
	this.id = "";
	this.index = [];
	this.root = "";
	this.rewrite = [];
	this.websocket = {};

	if (options != null && typeof (options) == "object") {
		for (var key in options) {
			var value = options[key];
			// fastcgi
			if (key == "fastcgi") {
				for (var i = 0; i < value.length; i++) {
					var fastcgi = new http.Fastcgi(value[i]);
					this.fastcgi.push(fastcgi);
				}
			}
			// rewrite
			else if (key == "rewrite") {
				for (var i = 0; i < value.length; i++) {
					var rewrite = new http.Rewrite(value[i]);
					this.rewrite.push(rewrite);
				}
			} else if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
				this[key] = value;
			}
		}
	}
};

http.Fastcgi = function (options) {
	this.id = "";
	this.isOn = true;
	this.pass = "";

	for (var key in options) {
		var value = options[key];
		if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
			this[key] = value;
		}
	}
};

http.Rewrite = function (options) {
	this.id = "";
	this.isOn = true;
	this.pattern = "";
	this.replace = "";

	for (var key in options) {
		var value = options[key];
		if (typeof (key) == "string" && typeof (this[key]) == typeof (value)) {
			this[key] = value;
		}
	}
};
