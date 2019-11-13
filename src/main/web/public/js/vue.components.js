/**
 * 自动补全
 */
Vue.component("auto-complete-box", {
	props: ["name", "placeholder", "options", "value", "maxlength", "autocomplete"],
	data: function () {
		return {
			newValue: this.value,
			visible: false,
			id: "auto-complete-box-" + Math.random().toString().replace(".", "-"),
			index: -1
		}
	},
	watch: {
		options: function (v) {
			var box = document.getElementById(this.id);
			var items = box.querySelectorAll(".item");
			this.index = -1;
			for (var i = 0; i < items.length; i++) {
				items[i].className = "item";
			}
		}
	},
	methods: {
		search: function () {
			if (this.autocomplete === false) {
				return;
			}
			this.index = -1;
			this.visible = (this.options.length > 0 && this.newValue.length > 0);
			this.$emit("change", this.newValue);
		},
		select: function (option) {
			this.newValue = option.value;
			this.visible = false;
			var box = document.getElementById(this.id);
			box.querySelector(".search").focus();
		},
		enter: function () {
			if (this.index > -1) {
				this.select(this.options[this.index]);
			}
			this.visible = false;
		},
		blur: function () {
			var that = this;
			setTimeout(function () {
				that.visible = false;
			}, 500);
		},
		down: function () {
			this.move(true);
		},
		up: function () {
			this.move(false);
		},
		move: function (isDown) {
			var box = document.getElementById(this.id);
			var items = box.querySelectorAll(".item");
			if (items.length > 0) {
				for (var i = 0; i < items.length; i++) {
					items[i].className = "item";
				}
				if (isDown) {
					this.index++;
					if (this.index >= items.length) {
						this.index = 0;
					} else if (this.index < 0) {
						items.index = 0;
					}
				} else {
					this.index--;
					if (this.index < 0) {
						this.index = items.length - 1;
					}
				}
				items[this.index].className = "item active";
				var offset = items[this.index].offsetTop;
				var box = document.getElementById(this.id);
				var menu = box.querySelector(".menu");
				if (offset + 20 >= menu.offsetHeight) {
					menu.scrollTop = 1000;
				} else {
					menu.scrollTop = 0;
				}
			} else {
				this.index = -1;
			}
		}
	},
	template: '<div class="autocomplete-box"> \
    <div class="ui fluid search selection dropdown" :id="id" style="padding-top:0;padding-bottom:0;height:2.7em;line-height:2.7em;z-index:1"> \
          <!--<span class="default text" v-if="newValue.length == 0">{{placeholder}}</span>--> \
          <input class="search fluid" :placeholder="placeholder" style="line-height:2.65em;padding-top:0;padding-bottom:0;z-index:10" :name="name" v-model="newValue" @input="search(newValue)" autocomplete="off" @keyup.down="down()" @keyup.up="up()" @keyup.enter="enter" @keypress.enter.prevent="1" @blur="blur()" :maxlength="maxlength"/> \
          <div class="ui menu blue" :style="{display: (!visible || this.options.length == 0) ? \'none\' : \'block\'}"> \
            <a class="item" v-for="option in options" @click.prevent="select(option)"> \
              {{option.name}} \
            </a> \
          </div> \
        </div> \
    </div>'
});

/**
 * 路径自动补全
 */
Vue.component("auto-complete-path-box", {
	props: ["name", "placeholder", "value", "maxlength", "autocomplete"],
	data: function () {
		return {
			"options": []
		};
	},
	methods: {
		change: function (v) {
			this.options = [];
			var that = this;
			Tea.action("/proxy/localPath")
				.params({
					"prefix": v
				})
				.success(function (resp) {
					that.options = resp.data.paths.$map(function (k, path) {
						return {
							"name": path,
							"value": path
						};
					});
				})
				.get();
		}
	},
	template: '<auto-complete-box \
        :name="this.name" \
        :value="this.value" \
        :placeholder="this.placeholder" \
        :options="options" \
        :maxlength="maxlength" \
        :autocomplete="autocomplete" \
        @change="change($event)"> \
            </auto-complete-box>'
});

