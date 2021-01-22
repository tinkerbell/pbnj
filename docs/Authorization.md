# Authorization

Authorization in PBnJ uses json web tokens (JWT).

- More info about JWT can be found at <https://jwt.io>

Authorization in PBnJ is optional.

- It is disabled by default.
  To enable it set the cli flag or env var or config value of `enableAuthz` to true.

Enabled Authorization must also include a symmetric or asymmetric key.

- HMAC(symmetric) and RS(asymmetric) keys are supported.
- Use the cli flag or env var or config value `HSKey` for symmetric HMAC.

```yaml
HSKey: "supersecret"
```

- Use the cli flag or env var or config value `RSPubKey` for asymmetric RS.

```yaml
RSPubKey: |
  -----BEGIN PUBLIC KEY-----
  MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnzyis1ZjfNB0bBgKFMSv
  vkTtwlvBsaJq7S5wA+kzeVOVpVWwkWdVha4s38XM/pa/yr47av7+z3VTmvDRyAHc
  aT92whREFpLv9cj5lTeJSibyr/Mrm/YtjCZVWgaOYIhwrXwKLqPr/11inWsAkfIy
  tvHWTxZYEcXLgAXFuUuaS3uF9gEiNQwzGTU1v0FqkqTBr4B8nW3HCN47XUu0t8Y0
  e+lf4s4OxQawWD79J9/5d3Ry0vbV3Am1FtGJiJvOwRsIfVChDpYStTcHTCMqtvWb
  V6L11BWkpzGXSW4Hv43qa+GSYOD2QU68Mb59oSk2OB+BtOLpJofmbGEGgvmwyCI9
  MwIDAQAB
  -----END PUBLIC KEY-----
```

When enabled, Authorization will protect the following RPC methods

- github.com.tinkerbell.pbnj.api.v1.
  - Machine/Power
  - Machine/BootDevice
  - BMC/NetworkSource
  - BMC/Reset
  - BMC/CreateUser
  - BMC/DeleteUser
  - BMC/UpdateUser

Clients must set the following gRPC metadata/header for requests

```evans grpc client
github.com.tinkerbell.pbnj.api.v1@127.0.0.1:50051> header authorization="bearer eyJhbG.eyJwcm9.3u8kMN"
```
