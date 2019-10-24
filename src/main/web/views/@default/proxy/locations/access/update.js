Tea.context(function () {
	this.submitSuccess = function () {
		window.location = "/proxy/locations/access?serverId=" + this.server.id + "&locationId=" + this.location.id;
	};

	// format policy
	if (this.policy.traffic.second.duration < 1) {
		this.policy.traffic.second.duration = 1;
	}
	if (this.policy.traffic.minute.duration < 1) {
		this.policy.traffic.minute.duration = 1;
	}
	if (this.policy.traffic.hour.duration < 1) {
		this.policy.traffic.hour.duration = 1;
	}
	if (this.policy.traffic.day.duration < 1) {
		this.policy.traffic.day.duration = 1;
	}
	if (this.policy.traffic.month.duration < 1) {
		this.policy.traffic.month.duration = 1;
	}

	this.trafficOn = this.policy.traffic.on;
	this.trafficTotalOn = this.policy.traffic.total.on;
	this.trafficSecondOn = this.policy.traffic.second.on;
	this.trafficMinuteOn = this.policy.traffic.minute.on;
	this.trafficHourOn = this.policy.traffic.hour.on;
	this.trafficDayOn = this.policy.traffic.day.on;
	this.trafficMonthOn = this.policy.traffic.month.on;

	this.accessOn = this.policy.access.on;
	this.accessAllowOn = this.policy.access.allowOn;
	this.accessDenyOn = this.policy.access.denyOn;
	if (this.policy.access.allow == null) {
		this.accessAllowIPs = [];
	} else {
		this.accessAllowIPs = this.policy.access.allow.$map(function (k, v) {
			return v.ip;
		});
	}
	if (this.policy.access.deny == null) {
		this.accessDenyIPs = [];
	} else {
		this.accessDenyIPs = this.policy.access.deny.$map(function (k, v) {
			return v.ip;
		});
	}
});