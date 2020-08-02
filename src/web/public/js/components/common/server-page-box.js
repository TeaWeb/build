/**
 * 代理服务特殊页面
 */
Vue.component("server-page-box", {
	props: ["comment", "pages", "prefix"],
	data: function () {
		var vOnes = [];
		if (this.pages != null) {
			vOnes = this.pages.$map(function (k, v) {
				v.status = (v.status.length > 0) ? v.status[0] : 0;
				v.newStatus = (v.newStatus > 0) ? v.newStatus.toString() : "";
				return v;
			});
		}
		return {
			editingIndex: -1,
			status: "",
			url: "",
			newStatus: "",
			vOnes: vOnes,
			vPrefix: (this.prefix == null) ? "" : this.prefix,
			isAdding: false,
			typicalPages: [
				{"name": "403页面", "url": "web/pages/403.html"},
				{"name": "404页面", "url": "web/pages/404.html"},
				{"name": "50x页面", "url": "web/pages/50x.html"},
				{"name": "暂时关闭英文页面", "url": "web/pages/shutdown_en.html"},
				{"name": "暂时关闭中文页面", "url": "web/pages/shutdown_zh.html"},
				{"name": "升级中中文页面", "url": "web/pages/shutdown_upgrade_zh.html"}
			],
			typicalPageVisible: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				if (that.$refs.statusInput != null) {
					that.$refs.statusInput.focus()
				}
			}, 50);
		},
		confirmAdd: function () {
			if (this.status == null || this.status.length == 0) {
				alert("请输入要匹配的响应状态码");
				this.focusStatus();
				return;
			}

			if (this.status.length != 3) {
				alert("请输入3位的状态码");
				this.focusStatus();
				return;
			}
			if (!this.status.match(/^[\dx]+$/)) {
				alert("状态码只能是数字或者字母x");
				this.focusStatus();
				return;
			}

			if (this.status[0] == "0") {
				alert("状态码第一位不能为0");
				this.focusStatus();
				return;
			}

			if (this.url.length == 0) {
				alert("请输入页面文件地址或者URL");
				this.focusURL();
				return;
			}

			if (this.newStatus != null && this.newStatus.length > 0) {
				if (this.newStatus.length != 3) {
					alert("请输入3位的新状态码");
					this.focusNewStatus();
					return;
				}
				if (!this.newStatus.match(/^\d+$/)) {
					alert("新状态码只能是数字");
					this.focusNewStatus();
					return;
				}

				if (this.newStatus[0] == "0") {
					alert("新状态码第一位不能为0");
					this.focusNewStatus();
					return;
				}
			}

			if (this.editingIndex > -1) {
				this.vOnes[this.editingIndex].status = this.status;
				this.vOnes[this.editingIndex].url = this.url;
				this.vOnes[this.editingIndex].newStatus = this.newStatus;
			} else {
				this.vOnes.push({
					status: this.status,
					url: this.url,
					newStatus: this.newStatus
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.status = "";
			this.url = "";
			this.newStatus = "";
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此特殊页面吗？")) {
				return;
			}
			this.cancel();
			this.vOnes.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.status = this.vOnes[index].status;
			this.url = this.vOnes[index].url;
			this.newStatus = this.vOnes[index].newStatus;
		},
		focusStatus: function () {
			var that = this;
			setTimeout(function () {
				if (that.$refs.statusInput != null) {
					that.$refs.statusInput.focus()
				}
			}, 50);
		},
		focusURL: function () {
			var that = this;
			setTimeout(function () {
				if (that.$refs.urlInput != null) {
					that.$refs.urlInput.focus()
				}
			}, 50);
		},
		focusNewStatus: function () {
			var that = this;
			setTimeout(function () {
				if (that.$refs.newStatusInput != null) {
					that.$refs.newStatusInput.focus()
				}
			}, 50);
		},
		showTypicalPages: function () {
			this.typicalPageVisible = !this.typicalPageVisible;
		},
		selectTypicalPage: function (page) {
			this.url = page.url;
			this.showTypicalPages();
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" style="padding:4px" :class="{blue:editingIndex == index}" v-for="(one,index) in vOnes"> \
				[{{one.status}}] -&gt; <span v-if="one.newStatus.length > 0">[{{one.newStatus}}]</span>{{one.url}}\
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" :name="vPrefix + \'StatusList\'" :value="one.status"/>\
				<input type="hidden" :name="vPrefix + \'URLList\'" :value="one.url"/> \
				<input type="hidden" :name="vPrefix + \'NewStatusList\'" :value="one.newStatus"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<table class="ui table definition"> \
				<tr> \
					<td class="title">响应状态码 *</td> \
					<td> \
						<input type="text" v-model="status" ref="statusInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="状态码" maxlength="3" style="width:5.2em"/> \
						<p class="comment">比如404，或者50x</p> \
					</td> \
				</tr> \
				<tr> \
					<td>URL *</td> \
					<td> \
						<input type="text" v-model="url"  ref="urlInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="页面文件路径或者完整的URL"/> \
						<p class="comment">页面文件是相对于TeaWeb目录的页面文件比如web/pages/40x.html，或者一个完整的URL。<a href="" @click.prevent="showTypicalPages()">推荐页面<i class="icon angle" :class="{down:!typicalPageVisible, up:typicalPageVisible}"></i></a> </p> \
						<div v-show="typicalPageVisible"> \
							<a class="ui label tiny" style="margin-bottom:2px;" @click.prevent="selectTypicalPage(page)" v-for="page in typicalPages" title="点击选中">{{page.name}}: <var>{{page.url}}</var></a> \
						</div> \
					</td> \
				</tr> \
				<tr> \
					<td class="title">新状态码</td> \
					<td> \
						<input type="text" v-model="newStatus" ref="newStatusInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="状态码" maxlength="3" style="width:5.2em"/> \
						<p class="comment">可以用来修改响应的状态码，不填表示不改变原有状态码。</p> \
					</td> \
				</tr> \
			</table> \
			<div style="margin-bottom:1em"> \
				<button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex == -1">确认添加</button><button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex > -1">确认保存</button>  &nbsp;<a href="" @click.prevent="cancel()">取消</a> \
			</div> \
		</div> \
		<button class="ui button tiny" type="button" @click.prevent="add()" v-if="!isAdding">+</button> \
		<p class="comment">{{comment}}</p> \
	</div>'
});