package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	params := make([]string, 0)

	var (
		useGo         bool = false
		help          bool = false
		hasGoFlag     bool = false
		hasGoGrpcFlag bool = false
	)
	flag.BoolVar(&useGo, "go", false, "use default go flags: '--go_out=. --go-grpc_out=.'")
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&help, "help", false, "help")
	flag.Parse()

	if useGo {
		hasGoFlag = true
		hasGoGrpcFlag = true
	} else if help {
		params = append(params, "-h")
		hasGoFlag = true
		hasGoGrpcFlag = true
		fmt.Println()
		fmt.Println("Go-specific options:")
		fmt.Println("  --go_out=OUT_DIR            Generate Go source file.")
		fmt.Println("  --go-grpc_out=OUT_DIR       Generate Go gRPC source file.")
		fmt.Println(" or:")
		fmt.Println("  -go                         use '--go_out=.' and '--go-grpc_out=.' options")
		fmt.Println()
	}

	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "-") {
			if strings.HasPrefix(arg, "--go_out") {
				hasGoFlag = true
			} else if strings.HasPrefix(arg, "--go-grpc_out") {
				hasGoGrpcFlag = true
			} else if strings.EqualFold(arg, "--help") {
				hasGoFlag = true
				hasGoGrpcFlag = true
				fmt.Println()
				fmt.Println("Go-specific options:")
				fmt.Println("  --go_out=OUT_DIR            Generate Go source file.")
				fmt.Println("  --go-grpc_out=OUT_DIR       Generate Go gRPC source file.")
				fmt.Println()
			} else if strings.EqualFold(arg, "--version") {
				hasGoFlag = true
				hasGoGrpcFlag = true
			}
			params = append(params, arg)
		} else {
			matches, err := filepath.Glob(arg)
			if err != nil {
				fmt.Printf("ERR: %v\n", err)
				os.Exit(1)
			}
			params = append(params, matches...)
		}
	}

	if useGo {
		params = append(params, "--go_out=.", "--go-grpc_out=.")
	}

	if !hasGoFlag {
		fmt.Println("-> WARN: no --go_out=OUT_DIR flag specified, no Go code will be generated.")
	}
	if !hasGoGrpcFlag {
		fmt.Println("-> WARN: no --go-grpc_out=OUT_DIR flag specified, no Go gRPC code will be generated.")
	}
	cmd := exec.Command("protoc.exe", params...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Printf("error: %v", err)
	}
}
