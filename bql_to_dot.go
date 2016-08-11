package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	app := cli.NewApp()
	app.Name = "bqldot"
	app.Usage = "bqldot command make DOT file for BQL graph"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "topology, t",
			Value: "",
			Usage: "name of the topology. by default, the topology name will be the BQL file name",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "name of output file. by default, output file name will be the BQL file name",
		},
	}
	app.Action = Run

	app.Run(os.Args)
}

// Run makes DOT file from BQL file.
func Run(c *cli.Context) error {
	if len(c.Args()) != 1 {
		cli.ShowSubcommandHelp(c)
		os.Exit(1)
	}

	bqlFilePath := c.Args()[0]
	bqlFileName := filepath.Base(bqlFilePath)
	topologyName := bqlFileName[:len(bqlFileName)-len(filepath.Ext(bqlFileName))]
	if n := c.String("topology"); n != "" {
		topologyName = n
	}
	outputName := topologyName + ".dot"
	if o := c.String("output"); o != "" {
		outputName = o
	}

	dot, err := convertToDot(bqlFilePath, topologyName)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return ioutil.WriteFile(outputName, []byte(dot), 0644)
}

func convertToDot(bqlFilePath, topologyName string) (string, error) {
	queries, err := func() (string, error) {
		f, err := os.Open(bqlFilePath)
		if err != nil {
			return "", err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}()
	if err != nil {
		return "", err
	}
	bp := parser.New()
	stmts, err := bp.ParseStmts(queries)
	if err != nil {
		return "", err
	}
	dot := ""
	for _, stmt := range stmts {
		dot += convertToDotLine(stmt)
	}
	return fmt.Sprintf(`digraph bql_graph {
  graph [label = "%s", labelloc=t];
%s}
`, topologyName, dot), nil
}

func convertToDotLine(stmt interface{}) string {
	switch stmt := stmt.(type) {
	case parser.CreateSourceStmt:
		return fmt.Sprintf("  %s [shape = box, label = \"%s\\nTYPE %s\"];\n",
			stmt.Name, stmt.Name, stmt.Type)
	case parser.CreateStreamAsSelectStmt:
		dot := fmt.Sprintf("  %s [shape = ellipse];\n", stmt.Name)
		dot += convertRelationEdge(string(stmt.Name), stmt.Select.Relations)
		return dot
	case parser.CreateStreamAsSelectUnionStmt:
		dot := fmt.Sprintf("  %s [shape = ellipse];\n", stmt.Name)
		for _, selStmt := range stmt.Selects {
			dot += convertRelationEdge(string(stmt.Name), selStmt.Relations)
		}
		return dot
	case parser.CreateSinkStmt:
		return fmt.Sprintf("  %s [shape = box, label = \"%s\\nTYPE %s\"];\n",
			stmt.Name, stmt.Name, stmt.Type)
	case parser.InsertIntoFromStmt:
		return fmt.Sprintf("  %s -> %s;\n", stmt.Input, stmt.Sink)
	case parser.CreateStateStmt:
		return convertStateNode(string(stmt.Name), string(stmt.Type), "")
	case parser.LoadStateStmt:
		return convertStateNode(string(stmt.Name), string(stmt.Type), stmt.Tag)
	case parser.LoadStateOrCreateStmt:
		return convertStateNode(string(stmt.Name), string(stmt.Type), stmt.Tag)
	default:
		// should output nothing when other statement type
		// because BQL file have potential to include operating statement
		// like SAVE, RESUME and so on.
		return ""
	}
}

func convertRelationEdge(out string, rels []parser.AliasedStreamWindowAST) string {
	dot := ""
	for _, rel := range rels {
		switch rel.Type {
		case parser.ActualStream:
			// not only interval, should add capacity and shedding?
			dot += fmt.Sprintf("  %s -> %s [label = \"RANGE %s %s\"];\n",
				rel.Name, out, rel.IntervalAST, rel.IntervalAST.Unit)
		case parser.UDSFStream:
			fmt.Println("not support UDSF, function name: " + rel.Name)
		}
	}
	return dot
}

func convertStateNode(name, sttype, tag string) string {
	if tag != "" {
		return fmt.Sprintf("  %s [shape = ellipse, label = \"%s\\nTYPE %s\\nTAG %s\"];\n",
			name, name, sttype, tag)
	}
	return fmt.Sprintf("  %s [shape = ellipse, label = \"%s\\nTYPE %s\"];\n",
		name, name, sttype)
}
