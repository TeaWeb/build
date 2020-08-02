Tea.context(function () {
	this.clients = [];
	this.isLoaded = false;

	this.$delay(function () {
		this.loadClients();
	});

	this.loadClients = function () {
		this.$post("$")
			.params({
				"serverId": this.server.id,
				"size": 100
			})
			.success(function (resp) {
				var that = this;
				this.clients = resp.data.clients.$map(function (k, v) {
					v.readSpeed = that.formatBytes(v.readSpeed) + "/s";
					v.writeSpeed = that.formatBytes(v.writeSpeed) + "/s";
					return v;
				});
			})
			.fail(function () {
				this.clients = [];
			})
			.done(function () {
				this.$delay(this.loadClients, 3000);
				this.isLoaded = true;
			});
	};

	this.formatBytes = function (bytes) {
		bytes = Math.ceil(bytes);
		if (bytes < 1024) {
			return bytes + " bytes";
		}
		if (bytes < 1024 * 1024) {
			return Math.ceil(bytes / 1024) + " k";
		}
		return (Math.ceil(bytes * 100 / 1024 / 1024) / 100) + " m";
	};

	this.disconnect = function (client) {
		if (!window.confirm("确定要断开客户端与服务器端之间的连接吗？")) {
			return;
		}
		this.$post(".clientDisconnect")
			.params({
				"serverId": this.server.id,
				"addr": client.clientAddr
			})
			.refresh();
	};
});