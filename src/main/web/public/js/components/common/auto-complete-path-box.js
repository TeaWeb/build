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