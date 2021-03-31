package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	var port int
	var ddPort int
	var rpcPort int
	var rpcCall string
	var err error

	if len(os.Args) < 3 {
		if p, d := os.Getenv("BYOND_REST_PORT"), os.Getenv("BYOND_PORT"); p != "" && d != "" {
			port, err = strconv.Atoi(p)
			ddPort, err = strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
		} else {
			log.Fatal("no port supplied, aborting")
		}
	} else {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		ddPort, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	}

	// a 5 length arg list implies that the RPC function was included
	if len(os.Args) != 5 {
		if p := os.Getenv("BYOND_REST_RPC_PORT"); p != "" {
			rpcPort, err = strconv.Atoi(p)
			if err != nil {
				log.Fatal("an error occurred reading the bot port")
			}
			if c := os.Getenv("BYOND_REST_RPC_CALL"); c != "" {
				rpcCall = os.Getenv(c)
			} else {
				log.Println("warning: no RPC call supplied, but a RPC port was supplied")
			}
		} else {
			log.Println("warning: no RPC port supplied, services related to this will not recieve updates")
		}
	} else {
		rpcPort, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		if rpcPort == 0 {
			log.Println("warning: no RPC port supplied, services related to this will not recieve updates")
		}

		rpcCall = os.Args[4]
	}

	state := new(State)
	state.ddport = ddPort

	go listenDD(port, rpcPort, rpcCall, state)
	go serveJSON(state)

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}
