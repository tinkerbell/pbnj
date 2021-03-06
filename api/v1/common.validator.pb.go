// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/common.proto

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

func (this *Host) Validate() error {
	if this.Host == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Host", fmt.Errorf(`value '%v' must not be an empty string`, this.Host))
	}
	return nil
}
func (this *ExternalAuthn) Validate() error {
	if this.Host != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Host); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Host", err)
		}
	}
	return nil
}
func (this *DirectAuthn) Validate() error {
	if this.Host != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Host); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Host", err)
		}
	}
	if this.Username == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Username", fmt.Errorf(`value '%v' must not be an empty string`, this.Username))
	}
	if this.Password == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Password", fmt.Errorf(`value '%v' must not be an empty string`, this.Password))
	}
	return nil
}
func (this *Authn) Validate() error {
	if oneOfNester, ok := this.GetAuthn().(*Authn_DirectAuthn); ok {
		if oneOfNester.DirectAuthn != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.DirectAuthn); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("DirectAuthn", err)
			}
		}
	}
	return nil
}
func (this *Vendor) Validate() error {
	return nil
}
