// Package main provides a CLI tool for creating distribution packages.
// It supports tar.gz, zip, and deb formats using pure Go libraries.
package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "tar":
		tarCmd := flag.NewFlagSet("tar", flag.ExitOnError)
		output := tarCmd.String("output", "", "Output file path (required)")
		binary := tarCmd.String("binary", "", "Binary file to package (required)")
		name := tarCmd.String("name", "", "Name of binary in archive (defaults to basename of binary)")
		if err := tarCmd.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if *output == "" || *binary == "" {
			tarCmd.Usage()
			os.Exit(1)
		}
		binaryName := *name
		if binaryName == "" {
			binaryName = filepath.Base(*binary)
		}
		if err := createTarGz(*output, *binary, binaryName); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created %s\n", *output)

	case "zip":
		zipCmd := flag.NewFlagSet("zip", flag.ExitOnError)
		output := zipCmd.String("output", "", "Output file path (required)")
		binary := zipCmd.String("binary", "", "Binary file to package (required)")
		name := zipCmd.String("name", "", "Name of binary in archive (defaults to basename of binary)")
		if err := zipCmd.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if *output == "" || *binary == "" {
			zipCmd.Usage()
			os.Exit(1)
		}
		binaryName := *name
		if binaryName == "" {
			binaryName = filepath.Base(*binary)
		}
		if err := createZip(*output, *binary, binaryName); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created %s\n", *output)

	case "deb":
		debCmd := flag.NewFlagSet("deb", flag.ExitOnError)
		output := debCmd.String("output", "", "Output file path (required)")
		binary := debCmd.String("binary", "", "Binary file to package (required)")
		name := debCmd.String("name", "render", "Package name")
		version := debCmd.String("version", "0.0.0", "Package version")
		arch := debCmd.String("arch", "amd64", "Architecture (amd64, arm64)")
		maintainer := debCmd.String("maintainer", "Unknown", "Maintainer name and email")
		description := debCmd.String("description", "", "Package description")
		if err := debCmd.Parse(os.Args[2:]); err != nil {
			os.Exit(1)
		}
		if *output == "" || *binary == "" {
			debCmd.Usage()
			os.Exit(1)
		}
		pkg := debPackage{
			Name:        *name,
			Version:     *version,
			Arch:        *arch,
			Maintainer:  *maintainer,
			Description: *description,
		}
		if err := createDeb(*output, *binary, pkg); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created %s\n", *output)

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: package <command> [options]

Commands:
  tar     Create a tar.gz archive
  zip     Create a zip archive
  deb     Create a Debian package

