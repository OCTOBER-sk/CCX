package packaging

// Packaging per spec 62 (lines 2136-2163) and spec 54 (packaging, 1915-1941) and spec 170 (installers, 5130-5145):
// Targets: Linux x64/ARM64, macOS x64/ARM64, Windows x64/ARM64 (spec 62, 2136-2145).
// Primary artifact: single native binary (spec 62, 2140-2143).
// Distribution channels: GitHub Releases, Homebrew, winget (spec 62, 2146-2155).
// Every package requires: target platform, binary path, checksum, signature, SBOM reference (spec 54, 1915-1925).

type Package struct {
	Target        string // platform tuple: linux-x64, linux-arm64, darwin-x64, darwin-arm64, windows-x64, windows-arm64
	Binary        string // native binary file name
	ChecksumPath  string // SHA-256 checksum file path
	SignaturePath string // ed25519/minisign signature path
	SBOMPath      string // SPDX/CycloneDX SBOM reference
}

// Platforms returns all 6 required target platforms per spec 62.
func Platforms() []string {
	return []string{
		"linux-x64", "linux-arm64",
		"darwin-x64", "darwin-arm64",
		"windows-x64", "windows-arm64",
	}
}

// VerifyPackage confirms packaging requirements per spec 54 and spec 62.
func VerifyPackage(p Package) bool {
	if p.Target == "" || p.Binary == "" {
		return false
	}
	valid := false
	for _, plat := range Platforms() {
		if p.Target == plat {
			valid = true
			break
		}
	}
	if !valid {
		return false // spec 62: target must be one of 6 platforms
	}
	return p.ChecksumPath != "" && p.SignaturePath != "" && p.SBOMPath != ""
}
