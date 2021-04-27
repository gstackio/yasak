package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"

	gopatch "github.com/gstackio/go-patch/patch"
	"github.com/spf13/cobra"
)

var debug, silent, version bool

var Version = "(development)"

func main() {

	var path string

	var locateCmd = &cobra.Command{
		Use:   "locate <yaml-file>",
		Short: "Locate the line and column where a value appears in a YAML file",
		Long: `Locate the line and column where a value appears in a YAML file.

    As input, the value is designated by a go-patch path. Go read
    <https://github.com/cppforlife/go-patch/blob/master/docs/examples.md> for
    reference.

    Line and column are separated by a TAB, whick make them easy to
    distinguish with the posix 'cut' utility.

    Whenever the designated value can't be found in the YAML file, a error is
    printed on stderr, and the tool
    exits with an error status of 1.

	When silent mode is enabled and the value can't be found, no error is
	printed, and the tool exits with no error, i.e. an exit status of 0.`,
		Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			file := args[0]

			if debug {
				fmt.Fprintf(os.Stderr, "DEBUG: file: '%s'\n", file)
				fmt.Fprintf(os.Stderr, "DEBUG: path: '%s'\n", path)
			}

			_, err := os.Stat(file)
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "ERROR: file does not exist: '%s' (%s)\n", file, err.Error())
				os.Exit(2)
			}

			// NOTE: for a 'fail-fast' attitude, we check the pointer before reading
			// the YAML file
			pointer, err := gopatch.NewPointerFromString(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: invalid path: '%s' (%s). "+
					"Please consult <https://github.com/cppforlife/go-patch/blob/master/docs/examples.md> for reference.\n",
					path, err.Error())
				os.Exit(2)
			}

			content, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: can't read file: '%s' (%s)\n", file, err.Error())
				os.Exit(2)
			}

			document := &yaml.Node{}
			err = yaml.Unmarshal(content, document)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: invalid YAML content in file: '%s' (%s)\n", file, err.Error())
				os.Exit(2)
			}

			node, err := locate(document, pointer, path)
			if err != nil && !silent {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			if node != nil {
				fmt.Printf("%d\t%d\t%s\n", node.Line, node.Column, node.LineComment)
			}
		},
	}

	locateCmd.Flags().StringVarP(&path, "path", "p", "", "YAML Path, expressed as a go-patch path.")
	locateCmd.MarkFlagRequired("path")

	var rootCmd = &cobra.Command{
		Version: Version,
		Use:   "yasak",
		Short: "Yet Another YAML Swiss Army Knife",
		Long: `Yet Another YAML Swiss Army Knife.

    A YAML processor built with love by Gstack and friends.`,
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "activate debug mode")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "activate silent mode")
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "display version")

	rootCmd.AddCommand(locateCmd)
	rootCmd.Execute()
	os.Exit(0)
}

