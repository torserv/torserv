# 🛡️ TorServe — Safety & Threat Model

**TorServe** is designed to help users publish static content over Tor with strong privacy defaults and zero configuration. This guide explains the safety principles behind TorServe and outlines key risks users should understand before using it.

---

## 🔐 Privacy Principles

TorServe follows these core principles:

- **No clearnet exposure** — Only serves over Tor's `.onion` network
- **Runs on localhost only** — Listens only on `127.0.0.1`
- **No logging** — No access logs, error logs, or analytics
- **No directory listing** — Prevents accidental exposure of extra files
- **All assets are local** — No JavaScript or image calls to external services
- **Static-only** — No dynamic upload forms, cookies, or sessions

---

## 📁 Content Safety

- **PDFs are Blocked**  
  PDF files are explicitly **not supported**. Even when sanitized, PDFs can carry embedded JavaScript, fonts, tracking pixels, or remote content.  
  If any `.pdf` file is detected in `public/`, the server will **refuse to start** or **reject it at runtime**.

- **Image Metadata Scrubbing**  
  TorServe removes EXIF and embedded metadata from:
  - `.jpg`, `.jpeg`, `.png`

  This protects against GPS coordinates, device info, and timestamps.

---

## 🧠 Common Risks

| Risk | How TorServe Reduces It |
|------|--------------------------|
| **Deanonymization via logs** | No logs are written anywhere. |
| **Leak via clearnet access** | Binds to `127.0.0.1`, not accessible from the LAN or internet. |
| **Metadata in files** | Auto-stripping of image metadata; blocks unsafe file types like PDFs. |
| **Timing analysis** | Adds 50–200ms randomized jitter to each request. |
| **Response size fingerprinting** | Pads all file responses to fixed block sizes. |
| **Cache leakage** | All caching headers are disabled. No ETag or Last-Modified. |

---

## 🚫 What TorServe Does *Not* Protect You From

- Malicious content you upload yourself
- Unsafe file types not recognized or filtered
- Compromised system (e.g., malware or local keyloggers)
- JavaScript or tracking in your own content
- DNS, browser, or OS leaks outside Tor (e.g., from Tor misconfiguration)

---

## ✅ Tips for Safe Usage

- **Use Tails or Whonix** if publishing from a sensitive environment
- **Test your files locally first** with Tor Browser before sharing the `.onion`
- **Never share `.onion` URLs on clearnet channels** (especially social media)
- **Use only image formats like `.png`, `.jpg`, `.gif`, `.bmp`, `.webp`**
- **Keep content non-interactive** (no forms, uploads, embedded scripts)

---

## 📌 Summary

TorServe tries to make anonymous publishing accessible and safe, but users must still act carefully. **Understand what you’re serving. Know your risks.**

If you're unsure — don't serve it.

---
