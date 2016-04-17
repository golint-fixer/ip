// Package ip filters HTTP traffic based on IP ranges.
// Supports IP v4/v6 and CIDR expressions.
// It also provides a middleware layer, useful for composition and multiplexing.
package ip

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"gopkg.in/vinxi/layer.v0"
)

// FilterFunc represents the filter function signature used to
// determine if should apply not the IP filtering.
type FilterFunc func(r *http.Request) bool

// ForbiddenResponder is used as default function to repond when the
// IP is not allowed. You can customize it via Filter.SetResponder(fn).
var ForbiddenResponder = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	w.Write([]byte("Forbidden: client IP not allowed"))
}

// Filter implements a IP range based authorization filter for incoming HTTP traffic.
type Filter struct {
	// layer stores the middleware layer.
	layer *layer.Layer
	// responser stores the responder function used when the IP is not allowed.
	responder http.HandlerFunc
	// filters stores a list of filters to determine if should apply the IP filter.
	filters []FilterFunc
	// ranges stores the allowed IP v4/v6 ranges.
	ranges []*net.IPNet
}

// New creates a new IP filter based on the given IP CIDR ranges.
func New(ranges ...string) *Filter {
	return &Filter{
		layer:     layer.New(),
		ranges:    parseRanges(ranges),
		responder: ForbiddenResponder,
	}
}

// SetResponder sets a custom function to reply in case that an IP not allowed.
func (f *Filter) SetResponder(fn http.HandlerFunc) {
	f.responder = fn
}

// Filter registers a new filter function.
// If the filter matches, the client IP won't be validated.
func (f *Filter) Filter(fn ...FilterFunc) {
	f.filters = append(f.filters, fn...)
}

// Use registers a new middleware handler in the middleware stack.
// The middleware will be executed only if the client IP is allowed.
func (f *Filter) Use(handler interface{}) *Filter {
	f.layer.Use(layer.RequestPhase, handler)
	return f
}

// UsePhase registers a new middleware handler in the middleware stack.
// The middleware will be executed only if the client IP is allowed.
func (f *Filter) UsePhase(phase string, handler interface{}) *Filter {
	f.layer.Use(phase, handler)
	return f
}

// UseFinalHandler registers a new middleware handler in the middleware stack.
// The middleware will be executed only if the client IP is allowed.
func (f *Filter) UseFinalHandler(handler http.Handler) *Filter {
	f.layer.UseFinalHandler(handler)
	return f
}

// Register registers the middleware handler.
func (f *Filter) Register(mw layer.Middleware) {
	mw.UsePriority("request", layer.Head, f.FilterHTTP)
}

// FilterHTTP filters an incoming HTTP request based on the client IP.
// If some filter passes, the request won't be limited.
func (f *Filter) FilterHTTP(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pass filters to determine if should apply the IP filter.
		// All the filters must pass to apply the IP filter.
		for _, filter := range f.filters {
			if !filter(r) {
				h.ServeHTTP(w, r)
				return
			}
		}

		// Check if the client IP is allowed.
		if !matchIPInRange(f.ranges, r.RemoteAddr) {
			// Otherwise reply with the default forbidden responder.
			f.responder(w, r)
			return
		}

		// Trigger the filter specific middleware layer
		// and forward the request.
		f.layer.Run("request", w, r, h)
	}
}

// parseRanges parses a range of CIDR expressions
// as a CIDR notation IP address and mask,
// like "192.168.100.1/24" or "2001:DB8::/48", as defined
// in RFC 4632 and RFC 4291.
func parseRanges(ranges []string) []*net.IPNet {
	cidrs := []*net.IPNet{}
	for _, expr := range ranges {
		_, cidr, err := net.ParseCIDR(expr)
		if err != nil {
			fmt.Errorf("Error parsing CIDR expression: %s\n", expr)
			continue
		}
		cidrs = append(cidrs, cidr)
	}
	return cidrs
}

// matchIPInRange compares if a given IP is contained in a range of IPs.
func matchIPInRange(ranges []*net.IPNet, IP string) bool {
	// Split by colons (also supports IPv6)
	parts := strings.Split(IP, ":")

	// Remove port from expression, if present
	if len(parts) > 1 {
		parts = parts[0 : len(parts)-1]
	}
	newIP := strings.Join(parts, ":")

	// For IPv6 expressions, remove brackets
	if len(newIP) > 1 && string(newIP[0]) == "[" {
		newIP = newIP[1 : len(newIP)-1]
	}

	// Parse IP v4/v6
	parsedIP := net.ParseIP(newIP)

	// Compare againts IP ranges
	for _, cidr := range ranges {
		if cidr.Contains(parsedIP) {
			return true
		}
	}
	return false
}
