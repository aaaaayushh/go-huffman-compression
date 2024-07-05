# Huffman Compression Tool

This is a Go-based implementation of the Huffman coding algorithm for file compression and decompression.

## Overview

This tool uses Huffman coding to compress and decompress files. Huffman coding is a lossless data compression algorithm that assigns variable-length codes to characters based on their frequency of occurrence. More frequent characters are assigned shorter codes, resulting in overall compression of the data.

## Features

- File compression using Huffman coding
- File decompression
- Command-line interface for easy use
- Efficient handling of large files

## Installation

To use this tool, you need to have Go installed on your system. If you don't have Go installed, you can download it from [golang.org](https://golang.org/).

Once Go is installed, you can clone this repository:

```
git clone https://github.com/yourusername/go-pwd-cracker.git
cd go-pwd-cracker
```

## Usage

### Compressing a file

To compress a file, use the following command:

```
go run main.go -path /path/to/input/file -output /path/to/output/file
```

Replace `/path/to/input/file` with the path to the file you want to compress, and `/path/to/output/file` with the desired path for the compressed file.

### Decompressing a file

The decompression is automatically performed after compression. The decompressed file will be created in the same directory as the compressed file, with the prefix "decoded_" added to the filename.

## How it works

1. The program reads the input file and builds a frequency map of characters.
2. It constructs a Huffman tree based on these frequencies.
3. The Huffman tree is used to generate a prefix table, mapping characters to their Huffman codes.
4. The prefix table is written to the output file as a header.
5. The input file is then compressed using the Huffman codes and written to the output file.
6. After compression, the program automatically decompresses the file to verify the process.

## Project Structure

- `main.go`: Contains the main logic for file I/O, compression, and decompression.
- `huff/huff.go`: Implements the Huffman tree data structure and related algorithms.
