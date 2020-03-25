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