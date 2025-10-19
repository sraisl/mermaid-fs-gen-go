package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateMermaidDiagram erzeugt ein Mermaid-Diagramm aus der Verzeichnisstruktur
func generateMermaidDiagram(directoryPath, outputFile string) (string, error) {
	var mermaidLines []string
	mermaidLines = append(mermaidLines, "```mermaid", "graph TD")

	// Rekursive Funktion zum Scannen der Verzeichnisse
	var scanDirectory func(path, prefix, parentID string) ([]string, error)
	scanDirectory = func(path, prefix, parentID string) ([]string, error) {
		var lines []string
		dirName := filepath.Base(path)
		currentID := fmt.Sprintf("node_%s", strings.ReplaceAll(strings.ReplaceAll(dirName, "/", "_"), ".", "_"))

		// Füge den aktuellen Knoten hinzu (nur für das Root-Verzeichnis)
		if parentID == "root" {
			lines = append(lines, fmt.Sprintf("    %s[%s]", currentID, dirName))
		}

		// Öffne das Verzeichnis
		dir, err := os.Open(path)
		if err != nil {
			lines = append(lines, fmt.Sprintf("    %s_perm[Permission Denied]", currentID))
			lines = append(lines, fmt.Sprintf("    %s --> %s_perm", currentID))
			return lines, nil // Fahre fort, trotz Fehler
		}
		defer dir.Close()

		// Lese Verzeichnisinhalte
		entries, err := dir.Readdir(-1)
		if err != nil {
			lines = append(lines, fmt.Sprintf("    %s_perm[Permission Denied]", currentID))
			lines = append(lines, fmt.Sprintf("    %s --> %s_perm", currentID))
			return lines, nil
		}

		// Sortiere Einträge für konsistente Ausgabe
		var dirs []os.FileInfo
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != ".git" {
				dirs = append(dirs, entry)
			}
		}

		// Verarbeite Unterverzeichnisse
		for _, entry := range dirs {
			subID := fmt.Sprintf("node_%s", strings.ReplaceAll(strings.ReplaceAll(entry.Name(), "/", "_"), ".", "_"))
			lines = append(lines, fmt.Sprintf("    %s[%s]", subID, entry.Name()))
			lines = append(lines, fmt.Sprintf("    %s --> %s", currentID, subID))

			// Rekursiver Aufruf für Unterverzeichnisse
			subPath := filepath.Join(path, entry.Name())
			subLines, err := scanDirectory(subPath, prefix+"  ", subID)
			if err != nil {
				return lines, err
			}
			lines = append(lines, subLines...)
		}

		return lines, nil
	}

	// Scanne das Verzeichnis
	lines, err := scanDirectory(directoryPath, "", "root")
	if err != nil {
		return "", err
	}
	mermaidLines = append(mermaidLines, lines...)
	mermaidLines = append(mermaidLines, "```")

	// Schreibe in die Ausgabedatei
	output := strings.Join(mermaidLines, "\n")
	err = os.WriteFile(outputFile, []byte(output), 0644)
	if err != nil {
		return "", err
	}

	return output, nil
}

func main() {
	// Hole Verzeichnis von Benutzer oder nutze aktuelles Verzeichnis
	fmt.Print("Enter directory path to scan (press Enter for current directory): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	directory := scanner.Text()
	if directory == "" {
		directory = "."
	}

	// Stelle sicher, dass der Pfad absolut ist
	absPath, err := filepath.Abs(directory)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Generiere das Diagramm
	diagram, err := generateMermaidDiagram(absPath, "directory_structure.md")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\nGenerated Mermaid diagram:")
	fmt.Println(diagram)
	fmt.Println("\nMermaid diagram has been saved to 'directory_structure.md'")
}
