package hello

import "sync"

type Inst struct {
	GreetSvc *greetService
}

var ins *Inst
var insOnce sync.Once

func GetInst() *Inst {
	insOnce.Do(func() {
		ins = &Inst{
			GreetSvc: newGreetService(),
		}
	})
	return ins
}
