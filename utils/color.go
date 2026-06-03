package utils

import (
	"github.com/fatih/color"
)

// Use a professional library for cross-platform color output
var (
	Red     = color.New(color.FgRed)                // Red 31
	Green   = color.New(color.FgGreen)              // Green 32
	Yellow  = color.New(color.FgYellow)             // Yellow 33
	Blue    = color.New(color.FgBlue, color.Bold)   // Blue 34
	Magenta = color.New(color.FgMagenta)            // Magenta 35
	Cyan    = color.New(color.FgHiCyan, color.Bold) // Cyan 36
	White   = color.New(color.FgWhite)              // White 37
)

