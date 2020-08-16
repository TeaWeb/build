package demo

// init the extension
func init() {
	// add request hook
	/**teaproxy.AddRequestHook(&teaproxy.RequestHook{
		BeforeRequest: func(req *teaproxy.Request, writer *teaproxy.ResponseWriter) (goNext bool) {
			writer.WriteString("hello")
			return false
		},
	})**/
}
