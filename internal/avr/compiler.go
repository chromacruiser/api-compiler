package avr

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// ConvertToHex takes the given elf file and converts it to a hex file`
func ConvertToHex(elf string) (string, error) {
	command := []string{
		"avr-objcopy", "-O", "ihex", elf, filepath.Join(filepath.Dir(elf), "project.hex"),
	}

	fmt.Println("command: ", command)

	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to convert to hex: %w", err)
	}

	fmt.Println(output)

	return filepath.Join(filepath.Dir(elf), "project.hex"), nil
}

// CompileAtmega328P compiles the given entrypoint to the elf file
func CompileAtmega328P(entrypoint string) (string, error) {
	command := []string{
		"avr-gcc", "-mmcu=atmega328p", "-o", filepath.Join(filepath.Dir(entrypoint), "project.elf"), entrypoint,
	}

	fmt.Println("command: ", command)

	cmd := exec.Command(command[0], command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to compile: %w", err)
	}

	fmt.Println(output)

	if err != nil {
		return "", fmt.Errorf("failed to compile: %w", err)
	}

	return filepath.Join(filepath.Dir(entrypoint), "project.elf"), nil
}
