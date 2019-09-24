Tea.context(function () {
	this.advancedVisible = false;

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	/**
	 * key
	 */
	var that = this;
	this.key = "${host}${requestURI}";
	this.formatKey = function () {
		var key = that.key;
		key = key.replace(/\${(.+?)}/g, "<em>${<a>$1</a>}</em>");
		return key;
	};

	/**
	 * 状态管理
	 */
	this.statusList = ["200"];

	this.statusAdding = false;
	this.addingStatus = "";

	this.addStatus = function () {
		this.statusAdding = true;
		this.$delay(function () {
			this.$find("form input[name='addingStatus']").focus();
		});
	};

	this.cancelAdding = function () {
		this.statusAdding = false;
	};

	this.addStatusConfirm = function (e) {
		if (this.addingStatus.length != 3) {
			alert("状态码必须是3位数字");
			this.$find("form input[name='addingStatus']").focus();
			return;
		}
		if (this.statusList.$contains(this.addingStatus)) {
			alert("状态码已存在");
			this.$find("form input[name='addingStatus']").focus();
			return;
		}
		this.statusList.push(this.addingStatus);
		this.statusAdding = false;
		this.addingStatus = "";

		return false;
	};

	this.deleteStatus = function (index) {
		this.statusList.$remove(index);
	};

	/**
	 * 类型
	 */
	this.cacheType = this.types[0].type;
	this.selectedType = this.types[0];

	this.changeCacheType = function () {
		var that = this;
		this.selectedType = this.types.$find(function (k, v) {
			return v.type == that.cacheType;
		});
	};

	/**
	 * Advanced
	 */
	this.showAdvanced = function () {
		this.advancedVisible = !this.advancedVisible;
	};

	/**
	 * Redis
	 */
	this.redisNetwork = "tcp";

	/**
	 * cache control
	 */
	this.addingCacheControl = "";
	this.cacheControlIsAdding = false;

	this.removeCacheControl = function (cacheControl) {
		this.skippedCacheControlValues.$removeValue(cacheControl);
	};

	this.addCacheControl = function () {
		this.cacheControlIsAdding = true;
		this.$delay(function () {
			this.$find("input[name='addingCacheControl']").focus();
		});
	};

	this.addCacheControlConfirm = function () {
		if (this.addingCacheControl.length == 0) {
			alert("请输入一个非空值");
			this.$find("input[name='addingCacheControl']").focus();
			return;
		}
		this.skippedCacheControlValues.push(this.addingCacheControl);
		this.addingCacheControl = "";
		this.cacheControlIsAdding = false;
	};

	this.cancelCacheControlAdding = function () {
		this.cacheControlIsAdding = false;
	};
});