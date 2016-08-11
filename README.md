# bqldot

make DOT file for graphviz, from BQL

## Install

### Require

* sensorbee http://sensorbee.io/
    * v0.5.1
* cli https://github.com/urfave/cli
    * v1.18.0

### Install

```bash
$ go get github.com/disktnk/bqldot
$ cd $GOPATH/src/github.com/disktnk/bqldot
$ go install
```

## Usage

```bash
$ bqldot path/to/bqlfile/foo.bql
```

"foo.dot" will be made.

* UDSF is not supported, not output as edge.

### Example

"sample.bql" is put on "sample" directory

```bash
$ bqldot -t topology_name sample.bql # "sample.dot" will be made
$ dot -Tgif sample.dot -o sample_graph.gif
```

Visualized as:

![sample_graph.gif](sample/sample_graph.gif)

## Reference

* graphviz http://www.graphviz.org/
    * The DOT language http://www.graphviz.org/content/dot-language
