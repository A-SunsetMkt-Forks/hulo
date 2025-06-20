package interpreter

type Debugger struct {
	interpreter *Interpreter
	breakpoints map[string]map[int]Breakpoint // file -> line -> breakpoint
	pauseChan   chan struct{}
	resumeChan  chan struct{}
}

type DebugState struct{}

func (d *Debugger) pause() {
	d.notifyClient(&DebugState{})

	<-d.resumeChan
}

func (d *Debugger) notifyClient(*DebugState) {}

type DebugCommand struct{}

func (d *Debugger) handleCommand(cmd DebugCommand) {}

func (d *Debugger) stepInto()        {}
func (d *Debugger) stepOver()        {}
func (d *Debugger) stepOut()         {}
func (d *Debugger) continue_()       {}
func (d *Debugger) inspectVariable() {}

type Breakpoint struct {
	File    string
	Line    int
	Column  int
	Enabled bool
}
