# MoveBin - Global Binary Installer

A cross-platform Go CLI tool to move a binary to a system-wide path (e.g., `/usr/local/bin` on Unix). Handles macOS, Linux, and Windows.

## Features

- Auto-detects OS and target directory.
- Sets executable permissions on Unix.
- Validates file existence and prevents overwrites.
- Creates destination dir if needed.

## Build & Install

1. Ensure Go 1.21+ is installed.
2. Save `movebin.go`.
3. Build: `go build movebin.go`.
4. Self-install: `sudo ./movebin ./movebin`.
5. Run: `movebin /path/to/binary`.

## Usage

```bash
movebin ./myapp
```

- macOS/Linux: Moves to `/usr/local/bin/myapp` (use `sudo` for perms).
- Windows: Moves to `C:\Program Files\bin\myapp.exe` (run as Admin).

Output: `Successfully moved /path/to/myapp to /usr/local/bin/myapp`.

## License

MIT.
