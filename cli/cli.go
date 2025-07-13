package cli

import "flag"

// Options holds CLI options.
type Options struct {
	Mode      string
	Token     string
	Pretty    bool
	ScoreOnly bool
}

// Parse parses command line flags.
func Parse() Options {
	mode := flag.String("mode", "multiscan", "scan mode")
	token := flag.String("token", "", "token address")
	pretty := flag.Bool("pretty", false, "pretty output")
	scoreOnly := flag.Bool("score-only", false, "print only risk score")
	flag.Parse()
	return Options{Mode: *mode, Token: *token, Pretty: *pretty, ScoreOnly: *scoreOnly}
}
