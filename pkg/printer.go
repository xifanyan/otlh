package otlh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Printer struct {
}

type JSONPrinterBuilder struct {
	*JSONPrinter
}

type JSONPrinter struct {
	indent string
	*Printer
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) JSON() *JSONPrinterBuilder {
	return &JSONPrinterBuilder{&JSONPrinter{indent: "  ", Printer: p}}
}

func (b *JSONPrinterBuilder) WithIndent(size int) *JSONPrinterBuilder {
	b.indent = strings.Repeat(" ", size)
	return b
}

func (b *JSONPrinterBuilder) Build() *JSONPrinter {
	return b.JSONPrinter
}

func (jb *JSONPrinter) Print(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, b, "", jb.indent)
	fmt.Println(buf.String())

	return nil
}
