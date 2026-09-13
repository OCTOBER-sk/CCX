package timeout

// Timeouts per spec section 31, lines 1165-1179.
// Separate timeouts (line 1169-1174):
//   connect timeout
//   request header timeout
//   idle stream timeout
//   total request timeout
//   shutdown grace period
// Line 1177: long thinking streams must not be killed merely because no text arrived.
// Line 1179: use event/keepalive activity to reset idle timers.

type Config struct {
	ConnectTimeout      int // connect timeout (line 1170)
	HeaderTimeout       int // request header timeout (line 1171)
	IdleStreamTimeout   int // idle stream timeout (line 1172)
	TotalRequestTimeout int // total request timeout (line 1173)
	ShutdownGrace       int // shutdown grace period (line 1174)
}

func Default() Config {
	return Config{
		ConnectTimeout:      10,  // seconds
		HeaderTimeout:       30,  // seconds
		IdleStreamTimeout:   60,  // seconds; reset by event/keepalive (line 1179)
		TotalRequestTimeout: 300, // seconds
		ShutdownGrace:       5,   // seconds
	}
}
