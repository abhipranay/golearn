package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
)

const ffStructName = "FeatureFlags"

func main() {
	traverse("./cmd/ast/feature_flag/ff.go")
}

func traverse(fileName string) {
	fDir := filepath.Dir(fileName)
	fName := filepath.Base(fileName)

	ffStructIdent := ast.NewIdent(ffStructName)
	fIdent := ast.NewIdent("a")

	var cc *ast.FuncDecl
	featureFlagKeyNames := make([]ast.Decl, 0)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fmt.Printf("working on file %v\n", fileName)
	ast.Inspect(file, func(n ast.Node) bool {
		// find feature flag struct
		spec, ok := n.(*ast.TypeSpec)
		if ok {
			name := spec.Name.Name
			if name == ffStructName {
				ff := spec.Type.(*ast.StructType)
				k, function := generateFunc(ffStructIdent, ff.Fields, fIdent)
				featureFlagKeyNames = append(featureFlagKeyNames, k)
				cc = function
			}
		}
		return true
	})
	featureFlagCode := generateFeatureFlagCode(ffStructIdent, featureFlagKeyNames, cc)
	//printer.Fprint(os.Stdout, fs, featureFlagCode)
	file.Decls = featureFlagCode
	buf := new(bytes.Buffer)
	err = format.Node(buf, fset, file)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else if fileName[len(fileName)-8:] != "_test.go" {
		err := ioutil.WriteFile(fmt.Sprintf("%s/gen_%s", fDir, fName), buf.Bytes(), 0664)
		if err != nil {
			panic(err)
		}
	}
}

func generateFunc(ffStructIdent *ast.Ident, structFields *ast.FieldList, fIdent *ast.Ident) (*ast.GenDecl, *ast.FuncDecl) {
	ffKeyIdent := getFeatureFlagKeyTypeIdent(ffStructIdent)
	resultIdent := ast.NewIdent("result")

	caseClauses := make([]ast.Stmt, 0)
	constDecls := &ast.GenDecl{
		Tok:   token.CONST,
		Specs: make([]ast.Spec, 0),
	}

	for i, field := range structFields.List {
		fieldName := field.Names[0].Name
		keyIdent := ast.NewIdent(fmt.Sprintf("Key%s", fieldName))
		vSpec := &ast.ValueSpec{
			Names: []*ast.Ident{
				keyIdent,
			},
		}
		if i == 0 {
			vSpec.Type = ffKeyIdent
			vSpec.Values = []ast.Expr{
				&ast.BasicLit{
					Value: "iota",
				},
			}
		}
		constDecls.Specs = append(constDecls.Specs, vSpec)
		cc := &ast.CaseClause{
			List: []ast.Expr{
				keyIdent,
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{resultIdent},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   fIdent,
							Sel: ast.NewIdent(field.Names[0].Name),
						},
					},
				},
			},
		}
		caseClauses = append(caseClauses, cc)
		//fmt.Printf("found: %v\n", field.Names[0].Name)
	}

	isEnabledFunc := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{fIdent},
					Type:  ffStructIdent,
				},
			},
		},
		Name: ast.NewIdent("IsFeatureEnabled"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("key")},
						Type:  ffKeyIdent,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("bool"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{resultIdent},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{ast.NewIdent("false")},
				},
				&ast.SwitchStmt{
					Tag: ast.NewIdent("key"),
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{resultIdent},
				},
			},
		},
	}
	return constDecls, isEnabledFunc
}

func getFeatureFlagKeyTypeIdent(ffStructIdent *ast.Ident) *ast.Ident {
	return ast.NewIdent(fmt.Sprintf("%sKey", ffStructIdent.Name))
}

func generateFeatureFlagCode(ffStructIdent *ast.Ident, keys []ast.Decl, isEnabledFunc *ast.FuncDecl) []ast.Decl {
	fKeyType := []ast.Decl{
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: getFeatureFlagKeyTypeIdent(ffStructIdent),
					Type: ast.NewIdent("int"),
				},
			},
		},
	}
	fKeyType = append(fKeyType, keys...)
	fKeyType = append(fKeyType, isEnabledFunc)

	return fKeyType
}
