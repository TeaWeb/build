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