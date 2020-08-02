/**
 * HTTP参数信息
 */
Vue.component("http-box", {
	props: ["prefix", "method", "url", "headers", "params", "timeout", "textBody"],
	data: function () {
		var timeout = this.timeout;
		if (timeout != null) {
			timeout = timeout.replace(/\D+/, "");
		}
		var tab = "params";
		if (this.textBody != null && this.textBody.length > 0) {
			tab = "text";
		}

		var headers = this.headers;
		if (headers == null) {
			headers = [];
		}

		var params = this.params;
		if (params == null) {
			params = [];
		}

		return {
			selectedTab: tab,
			moreOptionsVisible: false,

			vPrefix: this.prefix,
			vMethod: this.method,
			vURL: this.url,
			vHeaders: headers,
			vParams: params,
			vTimeout: timeout,
			vTextBody: this.textBody
		};
	},
	methods: {
		selectTab: function (tab) {
			this.selectedTab = tab;
			if (tab == "text") {
				var that = this;
				setTimeout(function () {
					var textInput = that.$refs.textInput;
					if (textInput != null) {
						textInput.focus();
					}
				});
			}
		},
		showMore: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible;
		}
	},
	watch: {
		vMethod: function (v) {
			if (v == "PUT") {
				this.selectTab("text");
			} else {
				this.selectTab("params");
			}
		}
	},
	template: '<tbody> \
		<tr> \
			<td class="title">URL *</td> \
			<td> \
				<input type="text" :name="prefix + \'_url\'" v-model="vURL" placeholder="http://" maxlength="500"/>\
			</td> \
		</tr> \
		<tr>\
			<td>请求方法 *</td> \
			<td> \
				<select :name="prefix + \'_method\'" v-model="vMethod" class="ui dropdown" style="width:8em">\
					<option value="GET">GET</option> \
					<option value="POST">POST</option> \
					<option value="PUT">PUT</option> \
				</select> \
			</td> \
		</tr> \
		<tr> \
			<td colspan="2"> \
				<a href="" style="font-weight: normal" @click.prevent="showMore()">更多请求选项<i class="icon angle" :class="{down:!moreOptionsVisible, up:moreOptionsVisible}"></i> </a> \
			</td> \
		</tr> \
		<tr v-show="moreOptionsVisible"> \
			<td>自定义Header</td> \
			<td> \
				<http-header-box :prefix="vPrefix" :headers="vHeaders"></http-header-box> \
			</td> \
		</tr> \
		<tr v-show="moreOptionsVisible" v-if="vMethod == \'POST\' || vMethod == \'PUT\'">\
			<td>自定义请求内容</td> \
			<td> \
				<div class="ui menu tabular attached"> \
					<a href="" class="item" :class="{active:selectedTab == \'params\'}" @click.prevent="selectTab(\'params\')" v-if="vMethod == \'POST\'">参数对</a> \
					<a href="" class="item" :class="{active:selectedTab == \'text\'}" @click.prevent="selectTab(\'text\')">文本</a> \
				</div>	\
				<div class="ui segment attached" v-if="selectedTab == \'params\'"> \
					<http-params :prefix="vPrefix" :params="vParams"></http-params> \
				</div> \
				<div class="ui segment attached" v-if="selectedTab == \'text\'">\
					<textarea rows="4" :name="prefix + \'_textBody\'" v-model="vTextBody" placeholder="要发送的内容文本" ref="textInput"></textarea> \
					<p class="comment">提醒：可能需要设置对应的<span class="ui label tiny">Content-Type</span>。</p> \
				</div> \
			</td> \
		</tr> \
		<tr v-show="moreOptionsVisible">\
			<td>超时时间</td>\
			<td> \
				<div class="ui fields inline"> \
				 	<div class="ui field">\
				 		<input type="text" :name="prefix + \'_timeout\'" style="width:4em" maxlength="6" v-model="vTimeout"/> \
				 	</div> \
				 	<div class="ui field">\
				 		s \
				 	</div> \
				</div> \
			</td> \
		</tr> \
	</tbody> \
	'
});