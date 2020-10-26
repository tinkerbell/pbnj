package oob

import "github.com/golang/protobuf/ptypes/any"

// User management methods
type User interface {
	Create() (result string, err *Error)
	Update() (result string, err *Error)
	Delete() (result string, err *Error)
}

// Machine management methods
type Machine interface {
	BootDevice() (result string, err Error)
	Power() (result string, err Error)
}

// BMC management methods
type BMC interface {
	Reset() (result string, err Error)
	NetworkSource() (result string, err Error)
}

// Error for all bmc actions
type Error struct {
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// A developer-facing human-readable error message in English. It should
	// both explain the error and offer an actionable resolution to it.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// Additional error information that the client code can use to handle
	// the error, such as retry delay or a help link.
	Details []*any.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
}
