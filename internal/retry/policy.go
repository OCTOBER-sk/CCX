package retry

// Retry policy per spec 27 (lines 1047-1075):
// Default: 400/401/403/404 → no retry; 429 → bounded backoff; 5xx → bounded retry; timeout/network → bounded retry.
// Honor Retry-After (line 1062-1064). Use exponential backoff + jitter + max attempts + max elapsed time (lines 1066-1073).
// Retries must be aware of request side effects (line 1075) — no replay after ambiguous upstream-accepted tool call (spec 87).

type Policy struct {
	MaxAttempts    int
	MaxElapsedMs   int
	BackoffBaseMs  int
	BackoffMaxMs   int
	JitterFactor   float64
}

func Default() Policy {
	return Policy{
		MaxAttempts:   3,
		MaxElapsedMs:  30000,
		BackoffBaseMs: 500,
		BackoffMaxMs:  8000,
		JitterFactor:  0.3,
	}
}

// Global retry rate budget per spec 264 (lines 6888-6903):
// Per provider retry rate budget to prevent 100 sessions × 3 retries = 300-request storm.
func GlobalRetryRateBudget() int {
	return 100 // max concurrent retries per provider (storm prevention)
}

func ShouldRetry(statusCode int, retryClass string) bool {
	switch statusCode {
	case 400, 401, 403, 404:
		return false // no retry per spec 1053
	case 429:
		return true // bounded backoff per spec 1054
	case 500, 502, 503, 504:
		return true // bounded retry per spec 1055
	default:
		if retryClass == "safe" || retryClass == "conditional" {
			return true
		}
		return false
	}
}
