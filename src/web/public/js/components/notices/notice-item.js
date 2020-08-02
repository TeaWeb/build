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