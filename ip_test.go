package ip

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/vinxi/utils.v0"
)

func TestIPFilterRanges(t *testing.T) {
	tests := []struct {
		IPs      string
		remoteIP string
		allowed  bool
	}{
		{"127.0.0.1/8", "127.0.0.1", true},
		{"127.0.0.1/8", "127.0.0.1:1234", true},
		{"127.0.0.1/8", "127.0.0.10", true},
		{"127.0.0.1/8", "127.10.1.1", true},
		{"127.0.0.1/8", "128.0.0.1", false},
		{"127.0.0.1/8", "127.0.1.1", true},
		{"127.0.0.1/8", "127.0.1.1:54231", true},
		{"192.168.0.1/27", "192.168.0.1", true},
		{"192.168.0.1/27", "192.168.0.26", true},
		{"192.168.0.1/27", "192.168.0.96", false},
		{"1.0.0.0/24", "1.0.0.1", true},
		{"1.0.0.0/24", "1.0.0.255", true},
		{"1.0.0.0/24", "1.0.1.1", false},
		{"::1/64", "[::1]:3223", true},
		{"2001:db8:a0b:12f0::1/64", "[2001:db8:a0b:12f0::1]:3223", true},
	}

	for _, test := range tests {
		filter := New(test.IPs)

		var allowed bool
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed = true
		})

		rw := utils.NewWriterStub()
		req := &http.Request{RemoteAddr: test.remoteIP}

		filter.FilterHTTP(handler)(rw, req)

		if test.allowed {
			rw.WriteHeader(200)
			rw.Write([]byte("foo"))
			st.Expect(t, allowed, true)
			st.Expect(t, rw.Code, 200)
			st.Expect(t, string(rw.Body), "foo")
		} else {
			st.Expect(t, allowed, false)
			st.Expect(t, rw.Code, 403)
			st.Expect(t, string(rw.Body), "Forbidden: client IP not allowed")
		}
	}
}
