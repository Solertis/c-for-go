package generator

import (
	"fmt"
	"io"

	tl "github.com/xlab/cgogen/translator"
)

func (gen *Generator) writeDefinesGroup(wr io.Writer, cdecl []tl.CDecl) {
	if len(cdecl) == 0 {
		return
	}
	writeStartConst(wr)
	for _, decl := range cdecl {
		if !decl.IsDefine {
			continue
		} else if len(decl.Expression) == 0 {
			continue
		}
		name := gen.tr.TransformName(tl.TargetDefine, decl.Name)
		fmt.Fprintf(wr, "// %s as defined in %s\n", name, tl.SrcLocation(decl.Pos))
		fmt.Fprintf(wr, "%s = %s\n", name, decl.Expression)
	}
	writeEndConst(wr)
}

func (gen *Generator) writeEnum(wr io.Writer, enumDecl tl.CDecl) {
	spec, ok := enumDecl.Spec.(*tl.CEnumSpec)
	if !ok {
		return
	}
	var typeName []byte
	if len(enumDecl.Name) > 0 {
		typeName = gen.tr.TransformName(tl.TargetTypedef, enumDecl.Name)
	} else if len(spec.Tag) > 0 {
		typeName = gen.tr.TransformName(tl.TargetTag, spec.Tag)
	}

	fmt.Fprintf(wr, "// %s as declared in %s\n", typeName, tl.SrcLocation(enumDecl.Pos))
	enumType, err := gen.tr.TranslateSpec(spec.Type)
	if err != nil {
		writeError(wr, err)
		return
	}

	writeType(wr, typeName, enumType)
	writeSpace(wr, 1)
	fmt.Fprintf(wr, "// %s enumeration from %s\n", typeName, tl.SrcLocation(enumDecl.Pos))
	writeStartConst(wr)
	first := true
	for _, en := range spec.Enumerators {
		enName := gen.tr.TransformName(tl.TargetDeclare, en.Name)
		if first {
			first = false
			fmt.Fprintf(wr, "%s %s = %s\n", enName, typeName, en.Expression)
			continue
		}
		fmt.Fprintf(wr, "%s = %s\n", enName, en.Expression)
	}
	writeEndConst(wr)
	writeSpace(wr, 1)
}

func writeStartConst(wr io.Writer) {
	fmt.Fprintln(wr, "const (")
}

func writeEndConst(wr io.Writer) {
	fmt.Fprintln(wr, ")")
}