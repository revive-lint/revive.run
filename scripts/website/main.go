// Program website generates the content of https://revive.run/ from the documentation of the revive-lint/revive repository.
//
// Run it from the repository root, with a checkout of revive-lint/revive in the directory given by -src:
//
//	go run ./scripts/website -src revive
//
// It reads README.md and RULES_DESCRIPTIONS.md from the checkout, writes Hugo content files under content/,
// and copies the files from assets/ to static/images/.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	src := flag.String("src", "revive", "path to a checkout of github.com/revive-lint/revive")
	flag.Parse()
	if err := run(*src); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(src string) error {
	rulesDoc, err := readTextFile(filepath.Join(src, "RULES_DESCRIPTIONS.md"))
	if err != nil {
		return err
	}

	sections := splitSections(rulesDoc)
	if len(sections) == 0 {
		return errors.New("no \"## \" sections found in RULES_DESCRIPTIONS.md")
	}
	seen := make(map[string]bool, len(sections))
	for _, sec := range sections {
		if seen[sec.name] {
			return fmt.Errorf("duplicate section %q in RULES_DESCRIPTIONS.md", sec.name)
		}
		seen[sec.name] = true
	}

	// All rules live on the single /r/ page, mirroring RULES_DESCRIPTIONS.md:
	// its "## " anchors keep legacy /r#<rule-name> links working natively.
	// The page lives under content/docs/ to be part of the docs sidebar.
	rules := frontMatter{
		title:       "Rules",
		description: "List of all available revive rules.",
		url:         "/r/",
		weight:      100,
		icon:        "rule",
	}
	if err := writePage(filepath.Join("content", "docs", "rules.md"), rules, rulesPageContent(sections)); err != nil {
		return err
	}

	readme, err := readTextFile(filepath.Join(src, "README.md"))
	if err != nil {
		return err
	}
	// A regular page served at /docs/ (see content/docs/_index.md) so that the
	// theme shows it in the sidebar and renders its table of contents.
	docs := frontMatter{
		title:       "Documentation",
		description: "Installation, usage, configuration, and CI integrations.",
		url:         "/docs/",
		weight:      1,
		icon:        "menu_book",
	}
	if err := writePage(filepath.Join("content", "docs", "documentation.md"), docs, transformReadme(readme)); err != nil {
		return err
	}

	copied, err := copyAssets(filepath.Join(src, "assets"), filepath.Join("static", "images"))
	if err != nil {
		return err
	}

	fmt.Printf("Generated docs/rules.md (%d sections) and docs/documentation.md; copied %d assets.\n", len(sections), copied)
	return nil
}

// resetDir recreates dir empty, so reruns never leave stale generated files behind.
func resetDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("removing %s: %w", dir, err)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	return nil
}

// readTextFile reads a file and normalizes its line endings to "\n".
func readTextFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

// writePage writes a Hugo content file: YAML front matter followed by body.
func writePage(path string, fm frontMatter, body string) error {
	if err := os.WriteFile(path, []byte(pageContent(fm, body)), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// copyAssets copies every regular file from srcDir to dstDir and returns the
// number of files copied.
func copyAssets(srcDir, dstDir string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := resetDir(dstDir); err != nil {
		return 0, err
	}

	copied := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, entry.Name()))
		if err != nil {
			return copied, err
		}
		if err := os.WriteFile(filepath.Join(dstDir, entry.Name()), data, 0o644); err != nil {
			return copied, fmt.Errorf("copying asset %s: %w", entry.Name(), err)
		}
		copied++
	}
	return copied, nil
}
