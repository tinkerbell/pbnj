// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/task.proto

package v1

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *StatusRequest) Validate() error {
	if this.TaskId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("TaskId", fmt.Errorf(`value '%v' must not be an empty string`, this.TaskId))
	}
	return nil
}
func (this *StatusResponse) Validate() error {
	if this.Error != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Error); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Error", err)
		}
	}
	return nil
}
func (this *Error) Validate() error {
	return nil
}
