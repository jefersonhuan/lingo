package transfer

import (
	"fmt"
	"lingo/utils"
	"time"
)

func (transfer *Transfer) Start() (finishedAt time.Time, err error) {
	fmt.Println(utils.ColorfulString("cyan", "Starting operations"))

	if err = utils.StepsFunctions(transfer.Source.Connect, transfer.Target.Connect); err != nil {
		return
	}

	for _, op := range []string{"Clone Databases", "Clone Collections", "Close Connections", "Finish"} {
		fmt.Println("-", op)
	}

	err = utils.StepsFunctions(transfer.clone)
	if err != nil {
		return
	}

	transfer.finish()

	return transfer.FinishedAt, nil
}

func (transfer *Transfer) finish() {
	transfer.Source.Disconnect()
	transfer.Target.Disconnect()

	transfer.FinishedAt = time.Now()
}
