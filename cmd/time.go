package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(timeCmd)
}

var timeCmd = &cobra.Command{
	Use:   "time [url]",
	Short: "Make and time a GET request",
	Long:  "Make and time a GET request to a given url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := makeMeasuredGetRequest(args[0])

		if err != nil {
			panic(err)
		}

		printMeasuredResponse(m)
	},
}

func printMeasuredResponse(m *measuredResponse) {
	lineBreakStyle := lipgloss.NewStyle().Width(80).Faint(true)
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Faint(true)

	durationPair := func(title string, value time.Duration) string {
		return fmt.Sprintf("%s: %s", titleStyle.Render(title), normalStyle.Render(value.String()))
	}

	textPair := func(title string, value string) string {
		return fmt.Sprintf("%s: %s", titleStyle.Render(title), normalStyle.Render(value))
	}

	lineBreak := func() {
		fmt.Println(lineBreakStyle.Render("--------------------------------------------------"))
	}

	lineBreak()
	fmt.Println(textPair("URL", m.Res.Request.URL.String()))
	fmt.Println(textPair("Response status", m.Res.Status))
	lineBreak()
	fmt.Println(durationPair("DNS Lookup", m.DNS))
	fmt.Println(durationPair("TCP Connection", m.Connect))
	fmt.Println(durationPair("TLS Handshake", m.TLS))
	fmt.Println(durationPair("Connect", m.TLS))
	fmt.Println(durationPair("First Byte", m.FirstByte))
	lineBreak()
	fmt.Println(durationPair("Name Lookup", m.FirstByte-m.TLS-m.Connect))
	fmt.Println(durationPair("Server Processing", m.FirstByte-m.TLS-m.Connect-m.DNS))
	lineBreak()
	fmt.Println(durationPair("Total", time.Since(m.Start)))
}

type measuredResponse struct {
	Res       *http.Response
	Start     time.Time
	DNS       time.Duration
	Connect   time.Duration
	TLS       time.Duration
	FirstByte time.Duration
}

func makeMeasuredGetRequest(url string) (*measuredResponse, error) {
	req, _ := http.NewRequest("GET", url, nil)

	measured := measuredResponse{}
	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			measured.DNS = time.Since(dns)
		},
		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			measured.Connect = time.Since(connect)
		},
		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			measured.TLS = time.Since(tlsHandshake)
		},
		GotFirstResponseByte: func() {
			measured.FirstByte = time.Since(start)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	measured.Start = start
	measured.Res = resp
	defer resp.Body.Close()

	return &measured, err
}
