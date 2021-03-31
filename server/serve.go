package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"

	"github.com/vulppine/byond-topic-go"
)

// State represents the current state of the server.
// It contains a mutex, the raw JSON representing the
// state, and the current status of the server, as
// well as the port that Dream Daemon is being hosted on.
//
// If the status changes from the last status that
// was in the State when a new state was recieved,
// the Dream Daemon listener will automatically
// call Update() over the given RPC port.
type State struct {
	m      sync.Mutex
	raw    []byte
	ddport int
	Status int `json:"status"`
}

func listenDD(port, rport int, rcall string, state *State) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	log.Printf("starting listener on port %d\n", port)

	var r *rpc.Client
	if rport != 0 {
		r, err = jsonrpc.Dial("tcp", fmt.Sprintf(":%d", rport))
		if err != nil {
			log.Println("error opening rpc port:", err)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go func(conn net.Conn) {
				state.m.Lock()
				// log.Println("connection recieved")

				s := state.Status

				// log.Println("reading sent state")
				b, err := io.ReadAll(conn)
				if err != nil {
					log.Println(err)
				} else {
					if len(b) != 0 {
						log.Println("updating state")
						state.raw = b
						json.Unmarshal(state.raw, state)
					}

					// log.Println(state)
				}

				if state.Status != s && r != nil {
					// log.Println("status changed, attempting to call Bot.StatusChange over rpc port")
					err = r.Call(rcall, string(state.raw), nil)
					if err == rpc.ErrShutdown {
						log.Println("attempting to open connection to rpc server again and retrying")
						r, err = jsonrpc.Dial("tcp", fmt.Sprintf(":%d", rport))

						if err == nil {
							r.Call(rcall, string(state.raw), nil)
						} else {
							log.Println("error opening rpc port:", err)
						}
					} else if err != nil {
						log.Println(err)
					}
				}

				// log.Println("attempting to unlock mutex")
				state.m.Unlock()
			}(c)
		}
	}
}

// you will want to route from the front end
// to at the very least, the REST endpoint /api/status
// from below for remote access
//
// enjoy the funny port number
func serveJSON(state *State) {
	if err := http.ListenAndServe(":3621", state); err != nil {
		log.Fatal(err)
	}
}

func (s *State) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.URL.Path != "/api/status" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// hardcoded
	if u, err := byondtopic.SendTopic(fmt.Sprintf(":%d", s.ddport), "update_rest"); u != "SUCCESS" && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(s.raw)))
	w.Write(s.raw)
}
