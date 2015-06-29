package attach

import (
	"encoding/json"
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

	var expected []*JavaThread
	expectedJSON := `
    [
      {
        "name": "main",
        "state": "RUNNABLE",
        "Stack": [
          {
            "method": "org.eclipse.swt.internal.gtk.OS.Call",
            "file": "Native Method",
            "line": -1
          },
          {
            "method": "org.eclipse.swt.widgets.Display.sleep",
            "file": "Display.java",
            "line": 4294
          },
          {
            "method": "org.eclipse.ui.application.WorkbenchAdvisor.eventLoopIdle",
            "file": "WorkbenchAdvisor.java",
            "line": 368
          },
          {
            "method": "org.eclipse.ui.internal.ide.application.IDEWorkbenchAdvisor.eventLoopIdle",
            "file": "IDEWorkbenchAdvisor.java",
            "line": 918
          }
        ]
      },
      {
        "name": "GC task thread#1 (ParallelGC)",
        "state": "",
        "Stack": []
      },
      {
        "name": "VM Periodic Task Thread",
        "state": "",
        "Stack": []
      }
    ]
	`
	err = json.Unmarshal([]byte(expectedJSON), &expected)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(threads, expected)
}
