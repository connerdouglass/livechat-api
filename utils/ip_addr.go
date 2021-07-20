package utils

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

// newAddrPattern is the pattern for parsing the IP address out of a net.Addr. This is
// needed because the net.Addr includes a port number at the end
var netAddrPattern = regexp.MustCompile(`^(.*):\d+$`)

// GetIpAddress gets the IP address from a set of headers and a net address
func GetIpAddress(
	header http.Header,
	addr net.Addr,
) string {

	// If there are headers, try to pull the CF-Connecting-IP header, which is forwarded
	// from Cloudflare in the event that Cloudflare is being used.
	if header != nil {
		ip := header.Get("CF-Connecting-IP")
		if len(ip) > 0 {
			return ip
		}
	}

	// If the address is nil, return an empty string
	if addr == nil {
		return ""
	}

	// Match against the pattern in order to pull the IP address out of the address
	submatch := netAddrPattern.FindStringSubmatch(addr.String())
	if len(submatch) < 2 {
		return ""
	}

	// Clean up the IP address. These only have an effect in the case of IPv6 addresses
	ip := submatch[1]
	ip = strings.Trim(ip, "[]")
	ip = strings.TrimPrefix(ip, "::ffff:")
	return ip

}
