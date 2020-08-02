Tea.context(function () {
	/**
	 * 缓存策略
	 */
	this.cacheEditing = false;
	this.selectedCachePolicy = this.cachePolicyFile;

	this.editCache = function () {
		this.cacheEditing = !this.cacheEditing;
		if (this.cacheEditing) {
			this.$delay(function () {
				window.scroll(0, 10000);
			});
		}
	};

	this.cancelCacheEditing = function () {
		this.cacheEditing = false;
	};

	this.saveCacheEditing = function () {
		this.$post(".updateCache")
			.params({
				"serverId": this.server.id,
				"policy": this.selectedCachePolicy
			})
			.success(function (resp) {
				this.cacheEditing = false;

				if (this.selectedCachePolicy.length == 0) {
					this.cachePolicy = null;
				} else {
					this.cachePolicy = resp.data.policy;
				}
			});
	};

	this.cacheKey = "";
	this.cleanToolVisible = false;

	this.showCleanTool = function () {
		this.cleanToolVisible = !this.cleanToolVisible;
		if (this.cleanToolVisible) {
			this.$delay(function () {
				this.$find("form input[name='cacheKey']").focus();
			});
		}
	};

	this.cleanCache = function () {
		if (this.cacheKey.length == 0) {
			alert("请输入要清除的Key");
			this.$find("form input[name='cacheKey']").focus();
			return;
		}
		this.$post("/proxy/cache/clean")
			.params({
				"filename": this.cachePolicy.filename,
				"key": this.cacheKey
			})
			.success(function () {
				alert("清除成功");
			});
	};
});