package tasks

// Code generated by go-ipld-prime gengo.  DO NOT EDIT.

import (
	ipld "github.com/ipld/go-ipld-prime"
)
var _ ipld.Node = nil // suppress errors when this dependency is not referenced
// Type is a struct embeding a NodePrototype/Type for every Node implementation in this package.
// One of its major uses is to start the construction of a value.
// You can use it like this:
//
// 		tasks.Type.YourTypeName.NewBuilder().BeginMap() //...
//
// and:
//
// 		tasks.Type.OtherTypeName.NewBuilder().AssignString("x") // ...
//
var Type typeSlab

type typeSlab struct {
	Any       _Any__Prototype
	Any__Repr _Any__ReprPrototype
	AuthenticatedRecord       _AuthenticatedRecord__Prototype
	AuthenticatedRecord__Repr _AuthenticatedRecord__ReprPrototype
	Bool       _Bool__Prototype
	Bool__Repr _Bool__ReprPrototype
	Bytes       _Bytes__Prototype
	Bytes__Repr _Bytes__ReprPrototype
	FinishedTask       _FinishedTask__Prototype
	FinishedTask__Repr _FinishedTask__ReprPrototype
	Float       _Float__Prototype
	Float__Repr _Float__ReprPrototype
	Int       _Int__Prototype
	Int__Repr _Int__ReprPrototype
	Link       _Link__Prototype
	Link__Repr _Link__ReprPrototype
	Link_FinishedTask       _Link_FinishedTask__Prototype
	Link_FinishedTask__Repr _Link_FinishedTask__ReprPrototype
	Link_List_StageDetails       _Link_List_StageDetails__Prototype
	Link_List_StageDetails__Repr _Link_List_StageDetails__ReprPrototype
	List       _List__Prototype
	List__Repr _List__ReprPrototype
	List_AuthenticatedRecord       _List_AuthenticatedRecord__Prototype
	List_AuthenticatedRecord__Repr _List_AuthenticatedRecord__ReprPrototype
	List_Logs       _List_Logs__Prototype
	List_Logs__Repr _List_Logs__ReprPrototype
	List_StageDetails       _List_StageDetails__Prototype
	List_StageDetails__Repr _List_StageDetails__ReprPrototype
	Logs       _Logs__Prototype
	Logs__Repr _Logs__ReprPrototype
	Map       _Map__Prototype
	Map__Repr _Map__ReprPrototype
	PopTask       _PopTask__Prototype
	PopTask__Repr _PopTask__ReprPrototype
	RecordUpdate       _RecordUpdate__Prototype
	RecordUpdate__Repr _RecordUpdate__ReprPrototype
	RetrievalTask       _RetrievalTask__Prototype
	RetrievalTask__Repr _RetrievalTask__ReprPrototype
	StageDetails       _StageDetails__Prototype
	StageDetails__Repr _StageDetails__ReprPrototype
	Status       _Status__Prototype
	Status__Repr _Status__ReprPrototype
	StorageTask       _StorageTask__Prototype
	StorageTask__Repr _StorageTask__ReprPrototype
	String       _String__Prototype
	String__Repr _String__ReprPrototype
	Task       _Task__Prototype
	Task__Repr _Task__ReprPrototype
	Tasks       _Tasks__Prototype
	Tasks__Repr _Tasks__ReprPrototype
	Time       _Time__Prototype
	Time__Repr _Time__ReprPrototype
	UpdateTask       _UpdateTask__Prototype
	UpdateTask__Repr _UpdateTask__ReprPrototype
}

// --- type definitions follow ---

// Any matches the IPLD Schema type "Any".  It has Union type-kind, and may be interrogated like map kind.
type Any = *_Any
type _Any struct {
	x _Any__iface
}
type _Any__iface interface {
	_Any__member()
}
func (_Bool) _Any__member() {}
func (_Int) _Any__member() {}
func (_Float) _Any__member() {}
func (_String) _Any__member() {}
func (_Bytes) _Any__member() {}
func (_Map) _Any__member() {}
func (_List) _Any__member() {}
func (_Link) _Any__member() {}

// AuthenticatedRecord matches the IPLD Schema type "AuthenticatedRecord".  It has Struct type-kind, and may be interrogated like map kind.
type AuthenticatedRecord = *_AuthenticatedRecord
type _AuthenticatedRecord struct {
	Record _Link_FinishedTask
	Signature _Bytes
}

// Bool matches the IPLD Schema type "Bool".  It has bool kind.
type Bool = *_Bool
type _Bool struct{ x bool }

// Bytes matches the IPLD Schema type "Bytes".  It has bytes kind.
type Bytes = *_Bytes
type _Bytes struct{ x []byte }

// FinishedTask matches the IPLD Schema type "FinishedTask".  It has Struct type-kind, and may be interrogated like map kind.
type FinishedTask = *_FinishedTask
type _FinishedTask struct {
	Status _Status
	StartedAt _Time
	ErrorMessage _String
	RetrievalTask _RetrievalTask__Maybe
	StorageTask _StorageTask__Maybe
	DealID _Int
	MinerMultiAddr _String
	ClientApparentAddr _String
	MinerLatencyMS _Int__Maybe
	TimeToFirstByteMS _Int__Maybe
	TimeToLastByteMS _Int__Maybe
	Events _Link_List_StageDetails
}