func fatalIf(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func locate(document *yaml.Node, pointer gopatch.Pointer, path string) (*yaml.Node, error) {
	tokens := pointer.Tokens()
	node := document.Content[0]

	if len(tokens) == 1 {
		return node, nil
	}

	for tokenIdx, token := range tokens[1:] {
		isLast := tokenIdx == len(tokens)-2
		currPath := gopatch.NewPointer(tokens[:tokenIdx+2])

		switch typedToken := token.(type) {
		case gopatch.IndexToken:
			if node.Kind != yaml.SequenceNode {
				return nil, gopatch.NewOpArrayMismatchTypeErr(currPath, node)
			}

			idx, err := gopatch.ArrayIndex{Index: typedToken.Index, Modifiers: typedToken.Modifiers, Array: reflect.ValueOf(node.Content), Path: currPath}.Concrete()
			if err != nil {
				return nil, err
			}

			if isLast {
				return node.Content[idx], nil
			}
			node = node.Content[idx]

		case gopatch.AfterLastIndexToken:
			errMsg := "Expected not to find after last index token in path '%s' (not supported in find operations)"
			return nil, fmt.Errorf(errMsg, path)

		case gopatch.MatchingIndexToken:
			if node.Kind != yaml.SequenceNode {
				return nil, gopatch.NewOpArrayMismatchTypeErr(currPath, node)
			}

			var matchingIndexes []int

			for _, mapping := range node.Content {
				if mapping.Kind != yaml.MappingNode {
					continue
				}
				for idx, keyNode := range mapping.Content {
					// key are even indexes, values are odd indexes
					if idx%2 == 1 {
						continue
					}
					if keyNode.Kind == yaml.ScalarNode && keyNode.Value == typedToken.Key {
						valNode := mapping.Content[idx+1]
						if valNode.Kind == yaml.ScalarNode && valNode.Value == typedToken.Value {
							matchingIndexes = append(matchingIndexes, idx)
						}
					}
				}
			}

			if typedToken.Optional && len(matchingIndexes) == 0 {
				// todo /blah=foo?:after, modifiers
				node = &yaml.Node{
					Kind:   yaml.MappingNode,
					Style:  0,
					Tag:    "!!map",
					Value:  "",
					Anchor: "",
					Alias:  nil,
					Content: []*yaml.Node{
						{
							Kind:        yaml.ScalarNode,
							Style:       yaml.LiteralStyle,
							Tag:         "!!str",
							Value:       typedToken.Key,
							Anchor:      "",
							Alias:       nil,
							Content:     []*yaml.Node{},
							HeadComment: "",
							LineComment: "",
							FootComment: "",
							Line:        -1,
							Column:      -1,
						},
						{
							Kind:        yaml.ScalarNode,
							Style:       yaml.LiteralStyle,
							Tag:         "!!str",
							Value:       typedToken.Value,
							Anchor:      "",
							Alias:       nil,
							Content:     []*yaml.Node{},
							HeadComment: "",
							LineComment: "",
							FootComment: "",
							Line:        -1,
							Column:      -1,
						},
					},
					HeadComment: "",
					LineComment: "",
					FootComment: "",
					Line:        -1,
					Column:      -1,
				}

				if isLast {
					return node, nil
				}
			} else {
				if len(matchingIndexes) != 1 {
					return nil, gopatch.OpMultipleMatchingIndexErr{Path: currPath, Idxs: matchingIndexes}
				}

				idx, err := gopatch.ArrayIndex{Index: matchingIndexes[0], Modifiers: typedToken.Modifiers, Array: reflect.ValueOf(node.Content), Path: currPath}.Concrete()
				if err != nil {
					return nil, err
				}

				if isLast {
					return node.Content[idx], nil
				}
				node = node.Content[idx]
			}

		case gopatch.KeyToken:
			if node.Kind != yaml.MappingNode {
				return nil, gopatch.OpMissingMapKeyErr{Key: typedToken.Key, Path: currPath, Obj: mappingKeys(node)}
			}

			var (
				child *yaml.Node
				found bool
			)

			for idx, keyNode := range node.Content {
				// key are even indexes, values are odd indexes
				if idx%2 == 1 {
					continue
				}
				if keyNode.Kind == yaml.ScalarNode && keyNode.Value == typedToken.Key {
					found = true
					child = node.Content[idx+1]
					break
				}
			}

			if !found && !typedToken.Optional {
				return nil, gopatch.OpMissingMapKeyErr{Key: typedToken.Key, Path: currPath, Obj: mappingKeys(node)}
			}

			node = child
			if isLast {
				return node, nil
			}
			if !found {
				// Determine what type of value to create based on next token
				switch tokens[tokenIdx+2].(type) {
				case gopatch.MatchingIndexToken:
					node = &yaml.Node{
						Kind:        yaml.SequenceNode,
						Style:       0,
						Tag:         "!!seq",
						Value:       "",
						Anchor:      "",
						Alias:       nil,
						Content:     []*yaml.Node{},
						HeadComment: "",
						LineComment: "",
						FootComment: "",
						Line:        -1,
						Column:      -1,
					}

				case gopatch.KeyToken:
					node = &yaml.Node{
						Kind:        yaml.MappingNode,
						Style:       0,
						Tag:         "!!map",
						Value:       "",
						Anchor:      "",
						Alias:       nil,
						Content:     []*yaml.Node{},
						HeadComment: "",
						LineComment: "",
						FootComment: "",
						Line:        -1,
						Column:      -1,
					}

				default:
					errMsg := "Expected to find key or matching index token at path '%s'"
					return nil, fmt.Errorf(errMsg, gopatch.NewPointer(tokens[:tokenIdx+3]))
				}
			}

		default:
			return nil, gopatch.OpUnexpectedTokenErr{Token: token, Path: currPath}
		}
	}

	return document.Content[0], nil
}

func mappingKeys(mapping *yaml.Node) (keysAsMap reflect.Value) {
	keys := map[string]interface{}{}
	for idx, keyNode := range mapping.Content {
		// key are even indexes, values are odd indexes
		if idx%2 == 1 {
			continue // skip values
		}
		if keyNode.Kind != yaml.ScalarNode {
			continue // skip non scalar keys
		}
		keys[keyNode.Value] = true
	}
	keysAsMap = reflect.ValueOf(keys)
	return
}
