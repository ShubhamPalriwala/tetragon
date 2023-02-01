package tracing

import (
	"context"
	"os"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/cilium/tetragon/pkg/observer"
	"github.com/cilium/tetragon/pkg/sensors"
	tus "github.com/cilium/tetragon/pkg/testutils/sensors"
)

func TestUprobeLoad(t *testing.T) {
	var sensorProgs = []tus.SensorProg{
		// uprobe
		0:  tus.SensorProg{Name: "generic_uprobe_event", Type: ebpf.Kprobe},
		1:  tus.SensorProg{Name: "generic_uprobe_process_event0", Type: ebpf.Kprobe},
		2:  tus.SensorProg{Name: "generic_uprobe_process_event1", Type: ebpf.Kprobe},
		3:  tus.SensorProg{Name: "generic_uprobe_process_event2", Type: ebpf.Kprobe},
		4:  tus.SensorProg{Name: "generic_uprobe_process_event3", Type: ebpf.Kprobe},
		5:  tus.SensorProg{Name: "generic_uprobe_process_event4", Type: ebpf.Kprobe},
		6:  tus.SensorProg{Name: "generic_uprobe_filter_arg1", Type: ebpf.Kprobe},
		7:  tus.SensorProg{Name: "generic_uprobe_filter_arg2", Type: ebpf.Kprobe},
		8:  tus.SensorProg{Name: "generic_uprobe_filter_arg3", Type: ebpf.Kprobe},
		9:  tus.SensorProg{Name: "generic_uprobe_filter_arg4", Type: ebpf.Kprobe},
		10: tus.SensorProg{Name: "generic_uprobe_filter_arg5", Type: ebpf.Kprobe},
		11: tus.SensorProg{Name: "generic_uprobe_process_filter", Type: ebpf.Kprobe},
	}

	var sensorMaps = []tus.SensorMap{
		// all uprobe programs
		tus.SensorMap{Name: "process_call_heap", Progs: []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
		tus.SensorMap{Name: "uprobe_calls", Progs: []uint{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},

		// generic_uprobe_process_filter,generic_uprobe_filter_arg*
		tus.SensorMap{Name: "filter_map", Progs: []uint{6, 7, 8, 9, 10, 11}},

		// generic_uprobe_filter_arg*,generic_retuprobe_event,base
		tus.SensorMap{Name: "tcpmon_map", Progs: []uint{6, 7, 8, 9, 10, 12}},

		// shared with base sensor
		tus.SensorMap{Name: "execve_map", Progs: []uint{6, 7, 8, 9, 10, 11, 12}},
	}

	nopHook := `
apiVersion: cilium.io/v1alpha1
metadata:
  name: "uprobe"
spec:
  uprobes:
  - path: "/bin/bash"
    symbol: "main"
`

	var sens []*sensors.Sensor
	var err error

	nopConfigHook := []byte(nopHook)
	err = os.WriteFile(testConfigFile, nopConfigHook, 0644)
	if err != nil {
		t.Fatalf("writeFile(%s): err %s", testConfigFile, err)
	}
	sens, err = observer.GetDefaultSensorsWithFile(t, context.TODO(), testConfigFile, tus.Conf().TetragonLib)
	if err != nil {
		t.Fatalf("GetDefaultObserverWithFile error: %s", err)
	}

	tus.CheckSensorLoad(sens, sensorMaps, sensorProgs, t)

	sensors.UnloadAll(tus.Conf().TetragonLib)
}