// Float matches the IPLD Schema type "Float".  It has float kind.
type Float = *_Float
type _Float struct{ x float64 }

// Int matches the IPLD Schema type "Int".  It has int kind.
type Int = *_Int
type _Int struct{ x int64 }

// Link matches the IPLD Schema type "Link".  It has link kind.
type Link = *_Link
type _Link struct{ x ipld.Link }

// Link_FinishedTask matches the IPLD Schema type "Link_FinishedTask".  It has link kind.
type Link_FinishedTask = *_Link_FinishedTask
type _Link_FinishedTask struct{ x ipld.Link }

// Link_List_StageDetails matches the IPLD Schema type "Link_List_StageDetails".  It has link kind.
type Link_List_StageDetails = *_Link_List_StageDetails
type _Link_List_StageDetails struct{ x ipld.Link }

// List matches the IPLD Schema type "List".  It has list kind.
type List = *_List
type _List struct {
	x []_Any__Maybe
}

// List_AuthenticatedRecord matches the IPLD Schema type "List_AuthenticatedRecord".  It has list kind.
type List_AuthenticatedRecord = *_List_AuthenticatedRecord
type _List_AuthenticatedRecord struct {
	x []_AuthenticatedRecord
}

// List_Logs matches the IPLD Schema type "List_Logs".  It has list kind.
type List_Logs = *_List_Logs
type _List_Logs struct {
	x []_Logs
}

// List_StageDetails matches the IPLD Schema type "List_StageDetails".  It has list kind.
type List_StageDetails = *_List_StageDetails
type _List_StageDetails struct {
	x []_StageDetails
}

// Logs matches the IPLD Schema type "Logs".  It has Struct type-kind, and may be interrogated like map kind.
type Logs = *_Logs
type _Logs struct {
	Log _String
	UpdatedAt _Time
}

// Map matches the IPLD Schema type "Map".  It has map kind.
type Map = *_Map
type _Map struct {
	m map[_String]MaybeAny
	t []_Map__entry
}
type _Map__entry struct {
	k _String
	v _Any__Maybe
}

// PopTask matches the IPLD Schema type "PopTask".  It has Struct type-kind, and may be interrogated like map kind.
type PopTask = *_PopTask
type _PopTask struct {
	Status _Status
	WorkedBy _String
}

// RecordUpdate matches the IPLD Schema type "RecordUpdate".  It has Struct type-kind, and may be interrogated like map kind.
type RecordUpdate = *_RecordUpdate
type _RecordUpdate struct {
	Records _List_AuthenticatedRecord
	SigPrev _Bytes
	Previous _Link__Maybe
}

// RetrievalTask matches the IPLD Schema type "RetrievalTask".  It has Struct type-kind, and may be interrogated like map kind.
type RetrievalTask = *_RetrievalTask
type _RetrievalTask struct {
	Miner _String
	PayloadCID _String
	CARExport _Bool
	Schedule _String__Maybe
	ScheduleLimit _String__Maybe
}

// StageDetails matches the IPLD Schema type "StageDetails".  It has Struct type-kind, and may be interrogated like map kind.
type StageDetails = *_StageDetails
type _StageDetails struct {
	Description _String__Maybe
	ExpectedDuration _String__Maybe
	Logs _List_Logs
	UpdatedAt _Time__Maybe
}

// Status matches the IPLD Schema type "Status".  It has int kind.
type Status = *_Status
type _Status struct{ x int64 }

// StorageTask matches the IPLD Schema type "StorageTask".  It has Struct type-kind, and may be interrogated like map kind.
type StorageTask = *_StorageTask
type _StorageTask struct {
	Miner _String
	MaxPriceAttoFIL _Int
	Size _Int
	StartOffset _Int
	FastRetrieval _Bool
	Verified _Bool
	Schedule _String__Maybe
	ScheduleLimit _String__Maybe
}

// String matches the IPLD Schema type "String".  It has string kind.
type String = *_String
type _String struct{ x string }

// Task matches the IPLD Schema type "Task".  It has Struct type-kind, and may be interrogated like map kind.
type Task = *_Task
type _Task struct {
	UUID _String
	Status _Status
	WorkedBy _String__Maybe
	Stage _String
	CurrentStageDetails _StageDetails__Maybe
	PastStageDetails _List_StageDetails__Maybe
	StartedAt _Time__Maybe
	RunCount _Int
	ErrorMessage _String
	RetrievalTask _RetrievalTask__Maybe
	StorageTask _StorageTask__Maybe
}

// Tasks matches the IPLD Schema type "Tasks".  It has list kind.
type Tasks = *_Tasks
type _Tasks struct {
	x []_Task
}

// Time matches the IPLD Schema type "Time".  It has int kind.
type Time = *_Time
type _Time struct{ x int64 }

// UpdateTask matches the IPLD Schema type "UpdateTask".  It has Struct type-kind, and may be interrogated like map kind.
type UpdateTask = *_UpdateTask
type _UpdateTask struct {
	Status _Status
	ErrorMessage _String__Maybe
	Stage _String__Maybe
	CurrentStageDetails _StageDetails__Maybe
	WorkedBy _String
	RunCount _Int
}

