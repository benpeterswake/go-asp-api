package asp

import (
	"context"
	"net/http"
	"testing"
)

func TestSignRequest(t *testing.T) {
	urlStr := "https://www.exmaple.com"
	r, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	r = r.WithContext(ctx)

	r, err = signRequest(r)
	if err != nil {
		t.Error(err)
		return
	}
}
