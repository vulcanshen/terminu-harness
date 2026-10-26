// Command terminu-harness is the harness's MCP server: the terminal the agent
// runs the app it writes in, and tdp's files.
package main

import (
	"fmt"
	"os"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/vulcanshen/terminu-harness/internal/assets"
	"github.com/vulcanshen/terminu-harness/internal/server"
)

// version is set when a release is built; a local build is dev.
var version = "dev"

const usage = `terminu-harness: the terminu harness's MCP server

Usage:
  terminu-harness mcp           serve MCP on stdin and stdout (Claude Code starts it)
  terminu-harness check [repo]  check the app's harness-assets/ (the current directory when left out)
  terminu-harness version       print the version (-v, --version)
  terminu-harness help          print this help (-h, --help)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "mcp":
		sv := server.New()
		err := mcpserver.ServeStdio(sv.MCP(version))
		sv.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "check":
		repo := "."
		if len(os.Args) > 2 {
			repo = os.Args[2]
		}
		r := assets.Check(repo, assets.Default(server.Source()))
		fmt.Print(r.Text())
		if len(r.Problems) > 0 {
			os.Exit(1)
		}
	case "version", "-v", "--version":
		fmt.Println(version)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}
