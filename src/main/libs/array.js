Array.$nil = {};

Array.prototype.$each = function (fn) {
	var that = this;
	if (that == null) {
		return true;
	}
	if (typeof (fn) != "function") {
		return true;
	}

	var length = that.length;
	for (var i = 0; i < length; i++) {
		fn.call(that, i, that[i]);
	}

	return true;
};

Array.prototype.$fill = function (value, length) {
	var that = this;
	if (that == null) {
		return false;
	}
	if (typeof (length) == "undefined") {
		length = that.length;
	}
	if (length < that.length) {
		return false;
	}
	if (length == that.length) {
		return true;
	}
	var size = length - that.length;
	for (var i = 0; i < size; i++) {
		that.push(value);
	}
	return true;
};

Array.prototype.$sum = function (fn) {
	var that = this;
	if (that == null) {
		return 0;
	}
	var sum = 0;
	that.$each(function (k, v) {
		if (typeof (fn) == "function") {
			v = fn.call(that, k, v);
		}
		if (typeof (v) == "number") {
			sum += v;
		} else if (typeof (v) == "string") {
			var n = parseFloat(v);
			if (!isNaN(n)) {
				sum += n;
			}
		}
	});
	return sum;
};

Array.prototype.$max = function (sortFunction) {
	var that = this;
	if (that == null) {
		return null;
	}
	if (that.length > 0) {
		var _array = that.$copy();
		_array.$rsort(sortFunction);
		return _array.$get(0);
	}
	return null;
};

Array.prototype.$copy = function () {
	var that = this;
	if (that == null) {
		return that;
	}

	var newArray = [];
	for (var i = 0; i < that.length; i++) {
		newArray.push(that[i]);
	}
	return newArray;
};

Array.prototype.$get = function (index) {
	var that = this;
	if (that == null) {
		return null;
	}
	if (index > that.length - 1) {
		return null;
	}
	return that[index];
};

Array.prototype.$rsort = function (sortFunction) {
	var that = this;
	if (that == null) {
		return false
	}

	this.$sort(sortFunction);
	that.reverse();
	return true;
};

Array.prototype.$sort = function (sortFunction) {
	var that = this;
	if (that == null) {
		return false
	}
	if (typeof (sortFunction) == "undefined") {
		sortFunction = function (v1, v2) {
			if (v1 > v2) {
				return 1;
			} else if (v1 == v2) {
				return 0;
			} else {
				return -1;
			}
		};
	}
	that.sort(sortFunction);
	return true;
};

Array.prototype.$map = function (fn) {
	var that = this;
	if (that == null) {
		return [];
	}
	var arr = [];
	for (var i = 0; i < that.length; i++) {
		var result = fn.call(that, i, that[i]);
		if (result === Array.$nil) {
			continue;
		}
		arr.push(result);
	}
	return arr;
};

Array.prototype.$contains = function (v) {
	var that = this;
	if (that == null) {
		return false;
	}
	for (var i = 0; i < that.length; i++) {
		if (that[i] == v) {
			return true;
		}
	}
	return false;
};

Array.prototype.$include = function (v) {
	var that = this;
	if (that == null) {
		return false;
	}
	return that.$contains(v);
};


/**
 * 取得某一个值在数组中第一次出现的位置
 */
Array.prototype.$indexOf = function (value, strict) {
	var that = this;
	if (that == null) {
		return [];
	}
	if (arguments.length == 0) {
		return -1;
	}
	if (typeof (strict) == "undefined") {
		strict = false;
	}
	for (var i = 0; i < that.length; i++) {
		if ((strict && value === that[i]) || (!strict && value == that[i])) {
			return i;
		}
	}
	return -1;
};

/**
 * 取得某一个值在数组中最后一次出现的位置
 */
Array.prototype.$lastIndexOf = function (value, strict) {
	var that = this;
	if (that == null) {
		return [];
	}
	if (arguments.length == 0) {
		return -1;
	}
	if (typeof (strict) == "undefined") {
		strict = false;
	}
	var index = -1;
	for (var i = 0; i < that.length; i++) {
		if ((strict && value === that[i]) || (!strict && value == that[i])) {
			index = i;
		}
	}
	return index;
};

/**
 * 一次性加入多个元素
 */
Array.prototype.$pushAll = function (array2) {
	var that = this;
	if (that == null) {
		return 0;
	}
	return Array.prototype.push.apply(that, array2);
};

