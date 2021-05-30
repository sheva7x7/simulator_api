package vehicles

import (
	"encoding/json"
	"fmt"
	ESClient "wenle/elasticsearch/esclient"
	"wenle/elasticsearch/mqtt"

	elastic "github.com/olivere/elastic/v7"
)

type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Vehicle struct {
	Car_no string `json:"car_no"`
	Latlon LatLon `json:"latlon"`
}

type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type Boundary struct {
	Fence Polygon `json:"fence"`
}

func UpdateVehicle(v Vehicle) (*Vehicle, error) {
	_, err := ESClient.UpdateItem("vehicles", v.Car_no, v)
	if err != nil {
		return nil, err
	}
	if CheckIfWithinBoundary(v.Car_no) {
		mqtt.Publish("car_in_boundary", v)
	} else {
		mqtt.Publish("car_outside_boundary", v)
	}
	return &v, err
}

func UpdateBoundary(boundary Boundary) (*Boundary, error) {
	_, err := ESClient.UpdateItem("boundary", "1", boundary)
	if err != nil {
		return nil, err
	}
	return &boundary, err
}

func GetBoundary() Boundary {
	item, _ := ESClient.GetItem("boundary", "1")
	var boundary Boundary
	json.Unmarshal(item.Source, &boundary)
	return boundary
}

func CheckIfWithinBoundary(id string) bool {
	boundary := GetBoundary()
	coordinates := boundary.Fence.Coordinates[0]
	searchSource := elastic.NewSearchSource()
	boolQuery := elastic.NewBoolQuery()
	matchQuery := elastic.NewMatchQuery("_id", id)
	geoPolygonQuery := elastic.NewGeoPolygonQuery("latlon")
	length := len(coordinates)
	for i := 0; i < length; i++ {
		coordinatesString := fmt.Sprintf("%f,%f", coordinates[i][0], coordinates[i][1])
		geoPoint, _ := elastic.GeoPointFromString(coordinatesString)
		geoPolygonQuery.AddGeoPoint(geoPoint)
	}
	boolQuery.Filter(geoPolygonQuery)
	boolQuery.Must(matchQuery)
	searchSource.Query(boolQuery)
	result, err := ESClient.QueryItem("vehicles", searchSource)
	if err != nil {
		fmt.Println(err)
	}
	if result.Hits.TotalHits.Value > 0 {
		return true
	}
	return false
}

func GetVehiclesWithinBoundary() ([]Vehicle, error) {
	boundary := GetBoundary()
	coordinates := boundary.Fence.Coordinates[0]
	searchSource := elastic.NewSearchSource()
	boolQuery := elastic.NewBoolQuery()
	geoPolygonQuery := elastic.NewGeoPolygonQuery("latlon")
	length := len(coordinates)
	for i := 0; i < length; i++ {
		coordinatesString := fmt.Sprintf("%f,%f", coordinates[i][0], coordinates[i][1])
		geoPoint, _ := elastic.GeoPointFromString(coordinatesString)
		geoPolygonQuery.AddGeoPoint(geoPoint)
	}
	boolQuery.Filter(geoPolygonQuery)
	searchSource.Query(boolQuery)
	result, err := ESClient.QueryItem("vehicles", searchSource)
	if err != nil {
		fmt.Println(err)
	}
	var vehicles []Vehicle
	for _, hit := range result.Hits.Hits {
		var item Vehicle
		json.Unmarshal(hit.Source, &item)
		vehicles = append(vehicles, item)
	}
	return vehicles, err
}
