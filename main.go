package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"sync"
	"time"
)

// Define all import variables
var (
	Token         string
	ProxyChannel  chan string
	ProxyList     []string
	VanityChecks  int
	VanityList    []string
	Client        http.Client
	Config        ConfigYaml
	SocketChannel chan *tls.Conn
)

func SnipingThread(Vanity string) {

	// Create clients and transport
	defaultclient := &http.Client{Timeout: 2 * time.Second}
	// Create http.Transport
	for {
		func() {
			// Get proxy from channel
			proxy := <-ProxyChannel
			proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxy))
			if err != nil {
				LogErr("invalid proxy type")
			}

			transport := http.Transport{Proxy: http.ProxyURL(proxyURL), DisableKeepAlives: false}
			defaultclient.Transport = &transport
			// Define errors and response
			var errs int
			var res int
		checkloop:
			for {
				// Measure start time
				start := time.Now()
				// Do vanity check
				res = VanityCheck(Vanity, defaultclient)
				// Calculate time taken for Vanity check
				elapsed := time.Since(start)
				// Check response code
				switch res {
				case 404: //Claim
					// Snipe depending on what mode is enabled
					fmt.Println("vanity is free attempting claim...")
					if Config.Main.SocketUsage == true {
						SnipeUsingSocket(Vanity, Config.Main.GuildID)
					} else {
						FastHttpTest(Vanity, Config.Main.GuildID)
					}
				case 200:
					// Continue checking
					VanityChecks++
					fmt.Println(fmt.Sprintf("%d vanity checks, took %s, current check: %s", VanityChecks, elapsed, Vanity))
				case 429:
					// If ratelimited break out of loop...
					fmt.Println("proxy ratelimited switching...")
					break checkloop
				default:
					// Bad response mainly 0, meaning the proxy is bad or there is an error with the request
					fmt.Println(fmt.Sprintf("got bad response with status code: %d", res))
					errs++
					// Maximum amount of errors before proxy is deemed invalid
					if errs > 20 {
						fmt.Println("exceeded error threshold")
						break checkloop
					}

				}
			}
			// Close Idle connections to prevent any memory from leaking
			transport.CloseIdleConnections()
			// Open a new goroutine which will wait for 30 seconds and then add back to channel (default discord ratelimit)
			go sleeper(proxy)

		}()
	}
}

func sleeper(proxy string) {
	// Sleep....
	time.Sleep(30 * time.Second)
	ProxyChannel <- proxy
}

// Create socket channels for claiming
func SetupSocketChannels() {
	for i := 0; i < Config.Main.SocketChannels; i++ {
		socket := CreateSocketChannel()
		SocketChannel <- socket
	}
	LogInfo(fmt.Sprintf("successfully setup %d socket channels", Config.Main.SocketChannels))

}

func main() {

	wg := &sync.WaitGroup{}

	// Setup sniper..
	SetupSniper()

	// Start threads
	for _, vanity := range VanityList {
		// Initiate threads
		for i := 0; i < Config.Main.Amplify; i++ {
			wg.Add(1)
			go SnipingThread(vanity)
		}

	}

	// Setup sockets if enabled
	if Config.Main.SocketUsage == true {
		SocketChannel = make(chan *tls.Conn, Config.Main.SocketChannels)
		go SetupSocketChannels()

	}

	ProxyChannel = make(chan string, len(ProxyList))

	// Add proxylist to channel
	for _, proxy := range ProxyList {
		ProxyChannel <- proxy
	}

	wg.Wait()
}
