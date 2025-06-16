# 🛠️ TorServe

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

---

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

## 🌍 Demo Page Languages Support

- English (default)
- 简体中文 (Simplified Chinese)
- Español
- Русский (Russian)
- <span dir="ltr">فارسی (Farsi)</span>

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---<span dir="ltr">فارسی (Farsi)</span>

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
🧅 Tor Binary Licensing

This project bundles the official Tor binary (unmodified) for convenience.
Tor is licensed under the BSD 3-Clause License.

All credit for Tor goes to the Tor Project. This project is not affiliated with or endorsed by the Tor Project.
