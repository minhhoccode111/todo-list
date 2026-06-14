package v1

import (
	"net/http"
	"testing"
)

// TestRefreshCookieSameSite guards the browser invariant that a SameSite=None
// cookie is only accepted together with Secure. The refresh-token cookie must
// therefore use None only when Secure is on, and fall back to Lax otherwise.
func TestRefreshCookieSameSite(t *testing.T) {
	t.Parallel()

	if got := refreshCookieSameSite(true); got != http.SameSiteNoneMode {
		t.Errorf("refreshCookieSameSite(true) = %v, want SameSiteNoneMode", got)
	}

	if got := refreshCookieSameSite(false); got != http.SameSiteLaxMode {
		t.Errorf("refreshCookieSameSite(false) = %v, want SameSiteLaxMode", got)
	}
}
