package hello

import "sync"

type Inst struct {
	helloLogic *helloLogic
	GreetSvc   *greetService
	Endpoints  *endpoints
}

var ins *Inst
var insOnce sync.Once

func GetInst() *Inst {
	insOnce.Do(func() {
		ins = &Inst{
			helloLogic: newHelloLogic(),
			GreetSvc:   newGreetService(),
			Endpoints:  newEndpoints(),
		}
	})
	return ins
}
