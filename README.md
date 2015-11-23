# cmpe273-assignment3

/****** POST ******/
curl -H "Content-Type: application/json" -X POST -d '{
    "starting_from_location_id": "565179fd91643f77ed15ed17",
    "location_ids": [
        "56517a6791643f77ed15ed18",
        "5651820991643f78bf22c90a",
        "56517a9191643f77ed15ed19",
        "56517ab291643f77ed15ed1a"
    ]
}' http://127.0.0.1:8080/trips

Created
************************
{"Id":"565259fa91643f7d29dbaa98","Status":"planning","Starting_from_location_id":"565179fd91643f77ed15ed17","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14}
************************


/****** GET ******/
 curl -H "Content-Type: application/json" -X GET http://127.0.0.1:8080/trips/565259fa91643f7d29dbaa98
OK
************************
{"Id":"565259fa91643f7d29dbaa98","Status":"planning","Starting_from_location_id":"565179fd91643f77ed15ed17","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14}
************************


/****** PUT ******/
curl -H "Content-Type: application/json" -H 'Accept: application/json' -X PUT 'http://127.0.0.1:8080/trips/565259fa91643f7d29dbaa98/request'
OK
************************
{"Id":"56525be191643f7d29dbaa99","Status":"requesting","Starting_from_location_id":"565179fd91643f77ed15ed17","Next_destination_location_id":"5651820991643f78bf22c90a","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":2}
************************

/****** 2nd PUT Call *****/
OK
************************
{"Id":"56525bf691643f7d29dbaa9a","Status":"requesting","Starting_from_location_id":"5651820991643f78bf22c90a","Next_destination_location_id":"56517ab291643f77ed15ed1a","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":2}
************************

/****** 3rd PUT Call *****/
OK
************************
{"Id":"56525cb891643f7d29dbaa9b","Status":"requesting","Starting_from_location_id":"56517ab291643f77ed15ed1a","Next_destination_location_id":"56517a9191643f77ed15ed19","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":2}
************************

/****** 4th Put call *****/
OK
************************
{"Id":"56525cc191643f7d29dbaa9c","Status":"requesting","Starting_from_location_id":"56517a9191643f77ed15ed19","Next_destination_location_id":"56517a6791643f77ed15ed18","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":2}
************************

/*******5th Put call ****/
OK
************************
{"Id":"56525ccc91643f7d29dbaa9d","Status":"requesting","Starting_from_location_id":"56517a6791643f77ed15ed18","Next_destination_location_id":"565179fd91643f77ed15ed17","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":2}
************************

/*******Subsequent calls ****/
This trip is finished
OK
************************
{"Id":"56525ccf91643f7d29dbaa9e","Status":"finished","Starting_from_location_id":"","Next_destination_location_id":"","Best_route_location_ids":["5651820991643f78bf22c90a","56517ab291643f77ed15ed1a","56517a9191643f77ed15ed19","56517a6791643f77ed15ed18"],"Total_uber_costs":77,"Total_uber_duration":5711,"Total_distance":29.14,"Eta":0}
************************

