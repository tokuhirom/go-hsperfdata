package hsperfdata

type DataEntryHeader2 struct {
	EntryLength  int32
	NameOffset   int32
	VectorLength int32
	DataType     byte
	Flags        byte
	DataUnits    byte
	DataVar      byte
	DataOffset   int32
}
