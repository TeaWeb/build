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