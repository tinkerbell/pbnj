FORMAT: 1A

# Power/Boot API

## /devices/{ip}/power
Manage the power state of a device.

+ Parameters
    + ip: 10.250.1.84 (string) - Address of the BMC

### GET

+ Response 200 (application/json)

        {
            "state": "on"
        }

### POST

+ Attributes
    + action (required, enum[string]) - The power action to perform.
        - on
        - off
        - cycle
        - reset

    + soft_timeout (optional, string) - How long to wait for device to gracefully shut down. If no value is given, a hard power off will be performed.

    + off_duration (optional, string) - How long to leave device off during a power cycle.
        + Default: 1s

+ Request Power On

        {
            "action": "on"
        }

+ Response 200 (application/json)

        {
            "state": "on"
        }

+ Request Power Off

        {
            "action": "off",
            "soft_timeout": "10s"
        }

+ Response 200 (application/json)

        {
            "state": "off"
        }

+ Response 202

    + Headers

        ```
        Location: /task/d7409dce-46ba-494c-a320-749c460357e1
        ```

+ Request Power Cycle

        {
            "action": "cycle",
            "soft_timeout": "10s",
            "off_duration": "2s"
        }

+ Response 200 (application/json)

        {
            "state": "on"
        }

+ Response 202

    + Headers

        ```
        Location: /task/d7409dce-46ba-494c-a320-749c460357e1
        ```

+ Request Reset

        {
            "action": "reset"
        }

+ Response 200 (application/json)

        {
            "state": "on"
        }


## /devices/{ip}/boot
Manage the boot device and options for a device.

+ Parameters
    + ip: 10.250.1.84 (string) - Address of the BMC

### PATCH
Set the boot device to use for the next boot.

+ Attributes
    + device (required, enum[string]) - The boot device to use.
        - none
        - bios
        - cdrom
        - disk
        - pxe

    + persistent (optional, boolean) - If true, use device for all future boots.
        + Default: false

+ Request

        {
            "device": "pxe"
        }

+ Response 204
