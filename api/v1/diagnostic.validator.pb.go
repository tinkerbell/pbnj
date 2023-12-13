// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/diagnostic.proto

package v1

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ScreenshotRequest) Validate() error {
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
	return nil
}
func (this *ScreenshotResponse) Validate() error {
	return nil
}
func (this *ClearSystemEventLogRequest) Validate() error {
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
	return nil
}
func (this *ClearSystemEventLogResponse) Validate() error {
	return nil
}
func (this *SystemEventLogRequest) Validate() error {
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
	return nil
}
func (this *SystemEventLogEntry) Validate() error {
	return nil
}
func (this *SystemEventLogResponse) Validate() error {
	for _, item := range this.Events {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Events", err)
			}
		}
	}
	return nil
}
func (this *SystemEventLogRawRequest) Validate() error {
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
	return nil
}
func (this *SystemEventLogRawResponse) Validate() error {
	return nil
}
