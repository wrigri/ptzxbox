# PTZ Xbox
Windows tools for controlling a PTZ camera with an Xbox controller

## Basic Usage

There are 3 executables that can be built from this package:
- `ptz_xbox.exe`
- `ptz_query.exe`
- `udp_test_server.exe`

In addition, it is necessary to have the `config.yaml` file in the same folder
as the executables in order for them run. An example/starter `config.yaml`
file is available in the top level of this package.

### `ptz_xbox.exe`
The `ptz_xbox.exe` executable is the primary application for translating Xbox
controller inputs to a PTZ Camera using VISCA over IP commands. It can either
be run directly from the command line or as a Windows service (see below for
more information about running it as a service).

### `config.yaml`
The `config.yaml` file is used to configure the applicaion. The `config.yaml`
in the repository can be used mostly as-is with the execption of the `ipAddress:`
field which will need to be updated to point to the address of your PTZ
camera. See below for additional information about the settings in `config.yaml`.

### `ptz_query.exe`
The `ptz_query.exe` executable is used to send a query to your PTZ camera to
get the current pan, tilt and zoom values. These values can be used in the
`config.yaml` file to assign presets to the A, B, X, and Y buttons.

### `udp_test_server.exe`
The `udp_test_server.exe` executable is only used for testing. To use this,
set the `ipAddress:` in `config.yaml` to `127.0.0.1` (localhost) and it will
print out any commands that are sent to it from `ptz_xbox.exe`. It can be
used to make sure that `ptz_xbox.exe` is reading the Xbox controller
correctly and sending the correct commands to the server.

## Controller Configuration
The controller is (currently) configured for the following commands (see
"Customizing the code" below if you would like to change it):

- Joysticks: Move the camera up/down/left/right at variable speed.
- D-Pad: Move the camera up/down/left/right at fixed speed.
- Triggers: Zoom in/out at variable speed.
- Bumpers: Zoom in/out at fixed speed.
- A, B, X, Y Buttons: Preset camera positions configured in `config.yaml`
- Menu buttton: Turn off and on auto focus. When held down, the bumpers
can adjust the focus, but auto focus will be turned back on when the
button is released.
- View button:  Reset camera and reload the `config.yaml`

## Building/Installing
All code is written in Go. The Go compiler can be obtained from https://go.dev/ 

Executables can be built and installed by running:
```
go install github.com/wrigri/ptzxbox/cmd/ptz_xbox@latest
go install github.com/wrigri/ptzxbox/cmd/ptz_query@latest
go install github.com/wrigri/ptzxbox/cmd/udp_test_server.exe@latest
```
By default this will install the executable files to your `%HOME%\go\bin`
folder. You will also need to copy the `config.yaml` file into that
directory for things to run.

You can also download the source from this repository and run the following
commands from the root directory of the repository to build:
```
go build .\cmd\ptz_xbox\
go build .\cmd\ptz_query\
go build .\cmd\udp_test_server\
```

## Configuration
To configure the application, you must have a `config.yaml` file in the same
folder as your executables. Here are the fields you can configure:

