package main

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"com.github/madbrain/gjq/lang"
	"github.com/mattn/go-colorable"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

func execute(text string) {
	reporter := &lang.DefaultReporter{}
	lexer := lang.NewLexer(text, reporter)
	parser := lang.NewParser(lexer, reporter)
	ast := parser.Parse()

	if reporter.HasErrors() {
		reporter.DisplayErrors(text)
		return
	}

	fmt.Println(lang.DisplayAst(ast))
}

func main() {
	history := simplehistory.New()

	editor := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "\x1B[36;22m> ")
		},
		Writer:  colorable.NewColorableStdout(),
		History: history,
		Highlight: []readline.Highlight{
			{Pattern: regexp.MustCompile("&"), Sequence: "\x1B[33;49;22m"},
			{Pattern: regexp.MustCompile(`"[^"]*"`), Sequence: "\x1B[35;49;22m"},
			{Pattern: regexp.MustCompile(`%[^%]*%`), Sequence: "\x1B[36;49;1m"},
			{Pattern: regexp.MustCompile("\u3000"), Sequence: "\x1B[37;41;22m"},
		},
		HistoryCycling: true,
		PredictColor:   [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
		ResetColor:     "\x1B[0m",
		DefaultColor:   "\x1B[33;49;1m",
	}

	editor.BindKey(keys.CtrlI, completion.CmdCompletionOrList{
		Completion: completion.File{},
		Postfix:    " ",
	})
	// If you do not want to list files with double-tab-key,
	// use `CmdCompletion` instead of `CmdCompletionOrList`

	fmt.Println("GJQ Shell. Type Ctrl-D to quit.")
	for {
		text, err := editor.ReadLine(context.Background())

		if err != nil {
			fmt.Printf("ERR=%s\n", err.Error())
			return
		}

		execute(text)

		history.Add(text)
	}
}