/**
 * 通知设置
 */
Vue.component("notice-item", {
	props: ["name", "item"],
	data: function () {
		return {
			"optionsVisible": false,
			"noticeChecked": this.item.on,
			"noticeLevel": this.item.level,
			"noticeSubject": this.item.subject,
			"noticeBody": this.item.body
		}
	},
	methods: {
		showOptions: function () {
			this.optionsVisible = !this.optionsVisible;
		}
	},
	template: '<div> \
	<div class="ui checkbox"> \
		<input type="checkbox" :name="name + \'NoticeOn\'" value="1" v-model="noticeChecked"/>	\
		<label> \
		 	<a href="" v-if="noticeChecked" @click.prevent="showOptions()">自定义<i class="icon angle" :class="{up:optionsVisible, down:!optionsVisible}"></i></a> \
		</label> \
	</div> \
	<table class="ui table definition" v-show="noticeChecked && optionsVisible">\
		<tr> \
			<td class="title">级别</td> \
			<td> \
				<select class="ui dropdown" :name="name + \'NoticeLevel\'" v-model="noticeLevel" style="width:5em"> \
				 	<option value="1">信息</option> \
				 	<option value="2">警告</option> \
				 	<option value="3">错误</option> \
				 	<option value="4">成功</option> \
				</select> \
			</td> \
		</tr>	\
		<tr> \
			<td>标题</td> \
			<td> \
				<input type="text" :name="name + \'NoticeSubject\'" v-model="noticeSubject"/> \
			</td> \
		</tr> \
		<tr> \
			<td>内容</td> \
			<td> \
				<textarea :name="name + \'NoticeBody\'" rows="2" v-model="noticeBody"></textarea> \
			</td> \
		</tr> \
	</table> \
		</div>'
});

/**
 * 通知设置显示
 */
Vue.component("notice-item-view", {
	props: ["item"],
	data: function () {
		var levelName = "";
		var levelColor = "olive";
		if (this.item != null) {
			switch (this.item.level) {
				case 0:
					levelName = "信息";
					break;
				case 1:
					levelName = "信息";
					break;
				case 2:
					levelName = "警告";
					levelColor = "yellow";
					break;
				case 3:
					levelName = "错误";
					levelColor = "red";
					break;
				case 4:
					levelName = "成功";
					levelColor = "green";
					break;
			}
		}
		return {
			"levelName": levelName,
			"levelColor": levelColor
		}
	},
	template: '<div> \
	  <span v-if="item == null">还没有设置</span> \
	  <div v-if="item != null"> \
	   	  <span class="ui label" v-if="!item.on">未开启</span>	\
	   	  <span v-if="item.on" :class="\'ui label tiny \' + levelColor">{{levelName}}级别：{{item.subject}}</span> \
	  </div> \
	</div>'
});

/**
 * Agent Group密钥管理
 */
