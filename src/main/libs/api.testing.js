var apis = {};

apis.API = function (path) {
	return new apis.APIObject(path);
};

apis.APIObject = function (path) {
	var _resp;

	this.attrs = {
		path: path,
		method: "",
		params: [],
		body: null,
		repeat: 1,
		timeout: 30,
		concurrent: 100,
		headers: [],
		remoteAddr: "",
		domain: "",
		files: [],
		cookies: [],
		author: "",
		description: "",
		onSuccess: null,
		onError: null,
		onDone: null,

		assertions: [],
		failures: []
	};

	this._assert = function (rule, args, message) {
		this.attrs.assertions.push({
			"rule": rule,
			"args": args,
			"message": message
		});

		return this;
	};

	this.method = function (method) {
		if (typeof (method) != "string") {
			throw Error("method():invalid method: not a string");
		}
		method = method.toUpperCase();
		var methods = ["GET", "POST", "PUT", "HEAD", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"];
		for (var i = 0; i < methods.length; i++) {
			if (methods[i] == method) {
				this.attrs.method = method;
				return this;
			}
		}
		throw Error("method():invalid method '" + method + "'");
	};

	this.param = function (name, value) {
		this.attrs.params.push({
			"name": name,
			"value": value
		});
		return this;
	};

	this.body = function (body) {
		if (body == null) {
			this.attrs.body = null;
			return this;
		}
		if (typeof (body) == "boolean") {
			if (body) {
				this.attrs.body = "1"
			} else {
				this.attrs.body = ""
			}
			return this;
		}
		if (typeof (body) == "number") {
			this.attrs.body = body.toString();
			return this;
		}
		if (typeof (body) == "string") {
			this.attrs.body = body;
			return this;
		}
		this.attrs.body = JSON.stringify(body);
		return this;
	};

	this.repeat = function (count) {
		if (typeof (count) != "number") {
			throw Error("repeat():invalid repeat count");
		}
		this.attrs.repeat = Math.floor(count);
		return this;
	};

	this.timeout = function (seconds) {
		if (typeof (timeout) != "number") {
			throw Error("timeout():invalid timeout");
		}
		this.attrs.timeout = seconds;
		return this;
	};

	this.concurrent = function (count) {
		if (typeof (count) != "number") {
			throw Error("concurrent():invalid concurrent count");
		}
		this.attrs.concurrent = Math.floor(count);
		return this;
	};

	this.header = function (name, value) {
		this.attrs.headers.push({
			"name": name,
			"value": value
		});
		return this;
	};

	this.remoteAddr = function (remoteAddr) {
		this.attrs.remoteAddr = remoteAddr;
		return this;
	};

	this.domain = function (domain) {
		this.attrs.domain = domain;
		return this;
	};

	this.file = function (field, path) {
		this.attrs.files.push({
			"field": field,
			"path": path
		});
		return this;
	};

	this.cookie = function (name, value) {
		this.attrs.cookies.push({
			"name": name,
			"value": value
		});
		return this;
	};

	this.author = function (author) {
		this.attrs.author = author;
		return this;
	};

	this.description = function (description) {
		this.attrs.description = description;
		return this;
	};

	this.addFailure = function (failure) {
		this.attrs.failures.push(failure);
		return this;
	};

	this.assertFail = function (assert) {
		if (assert.message != null && assert.message.toString().length > 0) {
			this.addFailure(assert.message);
		} else {
			this.addFailure(JSON.stringify(assert));
		}
	};

	this.assert = function (field, f, message) {
		if (typeof (f) != "function") {
			this.addFailure("assert() arguments.2 should be a function");
			return this;
		}
		this._assert("assert", [field, f], message);
		return this;
	};

	this._runAssertAssert = function (assert) {
		var field = assert.args[0];
		var f = assert.args[1];
		var value = f(this._fieldValue(field));
		if (!this._isTrue(value)) {
			this.assertFail(assert);
		}
	};

	this.assertFormat = function (format, message) {
		this._assert("format", [format], message);
		return this;
	};

	this._runAssertFormat = function (assert) {
		var format = assert.args[0];
		if (format == "json") {
			if (_resp.bodyJSON == null) {
				this.assertFail(assert);
			}
		}
	};

	this.assertHeader = function (name, value, message) {
		this._assert("header", [name, value], message);
		return this;
	};

	this._runAssertHeader = function (assert) {
		var name = assert.args[0];
		var value = assert.args[1];
		if (typeof (_resp.headers[name]) != "string") {
			this.assertFail(assert);
			return;
		}
		if (_resp.headers[name] != value) {
			this.assertFail(assert);
			return;
		}
	};

	this.assertStatus = function (status, message) {
		this._assert("status", [status], message);
		return this;
	};

	this._runAssertStatus = function (assert) {
		var status = assert.args[0];
		if (status != _resp.status) {
			this.assertFail(assert);
		}
	};

	this.assertEqual = function (field, value, message) {
		this._assert("equal", [field, value], message);
		return this;
	};

	this._runAssertEqual = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) != value) {
			this.assertFail(assert);
		}
	};

	this.assertNotEqual = function (field, value, message) {
		this._assert("notEqual", [field, value], message);
		return this;
	};

	this._runAssertNotEqual = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) == value) {
			this.assertFail(assert);
		}
	};

	this.assertGt = function (field, value, message) {
		this._assert("gt", [field, value], message);
		return this;
	};

	this._runAssertGt = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) <= value) {
			this.assertFail(assert);
		}
	};

	this.assertGte = function (field, value, message) {
		this._assert("gte", [field, value], message);
		return this;
	};

	this._runAssertGte = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) < value) {
			this.assertFail(assert);
		}
	};

	this.assertLt = function (field, value, message) {
		this._assert("lt", [field, value], message);
		return this;
	};

	this._runAssertLt = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) >= value) {
			this.assertFail(assert);
		}
	};

	this.assertLte = function (field, value, message) {
		this._assert("lte", [field, value], message);
		return this;
	};

	this._runAssertLte = function (assert) {
		var field = assert.args[0];
		var value = assert.args[1];
		if (this._fieldValue(field) > value) {
			this.assertFail(assert);
		}
	};

	this.assertTrue = function (field, message) {
		this._assert("true", [field], message);
		return this;
	};

	this._runAssertTrue = function (assert) {
		var field = assert.args[0];
		if (!this._isTrue(this._fieldValue(field))) {
			this.assertFail(assert);
		}
	};

	this.assertFalse = function (field, message) {
		this._assert("false", [field], message);
		return this;
	};

	this._runAssertFalse = function (assert) {
		var field = assert.args[0];
		if (this._isTrue(this._fieldValue(field))) {
			this.assertFail(assert);
		}
	};

	this.assertLength = function (field, length, message) {
		this._assert("length", [field, length], message);
		return this;
	};

	this._runAssertLength = function (assert) {
		var field = assert.args[0];
		var len = assert.args[1];
		var value = this._fieldValue(field);
		if (value == null && len != 0) {
			this.assertFail(assert);
			return;
		}
		if (typeof (value) == "object" && (value instanceof Array)) {
			if (value.length != len) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "number") {
			if (value.toString().length != len) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "string") {
			if (value.length != len) {
				this.assertFail(assert);
			}
			return;
		}

		this.assertFail(assert);
	};

	this.assertNotEmpty = function (field, message) {
		this._assert("notEmpty", [field], message);
		return this;
	};

	this._runAssertNotEmpty = function (assert) {
		var field = assert.args[0];
		var value = this._fieldValue(field);
		if (value == null) {
			this.assertFail(assert);
			return;
		}

		if (typeof (value) == "object" && (value instanceof Array)) {
			if (value.length == 0) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "number") {
			if (value.toString().length == 0) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "string") {
			if (value.length == 0) {
				this.assertFail(assert);
			}
			return;
		}

		this.assertFail(assert);
	};

	this.assertEmpty = function (field, message) {
		this._assert("empty", [field], message);
		return this;
	};

	this._runAssertEmpty = function (assert) {
		var field = assert.args[0];
		var value = this._fieldValue(field);
		if (value == null) {
			return;
		}

		if (typeof (value) == "object" && (value instanceof Array)) {
			if (value.length != 0) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "number") {
			if (value.toString().length != 0) {
				this.assertFail(assert);
			}
			return;
		}

		if (typeof (value) == "string") {
			if (value.length != 0) {
				this.assertFail(assert);
			}
			return;
		}

		this.assertFail(assert);
	};

	this.assertType = function (field, type, message) {
		this._assert("type", [field, type], message);
		return this;
	};

	this._runAssertType = function (assert) {
		var field = assert.args[0];
		var type = assert.args[1];
		var value = this._fieldValue(field);
		if (type == "bool" || type == "boolean") {
			if (typeof (value) != "boolean") {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "number") {
			if (typeof (value) != "number") {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "string") {
			if (typeof (value) != "string") {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "int") {
			if (typeof (value) != "number" || parseInt(value) != value) {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "float") {
			if (typeof (value) != "number" || parseInt(value) == value) {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "array") {
			if (typeof (value) != "object" || !(value instanceof Array)) {
				this.assertFail(assert);
			}
			return;
		}

		if (type == "object") {
			if (value == null) {
				this.assertFail(assert);
				return;
			}
			if (typeof (value) != "object") {
				this.assertFail(assert);
			}

			if (value instanceof Array) {
				this.assertFail(assert);
			}

			return;
		}

		if (type == "null") {
			if (value != null) {
				this.assertFail(assert);
			}
			return;
		}

		this.assertFail(assert);
	};

	this.assertBool = function (field, message) {
		return this.assertType(field, "bool", message);
	};

	this.assertNumber = function (field, message) {
		return this.assertType(field, "number", message);
	};

	this.assertString = function (field, message) {
		return this.assertType(field, "string", message);
	};

	this.assertInt = function (field, message) {
		return this.assertType(field, "int", message);
	};

	this.assertFloat = function (field, message) {
		return this.assertType(field, "float", message);
	};

	this.assertArray = function (field, message) {
		return this.assertType(field, "array", message);
	};

	this.assertObject = function (field, message) {
		return this.assertType(field, "object", message);
	};

	this.assertNull = function (field, message) {
		return this.assertType(field, "null", message);
	};

	this.assertExist = function (field, message) {
		return this._assert("exist", [field], message);
	};

	this._runAssertExist = function (assert) {
		var field = assert.args[0];
		if (!this._existField(field)) {
			this.assertFail(assert);
		}
	};

	this.assertNotExist = function (field, message) {
		return this._assert("notExist", [field], message);
	};

	this._runAssertNotExist = function (assert) {
		var field = assert.args[0];
		if (this._existField(field)) {
			this.assertFail(assert);
		}
	};

	this.onSuccess = function (callback) {
		this.attrs.onSuccess = callback;
		return this;
	};

	this.onError = function (callback) {
		this.attrs.onError = callback;
		return this;
	};

	this.onDone = function (callback) {
		this.attrs.onDone = callback;
		return this;
	};

	this.run = function () {
		runAPI.call(this, this.attrs, this._response);
		return this;
	};

	this._response = function (resp) {
		_resp = {
			"status": resp.status,
			"body": resp.body,
			"bodyJSON": null,
			"headers": {}
		};
		for (var name in resp.headers) {
			var values = resp.headers[name];
			if (values.length > 0) {
				_resp.headers[name] = values[0];
			}
		}

		// json
		try {
			_resp.bodyJSON = JSON.parse(_resp.body);
		} catch (e) {

		}

		// assert
		for (var i = 0; i < this.attrs.assertions.length; i++) {
			var assertion = this.attrs.assertions[i];
			var rule = assertion.rule;
			var method = "_runAssert" + rule[0].toUpperCase() + rule.substring(1);
			if (typeof (this[method]) == "function") {
				this[method](assertion);
			} else {
				this.addFailure("invalid assert '" + rule + "'");
			}
		}

		// success
		if (_resp.status < 400) {
			if (this.attrs.onSuccess && typeof (this.attrs.onSuccess) == "function") {
				this.attrs.onSuccess.call(this, _resp);
			}
		} else { // error
			if (this.attrs.onError && typeof (this.attrs.onError) == "function") {
				this.attrs.onError.call(this, _resp);
			}
		}

		// done
		if (this.attrs.onDone && typeof (this.attrs.onDone) == "function") {
			this.attrs.onDone.call(this, _resp);
		}

		return this.attrs.failures;
	};


	this._fieldValue = function (field) {
		if (_resp.bodyJSON == null) {
			return null;
		}
		if (typeof (field) != "string") {
			return null;
		}
		var pieces = field.split(".");
		var last = _resp.bodyJSON;
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

	this._existField = function (field) {
		if (_resp.bodyJSON == null) {
			return false;
		}
		if (typeof (field) != "string") {
			return false;
		}
		var pieces = field.split(".");
		var last = _resp.bodyJSON;
		for (var i = 0; i < pieces.length; i++) {
			var piece = pieces[i];
			if (last === null) {
				return false;
			}
			if (typeof (last) == "object" && typeof (last[piece]) != "undefined") {
				last = last[piece];
			} else {
				return false;
			}
		}
		return true;
	};

	this._isTrue = function (value) {
		if (typeof (value) == "boolean") {
			return value;
		}
		if (typeof (value) == "number") {
			return value > 0;
		}
		if (typeof (value) == "string" && value.length > 0) {
			return true;
		}
		return false;
	};
}

