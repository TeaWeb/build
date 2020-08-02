Tea.context(function () {
	var that = this;

	this.logs = [];
	this.fromId = "";

	this.isPlaying = true;
	this.isLoaded = false;

	// 搜索相关
	this.searchBoxVisible = teaweb.getBool("searchBoxVisible");
	this.searchIp = teaweb.getString("searchIp");
	this.searchDomain = teaweb.getString("searchDomain");
	this.searchOs = teaweb.getString("searchOs");
	this.searchBrowser = teaweb.getString("searchBrowser");
	this.searchMinCost = teaweb.getString("searchMinCost");
	this.searchKeyword = teaweb.getString("searchKeyword");
	this.searchBackendId = teaweb.getString("searchBackendId");
	this.searchLocationId = teaweb.getString("searchLocationId");
	if (this.searchLocationId.length > 0) {
		this.$delay(function () {
			this.location = this.server.locations.$find(function (k, v) {
				return v.id == that.searchLocationId;
			});
			this.changeLocation();
		});
	}
	this.searchRewriteId = teaweb.getString("searchRewriteId");
	this.searchFastcgiId = teaweb.getString("searchFastcgiId");

	this.$delay(function () {
		this.loadLogs();

		this.$find(".menu").bind("click", function () {
			window.scrollTo(0, 0);
		});
	});

	window.addEventListener("unload", function () {
		teaweb.set("searchIp", that.searchIp);
		teaweb.set("searchDomain", that.searchDomain);
		teaweb.set("searchOs", that.searchOs);
		teaweb.set("searchBrowser", that.searchBrowser);
		teaweb.set("searchMinCost", that.searchMinCost);
		teaweb.set("searchKeyword", that.searchKeyword);
		teaweb.set("searchBackendId", that.searchBackendId);
		teaweb.set("searchLocationId", that.searchLocationId);
		teaweb.set("searchRewriteId", that.searchRewriteId);
		teaweb.set("searchFastcgiId", that.searchFastcgiId);
	});

	/**
	 * 路径规则相关筛选
	 */
	this.location = null;

	this.changeLocation = function (locationId) {
		if (locationId != null) {
			this.searchRewriteId = "";
			this.searchFastcgiId = "";
		}

		if (this.searchLocationId == null || this.searchLocationId.length == 0) {
			this.location = null;
			return;
		}
		var that = this;
		this.location = this.server.locations.$find(function (k, v) {
			return v.id == that.searchLocationId;
		});

		this.changeFilter();
	};

	var loadSize = 10;
	this.loadLogs = function () {
		// 是否正在播放日志
		if (!this.isPlaying) {
			return;
		}

		var lastSize = 0;
		this.$get(".list")
			.params({
				"serverId": this.server.id,
				"fromId": this.fromId,
				"size": loadSize,
				"bodyFetching": this.bodyFetching ? 1 : 0,
				"logType": this.logType,
				"remoteAddr": this.searchIp,
				"domain": this.searchDomain,
				"osName": this.searchOs,
				"browser": this.searchBrowser,
				"cost": this.searchMinCost,
				"keyword": this.searchKeyword,
				"backendId": this.searchBackendId,
				"locationId": this.searchLocationId,
				"rewriteId": this.searchRewriteId,
				"fastcgiId": this.searchFastcgiId
			})
			.success(function (response) {
				// 日志
				lastSize = response.data.logs.length;
				if (lastSize == loadSize) {
					loadSize = 1000;
				} else {
					loadSize = 100;
				}

				if (lastSize == 0) {
					return;
				}

				if (response.data.lastId.length > 0) {
					this.fromId = response.data.lastId;
				}

				this.logs = response.data.logs.concat(this.logs);

				var max = 100;
				if (this.logs.length > max) {
					this.logs = this.logs.slice(0, max);
				}

				this.logs.$each(function (_, log) {
					if (typeof (log["isOpen"]) === "undefined") {
						log.isOpen = false;
					}

					// 浏览器图标
					log.browserIcon = "";
					if (log.extend != null) {
						var browserFamily = log.extend.client.browser.family.toLowerCase();
						if (["chrome", "firefox", "safari", "opera", "edge", "internet explorer"].$contains(browserFamily)) {
							log.browserIcon = browserFamily;
						} else if (browserFamily == "ie") {
							log.browserIcon = "internet explorer";
						} else if (browserFamily == "other") {
							log.extend.client.browser.family = "";
						}
					}
				});
			})
			.done(function () {
				this.$delay(function () {
					this.isLoaded = true;
				}, 100);

				// 每1秒刷新一次
				Tea.delay(function () {
					this.loadLogs();
				}, (lastSize < loadSize) ? 1000 : 100)
			});
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
			if (typeof (v.value) == "string" && v.value.trim().length > 0) {
				has = true;
			}
		});
		this.$find(".search-box form select").each(function (k, v) {
			if (typeof (v.value) == "string" && v.value.trim().length > 0) {
				has = true;
			}
		});
		return has;
	};

	this.resetSearchBox = function () {
		this.searchIp = "";
		this.searchDomain = "";
		this.searchOs = "";
		this.searchBrowser = "";
		this.searchMinCost = "";
		this.searchKeyword = "";
		this.searchBackendId = "";
		this.searchLocationId = "";
		this.searchRewriteId = "";
		this.searchFastcgiId = "";
		this.changeFilter();
	};

	this.showLogTab = function (log, index, tabName) {
		// 综合信息
		if (tabName == "summary") {

		}

		// 响应信息
		else if (tabName == "responseHeader") {
			this.$get(".responseHeader." + log.id + "." + log.day)
				.success(function (response) {
					log.responseHeaders = response.data.headers;
					log.responseHasBody = response.data.hasBody;

					if (log.responseHeaders == null) {
						log.responseHeaders = {};
					}

					this.$set(this.logs, index, log);
				});
		}

		// 请求信息
		else if (tabName == "request") {
			this.$get(".requestHeader." + log.id + "." + log.day)
				.success(function (response) {
					log.requestHeaders = response.data.headers;
					if (log.requestHeaders == null) {
						log.requestHeaders = {};
					}

					log.requestBody = response.data.body;
					log.hasRequestHeaders = false;
					for (var k in log.requestHeaders) {
						if (log.requestHeaders.hasOwnProperty(k)) {
							log.hasRequestHeaders = true;
							break;
						}
					}

					log.shouldHighlightRequest = false;
					var contentType = "";
					if (log.requestHeaders != null && log.requestBody != null) {
						contentType = log.requestHeaders["Content-Type"];
						if (contentType != null) {
							log.shouldHighlightRequest = [
								"application/json",
								"text/json"
							].$exist(function (k, v) {
								return contentType.toString().indexOf(v) > -1;
							});
						}
					}

					this.$set(this.logs, index, log);

					if (log.shouldHighlightRequest) {
						this.$delay(function () {
							var box = this.$find(".request-body-box")[0];
							if (box != null) {
								box.innerHTML = "";
							}

							var codeEditor = CodeMirror(box, {
								theme: "idea",
								lineNumbers: true,
								styleActiveLine: true,
								matchBrackets: true,
								value: "",
								readOnly: true,
								height: "auto",
								viewportMargin: Infinity,
								lineWrapping: true,
								highlightFormatting: true,
								indentUnit: 4,
								mode: "",
								indentWithTabs: true
							});

							var mimeType = "application/json";
							var info = CodeMirror.findModeByMIME(mimeType);
							if (info != null) {
								codeEditor.setOption("mode", info.mode);
								CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
								CodeMirror.autoLoadMode(codeEditor, info.mode);
							}

							codeEditor.setValue(log.requestBody);
						});
					}
				});

		}

		// 预览
		else if (tabName == "preview") {
			this.previewTab = "preview";

			log.responseBodyLoaded = false;
			log.responseRawBody = "";
			this.$set(this.logs, index, log);

			this.$get(".responseBody." + log.id + "." + log.day)
				.success(function (response) {
					log.responseRawBody = response.data.rawBody;
					log.responseIsImage = response.data.isImage;
					log.responseIsText = response.data.isText;
					log.responseBody = response.data.body;
					log.responseBodyContentType = response.data.contentType;
					log.responseBodyEncoding = response.data.encoding;

					var contentType = response.data.contentType;
					this.$set(this.logs, index, log);

					if (contentType == "text/plain") {
						log.responseIsText = false;
					}

					if (log.responseIsImage) {
						log.responseImageNatureSize = "";
						this.$delay(function () {
							var image = document.getElementById("log-response-image");
							if (image != null && image.naturalWidth != null && image.naturalWidth > 0 && image.naturalHeight != null && image.naturalHeight > 0) {
								log.responseImageNatureSize = image.naturalWidth + "x" + image.naturalHeight;
								this.$set(this.logs, index, log);
							}
						});
					}

					if (log.responseIsText && !log.responseBodyLoaded) {
						this.$delay(function () {
							var box = document.getElementById("log-response-text-editor");
							if (box != null) {
								box.innerHTML = "";
								var codeEditor = CodeMirror(box, {
									theme: "idea",
									lineNumbers: true,
									styleActiveLine: true,
									matchBrackets: true,
									value: "",
									readOnly: true,
									//height: "auto",
									viewportMargin: Infinity,
									lineWrapping: true,
									highlightFormatting: true,
									indentUnit: 4,
									mode: "",
									indentWithTabs: true
								});

								var info = CodeMirror.findModeByMIME(contentType);
								if (info != null) {
									codeEditor.setOption("mode", info.mode);
									CodeMirror.modeURL = "/codemirror/mode/%N/%N.js";
									CodeMirror.autoLoadMode(codeEditor, info.mode);
								}

								codeEditor.setValue(log.responseBody);
							}
						}, 100);
					}

					log.responseBodyLoaded = true;
				});
		}

		// Cookie
		else if (tabName == "cookie") {
			if (typeof (log.cookies) == "undefined") {
				this.$get(".cookies." + log.id + "." + log.day)
					.success(function (response) {
						log.cookies = response.data.cookies;
						log.countCookies = response.data.count;
						if (log.cookies == null) {
							log.cookies = {};
						}
						this.$set(this.logs, index, log);
					});
			}
		}

		// 终端信息
		else if (tabName == "client") {
			if (log.extend != null) {
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
		}

		log.tabName = tabName;
		this.$set(this.logs, index, log);
	};

	this.pause = function () {
		this.isPlaying = !this.isPlaying;

		if (this.isPlaying) {
			this.loadLogs();
		}
	};

	this.bodyFetching = false;
	this.startBodyFetching = function () {
		this.bodyFetching = !this.bodyFetching;
	};

	this.changeFilter = function () {
		this.fromId = "";
		this.logs = [];
		this.isLoaded = false;
	};

	/**
	 * 预览相关
	 */
	this.previewTab = "preview";

	this.selectPreviewTab = function (tab) {
		this.previewTab = tab;
	};
});