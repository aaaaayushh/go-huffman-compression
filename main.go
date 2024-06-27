package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-pwd-cracker/huff"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func buildFreqMap(file *os.File) map[rune]int {
	freqMap := make(map[rune]int)
	reader := bufio.NewReader(file)
	for {
		char, _, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		freqMap[char]++
	}
	return freqMap
}

func writePrefixTableToOutputFile(outputFile *os.File, prefixTable map[rune]string) {
	// Create a writer
	writer := bufio.NewWriter(outputFile)

	for char, prefix := range prefixTable {
		// Convert the character to its Unicode code point
		codePoint := int(char)

		// Write the code point, prefix, and a delimiter
		_, err := fmt.Fprintf(writer, "%d\t%s\n", codePoint, prefix)
		if err != nil {
			panic(err)
		}
	}
	_, err := writer.WriteString("***HEADER*END***\n")
	if err != nil {
		panic(err)
	}

	// Flush the writer to ensure all buffered operations have been applied to the underlying writer
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func main() {
	// declare a flag variable to accept file name as input
	filePath := flag.String("path", "", "path to file to be compressed")
	outputPath := flag.String("output", "", "output file path")
	flag.Parse()

	if *filePath == "" || *outputPath == "" {
		panic("File path is required")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer CloseFile(file)

	// declare a map to store the frequency of each character in the file
	freqMap := buildFreqMap(file)
	//
	//for char, freq := range freqMap {
	//	fmt.Println("Char:", char, "\tFrequency:", freq)
	//}

	// build a huffman tree with the freqMap
	huffTree := huff.BuildHuffmanTree(freqMap)
	prefixTable := BuildPrefixTable(huffTree.Root())

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		panic(err)
	}
	CloseFile(outputFile)

	// Open the file in append mode
	outputFile, err = os.OpenFile(*outputPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer CloseFile(outputFile)

	writePrefixTableToOutputFile(outputFile, prefixTable)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(outputFile)

	var bitBuffer uint8
	var bitCount uint8

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		bitString := prefixTable[char]
		for _, bit := range bitString {
			bitBuffer = (bitBuffer << 1) | uint8(bit-'0')
			bitCount++
			if bitCount == 8 {
				err := writer.WriteByte(bitBuffer)
				if err != nil {
					panic(err)
				}
				bitBuffer = 0
				bitCount = 0
			}
		}
	}
	if bitCount > 0 {
		bitBuffer <<= 8 - bitCount
		err := writer.WriteByte(bitBuffer)
		if err != nil {
			panic(err)
		}
	}
	err = writer.WriteByte(bitCount)
	if err != nil {
		panic(err)
	}

	// Flush the writer to ensure all buffered operations have been applied to the underlying writer
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
	DecodeFile(*outputPath)
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(err)
	}
}

func DecodeFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer CloseFile(file)
	reader := bufio.NewReader(file)

	prefixTable := make(map[string]rune)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		line = strings.TrimSpace(line)
		if line == "***HEADER*END***" {
			break
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		codePoint, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		char := rune(codePoint)
		prefix := parts[1]
		prefixTable[prefix] = char
	}

	var decodedData strings.Builder
	var currentPrefix strings.Builder

	// Read the compressed data
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// Process each bit in the byte
		for i := 7; i >= 0; i-- {
			bit := (b >> uint(i)) & 1
			currentPrefix.WriteByte('0' + bit)

			if char, ok := prefixTable[currentPrefix.String()]; ok {
				decodedData.WriteRune(char)
				currentPrefix.Reset()
			}
		}
	}

	// The last byte contains the number of valid bits in the previous byte
	validBits, err := reader.ReadByte()
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Create a new file for writing the decoded text
	decodedFilePath := filepath.Join(filepath.Dir(filePath), "decoded_"+filepath.Base(filePath))
	decodedFile, err := os.Create(decodedFilePath)
	if err != nil {
		panic(err)
	}
	defer CloseFile(decodedFile)

	// Write the decoded text to the file
	_, err = decodedFile.WriteString(decodedData.String()[:len(decodedData.String())-int(8-validBits)])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded text has been written to: %s\n", decodedFilePath)

}

func BuildPrefixTable(root huff.BaseNode) map[rune]string {
	prefixTable := make(map[rune]string)
	buildPrefixTableHelper(root, "", prefixTable)
	return prefixTable
}

// recursive function to assign huffman codes to each letter
func buildPrefixTableHelper(node huff.BaseNode, currentPrefix string, prefixTable map[rune]string) {
	/*
		 A type switch matches the dynamic type of the interface value 'x'. The dynamic type is matched against the types in
		'switch' cases. If a short variable assignment of the form 'v := x.(type)' is used as the switch guard and a switch
		case is used for a single type only, 'v' will have the type specified in the matching switch case.
	*/
	switch n := node.(type) {
	case *huff.LeafNode:
		prefixTable[n.Value()] = currentPrefix
	case *huff.InternalNode:
		buildPrefixTableHelper(n.Left(), currentPrefix+strconv.Itoa(n.LeftEdge), prefixTable)
		buildPrefixTableHelper(n.Right(), currentPrefix+strconv.Itoa(n.RightEdge), prefixTable)
	}
}
