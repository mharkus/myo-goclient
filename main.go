package main

import(
	"code.google.com/p/go.net/websocket"
	"strconv"
	"reflect"
	"fmt"
	"encoding/json"
)

type MyoData []map[string]interface{}

type MyoEvent struct {
	Accelerometer []float64
	Gyroscope []float64
	Myo float64
	Orientation struct {
		W float64
		X float64 
		Y float64
		Z float64
	}
	Timestamp int 
	Type string 
}

func addEvent(item map[string]interface{}) MyoEvent{
	myoEvent := MyoEvent{}
	myoEvent.Type = item["type"].(string)

	switch myoEvent.Type {
		case "orientation":
			addGyroscope(&myoEvent, item["gyroscope"])
			addAccelerometer(&myoEvent, item["accelerometer"])
			addOrientation(&myoEvent, item["orientation"])
		case "paired":
		case "connected":
		case "pose":
	}

	myoEvent.Myo = item["myo"].(float64)
	myoEvent.Timestamp,_ = strconv.Atoi(item["timestamp"].(string))
	
	return myoEvent
}

func addGyroscope(myoEvent * MyoEvent, item interface{}){
	switch reflect.TypeOf(item).Kind(){
	case reflect.Slice:
		s := reflect.ValueOf(item)
		myoEvent.Gyroscope = make([]float64, s.Len())
		for i := 0; i < s.Len(); i++ {
			myoEvent.Gyroscope[i] = s.Index(i).Interface().(float64)
		}
	}
}

func addAccelerometer(myoEvent * MyoEvent, item interface{}){
	switch reflect.TypeOf(item).Kind(){
	case reflect.Slice:
		s := reflect.ValueOf(item)
		myoEvent.Accelerometer = make([]float64, s.Len())
		for i := 0; i < s.Len(); i++ {
			myoEvent.Accelerometer[i] = s.Index(i).Interface().(float64)
		}
	}
}

func addOrientation(myoEvent * MyoEvent, item interface{}){
	switch reflect.TypeOf(item).Kind(){
	case reflect.Map:

		s := reflect.ValueOf(item)
		keys := s.MapKeys()
		for _,k := range keys {
			val := k.Interface().(string)
			result := s.Interface().(map[string]interface{})[k.Interface().(string)].(float64)
			switch val {
				case "w":
					myoEvent.Orientation.W = result	
				case "x":
					myoEvent.Orientation.X = result
				case "y":
					myoEvent.Orientation.Y = result
				case "z":
					myoEvent.Orientation.Z = result

			}
		}
	}
	
}


func main() {

	// connect to local websocket exposed by Myo Connect
	ws, err := websocket.Dial("ws://127.0.0.1:7204/myo/1", "", "http://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	var resp = make([]byte, 4096)
	var myoData MyoData

	for{
		// this blocks until websocket server sends some data
		n, err := ws.Read(resp)
		
		if err != nil {
			panic(err)
		}
		
		// convert to string
		item := string(resp[0:n])

		// then unmarshall the string to our MyoData struct which is just a 2-element array
		json.Unmarshal([]byte(item), &myoData);

		// manually traverse the json data and build our MyoEvent struct
		myoEvent := addEvent(myoData[1])

		fmt.Println(myoEvent);
	}
}

