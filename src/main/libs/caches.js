var caches = {};

caches.set = function (key, value, lifeSeconds) {
    if (!lifeSeconds) {
        lifeSeconds = 600;
    }
	if (value != null && value instanceof Array) {
		// 拷贝Value，防止同一个Value相互影响
		value = JSON.parse(JSON.stringify(value))
	}
    return callSetCache(key, value, lifeSeconds);
};

caches.get = function (key) {
    var value = callGetCache(key);
    if (value != null && value instanceof Array) {
		// 拷贝Value，防止同一个Value相互影响
		value = JSON.parse(JSON.stringify(value))
	}
    return value;
};