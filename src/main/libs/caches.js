var caches = {};

caches.set = function (key, value, lifeSeconds) {
    if (!lifeSeconds) {
        lifeSeconds = 600;
    }
    return callSetCache(key, value, lifeSeconds);
};

caches.get = function (key) {
    return callGetCache(key);
};