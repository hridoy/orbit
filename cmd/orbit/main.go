package main

import (
	"fmt"
	"os"
	"strings"

	"orbit/internal/buildinfo"
	"orbit/internal/cli"
)

func printHelp() {
	fmt.Printf("🪐 Orbit Sovereign CLI v%s (%s)\n", buildinfo.Version, buildinfo.GitCommit)
	fmt.Println("Usage: orbit <command> [arguments] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  push [dir]                Wrangler-style delta deploy to target application")
	fmt.Println("  status                    Display cluster, host, workload & database status")
	fmt.Println("  logs <app_id>             Tail or stream real-time logs for a workload")
	fmt.Println("  db create <name>          Provision a dedicated low-RAM native database")
	fmt.Println("  version                   Print agent and CLI version details")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --server=<url>            Orbit Agent server URL (default: http://127.0.0.1:8787)")
	fmt.Println("  --token=<token>           Ed25519 or Master Key (auto-discovered on local host)")
	fmt.Println("  --app=<app_id>            Target application slug")
	fmt.Println("  --follow, -f              Stream logs in real-time (SSE)")
	fmt.Println("  --engine=<engine>         Database engine: postgres or mongodb (default: postgres)")
	fmt.Println("  --profile=<profile>       RAM profile: micro, standard, or heavy (default: micro)")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	serverURL, token, appID, engine, profile, follow, positional := parseArgs(os.Args[2:])

	client := cli.NewClient(serverURL, token)

	switch command {
	case "version":
		fmt.Printf("🪐 Orbit CLI v%s (%s)\n", buildinfo.Version, buildinfo.GitCommit)
	case "status":
		if err := cli.RunStatus(client); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	case "push":
		dir := "."
		if len(positional) > 0 {
			dir = positional[0]
		}
		if appID == "" && len(positional) > 1 {
			appID = positional[1]
		}
		if err := cli.RunPush(client, dir, appID); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Push failed: %v\n", err)
			os.Exit(1)
		}
	case "logs":
		targetApp := appID
		if targetApp == "" && len(positional) > 0 {
			targetApp = positional[0]
		}
		if err := cli.RunLogs(client, targetApp, follow); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	case "db":
		if len(positional) < 2 || positional[0] != "create" {
			fmt.Println("Usage: orbit db create <name> [--engine=postgres|mongodb] [--profile=micro]")
			os.Exit(1)
		}
		dbName := positional[1]
		if err := cli.RunDBCreate(client, dbName, engine, profile); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	default:
		printHelp()
		os.Exit(1)
	}
}

func parseArgs(rawArgs []string) (serverURL, token, appID, engine, profile string, follow bool, positional []string) {
	for i := 0; i < len(rawArgs); i++ {
		arg := rawArgs[i]
		if strings.HasPrefix(arg, "--server=") {
			serverURL = strings.TrimPrefix(arg, "--server=")
		} else if arg == "--server" && i+1 < len(rawArgs) {
			serverURL = rawArgs[i+1]
			i++
		} else if strings.HasPrefix(arg, "--token=") {
			token = strings.TrimPrefix(arg, "--token=")
		} else if arg == "--token" && i+1 < len(rawArgs) {
			token = rawArgs[i+1]
			i++
		} else if strings.HasPrefix(arg, "--app=") {
			appID = strings.TrimPrefix(arg, "--app=")
		} else if arg == "--app" && i+1 < len(rawArgs) {
			appID = rawArgs[i+1]
			i++
		} else if strings.HasPrefix(arg, "--engine=") {
			engine = strings.TrimPrefix(arg, "--engine=")
		} else if arg == "--engine" && i+1 < len(rawArgs) {
			engine = rawArgs[i+1]
			i++
		} else if strings.HasPrefix(arg, "--profile=") {
			profile = strings.TrimPrefix(arg, "--profile=")
		} else if arg == "--profile" && i+1 < len(rawArgs) {
			profile = rawArgs[i+1]
			i++
		} else if arg == "--follow" || arg == "-f" {
			follow = true
		} else if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
		}
	}
	if engine == "" {
		engine = "postgres"
	}
	if profile == "" {
		profile = "micro"
	}
	return
}