- `debug:` Set this to `true` in order to print extra information out when
running `ptz_xbox.exe` from the command-line.
- `network/ipAddress:` Should be set to the IP Address of your PTZ camera
(or `127.0.0.1` for testing).
- `network/shortCmdPort:` If you are using a 'PTZOptics' camera, it can accept
a short version of commands on port `1259`.
- `network/longCmdPort:` This is the port for accepting standard 'VISCA over
IP' commands. The most common port appears to be `52381` but you may need to
double-check the documentation for your camera to confirm the port number
that the camera listens on.
- `network/useShortCommands:` If set to `true`, the application will use the
short versions of the commands on the `shortCmdPort`. Otherwise, if set to
`false`, the application will use the long versions of the commands on the
`longCmdPort`.
- `controller/deadZones:` This field has both a `low:` and `high:` setting.
These should be values between `0` and `1` which define the areas of the
Xbox joysticks that will be ignored. The `low:` value determines the area
in the center of the joystick that is ignored and the `high:` value determines
the area around the edge of the joysticks which is ignored.
- `camera:` There are several settings under `camera:`. This is where the
speed values are set for 'pan', 'tilt', and 'zoom'. Pan speeds are from `1`
to `24`. Tilt speeds are from `1` to `20` and Zoom speeds are from `0` , to
`7`. The 'max' values are used when the joystick is at the furthest from
center position, the 'min' values are used when the joystick is close to
the center and the 'default' values are used when using the dPad. For Zoom,
'min' and 'max' values are used by the triggers and the 'default' valued is
used by the bumpers. Available settings under `camera:` are `maxPanSpeed:`,
`minPanSpeed:`, `defaultPanSpeed:`, `maxTiltSpeed:`, `minTiltSpeed:`,
`defaultTiltSpeed`, `maxZoomSpeed:`, `minZoomSpeed:`, and `defaultZoomSpeed:`.
- `presets:` Under presets, there are sections for each button: `A:`, `B:`,
`X:`, and `Y:`. In addition, there are settings for `panSpeed:` and
`tiltSpeed:`. The speed options determine how fast the camera moves when
switching to a preset postion. Under each of the 'button' sections, you can
set `pan:`, `tilt:`, and `zoom:`. The values for each of these is a 4-digit
hexadecimal number. In addition, `panSpeed:` and `tiltSpeed:` can be set on
individual buttons and will override the values under `presets:`. Valid values
can be determined from your camera's documentation, but generally the best way
to set these is to move the camera to the desired postion and use
`ptz_query.exe` to get the values which can be copy/pasted into the
`config.yaml` file.

## Setting up a Windows Service
Instructions TBD (There are additional changes being made to the service part
of the code and this section will be updated after those changes are complete.)

## Customizing the code
If you would like to customize the code itself and change things like which
buttons send which commands, this will be a quick guide to the general layout
of the code. Except for some basic wrapper code, most of the code is in the
`internal` directory in different packages.
- `config` package: The `config` package is responsible for finding and
reading the `config.yaml` file and loading all of the values into a
`configData` struct, defined in that package. This is where any additional
config options would be handled.
- `controller` package: In the `controller` package are 3 files. 
    - `mainloop.go` is were the primary loop of the application runs. This
    reads the current state of the Xbox controller and determines if anything
    has changed. It runs about 20 times per second.
    - `actions.go` is where each input (button, D-Pad, joystick, trigger,
    etc.) is checked to see if it has changed. If it has, it calls a command
    in the `camctrl` package. This is the file you would change if you wanted
    a button to do something different.
    - `common.go` contains helper functions mostly for the `actions.go` file.
    It does things like converting the 'X' & 'Y' coordinates of the joysticks
    to a direction and distance from the center. It also implements the dead
    zones for the joysticks and converts the individual button presses of the
    D-Pad to a direction.
- `camctrl` package: In the `camctrl` package are two files.
    - `commands.go` is where named 'commands' are converted into the
    hexadecimal commands expected by the PTZ camera. For example, the `Home()`
    function sends the command `81 01 06 04 FF` to the camera. This is the
    file you would update if you wanted to add additional commands from the
    VISCA command set. There are commands defined in this file that are not 
    currently configured in the `controller/actions.go` file.
    - `common.go` contains helper functions for sending commands to the
    camera. Primarily the `sendCommand` function which constructs the command
    to send based on whether long or short commands are configured.
- `netcon` package: The `netcon` package handes the network conneciton to
the camera.
- `util` package: The `util` package includes some helper functions primarily
dealing with hexadecimal/byte conversions. For example, converting a value
like `0x3B7C` in the `config.yaml` file to a array of bytes like `03 0B 07 0C`
as expected by the camera in some commands.

The code for the 3 executables that can be built for this package can be 
found in:

- `cmd/ptz_query/main.go` - Source for ptz_query.exe.
- `cmd/ptz_xbox/main.go` - Source for ptz_xbox.exe.
- `cmd/udp_test_server/main.go` - Source for udp_test_server.exe.
