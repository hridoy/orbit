// Command releasesign creates and signs agent release manifests.
//
//	releasesign genkey                      print a new Ed25519 seed (secret) and public key (hex)
//	releasesign manifest -dir D -repo R -tag T -version V
//	                                        write D/manifest.json for orbit-agent-linux-* binaries
//	releasesign sign -dir D                 write D/manifest.sig, key from $ORBIT_RELEASE_KEY
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type asset struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

type manifest struct {
	Version string           `json:"version"`
	Assets  map[string]asset `json:"assets"`
}

func main() {
	if len(os.Args) < 2 {
		fail("usage: releasesign genkey|manifest|sign [flags]")
	}
	switch os.Args[1] {
	case "genkey":
		genkey()
	case "manifest":
		fs := flag.NewFlagSet("manifest", flag.ExitOnError)
		dir := fs.String("dir", "dist", "directory with the binaries")
		repo := fs.String("repo", "", "owner/name")
		tag := fs.String("tag", "", "release tag, e.g. v0.2.0")
		version := fs.String("version", "", "version without v prefix")
		_ = fs.Parse(os.Args[2:])
		if *repo == "" || *tag == "" || *version == "" {
			fail("manifest needs -repo, -tag, and -version")
		}
		writeManifest(*dir, *repo, *tag, *version)
	case "sign":
		fs := flag.NewFlagSet("sign", flag.ExitOnError)
		dir := fs.String("dir", "dist", "directory containing manifest.json")
		_ = fs.Parse(os.Args[2:])
		sign(*dir)
	default:
		fail("unknown command " + os.Args[1])
	}
}

func genkey() {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	check(err)
	fmt.Println("ORBIT_RELEASE_KEY (store as a GitHub Actions secret):", hex.EncodeToString(priv.Seed()))
	fmt.Println("ORBIT_RELEASE_PUBKEY (store as a GitHub Actions variable):", hex.EncodeToString(pub))
}

func writeManifest(dir, repo, tag, version string) {
	files, err := filepath.Glob(filepath.Join(dir, "orbit-agent-linux-*"))
	check(err)
	if len(files) == 0 {
		fail("no orbit-agent-linux-* binaries in " + dir)
	}
	sort.Strings(files)
	m := manifest{Version: version, Assets: map[string]asset{}}
	for _, f := range files {
		data, err := os.ReadFile(f)
		check(err)
		sum := sha256.Sum256(data)
		name := filepath.Base(f)
		arch := strings.TrimPrefix(name, "orbit-agent-linux-")
		m.Assets["linux/"+arch] = asset{
			URL:    fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, tag, name),
			SHA256: hex.EncodeToString(sum[:]),
		}
	}
	out, err := json.MarshalIndent(m, "", "  ")
	check(err)
	check(os.WriteFile(filepath.Join(dir, "manifest.json"), append(out, '\n'), 0644))
}

func sign(dir string) {
	seedHex := os.Getenv("ORBIT_RELEASE_KEY")
	seed, err := hex.DecodeString(strings.TrimSpace(seedHex))
	if err != nil || len(seed) != ed25519.SeedSize {
		fail("ORBIT_RELEASE_KEY must be a 32-byte hex seed")
	}
	raw, err := os.ReadFile(filepath.Join(dir, "manifest.json"))
	check(err)
	sig := ed25519.Sign(ed25519.NewKeyFromSeed(seed), raw)
	check(os.WriteFile(filepath.Join(dir, "manifest.sig"), []byte(hex.EncodeToString(sig)+"\n"), 0644))
}

func check(err error) {
	if err != nil {
		fail(err.Error())
	}
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "releasesign:", msg)
	os.Exit(1)
}
