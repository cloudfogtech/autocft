package model

type SystemConfig struct {
	CFAPIToken  string `json:"cfApiToken" env:"CF_API_TOKEN"`
	CFAccountID string `json:"cfAccountId" env:"CF_ACCOUNT_ID"`
	CFTunnelID  string `json:"cfTunnelId" env:"CF_TUNNEL_ID"`
	Basedir     string `json:"basedir" env:"BASEDIR"`
	Cron        string `json:"cron" env:"CRON"`
}

// IngressConfig define as zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress
// see https://developers.cloudflare.com/api/go/resources/zero_trust/subresources/tunnels/subresources/cloudflared/
type IngressConfig struct {
	Enabled bool `json:"enabled" tag:"autocft.enabled"`
	// Public hostname for this service.
	// <br>Required</br>
	Hostname string `json:"hostname" tag:"autocft.hostname"`
	// Protocol and address of destination server. Supported protocols: http://,
	// https://, unix://, tcp://, ssh://, rdp://, unix+tls://, smb://. Alternatively
	// can return an HTTP status code http_status:[code] e.g. 'http_status:404'.
	// <br>Required</br>
	Service string `json:"service" tag:"autocft.service"`
	// Requests with this path route to this public hostname.
	// Default: '/'
	Path string `json:"path" tag:"autocft.path"`
	// Configuration parameters for the public hostname specific connection settings
	// between cloudflared and origin server.
	Origin *IngressOriginConfig `json:"origin"`
}

type IngressOriginConfig struct {
	// Timeout for establishing a new TCP connection to your origin server. This
	// excludes the time taken to establish TLS, which is controlled by tlsTimeout.
	ConnectTimeout int64 `json:"connectTimeout" tag:"autocft.origin.connect-timeout" env:"ORIGIN_CONNECT_TIMEOUT"`
	// Disables chunked transfer encoding. Useful if you are running a WSGI server.
	DisableChunkedEncoding bool `json:"disableChunkedEncoding" tag:"autocft.origin.disable-chunked-encoding" env:"ORIGIN_DISABLE_CHUNKED_ENCODING"`
	// Attempt to connect to origin using HTTP2. Origin must be configured as https.
	HTTP2Origin bool `json:"http2Origin" tag:"autocft.origin.http2-origin" env:"ORIGIN_HTTP2_ORIGIN"`
	// Sets the HTTP Host header on requests sent to the local service.
	HTTPHostHeader string `json:"httpHostHeader" tag:"autocft.origin.http-host-header" env:"ORIGIN_HTTP_HEADER"`
	// Maximum number of idle keepalive connections between Tunnel and your origin.
	// This does not restrict the total number of concurrent connections.
	KeepAliveConnections int64 `json:"keepAliveConnections" tag:"autocft.origin.keep-alive-connections" env:"ORIGIN_KEEP_ALIVE_CONNECTIONS"`
	// Timeout after which an idle keepalive connection can be discarded.
	KeepAliveTimeout int64 `json:"keepAliveTimeout" tag:"autocft.origin.keep-alive-timeout" env:"ORIGIN_KEEP_ALIVE_TIME"`
	// Disable the “happy eyeballs” algorithm for IPv4/IPv6 fallback if your local
	// network has misconfigured one of the protocols.
	NoHappyEyeballs bool `json:"noHappyEyeballs" tag:"autocft.origin.no-happy-eyeballs" env:"ORIGIN_NO_HAPPY_EYEBALLS"`
	// Disables TLS verification of the certificate presented by your origin. Will
	// allow any certificate from the origin to be accepted.
	NoTLSVerify bool `json:"noTLSVerify" tag:"autocft.origin.no-tls-verify" env:"ORIGIN_NO_TLS_VERIFY"`
	// Hostname that cloudflared should expect from your origin server certificate.
	OriginServerName string `json:"originServerName" tag:"autocft.origin.origin-server-name" env:"ORIGIN_ORIGIN_SERVER_NAME"`
	// cloudflared starts a proxy server to translate HTTP traffic into TCP when
	// proxying, for example, SSH or RDP. This configures what type of proxy will be
	// started. Valid options are: "" for the regular proxy and "socks" for a SOCKS5
	// proxy.
	ProxyType string `json:"proxyType" tag:"autocft.origin.proxy-type" env:"ORIGIN_PROXY_TYPE"`
	// The timeout after which a TCP keepalive packet is sent on a connection between
	// Tunnel and the origin server.
	TCPKeepAlive int64 `json:"tcpKeepAlive" tag:"autocft.origin.tcp-keep-alive" env:"ORIGIN_TCP_KEEP_ALIVE"`
	// Timeout for completing a TLS handshake to your origin server, if you have chosen
	// to connect Tunnel to an HTTPS server.
	TLSTimeout int64 `json:"tlsTimeout" tag:"autocft.origin.tls-timeout" env:"ORIGIN_TLS_TIMEOUT"`
}

type SyncOptions struct {
	Dry bool
}

type SyncCount struct {
	WebManaged int
	Unchanged  int
	Updated    int
	Added      int
	Deleted    int
}
