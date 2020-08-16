var apis = {};

apis.APIResponseObject = function (resp) {
	resp.bodyJSON = JSON.parse(resp.body);

	this.header = function (name) {
		if (typeof (resp.headers[name]) == "undefined") {
			return "";
		}
		if (resp.headers[name].length == 0) {
			return "";
		}
		return resp.headers[name][0];
	};

	this.value = function (field) {
		if (resp.bodyJSON == null) {
			return null;
		}
		if (typeof (field) != "string") {
			return null;
		}
		var pieces = field.split(".");
		var last = resp.bodyJSON;
		for (var i = 0; i < pieces.length; i++) {
			var piece = pieces[i];
			if (last === null) {
				return null;
			}
			if (typeof (last) == "object" && typeof (last[piece]) != "undefined") {
				last = last[piece];
			} else {
				return null;
			}
		}
		return last;
	};
}

apis.NewAPIResponseObject = function (resp) {
	return new apis.APIResponseObject(resp);
};