Examples:
  package tar -binary bin/linux/amd64/render -output dist/render-linux-amd64.tar.gz
  package zip -binary bin/windows/amd64/render.exe -output dist/render-windows-amd64.zip
  package deb -binary bin/linux/amd64/render -output dist/render_1.0.0_amd64.deb -version 1.0.0`)
}

// createTarGz creates a gzipped tar archive containing the binary.
func createTarGz(output, binaryPath, binaryName string) (err error) {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	info, err := os.Stat(binaryPath)
	if err != nil {
		return fmt.Errorf("stat binary: %w", err)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close output file: %w", cerr)
		}
	}()

	gw := gzip.NewWriter(f)
	defer func() {
		if cerr := gw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close gzip writer: %w", cerr)
		}
	}()

	tw := tar.NewWriter(gw)
	defer func() {
		if cerr := tw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close tar writer: %w", cerr)
		}
	}()

	header := &tar.Header{
		Name:    binaryName,
		Size:    info.Size(),
		Mode:    0o755,
		ModTime: info.ModTime(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("write tar header: %w", err)
	}

	binaryFile, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("open binary: %w", err)
	}
	defer func() { _ = binaryFile.Close() }()

	if _, err := io.Copy(tw, binaryFile); err != nil {
		return fmt.Errorf("copy binary to tar: %w", err)
	}

	return nil
}

// createZip creates a zip archive containing the binary.
func createZip(output, binaryPath, binaryName string) (err error) {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close output file: %w", cerr)
		}
	}()

	zw := zip.NewWriter(f)
	defer func() {
		if cerr := zw.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close zip writer: %w", cerr)
		}
	}()

	binaryFile, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("open binary: %w", err)
	}
	defer func() { _ = binaryFile.Close() }()

	info, err := binaryFile.Stat()
	if err != nil {
		return fmt.Errorf("stat binary: %w", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("create zip header: %w", err)
	}
	header.Name = binaryName
	header.Method = zip.Deflate

	w, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create zip entry: %w", err)
	}

	if _, err := io.Copy(w, binaryFile); err != nil {
		return fmt.Errorf("copy binary to zip: %w", err)
	}

	return nil
}

// debPackage holds metadata for a Debian package.
type debPackage struct {
	Name        string
	Version     string
	Arch        string
	Maintainer  string
	Description string
}

// createDeb creates a Debian package containing the binary.
// A .deb file is an ar archive containing:
// - debian-binary (version string "2.0\n")
// - control.tar.gz (package metadata)
// - data.tar.gz (package contents)
func createDeb(output, binaryPath string, pkg debPackage) (err error) {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// Read binary
	binaryData, err := os.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("read binary: %w", err)
	}

	// Create control.tar.gz
	controlTar, err := createControlTar(pkg)
	if err != nil {
		return fmt.Errorf("create control tar: %w", err)
	}

	// Create data.tar.gz
	dataTar, err := createDataTar(pkg.Name, binaryData)
	if err != nil {
		return fmt.Errorf("create data tar: %w", err)
	}

	// Create ar archive
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close output file: %w", cerr)
		}
	}()

	// Write ar magic
	if _, err := f.Write([]byte("!<arch>\n")); err != nil {
		return fmt.Errorf("write ar magic: %w", err)
	}

	// Write debian-binary
	debianBinary := []byte("2.0\n")
	if err := writeArEntry(f, "debian-binary", debianBinary); err != nil {
		return fmt.Errorf("write debian-binary: %w", err)
	}

	// Write control.tar.gz
	if err := writeArEntry(f, "control.tar.gz", controlTar); err != nil {
		return fmt.Errorf("write control.tar.gz: %w", err)
	}

	// Write data.tar.gz
	if err := writeArEntry(f, "data.tar.gz", dataTar); err != nil {
		return fmt.Errorf("write data.tar.gz: %w", err)
	}

	return nil
}

// createControlTar creates the control.tar.gz containing package metadata.
func createControlTar(pkg debPackage) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Control file content
	control := fmt.Sprintf(`Package: %s
Version: %s
Section: utils
Priority: optional
Architecture: %s
Maintainer: %s
Description: %s
`, pkg.Name, pkg.Version, pkg.Arch, pkg.Maintainer, pkg.Description)

	header := &tar.Header{
		Name:    "./control",
		Size:    int64(len(control)),
		Mode:    0o644,
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(control)); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// createDataTar creates the data.tar.gz containing the binary.
func createDataTar(name string, binaryData []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Create directory entries
	dirs := []string{"./", "./usr/", "./usr/local/", "./usr/local/bin/"}
	for _, dir := range dirs {
		header := &tar.Header{
			Name:     dir,
			Mode:     0o755,
			Typeflag: tar.TypeDir,
			ModTime:  time.Now(),
		}
		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}
	}

	// Add binary
	header := &tar.Header{
		Name:    "./usr/local/bin/" + name,
		Size:    int64(len(binaryData)),
		Mode:    0o755,
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return nil, err
	}

	if _, err := tw.Write(binaryData); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// writeArEntry writes a single entry to an ar archive.
func writeArEntry(w io.Writer, name string, data []byte) error {
	// ar header format (60 bytes):
	// - name: 16 bytes (space padded)
	// - mtime: 12 bytes (decimal, space padded)
	// - owner: 6 bytes (decimal, space padded)
	// - group: 6 bytes (decimal, space padded)
	// - mode: 8 bytes (octal, space padded)
	// - size: 10 bytes (decimal, space padded)
	// - magic: 2 bytes ("`\n")

	// Pad name to 16 bytes
	paddedName := name
	if len(paddedName) < 16 {
		paddedName += strings.Repeat(" ", 16-len(paddedName))
	}

	mtime := fmt.Sprintf("%-12d", time.Now().Unix())
	owner := "0     "
	group := "0     "
	mode := "100644  "
	size := fmt.Sprintf("%-10d", len(data))
	magic := "`\n"

	header := paddedName + mtime + owner + group + mode + size + magic

	if _, err := w.Write([]byte(header)); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	// ar entries must be 2-byte aligned
	if len(data)%2 != 0 {
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

// For big-endian binary writing (unused but keeping for reference)
var _ = binary.BigEndian
