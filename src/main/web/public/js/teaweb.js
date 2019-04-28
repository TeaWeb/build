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
        if (typeof(value) == "string") {
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
        for (var i = 0; i < pieces.length; i ++) {
            var pattern = pieces[i];
            pattern = pattern.replace(/(\+|\*|\?|[|]|{|}|\||\\|\(|\)|\.)/g, "\\$1");
            var reg = new RegExp(pattern, "i");
            if (!reg.test(source)) {
                return false;
            }
        }
        return true;
    }
};
