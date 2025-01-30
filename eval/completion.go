package eval

import (
	"context"

	rl "github.com/nyaosorg/go-readline-ny"
)

type ExprCompletion struct {
}

func (C ExprCompletion) String() string {
	return "EXPR_COMPLETION"
}

func (C ExprCompletion) Call(ctx context.Context, B *rl.Buffer) rl.Result {
	// B as content and position
	// TODO this code is responsible to modify buffer and print information
	// example CmdCompletionOrList { Completion: completion.File {} }
	B.Out.WriteByte('\a')
	B.Out.WriteByte('\n')
	B.Out.WriteString(("TODO Complete\n"))
	return rl.CONTINUE
}
