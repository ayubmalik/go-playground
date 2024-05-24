package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/NordSecurity/nordvpn-linux/daemon/pb"
)

// Status returns ready to print status string.
func Status(resp *pb.StatusResponse) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Status: %s\n", resp.State))

	if resp.Hostname != "" {
		b.WriteString(fmt.Sprintf("Hostname: %s\n", resp.Hostname))
	}

	if resp.Ip != "" {
		b.WriteString(fmt.Sprintf("IP: %s\n", resp.Ip))
	}

	if resp.Country != "" {
		b.WriteString(fmt.Sprintf("Country: %s\n", resp.Country))
	}

	if resp.City != "" {
		b.WriteString(fmt.Sprintf("City: %s\n", resp.City))
	}

	if resp.Uptime != -1 {
		b.WriteString(
			fmt.Sprintf("Current technology: %s\n", resp.Technology.String()),
		)
		b.WriteString(
			fmt.Sprintf("Current protocol: %s\n", resp.Protocol.String()),
		)
	}

	// show transfer rates only if running
	if resp.Download != 0 || resp.Upload != 0 {
		b.WriteString(fmt.Sprintf("Transfer: %d received, %d sent\n", resp.Download, resp.Upload))
	}

	if resp.Uptime != -1 {
		// truncate to skip milliseconds from being displayed
		uptime := time.Duration(resp.Uptime).Truncate(1000 * time.Millisecond)
		up := uptime / 1000000000
		b.WriteString(fmt.Sprintf("Uptime: %d\n", up))
	}
	return b.String()
}
