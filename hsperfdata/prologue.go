package hsperfdata

// see sun/management/counter/perf/Prologue.java
type Prologue struct {
	Accessible   byte
	Used         int32
	Overflow     int32
	ModTimestamp int64
	EntryOffset  int32
	NumEntries   int32
}
