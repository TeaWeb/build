var widget = new widgets.Widget({
	"name": "代理状态",
	"code": "proxy_status@tea",
	"author": "TeaWeb",
	"version": "0.0.1"
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.options.name = "代理状态";
	chart.options.columns = 1;

	// ports
	var ports = [];
	if (context.server.listen != null) {
		context.server.listen.$each(function (k, v) {
			if (v.length > 0) {
				var index = v.indexOf(":");
				var port = "80";
				if (index > -1) {
					port = v.substring(index + 1);
				}
				if (!ports.$contains(port)) {
					ports.push(port);
				}
			}
		});
	}
	if (context.server.ssl.listen != null) {
		context.server.ssl.listen.$each(function (k, v) {
			if (v.length > 0) {
				var index = v.indexOf(":");
				var port = "443";
				if (index > -1) {
					port = v.substring(index + 1);
				}
				if (!ports.$contains(port)) {
					ports.push(port);
				}
			}
		});
	}
	if (ports.length > 0) {
		chart.options.name = "代理状态<em>（已绑定端口：" + ports.join(", ") + "）</em>";
	} else {
		chart.options.name = "代理状态<em>（还没有绑定网络地址）</em>";
	}

	chart.html = "<style type='text/css'> \
    .backends-box { \
		position: absolute; \
		width: 45%; \
		right: 0; \
		top: 0; \
		bottom: 0; \
		overflow-y: auto; \
     } \
    .backends-box .backend { \
         font-size: 0.8em; \
    } \
    .backends-box .backend .green { \
		display: inline-block; \
		width: 8px; \
		height: 8px; \
		background: #21ba45; \
		margin: 0 2px; \
    } \
    .backends-box .backend .grey { \
		display: inline-block; \
		width: 8px; \
		height: 8px; \
		background: grey; \
		margin: 0 2px; \
    } \
    .backends-box .backend .red { \
        display: inline-block; \
        width: 8px; \
        height: 8px; \
        background: #db2828; \
        margin: 0 2px; \
    } \
    .backends-box .backend .blue { \
        display: inline-block; \
        width: 8px; \
        height: 8px; \
        background: #2185d0; \
        margin: 0 2px; \
    } \
    .summary-state { \
        position: absolute; \
        left: 0; \
        top: 0; \
        bottom: 0; \
        width: 45%; \
	} \
    .summary-state .circle { \
		width: 120px; \
		height: 120px; \
		border-radius: 50%; \
		position:absolute; \
		top: 50%; \
		margin-top: -60px; \
		left: 50%; \
		margin-left: -60px; \
		text-align: center; \
		line-height: 120px; \
		color: white; \
		font-size: 70px; \
		opacity: 0.7; \
    } \
    .summary-state.red .circle  { \
        background: #db2828; \
    } \
    .summary-state.grey .circle { \
        background: grey; \
    } \
    .summary-state.green .circle { \
		background: #21ba45; \
    } \
    </style>";

	var hasDown = false;
	var hasOn = false;
	var backends = context.server.backends;
	for (var i = 0; i < backends.length; i++) {
		var backend = backends[i];
		if (backend.on) {
			hasOn = true;
		}
		if (backend.on && backend.isDown) {
			hasDown = true;
		}
	}
	var summaryState = "green";
	if (backends.length == 0) {
		summaryState = "";
	} else if (hasDown) {
		summaryState = "red";
	} else if (!hasOn) {
		summaryState = "grey";
	}

	chart.html += "<div>";
	if (backends.length == 0) {
		chart.html = "<div>";
		chart.html += "<p class='grey'><i class='icon paper plane'></i>暂时还没有配置后端服务</p>";
	} else {
		chart.html += "<div class='summary-state " + summaryState + "'><div class='circle'>" + ((backends.length > 0) ? backends.length : "") + "</div></div><div class='backends-box'>";
		for (var i = 0; i < backends.length; i++) {
			var backend = backends[i];
			chart.html += "<div class='backend'>";
			chart.html += "<span class='on-status " + (backend.on ? "blue" : "grey") + "' title='开启状态'></span><span class='down-status " + ((backend.on && !backend.isDown) ? "green" : "red") + "' title='连接状态'></span><span class='address'>" + backend.address + "</span>";
			chart.html += "</div>";
		}
	}
	chart.html += "</div></div>";
	chart.render();
};
