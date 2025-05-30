# 🖼️ Go Image Watermarker

A concurrent CLI tool written in Go that recursively scans a folder of images and applies a watermark (either text or a logo). It supports basic image preprocessing steps like resizing or grayscale conversion, using an optional image processing pipeline.

---

## 🚀 Features

- ✅ Recursively scans folders for `.jpg`, `.jpeg`, `.png`
- ✅ Concurrent image processing using goroutines and channels
- ✅ Supports text and logo (image) watermarks
- ✅ Optional image processing pipeline (e.g., resize, grayscale)
- ✅ CLI options to configure:
  - Watermark type and content
  - Position (top-left, center, bottom-right)
  - Opacity
  - Output directory
- ✅ Skips invalid or corrupted images gracefully

---

## 🧠 Architecture Overview

### 🧵 Concurrency Pattern: Fan-Out / Fan-In

- Main goroutine sends file paths to a **worker pool**
- Workers concurrently:
  - Load & preprocess the image
  - Apply watermark (text or image)
  - Save result to output directory

### 🔁 Optional Pipeline Design

```text
File Path → [ Load Image ]
          → [ Preprocess (resize/grayscale) ]
          → [ Watermark ]
          → [ Save ]

Each step can be modularized and piped using channels for cleaner separation.

⸻

🗂️ Project Structure

go-image-watermarker/
│
├── cmd/
│   └── main.go              # CLI entry point
│
├── internal/
│   ├── scanner/             # Recursively scan directories
│   ├── processor/           # Image preprocessing (resize, grayscale, etc.)
│   ├── watermark/           # Text or logo watermark functions
│   ├── saver/               # Saving output images
│   └── pipeline/            # (Optional) stage-based pipeline abstraction
│
├── assets/
│   └── sample-logo.png      # Example watermark logo
│
├── output/                  # Processed images go here
│
├── go.mod
└── README.md


⸻

🛠️ CLI Usage

go run cmd/main.go --input ./images \
                   --output ./output \
                   --watermark-type text \
                   --text "CONFIDENTIAL" \
                   --position bottom-right \
                   --opacity 0.5 \
                   --resize 800x800

Or for a logo:

go run cmd/main.go --input ./images \
                   --output ./output \
                   --watermark-type image \
                   --logo ./assets/sample-logo.png \
                   --opacity 0.3


⸻

🧪 Possible Enhancements
 • Parallel image decoding using buffer pools
 • EXIF metadata preservation
 • GUI frontend or web API
 • Integration with cloud storage (S3, GCS)
 • JPEG compression level control

⸻

📚 Dependencies
 • golang.org/x/image
 • github.com/disintegration/imaging
 • github.com/golang/freetype (for text rendering)
