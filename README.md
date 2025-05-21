# HEIC to JPEG Converter (with EXIF support)

This is a simple Go-based command-line tool to convert HEIC image files to JPEG format, preserving EXIF metadata.

---

## 📦 Features

- Converts `.heic` to `.jpg`
- Preserves EXIF metadata
- Customizable output filename
- Optional verbose logging

---

## 🛠 Requirements

- [Go](https://golang.org/dl/) 1.18 or higher
- macOS or Linux (HEIC support may require native libraries)

---

## 📥 Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/heic-to-jpg
cd heic-to-jpg
```
2. Download dependencies:

```bash
go mod tidy
```

3. Build the executable:

```bash
go build -o heicconvert
```

## 🚀 Usage
```bash
./heicconvert -input=input.heic -output=output.jpg [-verbose]
```

### Example:
```bash
./heicconvert -input=photo.heic -output=photo.jpg -verbose
```
