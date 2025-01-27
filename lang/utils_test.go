package lang

import "fmt"

type Message struct {
	span    Span
	message string
}

type TestReporter struct {
	errors []Message
}

func (reporter *TestReporter) Report(span Span, message string) {
	reporter.errors = append(reporter.errors, Message{span: span, message: message})
}

func (reporter TestReporter) Display(content string) {
	for _, error := range reporter.errors {
		fmt.Println(content)
		for i := 0; i < error.span.start; i += 1 {
			fmt.Print(" ")
		}
		for i := 0; i < max(error.span.Length(), 1); i += 1 {
			fmt.Print("^")
		}
		fmt.Println()
		fmt.Printf("%d:%d: %s\n", error.span.start, error.span.end, error.message)
	}
}
