package answers

import (
	"fmt"
	"github.com/levigross/grequests"
	"github.com/ravlio/highloadcup2018"
	"github.com/ravlio/hlcuptester"
	"log"
	"net"
	//"regexp"
	"github.com/stretchr/testify/assert"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"testing"
	"time"
	// "regexp"
	"os"
	// "runtime/pprof"
)

// import _ "github.com/pkg/profile"

var (
	db    *hl.DB
	srv   *hl.HTTPServer
	dpath string
	host  string
	port  string
)

func TestMain(m *testing.M) {
	/*f, err := os.Create("cpu.prof")
	  if err != nil {
	      log.Fatal(err)
	  }
	  pprof.StartCPUProfile(f)
	  defer pprof.StopCPUProfile()
	*/
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// defer profile.Start(profile.MemProfile).Stop()
	host = "localhost"
	port = "8181"
	// dpath="../../../data/"
	//dpath="/Users/ravlio/Downloads/elim_accounts_261218/"
	dpath = "/Users/maksim.bogdanov/Downloads/elim_accounts_261218/"
	//dpath="/Users/maksim.bogdanov/Downloads/test_accounts_210119/"
	//dpath="/tmp/data/"

	db = hl.NewDB(&hl.Options{Indexer: hl.OptionsIndexer{Enabled: true, WorkerCount: 4, BufferSize: 1}, Debug: true})
	err := db.Start()
	if err != nil {
		log.Fatal("error starting database")
	}

	err = db.LoadFromJSON(dpath + "data/")
	if err != nil {
		log.Fatalf("error loading json: %s", err.Error())
	}

	srv = &hl.HTTPServer{DB: db, Addr: net.JoinHostPort(host, port)}

	go func() {
		err := srv.Run()
		if err != nil {
			log.Fatalf("error running server: %v", err)
		}
	}()

	for {
		con, err := net.Dial("tcp", "localhost:"+port)
		if err == nil {
			con.Close()
			break
		}
		time.Sleep(time.Millisecond * 5)
	}

	defer func() {
		err := db.Stop()
		if err != nil {
			log.Fatal(err)
		}
		err = srv.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	os.Exit(m.Run())
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

func getPhase(t *testing.T, phase int) {
	ch, err := hlcuptester.Load(dpath, phase,
		//hlcuptester.PathRegexp(regexp.MustCompile(`/(\_contains|\_gt|\_lt)/`)),
		//hlcuptester.InsidePath("accounts/filter"),
		//hlcuptester.InsidePath("accounts/group"),
		hlcuptester.F(func(u string) bool {
			/*if strings.Contains(u, "accounts/group") {
				return true
			}

			return false*/
			if strings.Contains(u, "accounts/filter") {
				if strings.Contains(u, "_contains") || strings.Contains(u, "_gt") || strings.Contains(u, "_lt") {
					return false
				}
				return true
			}
			return false
		}),
	)

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	tim := time.Now()
	id := 0
	for rr := range ch {
		// fmt.Printf("URI: %s\n",rr.URI)

		id++

		if !assert.NoError(t, rr.Err) {
			t.FailNow()
		}

		uri := fmt.Sprintf("http://%s:%s%s", host, port, rr.URI)
		resp, err := grequests.Get(uri, nil)
		if !assert.NoError(t, err) {
			printResponseErr(resp, uri, id)
			t.FailNow()
		}
		if !assert.Equal(t, rr.ResponseStatus, resp.StatusCode) {
			printResponseErr(resp, uri, id)
			t.FailNow()
		}

		if rr.ResponseStatus == 400 || len(rr.ResponseBody) == 0 {
			if !assert.Empty(t, resp.String()) {
				t.FailNow()
			}
		} else {
			if !assert.JSONEq(t, rr.ResponseBody, resp.String()) {
				t.Logf("EXPECT JSON: %s\n", rr.ResponseBody)
				t.Logf("ACTUAL JSON: %s\n", resp.String())
				printResponseErr(resp, uri, id)
				t.FailNow()
			}
		}

		//fmt.Printf("URI: %s\nRequestBody: %s\nResponseCode: %d\nResponseBody:%s\n\n", rr.URI, rr.RequestBody, rr.ResponseStatus, rr.ResponseBody)
		//exp.GET(rr.URI).Expect().Status(rr.ResponseStatus).ContentType("application/json").Body().Equal(rr.ResponseBody)
		//break
	}
	println("ids", id, "dur", time.Since(tim).String())

}

func TestAll(t *testing.T) {
	t.Run("Phase #1 (GET)", func(t *testing.T) {
		getPhase(t, 1)
	})

	if t.Failed() {
		return
	}

	t.Run("Phase #2 (POST", func(t *testing.T) {
		println(1)
		tim := time.Now()
		ch, err := hlcuptester.Load(
			dpath,
			2,
			/*hlcuptester.InsidePath("accounts/new"),
			hlcuptester.PathRegexp(regexp.MustCompile(`/accounts/\d+/`)),
			hlcuptester.InsidePath("accounts/likes"),*/
		)

		if !assert.NoError(t, err) {
			t.FailNow()
		}

		id := 0
		for rr := range ch {
			// fmt.Printf("URI: %s\n",rr.URI)

			id++

			if !assert.NoError(t, rr.Err) {
				t.FailNow()
			}

			uri := fmt.Sprintf("http://%s:%s%s", host, port, rr.URI)
			ro := &grequests.RequestOptions{
				JSON: rr.RequestBody,
			}

			resp, err := grequests.Post(uri, ro)

			if !assert.NoError(t, err) {
				printPostResponseErr(resp, uri, id, rr.RequestBody)
				t.FailNow()
			}
			if !assert.Equal(t, rr.ResponseStatus, resp.StatusCode) {
				printPostResponseErr(resp, uri, id, rr.RequestBody)
				t.FailNow()
			}

			if !assert.Equal(t, rr.ResponseBody, resp.String()) {
				printPostResponseErr(resp, uri, id, rr.RequestBody)
				t.FailNow()
			}
		}
		println("ids", id, "dur", time.Since(tim).String())
	})

	if t.Failed() {
		return
	}

	time.Sleep(time.Second * 5)
	t.Run("Phase #3 (GET)", func(t *testing.T) {
		getPhase(t, 3)
	})

	time.Sleep(time.Second * 5)

}
