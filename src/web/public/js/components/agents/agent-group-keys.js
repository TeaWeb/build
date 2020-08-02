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