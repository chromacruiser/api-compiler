# API Compiler

API Compiler is an open-source HTTP service that exposes an endpoint for AVR code to be compiled. It is a tool for
developers who are working on AVR-based microcontroller projects and need an easy way to compile their code.

## Features

* Exposes a simple RESTful API for AVR code compilation
* Supports the AVR-GCC compiler for compiling C and C++ code
* Accepts a ZIP file containing the AVR project to be compiled
* Returns the compiled binary as an Intel HEX file
* Provides an example endpoint (`/example/compile`) that returns an example Intel HEX file for testing purposes

## Usage

To use API Compiler, follow these simple steps:

1. Clone the repository or download the source code.
2. Run the API server using the command `go run main.go`. The server will be available on port `8080`.
3. Send a POST request to the `/compile` endpoint with a ZIP file containing the AVR project to be compiled.
4. The compiled binary as an Intel HEX file will be returned as the response to the API request.

## License

API Compiler is released under the MIT license. See [LICENSE](./LICENSE) for more details.
