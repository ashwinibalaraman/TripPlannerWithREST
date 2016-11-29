package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

type ParamsForReqId struct {
	Start_latitude  string `json:"start_latitude"`
	Start_longitude string `json:"start_longitude"`
	End_latitude    string `json:"end_latitude"`
	End_longitude   string `json:"end_longitude"`
	Product_id      string `json:"product_id"`
}

type PostReqParams struct {
	Starting_from_location_id string
	Location_ids              []string
	Address                   string
}

type StatusStruct struct {
	Status string `json:"status"`
}

type PostResParams struct {
	Id                        bson.ObjectId `bson:"_id"`
	Status                    string        `bson:"status"`
	Starting_from_location_id string        `bson:"starting_from_location_id"`
	Best_route_location_ids   []string      `bson:"Best_route_location_ids"`
	Total_uber_costs          float64       `bson:"Total_uber_costs"`
	Total_uber_duration       float64       `bson:"Total_uber_duration"`
	Total_distance            float64       `bson:"Total_distance"`
}

type PutResParams struct {
	Id                           bson.ObjectId `bson:"_id"`
	Status                       string        `bson:"status"`
	Starting_from_location_id    string        `bson:"starting_from_location_id"`
	Next_destination_location_id string        `bson:"next_destination_location_id"`
	Best_route_location_ids      []string      `bson:"best_route_location_ids"`
	Total_uber_costs             float64       `bson:"total_uber_costs"`
	Total_uber_duration          float64       `bson:"total_uber_duration"`
	Total_distance               float64       `bson:"total_distance"`
	Eta                          float64       `bson:"eta"`
}
type Coordinate struct {
	Lat string `bson:"lat"`
	Lng string `bson:"lng"`
}
type ResParameters struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string        `bson:"name"`
	Address string        `bson:"address"`
	City    string        `bson:"city"`
	State   string        `bson:"state"`
	Zip     string        `bson:"zip"`
	Coord   Coordinate    `bson:"coordinate"`
}

type PutParameters struct {
	Address string     `bson:"address"`
	City    string     `bson:"city"`
	State   string     `bson:"state"`
	Zip     string     `bson:"zip"`
	Coord   Coordinate `bson:"coordinate"`
}

type PutReqParameters struct {
	Address string
	City    string
	State   string
	Zip     string
}

var Url string
var PutCounter int

func main() {
	//Url = "localhost"
	Url = "mongodb://ashwini:<passwd>@ds045064.mongolab.com:45064/cmpe273"
	mux := httprouter.New()

	mux.GET("/trips/:trip_id", getTrips)
	mux.POST("/trips", postTrips)
	mux.PUT("/trips/:trip_id/request", putTrips)
	//mux.DELETE("/locations/:location_id", deleteLocations)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
	PutCounter = 0

}

func getAuth(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Fprintf(rw, "Authorized!!!!")
	//Authorization_code = p.ByName("code")
}

func databaseFindAddressBook(id string) (ResParameters, error) {
	sess, err := mgo.Dial(Url)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("AddressBook")
	data := ResParameters{}

	if len(id) != 24 {
		err := errors.New("Id: " + id + " is not 24 characters long")
		return data, err
	}
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).Select(bson.M{}).One(&data)

	fmt.Println("Data for id: " + id + " ====" + data.Address)

	return data, err

}

func databaseFindTrips(id string) (PostResParams, error) {
	sess, err := mgo.Dial(Url)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Trips")

	data := PostResParams{}
	if len(id) != 24 {
		err := errors.New("Id: " + id + " is not 24 characters long")
		return data, err
	}
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).Select(bson.M{}).One(&data)

	fmt.Println("Data for id: " + id + " ====" + data.Starting_from_location_id)

	return data, err

}

func databaseInsert(ResParams PostResParams) error {
	sess, err := mgo.Dial(Url)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Trips")

	doc := ResParams
	err = collection.Insert(doc)

	return err
}

