package main

import (
	"bufio"
	"log"
	"os"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.design/x/clipboard"
)

type Game struct {
	file1Lines     []string
	file2Lines     []string
	discrepancies  []string
	compareClicked bool
	copyClicked    bool
	file1Name      string
	file2Name      string
}

func readLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	return lines
}

func (g *Game) compareFiles() {
	// Reset discrepancies
	g.discrepancies = []string{}

	// Create maps to track unique lines
	linesInFile1 := make(map[string]bool)
	linesInFile2 := make(map[string]bool)

	for _, line := range g.file1Lines {
		linesInFile1[line] = true
	}
	for _, line := range g.file2Lines {
		linesInFile2[line] = true
	}

	// Add lines unique to file1 with file name
	for line := range linesInFile1 {
		if !linesInFile2[line] && line != "" {
			g.discrepancies = append(g.discrepancies, g.file1Name+": "+line)
		}
	}

	// Add lines unique to file2 with file name
	for line := range linesInFile2 {
		if !linesInFile1[line] && line != "" {
			g.discrepancies = append(g.discrepancies, g.file2Name+": "+line)
		}
	}

	// Update window size based on number of discrepancies
	// Each discrepancy will take up 20px in height
	ebiten.SetWindowSize(800, 60+len(g.discrepancies)*20)
}

func (g *Game) copyToClipboard() {
	// Combine discrepancies into a single string
	result := ""
	for _, line := range g.discrepancies {
		result += line + "\n"
	}
	// Copy to clipboard
	clipboard.Write(clipboard.FmtText, []byte(result))
}

func (g *Game) Update() error {
	// Check for mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Check if the Compare button is clicked
		if x >= 10 && x <= 110 && y >= 10 && y <= 40 {
			g.compareClicked = true
			g.compareFiles()
		}

		// Check if the Copy button is clicked
		if x >= 120 && x <= 220 && y >= 10 && y <= 40 {
			g.copyClicked = true
			g.copyToClipboard()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw Compare button
	ebitenutil.DrawRect(screen, 10, 10, 100, 30, color.RGBA{0, 128, 255, 255})
	ebitenutil.DebugPrintAt(screen, "Compare", 30, 20)

	// Draw Copy button
	ebitenutil.DrawRect(screen, 120, 10, 100, 30, color.RGBA{0, 200, 100, 255})
	ebitenutil.DebugPrintAt(screen, "Copy", 140, 20)

	// Draw discrepancies
	yOffset := 60
	if g.compareClicked {
		for i, line := range g.discrepancies {
			ebitenutil.DebugPrintAt(screen, line, 10, yOffset+(20*i))
		}
	}

	// Confirmation for Copy button
	if g.copyClicked {
		ebitenutil.DebugPrintAt(screen, "Copied to clipboard!", 10, 50)
		g.copyClicked = false // Reset copyClicked after showing the message
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Keep the window width fixed and adjust the height dynamically
	return 800, 60 + len(g.discrepancies)*20
}

func main() {
	// Initialize clipboard
	err := clipboard.Init()
	if err != nil {
		log.Fatalf("Failed to initialize clipboard: %s", err)
	}

	// Define filenames
	file1Name := "file1.txt" // Change this to your desired filename
	file2Name := "file2.txt"   // Change this to your desired filename

	// Create game instance with custom file names
	game := &Game{
		file1Lines: readLines(file1Name),
		file2Lines: readLines(file2Name),
		file1Name:  file1Name,
		file2Name:  file2Name,
	}
	ebiten.SetWindowSize(800, 600) // Initial window size
	ebiten.SetWindowTitle("Unique Lines Comparison")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
