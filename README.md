## 🛠️ torserv

**torserv** is a hardened, zero-config static web server that automatically launches as a Tor hidden service. It enables anonymous web publishing with no setup, making it ideal for hostile or censored environments.

Designed to be unzip-and-run, it includes built-in privacy protections and a multilingual landing page explaining safe usage.

---

### 🧭 Quick Start

**Unzip → Execute binary → Get instant `.onion` address**

The Tor hidden service starts automatically. The `.onion` address is printed to the console.

---

## 🎯 Project Goals

* 🧳 Minimal setup: unzip → run → get `.onion` URL
* 🕳️ No clearnet exposure
* 🕵️ Privacy-first by default
* 🛡️ Safe for use in hostile or censored environments *(use extreme caution)*
* 🌐 Multilingual landing page with embedded safety guide
* 🖥️ Cross-platform binaries for Linux and Raspberry Pi

---

## ✨ Key Features

* ✅ Hardened static web server (Go)
* ✅ Automatic Tor hidden service (bundled `tor`)
* ✅ Multilingual `index.html` with usage guidance
* ✅ Image metadata scrubbing: JPG, WebP, BMP, GIF, PNG (EXIF)
* ✅ Optional `.onion` key rotation via `--new-key`
* ✅ No logs, analytics, or clearnet requests
* ✅ Binds only to `127.0.0.1`
* ✅ Chunked transfer + response padding for fingerprinting resistance

---

## 🧠 Safety Features

* **Header Hardening** – Strips User-Agent, Referer, ETag, etc.
* **Secure Defaults** – Localhost-only, no logs, no directory listing
* **Metadata Scrubbing** – EXIF stripped from many file types
* **Timing Jitter** – Random 50–200ms delay to mask response timing
* **Response Padding** – Normalized sizes prevent fingerprinting
* **No Caching** – Disables ETag, Last-Modified, public caching
* **Bundled Assets Only** – No outbound JavaScript or network calls
* **Multilingual Safety Guide** – Static html with tabbed translations
* **File/Directory Obfuscation** – Encrypted links in html hiding file names and directory structure

---

## 🖥️ Installation & Usage

### 🐧 Linux (x86\_64)

```bash
unzip torserv-linux-amd64.zip
cd TorServ
./torserv
```

---

### 🍓 Raspberry Pi (ARM64)

```bash
unzip torserv-rpi-arm64.zip
cd TorServ
./torserv
```

> The Tor hidden service will start and print a `.onion` address to the terminal.
> Use [Tor Browser](https://www.torproject.org/download/) to access it.

---

### 🚫 Windows

⚠️ **Note:** Windows release dropped due to aggressive antivirus false positives.
You may still build from source if needed.

---

## 🛠️ Build from Source

torserv is written in Go and requires the Tor binary to be present in a `tor/` directory inside the project.

---

### 📦 Requirements

* Go 1.20+
* Git
* `tor` binary (from [torproject.org](https://www.torproject.org/download/tor/))
* Optional: cross-compilers for other platforms

---

### 📁 Setup Instructions

```bash
git clone https://github.com/torserv/torserv.git
cd TorServ
mkdir tor/
```

---

### 🔍 Install Tor Binary

#### Linux (Debian-based)

```bash
sudo apt install tor
which tor
cp $(which tor) tor/
```

> Same applies on Raspberry Pi if using Raspbian/Debian.

#### Windows

Download the **Tor Expert Bundle** and place `tor.exe` into `tor\`.

---

### 🔧 Build Commands

#### 🐧 Linux (x86\_64 or ARM64, including Raspberry Pi)

```bash
go build -o release/linux/TorServ/torserv ./cmd/torserv
```

#### 🪟 Windows (Cross-compile from Linux/macOS)

```bash
sudo apt install gcc-mingw-w64

GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
CGO_ENABLED=1 go build -o release/windows/TorServ/torserv.exe ./cmd/torserv
```

#### 🍓 Raspberry Pi 4+ (ARM64, cross-compiled)

```bash
sudo apt install gcc-aarch64-linux-gnu

GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
CGO_ENABLED=1 go build -o release/rpi/TorServ/torserv ./cmd/torserv
```

---

### 📂 After Building

Your binary will be in:
`release/<platform>/TorServ/`

Copy it to the project root to run:

```bash
./torserv
```

torserv will auto-launch the Tor hidden service if `tor/` is present. If not, it will exit.

---

## 🌍 Demo Page Language Support

* English (default)
* 简体中文 (Simplified Chinese)
* Español
* Русский (Russian)
* <span dir="ltr">فارسی (Farsi)</span>

---

## 📜 License

This project is licensed under the MIT License (see LICENSE.md)

---

## ❤️ Support This Project

torserv is and always will be **Forever Free Open Source Software (FFOSS)**.

If it helps you or your mission, consider buying the dev a 🍔 or 🍺:

👉 [Donate via PayPal](https://paypal.me/torserv)

---

## 🧭 Support the Tor Project

You’ll need the [Tor Browser](https://www.torproject.org/download/) to access `.onion` sites.

If you care about privacy, consider supporting the [Tor Project](https://support.torproject.org/).

---

## 👋 Contributing

Welcoming:

* 🔐 Security audits
* 🐞 Bug reports
* 🌟 Feature requests
* 🌍 Translations
* 💻 Code contributions
* 🧠 Thoughtful feedback

---

### 🧅 Tor Binary Licensing

torserv bundles the unmodified official `tor` binary for convenience.
Tor is licensed under the **BSD 3-Clause License**.

This project is **not affiliated with or endorsed by the Tor Project**.
All credit for Tor belongs to [The Tor Project](https://www.torproject.org/).
