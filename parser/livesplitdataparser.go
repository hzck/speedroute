package parser

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"strconv"
	"strings"
)

type segments struct {
	Segments []struct {
		Name            string `xml:"Name"`
		BestSegmentTime string `xml:"BestSegmentTime>RealTime"`
	} `xml:"Segments>Segment"`
}

func LivesplitXMLtoJSON(data string) (string, error) {
	var v segments
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		return "", err
	}
	var g graph
	currentNode := "START"
	emptyRewardRefs := make([]rewardRef, 0)
	g.Rewards = make([]reward, 0)
	g.Nodes = append(g.Nodes, node{ID: currentNode, Rewards: emptyRewardRefs})
	g.StartID = currentNode
	for _, seg := range v.Segments {
		g.Nodes = append(g.Nodes, node{ID: seg.Name, Rewards: emptyRewardRefs})
		time, err := parseTime(seg.BestSegmentTime)
		if err != nil {
			return "", err
		}
		g.Edges = append(g.Edges, edge{From: currentNode, To: seg.Name, Weights: []weight{weight{Time: &time, Requirements: emptyRewardRefs}}})
		currentNode = seg.Name
	}
	g.EndID = currentNode
	result, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func parseTime(time string) (int, error) {
	times := strings.Split(time, ":")
	hours, err := strconv.ParseInt(times[0], 10, 0)
	if err != nil {
		return -1, err
	}
	minutes, err := strconv.ParseInt(times[1], 10, 0)
	if err != nil {
		return -1, err
	}
	ms, err := strconv.ParseFloat(times[2], 64)
	if err != nil {
		return -1, err
	}
	result := (int(hours)*60+int(minutes))*60*1000 + int(round(ms*1000))
	return result, nil
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
