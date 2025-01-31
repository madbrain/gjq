package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"com.github/madbrain/gjq/eval"
	"github.com/mattn/go-colorable"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
)

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

	editor.BindKey(keys.CtrlI, eval.ExprCompletion{})

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Expecting JSON file")
		return
	}

	evaluator := eval.Evaluator{}
	if err := evaluator.ReadJsonFile(args[0]); err != nil {
		fmt.Println("Error reading JSON file")
		return
	}

	ctxt := context.WithValue(context.Background(), "evaluator", &evaluator)

	fmt.Println("GJQ Shell. Type Ctrl-D to quit.")
	for {
		text, err := editor.ReadLine(ctxt)

		if err != nil {
			//fmt.Printf("ERR=%s\n", err.Error())
			return
		}

		evaluator.Evaluate(text)

		history.Add(text)
	}
}
