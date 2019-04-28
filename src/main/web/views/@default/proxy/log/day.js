Tea.context(function () {
	this.sourceLogs = this.logs;
	this.sourceLogs.$each(function (_, log) {
		if (typeof (log["isOpen"]) === "undefined") {
			log.isOpen = false;
		}

		// 浏览器图标
		var browserFamily = log.extend.client.browser.family.toLowerCase();
		log.browserIcon = "";
		if (["chrome", "firefox", "safari", "opera", "edge", "internet explorer"].$contains(browserFamily)) {
			log.browserIcon = browserFamily;
		} else if (browserFamily == "ie") {
			log.browserIcon = "internet explorer";
		} else if (browserFamily == "other") {
			log.extend.client.browser.family = "";
		}
	});

	this.formatCost = function (seconds) {
		var s = (seconds * 1000).toString();
		var pieces = s.split(".");
		if (pieces.length < 2) {
			return s;
		}

		return pieces[0] + "." + pieces[1].substr(0, 3);
	};

	this.showLog = function (index) {
		var log = this.logs[index];
		log.isOpen = !log.isOpen;
		if (log.isOpen) {
			this.isPlaying = false;
		}

		log.tabName = "summary";

		// 由于Vue的限制直接设置 log.isOpen 并不起作用
		this.$set(this.logs, index, log);

		// 关闭别的
		if (log.isOpen) {
			this.logs.$each(function (k, v) {
				if (v.id != log.id) {
					v.isOpen = false;
				}
			});
		}
	};

	this.showLogTab = function (log, index, tabName) {
		// 综合信息
		if (tabName == "summary") {

		}

		// 响应信息
		else if (tabName == "responseHeader") {
			if (typeof (log.responseHeaders) == "undefined") {
				this.$get(".responseHeader." + log.id + "." + log.day)
					.success(function (response) {
						log.responseHeaders = response.data.headers;
						log.responseBody = response.data.body;
						this.$set(this.logs, index, log);
					});
			}
		}

		// 请求信息
		else if (tabName == "request") {
			if (typeof (log.requestHeaders) == "undefined") {
				this.$get(".requestHeader." + log.id + "." + log.day)
					.success(function (response) {
						log.requestHeaders = response.data.headers;
						log.requestBody = response.data.body;
						log.hasRequestHeaders = false;
						for (var k in log.requestHeaders) {
							if (log.requestHeaders.hasOwnProperty(k)) {
								log.hasRequestHeaders = true;
								break;
							}
						}
						this.$set(this.logs, index, log);
					});
			}
		}

		// 预览
		else if (tabName == "preview") {
			if (typeof (log.responseHeaders) == "undefined") {
				log.previewImageLoaded = false;
				this.$get(".responseHeader." + log.id + "." + log.day)
					.success(function (response) {
						log.responseHeaders = response.data.headers;
						log.responseBody = response.data.body;

						if (typeof (log.responseHeaders["Content-Type"]) != "undefined" && log.responseHeaders["Content-Type"].length > 0 && log.responseHeaders["Content-Type"][0].match(/image\//)) {
							log.previewImageURL = log.requestScheme + "://" + log.host + log.requestURI;
						}

						this.$set(this.logs, index, log);
					})
					.done(function () {
						log.previewImageLoaded = true;
					});
			} else {
				if (typeof (log.responseHeaders["Content-Type"]) != "undefined" && log.responseHeaders["Content-Type"].length > 0 && log.responseHeaders["Content-Type"][0].match(/image\//)) {
					log.previewImageURL = log.requestScheme + "://" + log.host + log.requestURI;
				}
				log.previewImageLoaded = true;
			}
		}

		// 响应内容
		else if (tabName == "responseBody") {
			// @TODO
		}

		// Cookie
		else if (tabName == "cookie") {
			if (typeof (log.cookies) == "undefined") {
				this.$get(".cookies." + log.id + "." + log.day)
					.success(function (response) {
						log.cookies = response.data.cookies;
						log.countCookies = response.data.count;
						this.$set(this.logs, index, log);
					});
			}
		}

		// 终端信息
		else if (tabName == "client") {
			var client = log.extend.client;

			// 操作系统信息
			var osVersion = client.os.family;
			if (osVersion.length == 0 || osVersion == "Other") {
				log.osVersion = "";
			} else {
				if (client.os.major.length > 0) {
					osVersion += " " + client.os.major;
				}
				if (client.os.minor.length > 0) {
					osVersion += "." + client.os.minor;
				}
				if (client.os.patch.length > 0) {
					osVersion += "." + client.os.patch;
				}
				if (client.os.patchMinor.length > 0) {
					osVersion += "." + client.os.patchMinor;
				}
				log.osVersion = osVersion;
			}

			// 浏览器信息
			var browserVersion = client.browser.family;
			if (browserVersion.length == 0 || browserVersion == "Other") {
				log.browserVersion = "";
			} else {
				if (client.browser.major.length > 0) {
					browserVersion += " " + client.browser.major;
				}
				if (client.browser.minor.length > 0) {
					browserVersion += "." + client.browser.minor;
				}
				if (client.browser.patch.length > 0) {
					browserVersion += "." + client.browser.patch;
				}
				log.browserVersion = browserVersion;
			}

			// 地理位置信息
			var geo = log.extend.geo;
			var geoAddr = geo.region + " " + geo.state;
			if (![geo.city, geo.city + "市", geo.city + "州"].$contains(geo.state)) {
				geoAddr += " " + geo.city;
			}
			log.geoAddr = geoAddr.trim();

			if (log.geoAddr.length > 0) {
				[1].$loop(function (k, v, loop) {
					var mapBoxId = "map-box-" + log.id;
					if (document.getElementById(mapBoxId) == null) {
						setTimeout(function () {
							loop.next();
						}, 100);
						return;
					}
					var map = new BMap.Map("map-box-" + log.id);
					var decoder = new BMap.Geocoder();
					decoder.getPoint(log.geoAddr, function (point) {
						if (point == null) {
							point = new BMap.Point(geo.location.longitude, geo.location.latitude);
							var converter = new BMap.Convertor();
							converter.translate([point], 3, 5, function (data) {
								if (data.status == 0) {
									point = data.points[0];
								}

								var marker = new BMap.Marker(point, {
									icon: new BMap.Icon("/images/poi.png", new BMap.Size(20, 20), {
										anchor: new BMap.Size(10, 20),
										imageSize: new BMap.Size(20, 20)
									})
								});
								map.addOverlay(marker);
								map.centerAndZoom(point, 5);
							});
						} else {
							var marker = new BMap.Marker(point, {
								icon: new BMap.Icon("/images/poi.png", new BMap.Size(20, 20), {
									anchor: new BMap.Size(10, 20),
									imageSize: new BMap.Size(20, 20)
								})
							});
							map.addOverlay(marker);
							map.centerAndZoom(point, 5);
						}
					});
				});
			}
		}

		log.tabName = tabName;
		this.$set(this.logs, index, log);
	};
});