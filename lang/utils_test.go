package lang

type Message struct {
	span    Span
	message string
}

type TestReporter struct {
	messages []Message
}

func (reporter *TestReporter) Report(span Span, message string) {
	reporter.messages = append(reporter.messages, Message{span: span, message: message})
}
