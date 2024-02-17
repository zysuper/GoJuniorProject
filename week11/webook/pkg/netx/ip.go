package netx

import "net"

// GetOutboundIP 获得对外发送消息的 IP 地址
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "114.114.114.114:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
