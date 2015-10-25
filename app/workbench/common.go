package workbench

import (
	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/etcinit/gonduit"
	"github.com/etcinit/phabulous/app/factories"
)

type gonduitCaller func(*gonduit.Conn) (interface{}, error)

func spewOrPanic(factory *factories.GonduitFactory, caller gonduitCaller) {
	client, err := factory.Make()
	if err != nil {
		panic(err)
	}

	res, err := caller(client)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(res)
}
