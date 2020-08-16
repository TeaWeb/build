/**
 * 临时关闭页面
 */
Vue.component("server-shutdown-box", {
	props: ["shutdown", "prefix"],
	data: function () {
		return {
			vShutdown: (this.shutdown == null) ? {
				on: false,
				url: "",
				status: 200
			} : this.shutdown,
			typicalPages: [
				{"name": "404页面", "url": "web/pages/404.html"},
				{"name": "50x页面", "url": "web/pages/50x.html"},
				{"name": "暂时关闭英文页面", "url": "web/pages/shutdown_en.html"},
				{"name": "暂时关闭中文页面", "url": "web/pages/shutdown_zh.html"},
				{"name": "升级中中文页面", "url": "web/pages/shutdown_upgrade_zh.html"}
			],
			typicalPageVisible: false
		}
	},
	watch: {
		"vShutdown.on": function (v) {
			if (v) {
				this.focusURL();
			}
		}
	},
	methods: {
		focusURL: function () {
			var that = this;
			setTimeout(function () {
				if (that.$refs.urlInput != null) {
					that.$refs.urlInput.focus()
				}
			}, 50);
		},
		showTypicalPages: function () {
			this.typicalPageVisible = !this.typicalPageVisible;
		},
		selectTypicalPage: function (page) {
			this.vShutdown.url = page.url;
			this.showTypicalPages();
		}
	},
	template: '<div> \
		<div class="ui checkbox"> \
			<input type="checkbox" :name="prefix + \'On\'" :id="prefix + \'shutdownPageOn\'" v-model="vShutdown.on"/> \
			<label :for="prefix + \'shutdownPageOn\'">是否开启</label> \
		</div>\
		<div style="margin-top:0.4em" v-show="vShutdown.on"> \
			<table class="ui table"> \
				<tr> \
					<td class="title">页面URL</td> \
					<td> \
						<input type="text" :name="prefix + \'URL\'" ref="urlInput" v-model="vShutdown.url" placeholder="页面文件路径或一个完整URL" maxlength="100" style="width:30em"/> \
						<p class="comment">页面文件是相对于TeaWeb目录的页面文件比如web/pages/40x.html，或者一个完整的URL。<a href="" @click.prevent="showTypicalPages()">推荐页面<i class="icon angle" :class="{down:!typicalPageVisible, up:typicalPageVisible}"></i></a> </p> \
						<div v-show="typicalPageVisible"> \
							<a class="ui label tiny" style="margin-bottom:2px;" @click.prevent="selectTypicalPage(page)" v-for="page in typicalPages" title="点击选中">{{page.name}}: <var>{{page.url}}</var></a> \
						</div> \
					</td> \
				</tr> \
				<tr> \
					<td>状态码</td> \
					<td><input type="text" :name="prefix + \'Status\'" v-model="vShutdown.status" style="width:5.2em" placeholder="状态码" maxlength="3"/></td> \
				</tr> \
			</table> \
		</div> \
		<p class="comment">开启临时关闭页面时，所有请求的响应都会显示此页面。可用于临时升级网站使用。</p>\
	</div>'
});