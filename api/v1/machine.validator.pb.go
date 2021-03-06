// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/machine.proto

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

func (this *DeviceRequest) Validate() error {
	if this.Authn != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Authn); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Authn", err)
		}
	}
	if this.Vendor != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Vendor); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Vendor", err)
		}
	}
	if _, ok := BootDevice_name[int32(this.BootDevice)]; !ok {
		return github_com_mwitkow_go_proto_validators.FieldError("BootDevice", fmt.Errorf(`value '%v' must be a valid BootDevice field`, this.BootDevice))
	}
	return nil
}
func (this *DeviceResponse) Validate() error {
	return nil
}
func (this *PowerRequest) Validate() error {
	if this.Authn != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Authn); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Authn", err)
		}
	}
	if this.Vendor != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Vendor); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Vendor", err)
		}
	}
	if _, ok := PowerAction_name[int32(this.PowerAction)]; !ok {
		return github_com_mwitkow_go_proto_validators.FieldError("PowerAction", fmt.Errorf(`value '%v' must be a valid PowerAction field`, this.PowerAction))
	}
	if !(this.SoftTimeout > -1) {
		return github_com_mwitkow_go_proto_validators.FieldError("SoftTimeout", fmt.Errorf(`value '%v' must be greater than '-1'`, this.SoftTimeout))
	}
	if !(this.OffDuration > -1) {
		return github_com_mwitkow_go_proto_validators.FieldError("OffDuration", fmt.Errorf(`value '%v' must be greater than '-1'`, this.OffDuration))
	}
	return nil
}
func (this *PowerResponse) Validate() error {
	return nil
}
