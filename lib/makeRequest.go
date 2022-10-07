package lib

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

func makeMeasuredRequest(id int, url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("[%d] DNS lookup: %s\n", id, time.Since(dns))
		},
		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("[%d] Connect: %s\n", id, time.Since(connect))
		},
		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("[%d] TLS handshake: %s\n", id, time.Since(tlsHandshake))
		},
		GotFirstResponseByte: func() {
			fmt.Printf("[%d] First byte: %s\n", id, time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Total time: %v\n", time.Since(start))

	return resp, err
}
