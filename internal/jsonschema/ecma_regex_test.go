package jsonschema

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ecmaCases = []struct {
	re2  string
	ecma string
}{
	{`^[a-zA-Z0-9]{27}$`, `^[0-9A-Za-z]{27}$`},
	{`^$|^[a-z]+$`, `^$|^[a-z]+$`},
	{`^a.b$`, `^a[^\n]b$`},
	{`(?s)^a.b$`, `^a[\s\S]b$`},
	{`^(?i)ab$`, `^[Aa][Bb]$`},
	{`^\d+\.\d+$`, `^[0-9]+\.[0-9]+$`},
	{`^\w+$`, `^[0-9A-Z_a-z]+$`},
	{`^\s*$`, `^[\u0009\u000A\u000C\u000D ]*$`},
	{`^[[:alpha:]]+$`, `^[A-Za-z]+$`},
	{`\Aabc\z`, `^abc$`},
	{`^(?P<id>[a-z]+)$`, `^(?:[a-z]+)$`},
	{`^\Qa.b\E$`, `^a\.b$`},
	{`^[^:]{1,10}$`, `^[^:]{1,10}$`},
	{`^(ab)+?$`, `^(?:ab)+?$`},
	{`^a{2,}$`, `^a{2,}$`},
	{`^/x/$`, `^\/x\/$`},
	{`^é$`, `^\u00E9$`},
	{`^[a-]$`, `^[\-a]$`},
	{`^\bword\b$`, `^\bword\b$`},
	{`^(?:a|b)c$`, `^[ab]c$`},
}

var ecmaUntranslatable = []string{
	`(?m)^a$`,      // multiline anchors cannot be expressed inside a pattern
	"^\U0001F600$", // astral literal has no escape valid with and without "u"
	`(`,            // invalid RE2
}

func TestECMAPattern(t *testing.T) {
	for _, c := range ecmaCases {
		got, ok := ecmaPattern(c.re2)
		require.True(t, ok, c.re2)
		assert.Equal(t, c.ecma, got, c.re2)
	}
	for _, p := range ecmaUntranslatable {
		_, ok := ecmaPattern(p)
		assert.False(t, ok, p)
	}
}

// TestECMAPatternMatchesRE2 compiles every translated pattern with an ECMA-262
// engine (node), with and without the "u" flag, and checks that it accepts and
// rejects the same sample strings as Go's RE2 implementation.
func TestECMAPatternMatchesRE2(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not installed")
	}
	samples := []string{
		"", "a", "ab", "AB", "aB", "a\nb", "a\rb", "a b", "a.b", "axb", "1.5", "12",
		"abcdefghijklmnopqrstuvwxyz0", "word", "é", "-", "x:y", "a:", "ababab", "aa", "aaa",
		"/x/", "ac", "bc", "cc", " \t", " ", "abc", "id", "a_b",
	}
	type check struct {
		ECMA    string   `json:"ecma"`
		Samples []string `json:"samples"`
	}
	var checks []check
	var want [][]bool
	for _, c := range ecmaCases {
		got, _ := ecmaPattern(c.re2)
		re := regexp.MustCompile(c.re2)
		row := make([]bool, len(samples))
		for i, s := range samples {
			row[i] = re.MatchString(s)
		}
		checks = append(checks, check{ECMA: got, Samples: samples})
		want = append(want, row)
	}
	input, err := json.Marshal(checks)
	require.NoError(t, err)
	script := `
const checks = JSON.parse(require("fs").readFileSync(0, "utf8"));
const out = checks.map(c => {
  const u = new RegExp(c.ecma, "u");
  new RegExp(c.ecma);
  return c.samples.map(s => u.test(s));
});
process.stdout.write(JSON.stringify(out));
`
	cmd := exec.Command(node, "-e", script)
	cmd.Stdin = bytesReader(input)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	require.NoError(t, err, "every translated pattern must compile as ECMA-262")
	var got [][]bool
	require.NoError(t, json.Unmarshal(out, &got))
	for i, c := range ecmaCases {
		for j, s := range samples {
			assert.Equal(t, want[i][j], got[i][j], "RE2 %q vs ECMA %q on %q", c.re2, checks[i].ECMA, s)
		}
	}
}

func bytesReader(b []byte) *os.File {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	go func() {
		_, _ = w.Write(b)
		_ = w.Close()
	}()
	return r
}
