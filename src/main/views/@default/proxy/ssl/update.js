Tea.context(function () {
    this.httpsOn = (this.proxy.ssl != null && this.proxy.ssl.on);

    this.submitSuccess = function () {
        alert("修改成功");

        window.location = "/proxy/ssl?server=" + this.proxy.filename;
    };
});