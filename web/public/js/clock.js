// 原始文件：https://raw.githubusercontent.com/MichalPaszkiewicz/clockjs/gh-pages/clock.js
// 主页：http://www.michalpaszkiewicz.co.uk/clockjs/

Date.prototype.addHours= function(h){
	this.setHours(this.getHours()+h);
	return this;
}

Date.prototype.addMinutes = function(m){
	this.setMinutes(this.getMinutes()+m);
	return this;
}

Date.prototype.addSeconds = function(s){
	this.setSeconds(this.getSeconds()+s);
	return this;
}

var romanNumerals = [0,"I","II","III","IV","V","VI","VII","VIII","IX","X","XI","XII"];

var clock = function(id, options){
	var self = this;

	self.started = false;

	//initialise canvas && context
	self.canvas = document.getElementById(id);
	self.context = self.canvas.getContext("2d");

	var currentDate = new Date();

	//default options
	self.options = {
		radius: function(){ return Math.min(self.canvas.height, self.canvas.width) / 2 },
		colour:"rgba(255,0,0,0.2)",
		rim: function(){ return getValue("radius") * 0.2; },
		rimColour: function(){ return self.options.colour; },
		x: function(){ return self.canvas.width / 2 },
		y: function(){ return self.canvas.height / 2 },
		lineColour: function(){ return self.options.colour; },
		fillColour: function(){  return self.options.colour; },
		lineWidth: 1,
		centreCircle: true,
		centreCircleRadius: function(){ return getValue("radius") * 0.03; },
		centreCircleColour: function(){return getValue("colour");},
		centreCircleCutout: function(){ return getValue("radius") * 0.01; },
		addHours: 0,
		addMinutes: 0,
		addSeconds: 0,
		directionCoefficient: 1,
		markerType: "number",
		markerColour: function(){ return self.options.colour; },
		markerSize: function(){ return getValue("radius") * 0.02; },
		markerDistance: function(){ return getValue("radius") * 0.9; },
		markerDisplay: true,
	};

	//hands settings
	self.hands = {
		secondHand:{
			length: 1, width: 0.1,
			percentile:function(){
				return (currentDate.getSeconds() + currentDate.getMilliseconds() / 1000) / 60;
			}},
		minuteHand:{
			length: 0.8, width: 0.4,
			percentile:function(){
				return (currentDate.getMinutes() + currentDate.getSeconds() / 60) / 60;
			}},
		hourHand:{
			length: 0.5, width: 0.9,
			percentile:function(){
				return (currentDate.getHours() + currentDate.getMinutes() / 60) / 12;
			}}
	}

	//set specified options
	for (var key in options) {
		if (options.hasOwnProperty(key)) {
			self.options[key] = options[key];
		}
	}

	//get function - gets a function, otherwise value.
	var getValue = function(name, defaultName){
		if(name == null){
			if(defaultName == null){
				throw new Error("No value set for this option.");
			}
			if(typeof defaultName == "function"){
				return defaultName();
			}
			return (typeof self.options[defaultName] == "function") ? self.options[defaultName]() : self.options[defaultName];
		}
		if(self.options == null){
			throw new Error("Someone has deleted the clock's options. Uh-oh!");
		}
		if(typeof self.options[name] == "function"){
			var result = self.options[name]();
			if(result != null){
				return result;
			}
			if(typeof defaultName == "function"){
				return defaultName();
			}
			return (typeof self.options[defaultName] == "function") ? self.options[defaultName]() : self.options[defaultName];
		}
		else {
			return self.options[name];
		}
	}

	//for drawing a handleEvent
	var drawHand = function(x, y, radius, theta, lineWidth){
		self.context.lineWidth = 1;
		self.context.beginPath();
		self.context.moveTo(x,y);
		var offAmount = (lineWidth != null) ? lineWidth : 0.5;
		var one = {x: x + 2 * radius / 8 * Math.cos(theta + offAmount), y: y + 2 * radius / 8 * Math.sin(theta + offAmount)};
		var two = {x: x, y: y};
		var one2 = {x: x + 2 * radius / 8 * Math.cos(theta - offAmount), y: y + 2 * radius / 8 * Math.sin(theta - offAmount)};
		var finalx = x + radius * Math.cos(theta);
		var finaly = y + radius * Math.sin(theta);
		self.context.bezierCurveTo(one.x, one.y, two.x, two.y, finalx,finaly);
		self.context.bezierCurveTo(two.x, two.y, one2.x, one2.y, x,y);
		self.context.stroke();
		self.context.fill();
		self.context.lineWidth = 1;
	}

	//draw single marker on the clock
	var drawMarker = function(x, y, i){
		self.context.beginPath();
		self.context.fillStyle = getValue("markerColour", "colour");
		var markerSize = getValue("markerSize");

		switch(getValue("markerType")){
			case "numeral":
				markerSize *= 4;
				self.context.font = markerSize + "px sans-serif";
				self.context.textAlign = "center";
				self.context.fillStyle = getValue("markerColour");
				self.context.textBaseline = "middle";
				self.context.fillText(romanNumerals[i + 1],x,y);
				break;
			case "number":
				markerSize *= 4;
				self.context.font = markerSize + "px sans-serif";
				self.context.textAlign = "center";
				self.context.fillStyle = getValue("markerColour");
				self.context.textBaseline = "middle";
				self.context.fillText(i + 1,x,y);
				break;
			case "dot":
				self.context.arc(x,y,markerSize,0,2*Math.PI);
				self.context.fill();
				break;
			case "none":
			default:
				return;
		}
	}

	//for drawing the markers on the clock
	var drawMarkers = function(x, y){
		if(getValue("markerDisplay") == false){
			return;
		}
		var directionCoefficient = getValue("directionCoefficient");
		var markerDistance = getValue("markerDistance");
		var theta = directionCoefficient * 2 * Math.PI / 12 - Math.PI / 2
		for(var i = 0; i < 12; i++){
			var markerX =	x + markerDistance * Math.cos(theta);
			var markerY = y + markerDistance * Math.sin(theta);
			drawMarker(markerX, markerY, i);
			theta += directionCoefficient * 2 * Math.PI / 12;
		}
	}

	//update the date, change time zone etc.
	var updateDate = function(){
		//update date;
		currentDate = new Date();
		currentDate.addHours(getValue("addHours", function(){return 0;}));
		currentDate.addMinutes(getValue("addMinutes", function(){return 0;}));
		currentDate.addSeconds(getValue("addSeconds", function(){return 0;}));
	}

	//updates and draws clock
	self.draw = function(){
		// 刘祥超增加判断
		if (self.canvas.parentNode == null) {
			return;
		}

		self.canvas.height = 90; //self.canvas.parentNode.offsetHeight; // 刘祥超修改
		self.canvas.width = 200; //self.canvas.parentNode.offsetWidth; // 刘祥超修改

		var radius = getValue("radius");
		var x = getValue("x");
		var y = getValue("y");

		self.context.clearRect(0,0, self.canvas.width, self.canvas.height);

		//outer circle
		if(getValue("rim") != "none"){
			self.context.strokeStyle = getValue("rimColour");
			self.context.lineWidth = getValue("rim");
			self.context.beginPath();
			self.context.arc(x,y,radius - getValue("rim")/2,0,2*Math.PI);
			self.context.stroke();

			self.context.strokeStyle = getValue("lineColour");
			self.context.fillStyle = getValue("fillColour");
			self.context.lineWidth = getValue("lineWidth");
		}

		//markers
		drawMarkers(x, y);

		updateDate();

		var directionCoefficient = getValue("directionCoefficient", function(){return 1;});

		//draw all hands
		for (var key in self.hands) {
			if (self.hands.hasOwnProperty(key)) {
				var tempTheta = directionCoefficient * self.hands[key].percentile() * 2 * Math.PI - Math.PI / 2;
				var tempRadius = radius * self.hands[key].length;
				drawHand(x, y, tempRadius, tempTheta, self.hands[key].width);
			}
		}

		//centreCircle
		if(getValue("centreCircle")){
			self.context.beginPath();
			self.context.fillStyle = getValue("centreCircleColour", "colour");
			self.context.arc(x,y,getValue("centreCircleRadius"),0,2*Math.PI);
			self.context.fill();
			self.context.stroke();

			//cutout
			self.context.beginPath();
			self.context.arc(x,y,getValue("centreCircleCutout"),0,2*Math.PI);
			self.context.clip();
			self.context.clearRect(0,0,self.canvas.width, self.canvas.height);
		}
	};

	self.animate = function(){
		if(self.started == false){
			return;
		}

		self.draw();

		window.requestAnimationFrame(self.animate);
	};

	self.start = function(){
		self.started = true;

		self.animate();
	};

	self.stop = function(){
		self.started = false;
	}

	return self;
}

var clockMaker = function(){
	var maker = this;

	maker.started = false;

	maker.clocks = [];

	maker.addClock = function(clockItem, options){
		if(typeof clockItem == "string"){
			clockItem = new clock(clockItem, options);
		}

		maker.clocks.push({clock: clockItem, started: true});
	}

	maker.draw = function(){
		for(var i in maker.clocks){
			var currentClock = maker.clocks[i];
			if(currentClock.started == true){
				currentClock.clock.draw();
			}
		}
	}

	maker.animate = function(){
		if(maker.started == false){
			return;
		}

		maker.draw();

		window.requestAnimationFrame(maker.animate);
	}

	maker.start = function(){
		maker.started = true;

		for(var i in maker.clocks){
			var currentClock = maker.clocks[i];
			currentClock.started = true;
		}

		maker.animate();
	}

	maker.stop = function(){
		maker.started = false;
	}

	return maker;
}
