package attach

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFoo(t *testing.T) {
	assert := assert.New(t)
	threads, err := ParseStack(`

"main" #1 prio=6 os_prio=0 tid=0x00007f357c00a000 nid=0x45c6 runnable [0x00007f358442a000]
	java.lang.Thread.State: RUNNABLE
	at org.eclipse.swt.internal.gtk.OS.Call(Native Method)
	at org.eclipse.swt.widgets.Display.sleep(Display.java:4294)
	at org.eclipse.ui.application.WorkbenchAdvisor.eventLoopIdle(WorkbenchAdvisor.java:368)
	at org.eclipse.ui.internal.ide.application.IDEWorkbenchAdvisor.eventLoopIdle(IDEWorkbenchAdvisor.java:918)

"GC task thread#1 (ParallelGC)" os_prio=0 tid=0x00007f357c020800 nid=0x45c9 runnable

"VM Periodic Task Thread" os_prio=0 tid=0x00007f357c223000 nid=0x45d1 waiting on condition
`)
	assert.Nil(err)
	assert.Equal(len(threads), 3)

	assert.Equal(threads[0].name, "main")
	assert.Equal(threads[0].state, "RUNNABLE")
	assert.Equal(len(threads[0].stacks), 4)
	assert.Equal(threads[0].stacks[0].method, "org.eclipse.swt.internal.gtk.OS.Call")
	assert.Equal(threads[0].stacks[0].file, "Native Method")
	assert.Equal(threads[0].stacks[0].line, -1)
	assert.Equal(threads[0].stacks[1].method, "org.eclipse.swt.widgets.Display.sleep")
	assert.Equal(threads[0].stacks[1].file, "Display.java")
	assert.Equal(threads[0].stacks[1].line, 4294)
	assert.Equal(threads[0].stacks[2].method, "org.eclipse.ui.application.WorkbenchAdvisor.eventLoopIdle")
	assert.Equal(threads[0].stacks[2].file, "WorkbenchAdvisor.java")
	assert.Equal(threads[0].stacks[2].line, 368)
	assert.Equal(threads[0].stacks[3].method, "org.eclipse.ui.internal.ide.application.IDEWorkbenchAdvisor.eventLoopIdle")

	assert.Equal(threads[1].name, "GC task thread#1 (ParallelGC)")
	assert.Equal(threads[1].state, "")
	assert.Equal(len(threads[1].stacks), 0)

	assert.Equal(threads[2].name, "VM Periodic Task Thread")
	assert.Equal(len(threads[2].stacks), 0)
}
