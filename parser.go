package migration

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func ScanCodebase(root string) ([]StructInfo, []TypeAlias, error) {
	var structs []StructInfo
	var aliases []TypeAlias

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fileStructs, fileAliases, err := parseFile(path, content)

			structs = append(structs, fileStructs...)
			aliases = append(aliases, fileAliases...)
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return structs, aliases, nil
}

func parseFile(filePath string, fileContent []byte) ([]StructInfo, []TypeAlias, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, fileContent, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var structs []StructInfo
	var aliases []TypeAlias

	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			switch typeSpec := spec.(type) {
			case *ast.TypeSpec:
				if aliasType, isAlias := typeSpec.Type.(*ast.Ident); isAlias {
					aliases = append(aliases, TypeAlias{
						Name:           typeSpec.Name.Name,
						UnderlyingType: aliasType.Name,
					})
				} else if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
					if hasTableTag(genDecl.Doc) {
						structInfo := processStruct(typeSpec.Name.Name, structType)
						structs = append(structs, structInfo)
					}
				}
			}
		}

		return true
	})

	return structs, aliases, nil
}
func hasTableTag(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "go:generate") && strings.Contains(comment.Text, "table") {
			return true
		}
	}
	return false
}

func processStruct(name string, structType *ast.StructType) StructInfo {
	structInfo := StructInfo{Name: name}
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldType := getFieldType(field.Type)
		fieldInfo := FieldInfo{
			Name: field.Names[0].Name,
			Type: fieldType,
		}
		if field.Tag != nil {
			fieldInfo.Tag = field.Tag.Value
		}
		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}
	return structInfo
}
func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	case *ast.StarExpr:
		return getFieldType(t.X) // Recursively get the underlying type for pointers
	case *ast.ArrayType:
		return "[]" + getFieldType(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getFieldType(t.Key), getFieldType(t.Value))
	default:
		return fmt.Sprintf("%T", expr) // Fallback for unknown types
	}
}
