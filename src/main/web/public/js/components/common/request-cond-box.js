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