Vue.component("agent-group-keys", {
	props: ["keys"],
	data: function () {
		if (this.keys == null) {
			this.keys = [];
		}

		return {
			editingIndex: -1,
			name: "",
			dayFrom: "",
			dayTo: "",
			maxAgents: 0,
			isAdding: false,
			on: true
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			Tea.delay(function () {
				teaweb.datepicker("day-input-1", function (day) {
					that.dayFrom = day;
				});
				teaweb.datepicker("day-input-2", function (day) {
					that.dayTo = day;
				});
			});
		},
		confirmAdd: function () {
			if (this.name == null || this.name.length == 0) {
				alert("请输入密钥说明文字");
				return;
			}

			if (this.editingIndex > -1) {
				this.keys[this.editingIndex].name = this.name;
				this.keys[this.editingIndex].dayFrom = this.dayFrom;
				this.keys[this.editingIndex].dayTo = this.dayTo;
				this.keys[this.editingIndex].maxAgents = this.maxAgents;
				this.keys[this.editingIndex].on = this.on;
			} else {
				this.keys.push({
					key: "",
					name: this.name,
					dayFrom: this.dayFrom,
					dayTo: this.dayTo,
					maxAgents: this.maxAgents,
					on: this.on
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.name = "";
			this.dayFrom = "";
			this.dayTo = "";
			this.maxAgents = 0;
			this.on = true;
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此密钥吗？")) {
				return;
			}
			this.cancel();
			this.keys.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.key = this.keys[index].key;
			this.name = this.keys[index].name;
			this.dayFrom = this.keys[index].dayFrom;
			this.dayTo = this.keys[index].dayTo;
			this.maxAgents = this.keys[index].maxAgents;
			this.on = this.keys[index].on;
			var that = this;
			Tea.delay(function () {
				teaweb.datepicker("day-input-1", function (day) {
					that.dayFrom = day;
				});
				teaweb.datepicker("day-input-2", function (day) {
					that.dayTo = day;
				});
			});
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" :class="{blue:editingIndex == index}" v-for="(key,index) in keys"> \
				{{key.name}}：\
				<span v-if="key.key.length > 0">[{{key.key}}]</span> \
				<span v-if="key.key.length == 0">[保存后自动生成Key]</span> \
				<span v-if="key.dayFrom.length > 0">{{key.dayFrom}}</span> <span v-if="key.dayFrom.length > 0 || key.dayTo.length > 0">-</span> <span v-if="key.dayTo.length > 0">{{key.dayTo}}</span> \
				<span v-if="key.maxAgents >  0">/ {{key.maxAgents}}</span> \
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" name="keysName" :value="key.name"/>\
				<input type="hidden" name="keysKey" :value="key.key"/> \
				<input type="hidden" name="keysDayFrom" :value="key.dayFrom"/> \
				<input type="hidden" name="keysDayTo" :value="key.dayTo"/> \
				<input type="hidden" name="keysMaxAgents" :value="key.maxAgents"/> \
				<input type="hidden" name="keysOn" :value="key.on ? 1 : 0"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<table class="ui table definition"> \
				<tr> \
					<td class="title">密钥</td> \
					<td>自动生成</td> \
				</tr> \
				<tr>\
					<td>说明 *</td>\
					<td>\
						<input type="text" name="name" v-model="name"  @keyup.enter="confirmAdd" @keypress.enter.prevent="1"/>\
					</td>\
				</tr>\
				<tr> \
					<td>密钥开始生效日期</td> \
					<td> \
						<input type="text" name="dayFrom" v-model="dayFrom" id="day-input-1" style="width:8em" autocomplete="off" @keyup.enter="confirmAdd" @keypress.enter.prevent="1"/> \
						<p class="comment">非必填信息。</p> \
					</td> \
				</tr> \
				<tr> \
					<td>密钥失效日期</td> \
					<td> \
						<input type="text" name="dayTo" v-model="dayTo" id="day-input-2" style="width:8em" autocomplete="off" @keyup.enter="confirmAdd" @keypress.enter.prevent="1"/> \
						<p class="comment">该日期结束后，不能再注册新的Agent。非必填信息。</p> \
					</td> \
				</tr> \
				<tr> \
					<td>能注册的Agent数量限制</td> \
					<td> \
						<input type="text" size="8" value="0" style="width:10em" v-model="maxAgents" @keyup.enter="confirmAdd" @keypress.enter.prevent="1"/> \
						<p class="comment">超出此限制则不能再添加新的Agent。0表示不限制。</p> \
					</td> \
				</tr> \
				<tr> \
					<td>是否启用</td> \
					<td> \
						<div class="ui checkbox"> \
							<input type="checkbox" v-model="on"/> \
							<label></label>\
						</div> \
					</td> \
				</tr> \
			</table> \
			<div style="margin-bottom:1em"> \
				<button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex == -1">确认添加</button><button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex > -1">确认保存</button>  &nbsp;<a href="" @click.prevent="cancel()">取消</a> \
			</div> \
		</div> \
		<button class="ui button tiny" type="button" @click.prevent="add()" v-if="!isAdding">+</button> \
		<p class="comment">用于授权客户端快速注册Agent。</p> \
	</div>'
});

/**
 * Header参数
 */
Vue.component("http-header-box", {
	props: ["headers", "comment", "prefix"],
	data: function () {
		return {
			editingIndex: -1,
			name: "",
			value: "",
			isAdding: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				if (that.$refs.nameInput != null) {
					that.$refs.nameInput.focus()
				}
			}, 50);
		},
		confirmAdd: function () {
			if (this.name == null || this.name.length == 0) {
				alert("请输入Header名称");
				return;
			}

			if (this.editingIndex > -1) {
				this.headers[this.editingIndex].name = this.name;
				this.headers[this.editingIndex].value = this.value;
			} else {
				this.headers.push({
					name: this.name,
					value: this.value
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.name = "";
			this.value = "";
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此Header吗？")) {
				return;
			}
			this.cancel();
			this.headers.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.name = this.headers[index].name;
			this.value = this.headers[index].value;
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" style="padding:4px" :class="{blue:editingIndex == index}" v-for="(header,index) in headers"> \
				{{header.name}}: {{header.value}}\
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" :name="prefix + \'_headerNames\'" :value="header.name"/>\
				<input type="hidden" :name="prefix + \'_headerValues\'" :value="header.value"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<table class="ui table definition"> \
				<tr> \
					<td class="title">名称</td> \
					<td> \
						<input type="text" name="name" v-model="name" ref="nameInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="Header名"/> \
					</td> \
				</tr> \
				<tr> \
					<td>值</td> \
					<td> \
						<input type="text" name="value" v-model="value"  @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="Header值"/> \
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

/**
 * HTTP参数
 */
Vue.component("http-params", {
	props: ["params", "comment", "prefix"],
	data: function () {
		var paramList = this.params;
		if (paramList == null) {
			paramList = [];
		}
		return {
			paramList: paramList,
			editingIndex: -1,
			name: "",
			value: "",
			isAdding: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				if (that.$refs.nameInput != null) {
					that.$refs.nameInput.focus()
				}
			}, 50);
		},
		confirmAdd: function () {
			if (this.name == null || this.name.length == 0) {
				alert("请输入参数名称");
				return;
			}

			if (this.editingIndex > -1) {
				this.paramList[this.editingIndex].name = this.name;
				this.paramList[this.editingIndex].value = this.value;
			} else {
				this.paramList.push({
					name: this.name,
					value: this.value
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.name = "";
			this.value = "";
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此参数吗？")) {
				return;
			}
			this.cancel();
			this.params.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.name = this.paramList[index].name;
			this.value = this.paramList[index].value;
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" style="padding:4px" :class="{blue:editingIndex == index}" v-for="(param,index) in paramList"> \
				{{param.name}}: {{param.value}}\
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" :name="prefix + \'_paramNames\'" :value="param.name"/>\
				<input type="hidden" :name="prefix + \'_paramValues\'" :value="param.value"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<table class="ui table definition"> \
				<tr> \
					<td class="title">名称</td> \
					<td> \
						<input type="text" name="name" v-model="name" ref="nameInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="参数名"/> \
					</td> \
				</tr> \
				<tr> \
					<td>值</td> \
					<td> \
						<input type="text" name="value" v-model="value"  @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="参数值"/> \
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

/**
 * 单个值列表组件
 */
Vue.component("single-value-list", {
	props: ["values", "comment", "prefix", "valueName"],
	data: function () {
		var values = this.values;
		if (values == null) {
			values = [];
		} else {
			values = this.values.$map(function (k, v) {
				return {
					value: v
				};
			});
		}
		return {
			vValues: values,
			editingIndex: -1,
			value: "",
			isAdding: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				if (that.$refs.valueInput != null) {
					that.$refs.valueInput.focus()
				}
			}, 50);
		},
		confirmAdd: function () {
			if (this.value == null || this.value.length == 0) {
				alert("请输入" + this.valueName + "信息");
				if (this.$refs.valueInput != null) {
					this.$refs.valueInput.focus()
				}
				return;
			}

			if (this.editingIndex > -1) {
				this.vValues[this.editingIndex].value = this.value;
			} else {
				this.vValues.push({
					value: this.value
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.value = "";
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此" + this.valueName + "吗？")) {
				return;
			}
			this.cancel();
			this.vValues.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.value = this.vValues[index].value;
			var that = this;
			setTimeout(function () {
				if (that.$refs.valueInput != null) {
					that.$refs.valueInput.focus()
				}
			});
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" style="padding:4px" :class="{blue:editingIndex == index}" v-for="(param,index) in vValues"> \
				{{param.value}}\
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" :name="prefix + \'Values\'" :value="param.value"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<div style="margin-bottom:1em"> \
				<input type="text" name="value" v-model="value" ref="valueInput"  @keyup.enter="confirmAdd" @keypress.enter.prevent="1" :placeholder="valueName" style="width:10em"/> &nbsp; \
				<button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex == -1">确认添加</button><button class="ui button tiny" type="button" @click.prevent="confirmAdd()" v-if="editingIndex > -1">确认保存</button>  &nbsp;<a href="" @click.prevent="cancel()">取消</a> \
			</div> \
		</div> \
		<button class="ui button tiny" type="button" @click.prevent="add()" v-if="!isAdding">+</button> \
		<p class="comment">{{comment}}</p> \
	</div>'
});

/**
 * 请求匹配条件
 */
Vue.component("request-cond-box", {
	props: ["conds", "prefix", "operators", "comment", "variables"],
	data: function () {
		var conds = this.conds;
		if (conds == null) {
			conds = [];
		}

		var variables = this.variables;
		if (variables == null) {
			variables = [];
		}

		var that = this;
		setTimeout(function () {
			that.changeOp("eq")
		}, 100);
		return {
			vConds: conds,
			vOperators: this.operators,
			vVariables: variables,
			editingIndex: -1,
			vParam: "",
			vValue: "",
			vValues: [],
			vOperator: "eq",
			vOperatorDescription: "",
			vVariable: "",
			vVariableDescription: "",
			isAdding: false,
			variablesVisible: false
		}
	},
	watch: {
		vOperator: function (op) {
			this.changeOp(op);
		},
		vVariable: function (variable) {
			this.vOperatorDescription = "";
			this.vParam += variable;
			if (variable.length > 0) {
				var v = this.vVariables.$find(function (k, v1) {
					return v1.code == variable;
				});
				if (v) {
					this.vVariableDescription = v.description;
				}
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				if (that.$refs.paramInput != null) {
					that.$refs.paramInput.focus()
				}
			});
		},
		confirmAdd: function () {
			if (this.vParam == null || this.vParam.length == 0) {
				alert("请输入参数");
				if (this.$refs.paramInput != null) {
					this.$refs.paramInput.focus()
				}
				return;
			}

			if (this.isArrayOperator(this.vOperator)) {
				this.vValue = JSON.stringify(this.vValues);
			}

			if (this.editingIndex > -1) {
				this.vConds[this.editingIndex].param = this.vParam;
				this.vConds[this.editingIndex].operator = this.vOperator;
				this.vConds[this.editingIndex].value = this.vValue;
			} else {
				this.vConds.push({
					param: this.vParam,
					operator: this.vOperator,
					value: this.vValue
				});
			}
			this.cancel();
		},
		cancel: function () {
			this.editingIndex = -1;
			this.isAdding = false;
			this.vParam = "";
			this.vOperator = "eq";
			this.vValue = "";
			this.vValues = [];
			this.vVariable = "";
		},
		remove: function (index) {
			if (!window.confirm("确定要删除此匹配条件吗？")) {
				return;
			}
			this.cancel();
			this.vConds.$remove(index);
		},
		edit: function (index) {
			this.editingIndex = index;
			this.isAdding = true;
			this.vParam = this.vConds[index].param;
			this.vOperator = this.vConds[index].operator;
			this.vValue = this.vConds[index].value;

			if (this.isArrayOperator(this.vOperator)) {
				this.vValues = JSON.parse(this.vValue);
			}

			var that = this;
			setTimeout(function () {
				if (that.$refs.paramInput != null) {
					that.$refs.paramInput.focus()
				}
			})
		},
		changeOp: function (op) {
			var operator = this.vOperators.$find(function (k, v) {
				return v.op == op;
			});
			if (operator != null) {
				this.vOperatorDescription = operator.description;
			}
		},
		addValue: function () {
			this.vValues.push("");

			var that = this;
			setTimeout(function () {
				var inputs = that.$refs.valuesInput;
				if (inputs != null && (inputs instanceof Array)) {
					inputs[inputs.length - 1].focus();
				}
			});
		},
		removeValue: function (index) {
			this.vValues.$remove(index);
		},
		showVariables: function () {
			this.variablesVisible = !this.variablesVisible;
		},
		isArrayOperator: function (operator) {
			return ["in", "not in", "file ext", "mime type"].$contains(operator);
		},
		hasValue: function (operator) {
			return !["file exist", "file not exist"].$contains(operator);
		}
	},
	template: '<div> \
		<div style="margin-bottom: 1em">\
			<div class="ui label tiny" style="padding:4px;margin-top:3px;margin-bottom:3px" :class="{blue:editingIndex == index}" v-for="(cond,index) in vConds"> \
				{{cond.param}} <var>{{cond.operator}}</var> {{cond.value}}\
				&nbsp; <a href="" title="修改" @click.prevent="edit(index)"><i class="icon pencil small"></i></a>&nbsp; \
				<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i> </a> \
				<input type="hidden" :name="prefix + \'_condParams\'" :value="cond.param"/> \
				<input type="hidden" :name="prefix + \'_condOperators\'" :value="cond.operator"/> \
				<input type="hidden" :name="prefix + \'_condValues\'" :value="cond.value"/> \
			</div>\
		</div> \
		<div v-if="isAdding"> \
			<table class="ui table definition"> \
				<tr> \
					<td class="title">参数</td> \
					<td> \
						<input type="text" v-model="vParam" ref="paramInput" @keyup.enter="confirmAdd" @keypress.enter.prevent="1" placeholder="参数，类似于${arg.name}"/> \
						<div style="margin-top:0.6em"> \
							<a href="" @click.prevent="showVariables">内置变量 <i class="icon angle" :class="{down:!variablesVisible, up:variablesVisible}"></i></a> \
							<div v-show="variablesVisible" style="margin-top:0.4em"> \
								<select class="ui dropdown small" style="width:20em" v-model="vVariable"> \
									<option value="">[内置变量]</option> \
									<option v-for="variable in variables" :value="variable.code">{{variable.code}} - {{variable.name}}</option> \
								</select> \
								<p class="comment">{{vVariableDescription}}</p> \
							</div> \
						</div> \
					</td> \
				</tr> \
				<tr> \
					<td>操作符</td> \
					<td> \
						<select style="width:10em" class="ui dropdown" v-model="vOperator"> \
							<option v-for="operator in vOperators" :value="operator.op">{{operator.name}}</option>\
						</select>\
						<p class="comment">{{vOperatorDescription}}</p> \
					</td> \
				</tr> \
				<tr v-show="!isArrayOperator(vOperator) && hasValue(vOperator)"> \
					<td>对比值</td> \
					<td> \
						<textarea type="text"  v-model="vValue" rows="2" placeholder="对比值"/> \
					</td> \
				</tr> \
				<tr v-show="isArrayOperator(vOperator)"> \
					<td>对比值</td> \
					<td> \
						<div> \
							<div class="ui field" v-for="(v,index) in vValues"> \
								<input type="text" v-model="vValues[index]" style="width:10em" ref="valuesInput"/> \
								<a href="" title="删除" @click.prevent="removeValue(index)"><i class="icon remove small"></i></a> \
							</div> \
							<button class="ui button tiny" type="button" @click.prevent="addValue()">+</button>\
						</div> \
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

/**
 * 更多选项
 */
Vue.component("more-options-indicator", {
	data: function () {
		return {
			visible: false
		};
	},
	methods: {
		changeVisible: function () {
			this.visible = !this.visible;
			if (Tea.Vue != null) {
				Tea.Vue.moreOptionsVisible = this.visible;
			}
		}
	},
	template: '<a href="" style="font-weight: normal" @click.prevent="changeVisible()">更多选项 <i class="icon angle" :class="{down:!visible, up:visible}"></i> </a>'
});

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