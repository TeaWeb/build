package forms

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"net/http"
)

// 环境变量
type EnvBox struct {
	Element `yaml:",inline"`
}

func NewEnvBox(title string, subtitle string) *EnvBox {
	return &EnvBox{
		Element{
			Title:    title,
			Subtitle: subtitle,
		},
	}
}

func (this *EnvBox) Super() *Element {
	return &this.Element
}

func (this *EnvBox) Compose() string {
	this.Javascript = `
/**
 * 环境变量
 */
this.env = `
	if types.IsSlice(this.Value) {
		this.Javascript += stringutil.JSONEncode(this.Value)
	} else {
		this.Javascript += "[]"
	}

	this.Javascript += `;
this.envAdding = false;
this.envAddingName = "";
this.envAddingValue = "";

this.addEnv = function () {
	this.envAdding = !this.envAdding;
	this.$delay(function () {
		this.$find("form input[name='envAddingName']").focus();
	});
};

this.confirmAddEnv = function () {
	if (this.envAddingName.length == 0) {
		alert("请输入变量名");
		this.$find("form input[name='envAddingName']").focus();
		return;
	}
	this.env.push({
		"name": this.envAddingName,
		"value": this.envAddingValue
	});
	this.envAdding = false;
	this.envAddingName = "";
	this.envAddingValue = "";
};

this.removeEnv = function (index) {
	this.env.$remove(index);
};

this.cancelEnv = function () {
	this.envAdding = false;
};`

	return `
<div class="ui field">
	<span class="ui label small" v-for="(var1, index) in env">
		<input type="hidden" name="` + this.Namespace + "_" + this.Code + `_envNames" :value="var1.name"/>
		<input type="hidden" name="` + this.Namespace + "_" + this.Code + `_envValues" :value="var1.value"/>
		<em>{{var1.name}}</em>: {{var1.value}}
		<a href="" @click.prevent="removeEnv(index)"><i class="icon remove"></i></a>
	</span>
</div>
<div v-if="envAdding" class="ui fields inline">
	<div class="ui field">
		<input type="text" name="envAddingName" v-model="envAddingName" placeholder="变量名" style="width:9em" @keyup.enter="confirmAddEnv" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<input type="text" name="envAddingValue" v-model="envAddingValue" placeholder="变量值" style="width:15em" @keyup.enter="confirmAddEnv" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<button class="ui button" type="button" @click="confirmAddEnv()">添加</button>
	</div>
	<div class="ui field" style="padding-left:0;padding-right:0">
		<a href="" @click.prevent="cancelEnv()" title="删除"><i class="icon remove"></i></a>
	</div>
</div>
<div class="ui field">
	<button class="ui button small" type="button" @click="addEnv()">+</button>
</div>`
}

func (this *EnvBox) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	names, found := req.Form[this.Namespace+"_"+this.Code+"_envNames"]
	envs := []*shared.Variable{}
	if found {
		values, found := req.Form[this.Namespace+"_"+this.Code+"_envValues"]
		if found {
			for index, name := range names {
				if index < len(values) {
					envs = append(envs, &shared.Variable{
						Name:  name,
						Value: values[index],
					})
				} else {
					envs = append(envs, &shared.Variable{
						Name:  name,
						Value: "",
					})
				}
			}
		}
	}
	return envs, false, nil
}
