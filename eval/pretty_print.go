package eval

import "fmt"

const TAB_SIZE = 2

func PrettyPrintJson(v any) {
	prettyPrint(v, 0)
	fmt.Println()
}

func prettyPrint(v any, indent int) {
	switch t := v.(type) {
	case nil:
		fmt.Print("null")
	case int:
		fmt.Printf("%d", t)
	case string:
		fmt.Printf("\"%s\"", t)
	case []any:
		fmt.Printf("[")
		if len(t) > 0 {
			for i, value := range t {
				if i > 0 {
					fmt.Print(", ")
				}
				prettyPrint(value, indent+TAB_SIZE)
			}
		}
		fmt.Printf("]")
	case map[string]any:
		fmt.Printf("{")
		if len(t) > 0 {
			fmt.Println()
			i := 0
			for key, value := range t {
				if i > 0 {
					fmt.Printf(",\n")
				}
				printIndent(indent + TAB_SIZE)
				fmt.Printf("%s: ", key)
				prettyPrint(value, indent+TAB_SIZE)
				i++
			}
			fmt.Printf("\n")
			printIndent(indent)
		}
		fmt.Printf("}")
	default:
		fmt.Printf("@%+v@", v)
	}
}

func printIndent(amount int) {
	for i := 0; i < amount; i++ {
		fmt.Print(" ")
	}
}
