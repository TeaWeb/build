var ENGINE = {
    "version": 0,
    "probes": [function () {
	var probe = new ProcessProbe();
	probe.id = "tea_1545120945542156984";
	probe.name = "ElasticSearch";
	probe.site = "https://www.elastic.co";
	probe.docSite = "https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html";
	probe.developer = "Elasticsearch B.V.";
	probe.commandName = "java";
	probe.commandPatterns = [ "Elasticsearch" ];
	probe.commandVersion = "${commandFile} --version";
	probe.onProcess(function (p) {
		var args = parseArgs(p.cmdline);
		var homeDir = "";
		for (var i = 0; i < args.length; i ++) {
			var arg = args[i];
			var index = arg.indexOf("-Des.path.home=");
			if (index < 0) {
				continue;
			}
			homeDir = arg.substring("-Des.path.home=".length);
		}
		if (homeDir.length > 0) {
			p.dir = homeDir;
			p.file = homeDir + "/bin/elasticsearch";
		}
		return true;
	});
	probe.onParseVersion(function (v) {
		return v;
	});
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545120945562156984"; // 探针ID，
	probe.name = "Apache Http Server"; // App名称
	probe.site = "http://httpd.apache.org/"; // App官方网站
	probe.docSite = "http://httpd.apache.org/docs/current/"; // 官方文档网址
	probe.developer = "The Apache Software Foundation"; // App开发者公司、团队或者个人名称
	probe.commandName = "httpd"; // App启动的命令名称
	probe.commandPatterns = []; // 进程匹配规则
	probe.commandVersion = "${commandFile} -v"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		return v;
	});

	// 运行探针
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545123256877858190"; // 探针ID，
	probe.name = "PHP-FPM"; // App名称
	probe.site = "http://php.net/"; // App官方网站
	probe.docSite = "http://php.net/docs.php"; // 官方文档网址
	probe.developer = "The PHP Group"; // App开发者公司、团队或者个人名称
	probe.commandName = "php-fpm"; // App启动的命令名称
	probe.commandPatterns = []; // 进程匹配规则
	probe.commandVersion = "${commandFile} -v"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		var match = v.match(/PHP \d+\.\d+\.\d+/);
		if (match) {
			return match[0];
		}
		return v;
	});

	// 运行探针
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545123416213326756"; // 探针ID，
	probe.name = "Redis"; // App名称
	probe.site = "https://redis.io/"; // App官方网站
	probe.docSite = "https://redis.io/documentation"; // 官方文档网址
	probe.developer = "redislabs"; // App开发者公司、团队或者个人名称
	probe.commandName = "redis-server"; // App启动的命令名称
	probe.commandPatterns = [""]; // 进程匹配规则
	probe.commandVersion = "{commandFile} -v"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		return v;
	});

	// 运行探针
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545123531232625878"; // 探针ID，
	probe.name = "MongoDB"; // App名称
	probe.site = "https://www.mongodb.com/"; // App官方网站
	probe.docSite = "https://docs.mongodb.com/"; // 官方文档网址
	probe.developer = "MongoDB, Inc"; // App开发者公司、团队或者个人名称
	probe.commandName = "mongod"; // App启动的命令名称
	probe.commandPatterns = ["/mongod"]; // 进程匹配规则
	probe.commandVersion = "${commandFile} --version"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		var match = v.match(/version (v\S+)/);
		if (match) {
			return match[1];	
		}
		return v;
	});

	// 运行探针
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545123651698862290"; // 探针ID，
	probe.name = "nginx"; // App名称
	probe.site = "http://nginx.org/"; // App官方网站
	probe.docSite = "http://nginx.org/en/docs/"; // 官方文档网址
	probe.developer = "nginx.org"; // App开发者公司、团队或者个人名称
	probe.commandName = "nginx"; // App启动的命令名称
	probe.commandPatterns = []; // 进程匹配规则
	probe.commandVersion = "${commandFile} -v"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		var index = v.indexOf("nginx version:");
		if (index > -1) {
			return v.substring("nginx version:".length);
		}
		return v;
	});

	// 运行探针
	probe.run();
},
function () {
	var probe = new ProcessProbe(); // 构造对象
	probe.author = ""; // 探针作者
	probe.id = "tea_1545123805152211975"; // 探针ID，
	probe.name = "MySQL"; // App名称
	probe.site = "https://www.mysql.com/"; // App官方网站
	probe.docSite = "https://dev.mysql.com/doc/"; // 官方文档网址
	probe.developer = "Oracle Corporation"; // App开发者公司、团队或者个人名称
	probe.commandName = "mysqld_safe"; // App启动的命令名称
	probe.commandPatterns = ["mysqld_safe$"]; // 进程匹配规则
	probe.commandVersion = "${commandDir}/mysqld -V"; // 获取版本信息的命令

	// 进程筛选
	probe.onProcess(function (p) {
		return true;
	});

	// 版本信息分析
	probe.onParseVersion(function (v) {
		return v;
	});

	// 运行探针
	probe.run();
}]
};