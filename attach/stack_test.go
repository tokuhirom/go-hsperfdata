package attach

import (
	"testing"
)

func TestFoo(t *testing.T) {
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
	if err != nil {
		t.Error(err)
	}
	if len(threads) != 3 {
		t.Errorf("invalid count: %d", len(threads))
	}
	if threads[0].name != "main" {
		t.Errorf("bad name: %v", threads[0].name)
	}
	if threads[1].name != "GC task thread#1 (ParallelGC)" {
		t.Errorf("bad name: %v", threads[1].name)
	}
	if threads[2].name != "VM Periodic Task Thread" {
		t.Errorf("bad name: %v", threads[2].name)
	}
}
