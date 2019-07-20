package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ravlio/hlcuptester"
	"log"
	"net"
	"os"
	"reflect"

	"strings"
	"time"

	_ "net/http/pprof"
)

// import _ "github.com/pkg/profile"

var host = flag.String("host", "127.0.0.1", "hostname")
var port = flag.String("port", "8181", "port")
var dpath = flag.String("dpath", "/Users/maksim.bogdanov/Downloads/elim_accounts_261218/", "data path")
var phase = flag.Int("phase", 0, "phase number")

func main() {
	flag.Parse()
	for {
		con, err := net.Dial("tcp", net.JoinHostPort(*host, *port))
		if err == nil {
			con.Close()
			break
		}
		time.Sleep(time.Millisecond * 5)
	}

	if *phase > 0 {
		if *phase == 1 || *phase == 3 {
			getPhase(*phase)
		} else if *phase == 2 {
			postPhase()
		}
		return
	}

	getPhase(1)
	postPhase()
	getPhase(3)

}

func equalJSON(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func printResponseErr(resp *grequests.Response, uri string, id int) {
	fmt.Printf("URI [#%d]: %s\n", id, uri)
	if e := resp.Header.Get("X-Error"); e != "" {
		fmt.Printf("error: %s\n", e)
	}

	if t := resp.Header.Get("X-Error-Trace"); t != "" {
		fmt.Printf("trace: %s\n", t)
	}
}

func printPostResponseErr(resp *grequests.Response, uri string, id int, body string) {
	fmt.Printf("URI [#%d]: %s\n", id, uri)
	fmt.Printf("POST Body: \n%s\n", body)
	if e := resp.Header.Get("X-Error"); e != "" {
		fmt.Printf("error: %s\n", e)
	}

	if t := resp.Header.Get("X-Error-Trace"); t != "" {
		fmt.Printf("trace: %s\n", t)
	}
}

func postPhase() {
	tim := time.Now()
	ch, err := hlcuptester.Load(
		*dpath,
		2,
		/*hlcuptester.InsidePath("accounts/new"),
		hlcuptester.PathRegexp(regexp.MustCompile(`/accounts/\d+/`)),
		hlcuptester.InsidePath("accounts/likes"),*/
	)

	if err != nil {
		log.Fatal(err)
	}

	id := 0
	for rr := range ch {
		// fmt.Printf("URI: %s\n",rr.URI)

		id++

		if rr.Err != nil {
			log.Fatal(err)
		}

		uri := fmt.Sprintf("http://%s:%s%s", *host, *port, rr.URI)
		ro := &grequests.RequestOptions{
			JSON: rr.RequestBody,
		}

		resp, err := grequests.Post(uri, ro)

		if err != nil {
			printPostResponseErr(resp, uri, id, rr.RequestBody)
			log.Fatal(err)
		}
		if rr.ResponseStatus != resp.StatusCode {
			printPostResponseErr(resp, uri, id, rr.RequestBody)
			log.Fatal(err)
		}

		if rr.ResponseBody != resp.String() {
			printPostResponseErr(resp, uri, id, rr.RequestBody)
			log.Fatal(err)
		}
	}
	println("ids", id, "dur", time.Since(tim).String())
}

func getPhase(phase int) {
	ch, err := hlcuptester.Load(*dpath, phase,
		//hlcuptester.PathRegexp(regexp.MustCompile(`/(\_contains|\_gt|\_lt)/`)),
		//hlcuptester.InsidePath("accounts/filter"),
		//hlcuptester.InsidePath("accounts/group"),
		hlcuptester.F(func(u string) bool {
			if strings.Contains(u, "accounts/group") {
				return true
			}

			//return false
			if strings.Contains(u, "accounts/filter") {
				if strings.Contains(u, "_contains") || strings.Contains(u, "_gt") || strings.Contains(u, "_lt") {
					return false
				}
				return true
			}
			return false
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	tim := time.Now()
	id := 0
	for rr := range ch {
		// fmt.Printf("URI: %s\n",rr.URI)

		id++

		if rr.Err != nil {
			log.Fatal(err)
		}

		uri := fmt.Sprintf("http://%s:%s%s", *host, *port, rr.URI)
		resp, err := grequests.Get(uri, &grequests.RequestOptions{DialKeepAlive: time.Second * 30})
		if err != nil {
			printResponseErr(resp, uri, id)

			log.Fatal(err)
		}
		if rr.ResponseStatus != resp.StatusCode {
			printResponseErr(resp, uri, id)
			log.Fatal(err)
		}

		if rr.ResponseStatus == 400 || len(rr.ResponseBody) == 0 {
			if resp.String() != "" {
				log.Fatal(err)
			}
		} else {
			if !equalJSON(rr.ResponseBody, resp.String()) {
				log.Printf("EXPECT JSON: %s\n", rr.ResponseBody)
				log.Printf("ACTUAL JSON: %s\n", resp.String())
				printResponseErr(resp, uri, id)
				os.Exit(1)
			}
		}

	}
	println("ids", id, "dur", time.Since(tim).String())

}
