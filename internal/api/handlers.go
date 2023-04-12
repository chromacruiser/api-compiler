package api

import (
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"github.com/chromacruiser/api-compiler/internal/avr"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed fixtures/*
var fixtures embed.FS

type Handlers struct {
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h Handlers) ExampleCompile(ctx echo.Context) error {
	// Create a temporary directory to store the project
	tempDir, err := os.MkdirTemp("", "project-*")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	defer os.RemoveAll(tempDir)

	// Copy fixture into temporary directory
	err = copyFiles("fixtures/basic", tempDir)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	elf, err := avr.CompileAtmega328P(filepath.Join(tempDir, "main.c"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	hex, err := avr.ConvertToHex(elf)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.File(hex)
}

func (h Handlers) Compile(ctx echo.Context) error {
	body := &CompileMultipartRequestBody{}

	err := ctx.Bind(body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "Bad request"})
	}

	// Read the project from the request body
	project, err := body.Project.Bytes()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to read project from request body"})
	}

	// Create a temporary directory to store the project
	tempDir, err := os.MkdirTemp("", "project-*")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	defer os.RemoveAll(tempDir)

	// Unzip the project into the temporary directory
	unzipProject(project, tempDir)

	elf, err := avr.CompileAtmega328P(filepath.Join(tempDir, "main.c"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	hex, err := avr.ConvertToHex(elf)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.File(hex)
}


func copyFiles(source string, destination string) error {
	err := fs.WalkDir(fixtures, source, func(path string, d fs.DirEntry, err error) error {
		writePath := strings.TrimLeft(strings.TrimPrefix(path, source), "/")

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}
		if d.IsDir() {
			// Create the directory in the temporary directory
			err := os.MkdirAll(filepath.Join(destination, writePath), d.Type())
			if err != nil {
				return fmt.Errorf("error creating directory: %w", err)
			}
		} else {
			// Read the contents of the file from embed.FS and write them to a file in the temporary directory
			fileContents, err := fixtures.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}
			err = os.WriteFile(filepath.Join(destination, writePath), fileContents, d.Type())
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error copying files: %w", err)
	}

	return nil
}

func unzipProject(project []byte, destination string) {
	reader, _ := zip.NewReader(bytes.NewReader(project), -1)

	for _, file := range reader.File {
		path := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, _ := file.Open()
		defer fileReader.Close()

		targetFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		defer targetFile.Close()

		io.Copy(targetFile, fileReader)
	}
}
