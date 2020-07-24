package algorithm

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"testing"

	m "github.com/hzck/speedroute/model"
	"github.com/hzck/speedroute/parser"
)

// TestAllFiles goes through tests/ folder and tries to route a path, verifying towards files
// in validations/ folder. Folder structure is the same in both tests/ and validations/ to
// know which validation belongs to which test case.
func TestAllFiles(t *testing.T) {
	for _, testPath := range getDirFileNames(t, "tests/") {
		graph, err := parser.CreateGraphFromFile(testPath)
		if err != nil {
			failAndPrint(t, err.Error())
		}
		validation, err := os.Open("validations/" + testPath[6:])
		defer func() {
			if closeErr := validation.Close(); closeErr != nil && err == nil {
				failAndPrint(t, closeErr.Error())
			}
		}()
		if err != nil {
			assertFailingPath(t, graph, testPath)
			continue
		}
		var path []string
		scanner := bufio.NewScanner(validation)
		for scanner.Scan() {
			path = append(path, scanner.Text())
		}
		assertCorrectPath(t, path, Route(graph), testPath)
	}
}

func getDirFileNames(t *testing.T, dirName string) []string {
	var fileNames []string
	dir, err := os.Open(dirName)
	defer func() {
		if closeErr := dir.Close(); closeErr != nil && err == nil {
			failAndPrint(t, "Can't close dir: "+closeErr.Error())
		}
	}()
	if err != nil {
		failAndPrint(t, "Can't open dir: "+err.Error())
		return nil
	}
	files, _ := dir.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			fileNames = append(fileNames, getDirFileNames(t, dirName+file.Name()+"/")...)
			continue
		}
		fileNames = append(fileNames, dirName+file.Name())
	}
	return fileNames
}

func assertFailingPath(t *testing.T, graph *m.Graph, testPath string) {
	if Route(graph) != nil {
		failAndPrint(t, "assertFailingPath: "+testPath)
	}
}

func assertCorrectPath(t *testing.T, nodes []string, path []*m.Edge, testPath string) {
	if len(path) != len(nodes)-1 {
		failAndPrint(t, "assertCorrectPath - len: "+testPath)
	}
	for i, edge := range path {
		if edge.From().ID() != nodes[i] || edge.To().ID() != nodes[i+1] {
			failAndPrint(t, "assertCorrectPath - correct "+nodes[i]+"->"+nodes[i+1]+", is "+edge.From().ID()+"->"+edge.To().ID()+": "+testPath)
		}
	}
}

func failAndPrint(t *testing.T, testPath string) {
	t.Fail()
	fmt.Println("Failing " + testPath)
}

// TestBenchmarkGraph verifies that the benchmarking randomized graph works.
func TestBenchmarkGraph(t *testing.T) {
	size := 4
	graph := createBenchmarkGraph(size)
	path := Route(graph)
	if len(path) != size*2 {
		t.Fail()
	}
}

// BenchmarkAlgorithm consists of an int of the magnitude the graph should be benchmarked against.
func BenchmarkAlgorithm(b *testing.B) {
	graph := createBenchmarkGraph(30)
	b.ResetTimer()
	Route(graph)
}

// createBenchmarkGraph creates a graph from an int parameter specifying the magnitude of the graph.
// Nodes: n²+n+1. Edges: n³+n². Possible paths: (can't remember).
func createBenchmarkGraph(size int) *m.Graph {
	transitionNode := m.CreateNode("tempId", false)
	startNode := transitionNode
	for i := 0; i < size; i++ {
		var nodes []*m.Node
		tempNode := m.CreateNode("tempId", false)
		for k := 0; k < size; k++ {
			newNode := m.CreateNode("tempId", false)
			nodes = append(nodes, newNode)
			createWeightedEdge(transitionNode, newNode, rand.Intn(10000)+1)
			createWeightedEdge(newNode, tempNode, rand.Intn(10000)+1)
		}
		transitionNode = tempNode
		for _, a := range nodes {
			for _, b := range nodes {
				if a != b {
					createWeightedEdge(a, b, rand.Intn(10000)+1)
				}
			}
		}
	}
	return m.CreateGraph(startNode, transitionNode)
}

func createWeightedEdge(from, to *m.Node, time int) {
	edge := m.CreateEdge(from, to)
	edge.AddWeight(m.CreateWeight(time))
}
