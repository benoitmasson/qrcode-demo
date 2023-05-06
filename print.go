package main

import (
	"fmt"
	"strings"

	"gocv.io/x/gocv"
)

func printQRCode(qrcode [][]bool) {
	if len(qrcode) == 0 {
		return
	}

	fmt.Println("\033[7;m", strings.Repeat(" ", 2*len(qrcode[0])+3)) // turn on inverse mode, start with blank line

	for i := 0; i < len(qrcode); i++ {
		fmt.Print("  ") // start line with blank characters
		for j := 0; j < len(qrcode[i]); j++ {
			val := qrcode[i][j]
			char := " " // blank, displayed white
			if val {
				char = "█" // filled, displayed black
			}
			fmt.Print(char, char) // double print to achieve 1:1 scale
		}
		fmt.Println("  ") // end line with blank characters
	}

	fmt.Println(strings.Repeat(" ", 2*len(qrcode[0])+3), "\033[0;m") // end with blank line, turn off inverse mode
}

func printQRCodeMat(qrcode gocv.Mat) {
	fmt.Println("\033[7;m", strings.Repeat(" ", 2*qrcode.Cols()+3)) // turn on inverse mode, start with blank line

	for i := 0; i < qrcode.Rows(); i++ {
		fmt.Print("  ") // start line with blank characters
		for j := 0; j < qrcode.Cols(); j++ {
			val := qrcode.GetUCharAt(i, j)
			char := " " // blank, displayed white
			if val == 0 {
				char = "█" // filled, displayed black
			}
			fmt.Print(char, char) // double print to achieve 1:1 scale
		}
		fmt.Println("  ") // end line with blank characters
	}

	fmt.Println(strings.Repeat(" ", 2*qrcode.Cols()+3), "\033[0;m") // end with blank line, turn off inverse mode
}
