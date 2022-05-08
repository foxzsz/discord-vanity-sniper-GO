package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// Reuse one fasthttp client for claiming
var (
	fastclient = &fasthttp.Client{
		//Dial:                fasthttpproxy.FasthttpHTTPDialer(proxy),
		//MaxConnsPerHost:     1000000,
		ReadBufferSize:      8192000,
		ReadTimeout:         time.Duration(1) * time.Second,
		MaxIdleConnDuration: time.Duration(1) * time.Second,
		// prevent accs from getting locked
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			CipherSuites:       []uint16{0x1301, 0x1303, 0x1302, 0xc02b, 0xc02f, 0xcca9, 0xcca8, 0xc02c, 0xc030, 0xc00a, 0xc009, 0xc013, 0xc014, 0x009c, 0x009d, 0x002f, 0x0035},
			InsecureSkipVerify: true,
			CurvePreferences:   []tls.CurveID{tls.CurveID(0x001d), tls.CurveID(0x0017), tls.CurveID(0x0018), tls.CurveID(0x0019), tls.CurveID(0x0100), tls.CurveID(0x0101)},
		},
	}
)

const UserAgent = "Discord/32114 CFNetwork/1331.0.7 Darwin/21.4.0"
const XTrack = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2Ojk3LjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTcuMCIsImJyb3dzZXJfdmVyc2lvbiI6Ijk3LjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTk5OSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0="
const XSuper = "eyJvcyI6ImlPUyIsImJyb3dzZXIiOiJEaXNjb3JkIGlPUyIsImRldmljZSI6ImlQYWQxMywxNiIsInN5c3RlbV9sb2NhbGUiOiJlbi1JTiIsImNsaWVudF92ZXJzaW9uIjoiMTI0LjAiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJkZXZpY2VfYWR2ZXJ0aXNlcl9pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImRldmljZV92ZW5kb3JfaWQiOiJBMTgzNkNFRC1BRDI5LTRGRTAtQjVDNC0zODQ0NDU0MEFFQTciLCJicm93c2VyX3VzZXJfYWdlbnQiOiIiLCJicm93c2VyX3ZlcnNpb24iOiIiLCJvc192ZXJzaW9uIjoiMTUuNC4xIiwiY2xpZW50X2J1aWxkX251bWJlciI6MzIyNDcsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9"

// Default vanity check function
func VanityCheck(Vanity string, Client *http.Client) (status int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v9/invites/%s", Vanity), nil)
	if err != nil {
		LogErr(fmt.Sprintf("error while doing req: %s", err))
		return 0
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := Client.Do(req)
	if err != nil {
		LogErr(fmt.Sprintf("error while doing req: %s", err))
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

// Dial discord and return *tls.Conn
func CreateSocketChannel() *tls.Conn {
	conns, err := tls.Dial("tcp", "discord.com"+":443", nil)
	if err != nil {
		fmt.Println("error opening connection")
	}
	return conns

}

// Claim using fasthttp ( faster and better than sockets)
func FastHttpClaim(vanity string, guildid int) {
	start := time.Now()
	req := fasthttp.AcquireRequest()
	payload := []byte(fmt.Sprintf(`{"code": "%s"}`, vanity))
	req.SetBody(payload)
	req.Header.SetMethod("PATCH")
	// Headers
	for k, v := range map[string]string{
		"Host":               "discord.com",
		"User-Agent":         UserAgent,
		"Accept":             "*/*",
		"Accept-Language":    "en-US,en;q=0.5",
		"Content-Type":       "application/json",
		"Authorization":      Config.Main.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
	} {
		req.Header.Set(k, v)
	}

	req.SetRequestURI(fmt.Sprintf("https://discord.com/api/v9/guilds/%d/vanity-url", guildid))
	res := fasthttp.AcquireResponse()
	if err := fastclient.Do(req, res); err != nil {
		LogErr(fmt.Sprintf("error while claiming %s", err))
		SendFail(vanity, "", "0")
	}
	fasthttp.ReleaseRequest(req)
	elapsed := fmt.Sprint(time.Since(start))
	statuscode := fmt.Sprint(res.StatusCode())
	fasthttp.ReleaseResponse(res)
	if statuscode == "200" {
		SendSuccess(vanity, elapsed)

	} else if statuscode == "429" {
		SendRatelimit(vanity, elapsed)

	} else {
		SendFail(vanity, elapsed, statuscode)
	}

}

// Attempt to claim using sockets
func SnipeUsingSocket(vanity string, guildid int) {
	// Measure time
	start := time.Now()
	// Get socket from channel
	conn := <-SocketChannel
	// Create payload
	payload := fmt.Sprintf(`{"code": "%s"}`, vanity)
	// Create TCP request data
	data := fmt.Sprintf("PATCH /api/v9/guilds/%d/vanity-url HTTP/1.1\r\nContent-Type: application/json\r\nHost: discord.com\r\nAuthorization: %s\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0\r\nContent-Length: "+strconv.Itoa(len(payload))+"\r\n\r\n"+payload, guildid, Config.Main.Token)
	// Create bytes which will be used to recieve data
	authbytes := make([]byte, 4096)
	// Write and read response from request
	_, err := conn.Write([]byte(data))
	if err != nil {
		SendFail(vanity, "", "0")
	}
	_, err = conn.Read(authbytes)
	if err != nil {
		SendFail(vanity, "", "0")
	}
	// Calculate time taken
	elapsed := fmt.Sprint(time.Since(start))
	// Check the response ( and send webhook)
	CheckResponse(authbytes, vanity, elapsed)
	// Send Tcp socket back to channel
	SocketChannel <- conn
}

// Check the response for socket claims
func CheckResponse(Response []byte, vanity string, elapsed string) {
	statuscode := string(Response[9:12])
	if statuscode == "200" {
		SendSuccess(vanity, elapsed)

	} else if statuscode == "429" {
		SendRatelimit(vanity, elapsed)

	} else {
		SendFail(vanity, elapsed, statuscode)
	}
}
