HEAD
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
- 🖥️ Cross-platform binaries: Linux, Windows, Raspberry Pi

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

## 📁 File Structure

```
```
TorServ/
├── torserve        # Main binary (or torserve.exe on Windows)
├── torrc             # Minimal Tor config
├── public/         # Static content root
│   ├── index.html
│   ├── fonts/      # Web fonts for multilingual UI
│   └── ...            # Other translations or assets
├── LICENSE
├── README.md
```

---

## 🧠 Safety Features

- **Header hardening** – Strips User-Agent, Referer, ETag, and other identifying headers
- **Secure defaults** – No logs, no clearnet, localhost only (127.0.0.1), no directory listings
- **Metadata scrubbing** – Removes EXIF from images; PDF files are blocked entirely
- **Timing jitter** – Adds 50–200ms randomized delay to obscure traffic timing
- **Response padding** – Normalizes response sizes to resist size fingerprinting
- **No caching** – Disables ETag, Last-Modified, and other caching headers
- **No outbound JavaScript** – All assets are bundled and served locally
- **Ephemeral mode** – Optional: generates a temporary `.onion` and discards key on shutdown
- **Multilingual safety guide** – Included React-based UI explains risks and usage across multiple languages

---

## 🖥️ Installation & Usage

---

### 🐧 Linux

```bash
unzip torserve-linux-amd64.zip
cd TorServ
./torserve
```

---

### 🪟 Windows

```bash
Unzip torserve-windows-amd64.zip
cd TorServ
Double-click torserve.exe
```

---

### 🍓 Raspberry Pi (ARMv7)

```bash
unzip torserve-rpi-arm64.zip
cd TorServ
./torserve
```

> The Tor hidden service will start and print a `.onion` address in the terminal. Use the Tor Browser to access your new hidden service. Download the browser here – [Tor Project](https://support.torproject.org/) 

---

## 🌍 Demo Page Languages Support

- English (default)
- 简体中文 (Simplified Chinese)
- Español
- فارسی (Farsi)
- Русский (Russian)

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

Open a GitHub issue or pull request to get involved.
=======
## Hi there 👋

<!--
**torserv/torserv** is a ✨ _special_ ✨ repository because its `README.md` (this file) appears on your GitHub profile.

Here are some ideas to get you started:

- 🔭 I’m currently working on ...
- 🌱 I’m currently learning ...
- 👯 I’m looking to collaborate on ...
- 🤔 I’m looking for help with ...
- 💬 Ask me about ...
- 📫 How to reach me: ...
- 😄 Pronouns: ...
- ⚡ Fun fact: ...
-->
f22c32046d59fc4b415cdf0691f591d1eafde19c
