package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Define all import variables
var (
	Token             string
	ProxyChannel      chan string
	ProxyList         []string
	VanityChecks      int
	VanityList        []string
	Client            http.Client
	Config            ConfigYaml
	SocketChannel     chan *tls.Conn
	currentlysleeping int
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
			var start time.Time
			var elapsed time.Duration
		checkloop:
			for {
				// Measure start time
				if Config.Main.Debug == true {
					start = time.Now()
				}
				// Do vanity check
				res = VanityCheck(Vanity, defaultclient)
				// Calculate time taken for Vanity check
				if Config.Main.Debug == true {
					elapsed = time.Since(start)
				}

				// Check response code
				switch res {
				case 404: //Claim
					// Snipe depending on what mode is enabled
					if Config.Main.Debug == true {
						fmt.Println("vanity is free attempting claim...")
					}
					if Config.Main.SocketUsage == true {
						SnipeUsingSocket(Vanity, Config.Main.GuildID)
					} else {
						FastHttpClaim(Vanity, Config.Main.GuildID)
					}
				case 200:
					// Continue checking
					VanityChecks++
					if Config.Main.Debug == true {
						fmt.Println(fmt.Sprintf("%d vanity checks, took %s, current check: %s", VanityChecks, elapsed, Vanity))
					}
				case 429:
					// If ratelimited break out of loop...
					if Config.Main.Debug == true {
						fmt.Println("proxy ratelimited switching...")
					}
					break checkloop
				default:
					// Bad response mainly 0, meaning the proxy is bad or there is an error with the request
					if Config.Main.Debug == true {
						fmt.Println(fmt.Sprintf("got bad response with status code: %d", res))
					}
					errs++
					// Maximum amount of errors before proxy is deemed invalid
					if errs > 20 {
						if Config.Main.Debug == true {
							fmt.Println("exceeded error threshold")
						}
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
	currentlysleeping++
	time.Sleep(30 * time.Second)
	currentlysleeping--
	ProxyChannel <- proxy
}

func updater() {
	var prev int
	for {
		checkspersec := (VanityChecks - prev) / 5
		prev = VanityChecks
		Clear()

		LogInfo(fmt.Sprintf("total vanity checks %d, currently sleeping %d, checks per second %d", VanityChecks, currentlysleeping, checkspersec))
		time.Sleep(5 * time.Second)
	}
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

	ProxyChannel = make(chan string, len(ProxyList))

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

	// Add proxylist to channel
	for _, proxy := range ProxyList {
		ProxyChannel <- proxy
	}

	if Config.Main.Debug == false {
		go updater()

	}

	wg.Wait()
}
