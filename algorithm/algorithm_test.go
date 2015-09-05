package algorithm

import (
	"bufio"
	"fmt"
	"github.com/hzck/speedroute/json"
	m "github.com/hzck/speedroute/model"
	"math/rand"
	"os"
	"testing"
)

func TestAllFiles(t *testing.T) {
	for _, testPath := range getDirFileNames("tests/") {
		graph := json.CreateGraphFromFile(testPath)
		validation, err := os.Open("validations/" + testPath[6:])
		defer validation.Close()
		if err != nil {
			assertFailingPath(t, graph, testPath)
		} else {
			var path []string
			scanner := bufio.NewScanner(validation)
			for scanner.Scan() {
				path = append(path, scanner.Text())
			}
			assertCorrectPath(t, path, Route(graph), testPath)
		}
	}
}

func getDirFileNames(dirName string) []string {
	var fileNames []string
	dir, _ := os.Open(dirName)
	defer dir.Close()
	files, _ := dir.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			for _, fileName := range getDirFileNames(dirName + file.Name() + "/") {
				fileNames = append(fileNames, fileName)
			}
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
		if edge.From().Id() != nodes[i] || edge.To().Id() != nodes[i+1] {
			failAndPrint(t, "assertCorrectPath - correct "+nodes[i]+"->"+nodes[i+1]+", is "+edge.From().Id()+"->"+edge.To().Id()+": "+testPath)
		}
	}
}

func failAndPrint(t *testing.T, testPath string) {
	t.Fail()
	fmt.Println("Failing " + testPath)
}

func TestBenchMarkGraph(t *testing.T) {
	size := 2
	graph := createBenchmarkGraph(size)
	path := Route(graph)
	if len(path) != size*2 {
		t.Fail()
	}
}

func BenchmarkAlgorithm(b *testing.B) {
	graph := createBenchmarkGraph(30)
	b.ResetTimer()
	Route(graph)
}

func createBenchmarkGraph(size int) *m.Graph {
	graph := m.CreateGraph()
	transitionNode := m.CreateNode("tempId", false)
	graph.AddStartNode(transitionNode)
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
	graph.AddEndNode(transitionNode)
	return graph
}

func createWeightedEdge(from, to *m.Node, time int) *m.Edge {
	edge := m.CreateEdge(from, to)
	edge.AddWeight(m.CreateWeight(rand.Intn(10000) + 1))
	return edge
}
