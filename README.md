# Go USB Device Auth

This tool monitors USB add events using [go-udev](https://github.com/jochenvg/go-udev) and can automatically **block specific USB interfaces** based on their `bInterfaceClass`, `bInterfaceSubClass`, and `bInterfaceProtocol` values.  
It is based on the [Linux kernel USB authorization](https://www.kernel.org/doc/Documentation/usb/authorization.txt) and works by writing `0` to the `authorized` attribute of matching interfaces in sysfs.

---

## Requirements
- Linux system with `udev` support.
- Root privileges (needed to modify `/sys/bus/usb/devices/.../authorized`).

## Configuration

The default settings block USB mass storage (class=8, subclass=6, protocol=50).
You can modify these constants in the code to target other devices:
```go
const (
    BlockIFClass    = 8
    BlockIFSubClass = 6
    BlockIFProtocol = 50
)
```

## Installation & Usage

1. Clone this repository:
   ```bash
   git clone https://github.com/AryanBakliwal/go-usb-auth.git
   cd go-usb-auth
    ```

2. Build the program:
   ```bash
   go build
   ```

3. Run with root privileges:
   ```bash
   sudo ./go-usb-auth
   ```
<br>

---

<br>
If you find this project useful, consider giving it a star!