package parser

import (
	"encoding/json"
	"reflect"
	"testing"
)

const data16 string = `
<?xml version="1.0" encoding="UTF-8"?>
<Run version="1.6.0">
  <Segments>
    <Segment>
      <Name>This is 1</Name>
      <BestSegmentTime>
        <RealTime>00:00:04.9200000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>2 fo' sho'</Name>
      <BestSegmentTime>
        <RealTime>00:00:04.6610000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>Segment 3</Name>
      <BestSegmentTime>
        <RealTime>00:00:02.9120000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>4</Name>
      <BestSegmentTime>
        <RealTime>00:00:25.3680000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>5</Name>
      <BestSegmentTime>
        <RealTime>12:20:01.6740000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>6</Name>
      <BestSegmentTime>
        <RealTime>02:02:03.4480000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>7</Name>
      <BestSegmentTime>
        <RealTime>00:00:00.3525000</RealTime>
      </BestSegmentTime>
    </Segment>
    <Segment>
      <Name>8 going to the end</Name>
      <BestSegmentTime>
        <RealTime>00:00:03.7524000</RealTime>
      </BestSegmentTime>
    </Segment>
  </Segments>
</Run>
`

const answer16 string = `
{
  "rewards": [],
  "nodes": [
    {
      "id": "START",
      "rewards": []
    },
    {
      "id": "This is 1",
      "rewards": []
    },
    {
      "id": "2 fo' sho'",
      "rewards": []
    },
    {
      "id": "Segment 3",
      "rewards": []
    },
    {
      "id": "4",
      "rewards": []
    },
    {
      "id": "5",
      "rewards": []
    },
    {
      "id": "6",
      "rewards": []
    },
    {
      "id": "7",
      "rewards": []
    },
    {
      "id": "8 going to the end",
      "rewards": []
    }
  ],
  "edges": [
    {
      "from": "START",
      "to": "This is 1",
      "weights": [
        {
          "time": "00:00:04.9200000",
          "requirements": []
        }
      ]
    },
    {
      "from": "This is 1",
      "to": "2 fo' sho'",
      "weights": [
        {
          "time": "00:00:04.6610000",
          "requirements": []
        }
      ]
    },
    {
      "from": "2 fo' sho'",
      "to": "Segment 3",
      "weights": [
        {
          "time": "00:00:02.9120000",
          "requirements": []
        }
      ]
    },
    {
      "from": "Segment 3",
      "to": "4",
      "weights": [
        {
          "time": "00:00:25.3680000",
          "requirements": []
        }
      ]
    },
    {
      "from": "4",
      "to": "5",
      "weights": [
        {
          "time": "12:20:01.6740000",
          "requirements": []
        }
      ]
    },
    {
      "from": "5",
      "to": "6",
      "weights": [
        {
          "time": "02:02:03.4480000",
          "requirements": []
        }
      ]
    },
    {
      "from": "6",
      "to": "7",
      "weights": [
        {
          "time": "00:00:00.3525000",
          "requirements": []
        }
      ]
    },
    {
      "from": "7",
      "to": "8 going to the end",
      "weights": [
        {
          "time": "00:00:03.7524000",
          "requirements": []
        }
      ]
    }
  ],
  "startId": "START",
  "endId": "8 going to the end"
}
`

// TestParseVersion16 tests that version livesplit data version 1.6 is parseable.
func TestParseVersion16(t *testing.T) {
	var xmljson, correct graph
	result, err := LivesplitXMLtoJSON(data16)
	validateError(t, err)
	validateError(t, json.Unmarshal([]byte(result), &xmljson))
	validateError(t, json.Unmarshal([]byte(answer16), &correct))
	if !reflect.DeepEqual(xmljson, correct) {
		t.Errorf("The parsed result is not correct:\n" + result + "\nVS\n" + answer16)
	}
	t.Log(xmljson)
	t.Log(correct)
}

func validateError(t *testing.T, err error) {
	if err != nil {
		t.Errorf(err.Error())
	}
}
