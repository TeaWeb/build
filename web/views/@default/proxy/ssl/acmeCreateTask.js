Tea.context(function () {
	this.step = "prepare";

	/**
	 * 准备
	 */
	this.confirmPrepared = function () {
		this.step = "config";
	};

	/**
	 * 配置
	 */
	this.userId = "";
	this.domains = [];
	this.domainIsAdding = false;
	this.addingDomain = "";

	this.removeDomain = function (domain) {
		this.domains.$removeValue(domain);
	};

	this.addDomain = function () {
		this.domainIsAdding = true;
		this.$delay(function () {
			this.$find("input[name='addingDomain']").focus();
		});
	};

	this.confirmAddDomain = function () {
		if (this.addingDomain.length == 0) {
			alert("请输入要添加的域名");
			this.$find("input[name='addingDomain']").focus();
			return;
		}
		this.domains.push(this.addingDomain);
		this.addingDomain = "";
		this.domainIsAdding = false;
	};

	this.cancelAddingDomain = function () {
		this.domainIsAdding = false;
		this.addingDomain = "";
	};

	this.submitConfig = function () {
		if (this.domains.length == 0) {
			alert("请至少添加一个域名");
			this.addDomain();
			return;
		}

		this.step = "auth";
		this.dnsRecords = [];
		this.dnsError = "";
		this.dnsCheckingError = "";

		this.dnsRecords = [];

		this.$post(".acmeRecords")
			.params({
				"userId": this.userId,
				"domains": this.domains.join(",")
			})
			.success(function (response) {
				this.dnsRecords = response.data.records;
			})
			.fail(function (response) {
				this.dnsError = response.message;
			});
	};

	/**
	 * DNS
	 */
	this.dnsRecords = [];
	this.dnsError = "";
	this.dnsIsChecking = false;
	this.dnsCheckingError = "";
	this.dnsHelpVisible = false;

	this.checkDNS = function () {
		this.dnsIsChecking = true;
		this.dnsCheckingError = "";

		this.$post(".acmeDnsChecking")
			.params({
				"userId": this.userId,
				"domains": this.domains.join(","),
				"records": JSON.stringify(this.dnsRecords),
				"serverId": this.server.id
			})
			.success(function () {
				this.step = "finish";
			})
			.fail(function (response) {
				this.dnsCheckingError = response.message;
			})
			.done(function () {
				this.dnsIsChecking = false;
			});
	};

	this.showDnsHelp = function () {
		this.dnsHelpVisible = !this.dnsHelpVisible;
	};
});