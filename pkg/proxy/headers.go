package proxy

const (

	// HeaderXProxyProcess is the header that tell where the request needs to be
	// processed or forwarded.
	//
	// By default, the proxy will forward the request to the service.
	// If the header is present, the proxy will process the request on its own.
	HeaderXProxyProcess = "X-Proxy-Process"
)
