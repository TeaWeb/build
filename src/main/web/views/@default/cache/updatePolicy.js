Tea.context(function () {
	this.advancedVisible = false;

	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	// life
	if (this.policy.life == "0s") {
		this.life = "";
		this.lifeUnit = "s";
	} else {
		this.life = this.policy.life.substring(0, this.policy.life.length - 1);
		this.lifeUnit = this.policy.life[this.policy.life.length - 1];
	}

	// max size
	if (this.policy.maxSize.length > 0 && this.policy.maxSize != "0.00m") {
		this.maxSize = this.policy.maxSize.substring(0, this.policy.maxSize.length - 1);
		this.maxSizeUnit = this.policy.maxSize[this.policy.maxSize.length - 1];
	} else {
		this.maxSize = "";
		this.maxSizeUnit = "m";
	}

	// capacity
	if (this.policy.capacity.length > 0 && this.policy.capacity != "0.00g") {
		this.capacity = this.policy.capacity.substring(0, this.policy.capacity.length - 1);
		this.capacityUnit = this.policy.capacity[this.policy.capacity.length - 1];
	} else {
		this.capacity = "";
		this.capacityUnit = "g";
	}

	this.submitSuccess = function () {
		window.location = "/cache/policy?filename=" + this.policy.filename;
	};

	/**
	 * key
	 */
	var that = this;
	this.formatKey = function () {
		var key = that.policy.key;
		key = key.replace(/\${(.+?)}/g, "<em>${<a>$1</a>}</em>");
		return key;
	};

	/**
	 * 状态管理
	 */
	this.statusList = this.policy.status;

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
	this.cacheType = this.policy.type;

	this.changeCacheType = function () {
		var that = this;
		this.selectedType = this.types.$find(function (k, v) {
			return v.type == that.cacheType;
		});
	};
	this.changeCacheType();

	/**
	 * advanced
	 */
	this.showAdvanced = function () {
		this.advancedVisible = !this.advancedVisible;
	};

	/**
	 * Redis
	 */
	this.redisNetwork = "tcp";
	if (this.policy.options.network) {
		this.redisNetwork = this.policy.options.network;
	}

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