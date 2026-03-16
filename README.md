# PVM – PHP Version Manager for Windows

**PVM** is a PHP version manager for Windows, designed to make installing, managing, and switching between multiple PHP versions simple. PVM includes both a **CLI tool** (`pvm.exe`) written in Go and a **GUI installer** (`pvm-setup.exe`) for easy installation.

---

## Features

- Install and manage multiple PHP versions on Windows.
- Switch PHP versions globally with a single command.
- Support for Thread-Safe (TS) and Non-Thread-Safe (NTS) builds.
- Automatic environment variable (`PVM_HOME`) and PATH configuration.
- GUI installer with symlink management.
- Lightweight and portable CLI tool written in Go.
- Composer integration with wrapper scripts.

---

## Installation

### Using the GUI Installer

1. Download the latest installer from `dist/{version}/pvm-setup.exe`.  
2. Run `pvm-setup.exe` and follow the instructions:  
   - Choose the folder where PHP versions will be linked.  
   - Environment variables and PATH are configured automatically.  
3. Verify installation:

```powershell
pvm help
````

### Using the CLI Directly

1. Run `pvm.exe` from its location or copy it to a folder in your PATH:

```powershell
C:\path\to\pvm.exe help
```

2. Common commands:

```text
pvm install <version>     # Install a specific PHP version
pvm use <version>         # Switch to a specific PHP version
pvm list                  # List installed PHP versions
pvm uninstall <version>   # Remove a PHP version
```

---

## Project Structure

```
pvm-master/
├─ bin/                 # CLI wrappers and helper scripts
├─ commands/            # Go command packages
├─ common/              # Shared utilities
├─ theme/               # CLI color and theme utilities
├─ dist/                # Installer builds
│   └─ {version}/pvm-setup.exe
├─ LICENSE              # License file
├─ LICENSE.txt          # License text
├─ main.go              # Main Go executable
├─ go.mod               # Go module file
├─ go.sum               # Go dependencies
├─ pvm.ico              # Application icon
├─ pvm-setup.iss        # Inno Setup installer script
├─ README.md            # This file
├─ SUPPORT.md           # Support information
```

---

## Contributing

* Report bugs or request features via GitHub Issues.
* Submit pull requests to improve functionality.

---

## License

PVM is licensed under the **MIT License**. See the `LICENSE` file for details.

---

## Links

* **Project URL:** [https://github.com/ashoknepalofficial/pvm](https://github.com/ashoknepalofficial/pvm)
* **Releases:** `dist/{version}/pvm-setup.exe`
* **Support:** See `SUPPORT.md` for help and contact information


---

This version:

- Is **concise and professional**.  
- Mentions **CLI, GUI installer, and Go executable**.  
- Shows **folder structure** and usage instructions.  
- Ready to be copied to your GitHub repo.  
