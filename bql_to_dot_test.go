package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/parser"
	"testing"
)

func TestConvertToDotLine(t *testing.T) {
	Convey("Given a test set", t, func() {
		type tcase struct {
			name string
			bql  string
			dot  string
		}
		tcases := []tcase{
			tcase{"create source", `CREATE SOURCE s TYPE src;`,
				"  s [shape = box, label = \"s\\nTYPE src\"];\n"},
			tcase{"create state", `CREATE STATE s TYPE st;`,
				"  s [shape = ellipse, label = \"s\\nTYPE st\"];\n"},
			tcase{"load state", `LOAD STATE s TYPE st;`,
				"  s [shape = ellipse, label = \"s\\nTYPE st\"];\n"},
			tcase{"load state with tag", `LOAD STATE s TYPE st TAG v1;`,
				"  s [shape = ellipse, label = \"s\\nTYPE st\\nTAG v1\"];\n"},
			tcase{"load or create state", `LOAD STATE s TYPE st OR CREATE IF NOT SAVED;`,
				"  s [shape = ellipse, label = \"s\\nTYPE st\"];\n"},
			tcase{"load or create state with tag", `LOAD STATE s TYPE st TAG v1 OR CREATE IF NOT SAVED;`,
				"  s [shape = ellipse, label = \"s\\nTYPE st\\nTAG v1\"];\n"},
			tcase{"create stream", `CREATE STREAM b AS SELECT RSTREAM * FROM s [RANGE 1 TUPLES]`,
				"  b [shape = ellipse];\n  s -> b [label = \"RANGE 1 TUPLES\"];\n"},
			tcase{"create stream with join",
				`CREATE STREAM b AS SELECT RSTREAM * FROM s [RANGE 0.2 SECONDS], s2 [RANGE 500 MILLISECONDS];`,
				"  b [shape = ellipse];\n  s -> b [label = \"RANGE 0.2 SECONDS\"];\n  s2 -> b [label = \"RANGE 500 MILLISECONDS\"];\n"},
			tcase{"create stream with union all",
				`CREATE STREAM b AS SELECT RSTREAM * FROM s [RANGE 10 TUPLES] UNION ALL SELECT RSTREAM * FROM s2 [RANGE 10 TUPLES];`,
				"  b [shape = ellipse];\n  s -> b [label = \"RANGE 10 TUPLES\"];\n  s2 -> b [label = \"RANGE 10 TUPLES\"];\n"},
			tcase{"create sink", `CREATE SOURCE s TYPE si;`,
				"  s [shape = box, label = \"s\\nTYPE si\"];\n"},
			tcase{"insert into sink", `INSERT INTO s FROM b;`,
				"  b -> s;\n"},
		}
		bp := parser.New()
		for _, tcase := range tcases {
			tcase := tcase
			Convey(fmt.Sprintf("When given a test BQL line: %v", tcase.name), func() {
				stmt, _, err := bp.ParseStmt(tcase.bql)
				So(err, ShouldBeNil)
				dot := convertToDotLine(stmt)
				Convey("Then the BQL line should be converted correctly", func() {
					So(dot, ShouldEqual, tcase.dot)
				})
			})
		}
	})
}
