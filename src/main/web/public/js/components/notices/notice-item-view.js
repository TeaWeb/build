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