func getTrips(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	sess, err := mgo.Dial(Url)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("cmpe273").C("Trips")
	data := PostResParams{}
	id := p.ByName("trip_id")
	fmt.Println("id:", id)
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).Select(bson.M{}).One(&data)

	if err != nil {
		if err.Error() == "not found" {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(rw, http.StatusText(http.StatusNotFound))
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, err.Error())
		}
	} else {
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, http.StatusText(http.StatusOK))
		b, _ := json.Marshal(data)
		fmt.Println(string(b))
		fmt.Fprintf(rw, "\n************************\n")
		fmt.Fprintf(rw, string(b))
		fmt.Fprintf(rw, "\n************************\n")
	}
}

// edge struct holds the bare data needed to define a graph.
type Edge struct {
	vert1, vert2 string
	cost         float64
	duration     float64
	distance     float64
}

var Graph []Edge
var edgeCount int

func postTrips(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	ReqParams := PostReqParams{}

	json.NewDecoder(req.Body).Decode(&ReqParams)

	starting_from_location_id := ReqParams.Starting_from_location_id
	fmt.Println(starting_from_location_id)
	location_ids := ReqParams.Location_ids

	edgeCount = 0
	n := len(location_ids) + 1
	edges := ((n * (n - 1)) / 2)
	fmt.Println("EDGE COUNT", edgeCount)
	Graph = make([]Edge, edges)
	for i := range location_ids {
		_, _, _, err := callUberApi(starting_from_location_id, location_ids[i], true)
		if err != nil {
			handleError(rw, "Error: ", err)
			return
		}
		i++
	}
	for j := 0; j < len(location_ids); j++ {
		for i := range location_ids {
			i++
			if i+j < len(location_ids) {
				_, _, _, err := callUberApi(location_ids[j], location_ids[i+j], true)
				if err != nil {
					handleError(rw, "Error: ", err)
					return
				}
			}
			//j++
		}
	}

	fmt.Println("-------------\n\n\n\n")
	fmt.Println(Graph)
	fmt.Println("-------------\n\n\n\n")
	lastVal := len(location_ids) - 1
	nodeOrder := DijkstraAlgoCall(starting_from_location_id, location_ids[lastVal], Graph)

	var totalDistance float64
	var totalPrice float64
	var totalDuration float64

	for i := range nodeOrder {
		if i <= len(nodeOrder)-2 {
			for j := 0; j < len(Graph); j++ {
				if (Graph[j].vert1 == nodeOrder[i] || Graph[j].vert2 == nodeOrder[i]) && (Graph[j].vert1 == nodeOrder[i+1] || Graph[j].vert2 == nodeOrder[i+1]) {
					//fmt.Println(Graph[j].cost)
					totalPrice = totalPrice + Graph[j].cost
					totalDuration = totalDuration + Graph[j].duration
					totalDistance = totalDistance + Graph[j].distance
				}
			}

		}

	}
	dur, dist, cost, err := callUberApi(location_ids[len(location_ids)-1], starting_from_location_id, false)
	if err != nil {
		handleError(rw, "Error: ", err)
		return
	}
	//nodeOrder[len(location_ids)] = starting_from_location_id
	totalPrice = totalPrice + cost
	totalDuration = totalDuration + dur
	totalDistance = totalDistance + dist
	fmt.Println("duration: ")
	fmt.Println(totalDuration)
	fmt.Println("distance: ")
	fmt.Println(totalDistance)
	fmt.Println("price: ")
	fmt.Println(totalPrice)

	//remove the first element in the nodeOrder as it is the start_location_id
	nodeOrder = append(nodeOrder[:0], nodeOrder[1:]...)
	var ResParams PostResParams
	ResParams = PostResParams{
		Id:     bson.NewObjectId(),
		Status: "planning",
		Starting_from_location_id: starting_from_location_id,
		Best_route_location_ids:   nodeOrder,
		Total_uber_duration:       totalDuration,
		Total_distance:            totalDistance,
		Total_uber_costs:          totalPrice,
	}

	err = databaseInsert(ResParams)
	if err != nil {
		handleError(rw, "DB Error: ", err)
		return
	} else {
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, http.StatusText(http.StatusCreated))
		b, _ := json.Marshal(ResParams)
		fmt.Println(string(b))
		fmt.Fprintf(rw, "\n************************\n")
		fmt.Fprintf(rw, string(b))
		fmt.Fprintf(rw, "\n************************\n")
	}
}
func buildGraph(startId string, endId string, duration interface{}, distance interface{}, low_estimate interface{}) {
	var totalDuration float64
	var totalDistance float64
	var totalPrice float64

	if x, ok := duration.(float64); ok {
		totalDuration = totalDuration + x
		Graph[edgeCount].duration = x
	}
	if y, ok := distance.(float64); ok {
		totalDistance = totalDistance + y
		Graph[edgeCount].distance = y
	}
	if z, ok := low_estimate.(float64); ok {
		totalPrice = totalPrice + z
		Graph[edgeCount].cost = z
	}
	Graph[edgeCount].vert1 = startId
	Graph[edgeCount].vert2 = endId
	edgeCount++

}
func callUberApi(startId string, endId string, build_graph bool) (duration_f float64, distance_f float64, low_estimate_f float64, err error) {

	start_loc_data, err := databaseFindAddressBook(startId)
	end_loc_data, err := databaseFindAddressBook(endId)

	if err != nil {
		return 0, 0, 0, err
	}

	queryString := "start_latitude=" + start_loc_data.Coord.Lat + "&" + "start_longitude=" + start_loc_data.Coord.Lng + "&" + "end_latitude=" + end_loc_data.Coord.Lat + "&" + "end_longitude=" + end_loc_data.Coord.Lng
	fmt.Println(queryString)

	uberApiUrl := "https://sandbox-api.uber.com/v1/estimates/price?" + queryString
	//resp, _ := http.Get(uberApiUrl)

	requ, err := http.NewRequest("GET", uberApiUrl, nil)
	requ.Header.Set("Authorization", "Token VGbXSQDwDPSyZSav4Sn_ui8vtu6DNyzf5IaDL3go")

	if err != nil {
		return 0, 0, 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(requ)
	if err != nil {
		return 0, 0, 0, err
	}
	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	//fmt.Println(data)
	prices := data.(map[string]interface{})["prices"]
	result := prices.([]interface{})[0]
	distance := result.(map[string]interface{})["distance"]
	low_estimate := result.(map[string]interface{})["low_estimate"]
	duration := result.(map[string]interface{})["duration"]
	if build_graph == true {
		buildGraph(startId, endId, duration, distance, low_estimate)
	}
	if x, ok := duration.(float64); ok {
		duration_f = x
	}
	if y, ok := distance.(float64); ok {
		distance_f = y
	}
	if z, ok := low_estimate.(float64); ok {
		low_estimate_f = z
	}
	return duration_f, distance_f, low_estimate_f, err

}

func handleError(rw http.ResponseWriter, msg string, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(rw, msg+err.Error())
	fmt.Println(msg + err.Error())
}

var ResParams PutResParams

func putTrips(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	access_token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsImhpc3RvcnkiLCJyZXF1ZXN0Il0sInN1YiI6IjMwOWZmYzE3LTk5NDAtNDkwOS1iNjE3LTdjNDNmZDNiMjI3ZSIsImlzcyI6InViZXItdXMxIiwianRpIjoiNTAzMTU3NmQtMjljYi00MTUxLWJmM2YtZGJjMDljODRiMGI2IiwiZXhwIjoxNDUwNTgwODQ4LCJpYXQiOjE0NDc5ODg4NDcsInVhY3QiOiJNRG8zamRwbldNamN3QTEya29xTThpM1N4UDViUEEiLCJuYmYiOjE0NDc5ODg3NTcsImF1ZCI6InNxUEhQSmJkODlxaEVPVlphR053WnhDdkFjN284Y0swIn0.oMyG0aUUcZPx_vm-yK1OeswkH36Ait4jS5qU9e4H_MBT9zdXnn2MYxurHkqWA388PReKc-gcz2PWMnuCeV2PHSN3TNOZKcCSgRdLj8Hudk9zMwitF1ovLv23tresJiSaqfJQbYiVqepR_s8SKTXCU72HltRvTe-CSteuOKXF5NlmOfpyK2E0eJiV0MyHkW9yya8_ylbt9bgQSANjEAEf3QElAsmxUEk6g4pY7QAfvn4rMI5yVKsFii922M0MBz-zA7Fnm7C-iOFvtcmPmn6P0FQvjgDcinABKtzJoVf3vdHkECfMm0ZhLHvEdLOEjtCvLhkXhXvJEwvxbaBLibyPhQ"
	id := p.ByName("trip_id")
	fmt.Println("id:", id)

	data, err := databaseFindTrips(id)
	if err != nil {
		handleError(rw, "DB Error for "+id+":", err)
		return
	}

	var start_loc_id string
	var end_loc_id string
	best_route := data.Best_route_location_ids
	productId := ""
	status := "requesting"

	var start_loc_data ResParameters
	var end_loc_data ResParameters

	if PutCounter > len(best_route) {
		fmt.Println("This trip is finished\n")
		fmt.Fprintf(rw, "This trip is finished\n")
		status = "finished"

		ResParams = PutResParams{
			Id:     bson.NewObjectId(),
			Status: status,
			Starting_from_location_id:    start_loc_id,
			Next_destination_location_id: end_loc_id,
			Best_route_location_ids:      best_route,
			Total_uber_duration:          data.Total_uber_duration,
			Total_distance:               data.Total_distance,
			Total_uber_costs:             data.Total_uber_costs,
			Eta:                          0,
		}

		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, http.StatusText(http.StatusOK))
		fmt.Fprintf(rw, "\n************************\n")
		b, _ := json.Marshal(ResParams)
		fmt.Println(string(b))
		fmt.Fprintf(rw, string(b))
		fmt.Fprintf(rw, "\n************************\n")

		return
	}
	if PutCounter <= len(best_route) {
		if PutCounter == 0 {
			start_loc_id = data.Starting_from_location_id
		} else {
			start_loc_id = best_route[PutCounter-1]
		}

		if PutCounter == len(best_route) {
			//status = "finished"
			start_loc_id = best_route[PutCounter-1]
			end_loc_id = data.Starting_from_location_id
		} else {
			end_loc_id = best_route[PutCounter]
		}

		start_loc_data, err = databaseFindAddressBook(start_loc_id)
		if err != nil {
			handleError(rw, "DB Error for "+start_loc_id+":", err)
			return
		}
		end_loc_data, err = databaseFindAddressBook(end_loc_id)

		if err != nil {
			handleError(rw, "DB Error for "+end_loc_id+":", err)
			return
		}
		/***get product id for the location: GET Request**/
		queryString := "latitude=" + start_loc_data.Coord.Lat + "&longitude=" + start_loc_data.Coord.Lng
		getUberApiUrl := "https://sandbox-api.uber.com/v1/products?" + queryString
		getReq, err := http.NewRequest("GET", getUberApiUrl, nil)
		if err != nil {
			handleError(rw, "Error: ", err)
			return
		}
		getReq.Header.Set("Content-Type", "application/json")
		getReq.Header.Set("Authorization", "Bearer "+access_token)

		client := &http.Client{}
		getResp, err := client.Do(getReq)
		if err != nil {
			handleError(rw, "Get Error: ", err)
			return
		}
		var getData interface{}
		err = json.NewDecoder(getResp.Body).Decode(&getData)
		if err != nil {
			handleError(rw, "Decode Error: ", err)
			return
		}

		fmt.Println("Status", getResp.StatusCode)
		products := getData.(map[string]interface{})["products"]
		product := products.([]interface{})[0]
		product_id := product.(map[string]interface{})["product_id"]
		if str, ok := product_id.(string); ok {
			productId = str
		}
		/*****get productId END*****/

		/**Get requestId and eta for the trip: POST Request*****/
		jsonForReqId := ParamsForReqId{
			Product_id:      productId,
			Start_latitude:  start_loc_data.Coord.Lat,
			Start_longitude: start_loc_data.Coord.Lng,
			End_latitude:    end_loc_data.Coord.Lat,
			End_longitude:   end_loc_data.Coord.Lng,
		}
		//queryString := "product_id=" + product_id + "&start_latitude=" + start_loc_data.Coord.Lat + "&start_longitude=" + start_loc_data.Coord.Lng + "&end_latitude=" + end_loc_data.Coord.Lat + "&end_longitude=" + end_loc_data.Coord.Lng
		uberApiUrl := "https://sandbox-api.uber.com/v1/requests?" //+ queryString
		b, err := json.Marshal(jsonForReqId)
		fmt.Println("jsonForReqId: ", jsonForReqId)
		fmt.Println("json byte string: ", string(b))
		if err != nil {
			handleError(rw, "Marshal error: ", err)
			return
		}
		requ, _ := http.NewRequest("POST", uberApiUrl, bytes.NewBuffer(b))
		requ.Header.Set("Content-Type", "application/json")
		requ.Header.Set("Authorization", "Bearer "+access_token)
		//fmt.Println("request is:", requ)
		client = &http.Client{}
		resp, err := client.Do(requ)
		if err != nil {
			handleError(rw, "Post Error: ", err)
			return
		}
		var httpData interface{}
		err = json.NewDecoder(resp.Body).Decode(&httpData)
		if err != nil {
			handleError(rw, "Decode Error: ", err)
			return
		}

		fmt.Println("Status", resp.StatusCode)
		request_id := httpData.(map[string]interface{})["request_id"]
		eta := httpData.(map[string]interface{})["eta"]

		fmt.Println("REQUEST ID ", request_id)
		fmt.Println("ETA  ", eta)
		/*****Get requestId and eta for the trip END*****/

		/***** Change the status for the trip: PUT Request*****/
		statusJson := StatusStruct{Status: "accepted"}
		b, err = json.Marshal(statusJson)
		if err != nil {
			handleError(rw, "Marshal Error: ", err)
			return
		}
		if str, ok := request_id.(string); ok {
			uberApiUrl = "https://sandbox-api.uber.com/v1/sandbox/requests/" + str
		}

		fmt.Println(uberApiUrl)
		requ, err = http.NewRequest("PUT", uberApiUrl, bytes.NewBuffer(b))
		if err != nil {
			handleError(rw, "Put Error: ", err)
			return
		}
		requ.Header.Set("Content-Type", "application/json")
		requ.Header.Set("Authorization", "Bearer "+access_token)

		resp2, _ := client.Do(requ)
		var data2 interface{}
		err = json.NewDecoder(resp2.Body).Decode(&data2)

		if resp2.StatusCode == 204 {
			var eta_num float64
			if num, ok := eta.(float64); ok {
				eta_num = num
				fmt.Println("it is string", num, ok)
			}

			ResParams = PutResParams{
				Id:     bson.NewObjectId(),
				Status: status,
				Starting_from_location_id:    start_loc_id,
				Next_destination_location_id: end_loc_id,
				Best_route_location_ids:      best_route,
				Total_uber_duration:          data.Total_uber_duration,
				Total_distance:               data.Total_distance,
				Total_uber_costs:             data.Total_uber_costs,
				Eta:                          eta_num,
			}
			PutCounter++
			rw.WriteHeader(http.StatusCreated)
			fmt.Fprintf(rw, http.StatusText(http.StatusOK))
			b, _ := json.Marshal(ResParams)
			fmt.Println(string(b))
			fmt.Fprintf(rw, "\n************************\n")
			fmt.Fprintf(rw, string(b))
			fmt.Fprintf(rw, "\n************************\n")
		} else {
			rw.WriteHeader(resp2.StatusCode)
		}
	}

}
