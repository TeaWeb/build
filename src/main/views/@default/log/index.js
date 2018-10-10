Tea.context(function () {
    var that = this;

    this.logs = [];
    this.sourceLogs = [];
    this.fromId = "";
    this.total = 0;
    this.countSuccess = 0;
    this.countFail = 0;
    this.qps = 0;

    this.started = false;

    // 搜索相关
    this.searchBoxVisible = teaweb.getBool("searchBoxVisible");
    this.searchIp = teaweb.getString("searchIp");
    this.searchDomain = teaweb.getString("searchDomain");
    this.searchOs = teaweb.getString("searchOs");
    this.searchBrowser = teaweb.getString("searchBrowser");
    this.searchMinCost = teaweb.getString("searchMinCost");
    this.searchKeyword = teaweb.getString("searchKeyword");

    this.$delay(function () {
        this.loadLogs();

        this.$watch("searchIp", function (value) {
            that.filterLogs()
        });
        this.$watch("searchDomain", function (value) {
            that.filterLogs()
        });
        this.$watch("searchOs", function (value) {
            that.filterLogs()
        });
        this.$watch("searchBrowser", function (value) {
            that.filterLogs()
        });
        this.$watch("searchMinCost", function (value) {
            that.filterLogs()
        });
        this.$watch("searchKeyword", function (value) {
            that.filterLogs()
        });
    });

    window.addEventListener("unload", function () {
        teaweb.set("searchIp", that.searchIp);
        teaweb.set("searchDomain", that.searchDomain);
        teaweb.set("searchOs", that.searchOs);
        teaweb.set("searchBrowser", that.searchBrowser);
        teaweb.set("searchMinCost", that.searchMinCost);
        teaweb.set("searchKeyword", that.searchKeyword);
    });

    var loadSize = 100;
    this.loadLogs = function () {
        var lastSize = 0;
        this.$get(".get")
            .params({
                "fromId": this.fromId,
                "size": loadSize
            })
            .success(function (response) {
                lastSize = response.data.logs.length;
                if (lastSize == loadSize) {
                    loadSize = 1000;
                } else {
                    loadSize = 100;
                }

                this.total = Math.ceil(response.data.total * 100 / 10000) / 100;
                this.countSuccess = Math.ceil(response.data.countSuccess * 100 / 10000) / 100;
                this.countFail = Math.ceil(response.data.countFail * 100 / 10000) / 100;
                this.qps = response.data.qps;

                this.sourceLogs = response.data.logs.concat(this.sourceLogs);
                this.sourceLogs.$each(function (_, log) {
                    if (typeof(log["isOpen"]) === "undefined") {
                        log.isOpen = false;
                    }
                });

                if (this.sourceLogs.length > 0) {
                    this.fromId = this.sourceLogs.$first().id;

                    if (this.sourceLogs.length > 100) {
                        this.sourceLogs = this.sourceLogs.slice(0, 100);
                    }
                }

                this.filterLogs();
            })
            .done(function () {
                this.started = true;

                // 每1秒刷新一次
                Tea.delay(function () {
                    this.loadLogs();
                }, (lastSize < loadSize) ? 1000 : 100)
            });
    };

    this.showLog = function (index) {
        var log = this.logs[index];
        log.isOpen = !log.isOpen;

        // 由于Vue的限制直接设置 log.isOpen 并不起作用
        this.$set(this.logs, index, log);
    };

    this.formatCost = function (seconds) {
        var s = (seconds * 1000).toString();
        var pieces = s.split(".");
        if (pieces.length < 2) {
            return s;
        }

        return pieces[0] + "." + pieces[1].substr(0, 3);
    };

    this.showSearchBox = function () {
        this.searchBoxVisible = true;
        teaweb.set("searchBoxVisible", true);
    };

    this.hideSearchBox = function () {
        this.searchBoxVisible = false;
        teaweb.set("searchBoxVisible", false);
    };

    this.hasSearchConditions = function () {
        var has = false;
        this.$find(".search-box form input").each(function (k, v) {
           if (typeof(v.value) == "string" && v.value.trim().length > 0) {
               has = true;
           }
        });
        return has;
    };

    this.resetSearchBox = function () {
        that.searchIp = "";
        that.searchDomain = "";
        that.searchOs = "";
        that.searchBrowser = "";
        that.searchMinCost = "";
        that.searchKeyword = "";
    };

    this.filterLogs = function () {
        this.logs = this.sourceLogs.$filter(function (_, log) {
            if (!teaweb.match(log.remoteAddr, that.searchIp)) {
                return false;
            }

            if (!teaweb.match(log.host, that.searchDomain)) {
                return false;
            }

            if (typeof(log.extend.client.os.family) != "undefined" && !teaweb.match(log.extend.client.os.family, that.searchOs)) {
                return false;
            }

            if (typeof(log.extend.client.browser.family) != "undefined" && !teaweb.match(log.extend.client.browser.family, that.searchBrowser)) {
                return false;
            }

            if (that.searchMinCost.length > 0) {
                var cost = parseFloat(that.searchMinCost);
                if (isNaN(cost) || log.requestTime < cost * 0.001) {
                    return false;
                }
            }

            if (that.searchKeyword != null && that.searchKeyword.length > 0) {
                var values = [
                    log.requestPath,
                    log.requestURI,
                    log.userAgent,
                    log.remoteAddr,
                    log.requestMethod,
                    log.statusMessage,
                    log.timeLocal,
                    log.timeISO8601,
                    log.host,
                    log.request,
                    log.contentType,
                    JSON.stringify(log.extend)
                ];

                var found = false;
                for (var i = 0; i < values.length; i ++) {
                    if (teaweb.match(values[i], that.searchKeyword)) {
                        found = true;
                        break;
                    }
                }

                if (!found) {
                    return false;
                }
            }

            return true;
        });
    };

});