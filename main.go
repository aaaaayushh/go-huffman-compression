package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-pwd-cracker/huff"
	"io"
	"os"
	"strconv"
)

func main() {
	// declare a flag variable to accept file name as input
	filePath := flag.String("path", "", "path to file to be compressed")
	flag.Parse()
	if *filePath == "" {
		panic("File path is required")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// declare a map to store the frequency of each character in the file
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

	// build a huffman tree with the freqMap
	huffTree := huff.BuildHuffmanTree(freqMap)
	prefixTable := BuildPrefixTable(huffTree.Root())
	for char, code := range prefixTable {
		fmt.Printf("Character %c: Code %s\n", char, code)
	}
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
