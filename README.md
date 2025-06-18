## 🛠️ TorServe

**TorServe** is a hardened, zero-config static web server that automatically launches as a Tor hidden service. It allows users to anonymously publish web content with no setup, making it ideal for use in hostile or censored environments.

This tool is designed to be unzip-and-run, with built-in privacy protections and a multilingual landing page explaining safe usage.

🧭 **Quick Start**  
Unzip → Execute binary → Instant and safe Tor Hidden Service  
The `.onion` address for your site is generated automatically and printed in the console.

---

## 🎯 Project Goals

- 🧳 Minimal setup: unzip → run → get .onion URL
- 🕳️ No clearnet exposure
- 🕵️ Privacy-first by default
- 🛡️ Safe for use in hostile or censored environments *(use extreme caution for now)*
- 🌐 Multilingual landing page with embedded safety guide
- 🖥️ Cross-platform binaries: Linux, Raspberry Pi

---

## ✨ Key Features

- ✅ Hardened static web server written in Go
- ✅ Automatic Tor hidden service (bundled Tor)
- ✅ Multilingual `index.html` with safety instructions
- ✅ Image metadata scrubbing, jpg, webp, bmp, gif, png (EXIF)
- ✅ Optional rotating `.onion` support using –new-key command line arg
- ✅ No logs, no analytics, no clearnet connections
- ✅ 127.0.0.1 only (localhost)
- ✅ Chunked transfer + response padding to resist fingerprinting

---

## 🧠 Safety Features

- **Header hardening** – Strips User-Agent, Referer, ETag, and other identifying headers
- **Secure defaults** – No logs, no clearnet, localhost only (127.0.0.1), no directory listings
- **Metadata scrubbing** – Removes EXIF from images; PDF files are blocked entirely
- **Timing jitter** – Adds 50–200ms randomized delay to obscure traffic timing
- **Response padding** – Normalizes response sizes to resist size fingerprinting
- **No caching** – Disables ETag, Last-Modified, and other caching headers
- **No outbound JavaScript** – All assets are bundled and served locally
- **Multilingual safety guide** – Included React-based UI explains risks and usage across multiple languages

---

## 🖥️ Installation & Usage

### 🐧 Linux

```bash
unzip torserve-linux-amd64.zip
cd TorServ
./torserv
```

---

### 🪟 Windows

Binary Support dropped. Users may build from source if desired.

---

### 🍓 Raspberry Pi (ARMv7)

```bash
unzip torserve-rpi-arm64.zip
cd TorServ
./torserv
```

> The Tor hidden service will start and print a `.onion` address in the terminal. Use the Tor Browser to access your new hidden service. Download the browser here – [Tor Project](https://torproject.org/download/) 

---

## 🛠️ Build from Source

TorServe is written in Go and requires Tor to be available in the `tor/` directory. You can build for Linux, Windows, or Raspberry Pi with minimal setup.

---

### 📦 Requirements

* Go 1.20+ (`go version`)
* Git
* Cross-compiler if building for another OS (e.g., `mingw-w64` for Windows)
* `tor` binary from [torproject.org](https://www.torproject.org/download/tor/)

---

### 📁 Directory Setup

After cloning:

```bash
git clone https://github.com/torserv/torserv.git
cd TorServ
mkdir tor/
```

---

### 🔍 Get the Tor Binary

#### On Linux debian:

Install Tor if you haven’t already:

```bash
sudo apt install tor
```

Then locate the binary:

```bash
which tor
```

Copy it into the `tor/` directory:

```bash
cp $(which tor) tor/
```

> On Raspberry Pi, the same commands apply if you're using a Debian-based distro.

#### On Windows:

Download the **Expert Bundle** from [torproject.org](https://www.torproject.org/download/tor/).
Extract `tor.exe` into the `tor\` folder inside the project directory.

---

### 🖥️ Build for Linux (x86\_64)

```bash
go build -o release/linux/TorServ/torserve ./cmd/torserv
```

---

### 🪟 Build for Windows (x86\_64, requires mingw-w64)

> Install mingw-w64:
> `sudo apt install gcc-mingw-w64`

```bash
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
CGO_ENABLED=1 go build -o release/windows/TorServ/torserve.exe ./cmd/torserv
```
---

### 🍓 Build for Raspberry Pi 4+ (ARM64)

> Install the ARM cross-compiler if needed:
> `sudo apt install gcc-aarch64-linux-gnu`

```bash
GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
CGO_ENABLED=1 go build -o release/rpi/TorServ/torserve ./cmd/torserv
```
---

### 📂 After Building

Your binary will be in `release/<platform>/TorServ/`.
To run copy the bin to the project root and:

```bash
./torserve
```

TorServe will auto-launch the Tor hidden service if `tor/` is present or fail if not.

---

## 🌍 Demo Page Languages Support

- English (default)
- 简体中文 (Simplified Chinese)
- Español
- Русский (Russian)
- <span dir="ltr">فارسی (Farsi)</span>

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

## ❤️ Support This Project

TorServe is and always will be **Forever Free Open Source Software (FFOSS)**.

If TorServe helps you or your cause, please consider buying the dev a 🍔 or 🍺 — your support goes toward:

- Development
- Bug bounties
- Security audits
- Beer and cheeseburgers

👉 [Donate via PayPal](https://paypal.me/torserv)

---

## 🧭 Support Tor Browser

You can only access `.onion` sites using the [Tor Browser](https://www.torproject.org/download/).

Please consider supporting the [Tor Project](https://support.torproject.org/) — they make privacy tools like TorServe possible.

---

## 👋 Contributing

I welcome:
- Security audits
- Bug reports
- Feature suggestions
- Translations
- Code contributions
- Praise or constructive criticism

---

🧅 Tor Binary Licensing

This project bundles the official Tor binary (unmodified) for convenience.
Tor is licensed under the BSD 3-Clause License.

All credit for Tor goes to the Tor Project. This project is not affiliated with or endorsed by the Tor Project.
