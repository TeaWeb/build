Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/proxy/detail?serverId=" + this.server.id;
	};

	/**
	 * 域名
	 */
	this.nameAdding = false;
	this.addingNameName = "";
	this.nameEditingIndex = -1;

	this.addName = function () {
		this.nameAdding = true;
		this.nameEditingIndex = -1;
		this.$delay(function () {
			this.$find("form input[name='addingNameName']").focus();
		});
	};

	this.editName = function (index) {
		this.nameEditingIndex = index;
		this.addingNameName = this.server.name[index];
		this.nameAdding = true;
		this.$delay(function () {
			this.$find("form input[name='addingNameName']").focus();
		});
	};

	this.confirmAddName = function () {
		this.addingNameName = this.addingNameName.trim();
		if (this.addingNameName.length == 0) {
			alert("域名不能为空");
			this.$find("form input[name='addingNameName']").focus();
			return;
		}
		if (this.nameEditingIndex > -1) {
			this.server.name[this.nameEditingIndex] = this.addingNameName;
		} else {
			this.server.name.push(this.addingNameName);
		}
		this.cancelNameAdding();
	};

	this.cancelNameAdding = function () {
		this.nameAdding = false;
		this.addingNameName = "";
		this.nameEditingIndex = -1;
	};

	this.removeName = function (index) {
		this.cancelNameAdding();
		this.server.name.$remove(index);
	};

	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * 监听地址
	 */
	this.listenAdding = false;
	this.addingListenName = "";
	this.editingListenIndex = -1;

	this.addListen = function () {
		this.listenAdding = true;
		this.editingListenIndex = -1;
		this.$delay(function () {
			this.$find("form input[name='addingListenName']").focus();
		});
	};

	this.editListen = function (index) {
		this.listenAdding = true;
		this.editingListenIndex = index;
		this.$delay(function () {
			this.$find("form input[name='addingListenName']").focus();
		});
		this.addingListenName = this.server.listen[index];
	};

	this.confirmAddListen = function () {
		this.addingListenName = this.addingListenName.trim();
		if (this.addingListenName.length == 0) {
			alert("绑定地址不能为空");
			this.$find("form input[name='addingListenName']").focus();
			return;
		}
		if (this.editingListenIndex > -1) {
			this.server.listen[this.editingListenIndex] = this.addingListenName;
		} else {
			this.server.listen.push(this.addingListenName);
		}
		this.cancelListenAdding();
	};

	this.cancelListenAdding = function () {
		this.listenAdding = false;
		this.addingListenName = "";
		this.editingListenIndex = -1;
	};

	this.removeListen = function (index) {
		this.server.listen.$remove(index);
		this.cancelListenAdding();
	};

	/**
	 * index
	 */
	this.indexAdding = false;
	this.addingIndexName = "";

	this.addIndex = function () {
		this.indexAdding = true;
		this.$delay(function () {
			this.$find("form input[name='addingIndexName']").focus();
		});
	};

	this.confirmAddIndex = function () {
		this.addingIndexName = this.addingIndexName.trim();
		if (this.addingIndexName.length == 0) {
			alert("首页文件名不能为空");
			this.$find("form input[name='addingIndexName']").focus();
			return;
		}
		this.server.index.push(this.addingIndexName);
		this.cancelIndexAdding();
	};

	this.cancelIndexAdding = function () {
		this.indexAdding = !this.indexAdding;
		this.addingIndexName = "";
	};

	this.removeIndex = function (index) {
		this.server.index.$remove(index);
	};

	/**
	 * 单位
	 */
	this.maxBodyUnits = [{
		"code": "k",
		"name": "K"
	}, {
		"code": "m",
		"name": "M"
	}, {
		"code": "g",
		"name": "G"
	}];
	this.maxBodyUnit = "m";
	if (this.server.maxBodySize.length > 0) {
		this.maxBodyUnit = this.server.maxBodySize[this.server.maxBodySize.length - 1];
		this.server.maxBodySize = this.server.maxBodySize.substring(0, this.server.maxBodySize.length - 1);
	}

	/**
	 * 访问日志
	 */
	this.enableAccessLog = !this.server.disableAccessLog;
	this.enableStat = !this.server.disableStat;

	/**
	 * 小静态文件加速
	 */
	this.cacheStatic = this.server.cacheStatic;

	/**
	 * 压缩级别
	 */
	this.gzipLevels = Array.$range(1, 9);
	this.gzipMinUnits = [
		{
			"code": "b",
			"name": "B"
		},
		{
			"code": "k",
			"name": "K"
		}, {
			"code": "m",
			"name": "M"
		}];
	this.gzipMinUnit = "k";
	if (this.server.gzip.minLength.length > 0) {
		this.gzipMinUnit = this.server.gzip.minLength[this.server.gzip.minLength.length - 1];
		this.server.gzip.minLength = this.server.gzip.minLength.substring(0, this.server.gzip.minLength.length - 1);
	}

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		var that = this;
		[".names-box", ".indexes-box", ".listens-box"].$each(function (k, box) {
			var box = that.$find(box)[0];
			if (!box) {
				return;
			}
			Sortable.create(box, {
				draggable: ".label",
				handle: ".handle",
				onStart: function () {

				},
				onUpdate: function (event) {
				}
			});
		});
	};

	/**
	 * 高级选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * 通知设置
	 */
	this.noticeVisible = false;

	this.showNotice = function () {
		this.noticeVisible = !this.noticeVisible;
	};
});