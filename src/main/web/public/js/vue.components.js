/**
 * 自动补全
 */
Vue.component("auto-complete-box", {
	props: ["name", "placeholder", "options", "value", "maxlength"],
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
	props: ["name", "placeholder", "value", "maxlength"],
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

