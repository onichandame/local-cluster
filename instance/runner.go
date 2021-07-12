package instance

import (
	"context"
	"os/exec"
	"sync"
)

type Runner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

var RunnersMap = map[uint]*Runner{}
var RMLock sync.Mutex
