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
		TLSConfig:           &tls.Config{InsecureSkipVerify: true},
	}
)

// Default vanity check function
func VanityCheck(Vanity string, Client *http.Client) (status int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v9/invites/%s", Vanity), nil)
	if err != nil {
		LogErr(fmt.Sprintf("error while doing req: %s", err))
		return 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0")
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
func FastHttpTest(vanity string, guildid int) {
	start := time.Now()
	req := fasthttp.AcquireRequest()
	payload := []byte(fmt.Sprintf(`{"code": "%s"}`, vanity))
	req.SetBody(payload)
	req.Header.SetMethod("PATCH")
	req.Header.Set("Authorization", Config.Main.Token)
	req.Header.SetContentType("application/json")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(fmt.Sprintf("https://discord.com/api/v9/guilds/%d/vanity-url", guildid))
	res := fasthttp.AcquireResponse()
	if err := fastclient.Do(req, res); err != nil {
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
