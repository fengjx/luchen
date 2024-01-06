package luchen

import (
	"net"
	"net/http"
	"strings"
)

func getClientIP(req *http.Request) string {
	forwardIPs := getHeader(req, "X-Forwarded-For")
	if forwardIPs == "" {
		return ""
	}
	ips := strings.Split(forwardIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "unknown" || ip == "unknow" {
			continue
		}
		parseIP := net.ParseIP(ip)
		if parseIP.IsLoopback() || parseIP.IsPrivate() {
			continue
		}
		return ip
	}
	return ""
}

func getHeader(req *http.Request, key string) string {
	vals := req.Header[key]
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
