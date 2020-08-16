Vue.component("proxy-change-notify-box", {
	props: [],
	created: function () {
		var that = this;
		setTimeout(function () {
			that.refreshStatus();
		}, 100)
	},
	data: function () {
		return {
			statusChanged: false
		};
	},
	methods: {
		refreshStatus: function () {
			var that = this;
			Tea.action("/proxy/status")
				.get()
				.success(function (response) {
					that.statusChanged = response.data.changed;
				})
				.done(function () {
					this.$delay(function () {
						that.refreshStatus();
					}, 3000);
				});
		},
		restart: function () {
			var that = this;
			Tea.action("/proxy/restart")
				.get()
				.success(function () {
					that.statusChanged = false;
				});
		}
	},
	template: '<div> \
	<div class="ui icon message warning" v-if="statusChanged" style="margin-top:0.5em">\
		<i class="exclamation circle icon large"></i>\
		代理服务已被修改，<a href="" @click.prevent="restart()">点此重启后生效</a>\
	</div>\
</div>'
});