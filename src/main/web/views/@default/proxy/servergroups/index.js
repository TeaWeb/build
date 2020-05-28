Tea.context(function () {
	this.updateGroup = function (groupId) {
		teaweb.popup("/proxy/servergroups/update?groupId=" + groupId, {
			height: "10em",
			callback: function () {
				teaweb.successReload("保存成功");
			}
		});
	};

	this.deleteGroup = function (groupId) {
		teaweb.confirm("确定要删除此分组吗？", function () {
			this.$post(".delete")
				.params({
					"groupId": groupId
				})
				.refresh();
		});
	};

	this.addServer = function (groupId) {
		teaweb.popup("/proxy/servergroups/addServer?groupId=" + groupId, {
			callback: function () {
				teaweb.successReload("保存成功");
			}
		});
	};

	this.deleteServer = function (groupId, serverId) {
		teaweb.confirm("确定要移除此服务吗？", function () {
			this.$post(".deleteServer")
				.params({
					"groupId": groupId,
					"serverId": serverId,
				})
				.refresh();
		});
	};
});