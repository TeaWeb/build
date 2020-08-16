package scripts

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/robertkrimen/otto"
	"testing"
)

func TestEngine_Run(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`var widget = new widgets.Widget({
	"name": "测试Widget",
	"code": "test_stat@tea",
	"author": "我是测试的",
	"version": "0.0.1"
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.html = "测试HTML Chart";
	chart.render();
};`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(engine.Charts())
}

func TestEngine_Cache(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`
var widget = new widgets.Widget({});
widget.run = function () {
	caches.set("a", "b");
	console.log(caches.get("a"));
};
`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEngine_Log(t *testing.T) {
	engine := NewEngine()
	engine.SetContext(&Context{
		Server: &teaconfigs.ServerConfig{
			Id: "123",
		},
	})
	err := engine.RunCode(`
var widget = new widgets.Widget({});
widget.run = function () {
	var query = new logs.Query();
	query.attr("status", [200]);
	query.count();
};
`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEngine_Menu(t *testing.T) {
	vm := otto.New()
	_, err := vm.Run(`
var a = function () {
	this.menus = [];
};

var b = function () {
	this.menus = [];

	this.addMenu = function () {
		var m = new menu();
		m.a = this;
		this.menus.push(m);
	};

	this.render = function () {
		console.log(this.menus.length);
	};
};

var menu = function () {
	this.a = null;
};

a.prototype = new b();

var a1 = new a();
a1.addMenu();
a1.render();
var a2 = new a();
//a2.addMenu();
a2.render();
`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
