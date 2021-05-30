package esclient

import (
	"context"
	"fmt"
	"os"

	elastic "github.com/olivere/elastic/v7"
)

var client *elastic.Client
var ctx context.Context

func GetESClient() (*elastic.Client, error) {
	if client != nil {
		return client, nil
	}
	client, err := StartClient()
	return client, err
}

func StartClient() (*elastic.Client, error) {
	fmt.Println("start es client", os.Getenv("ES_CLIENT_URL"))
	ctx = context.Background()
	ESClient, clientErr := elastic.NewClient(elastic.SetURL(os.Getenv("ES_CLIENT_URL")),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	if clientErr != nil {
		fmt.Println("Initialization failed")
	}
	client = ESClient
	fmt.Println("ES initialized... 52")

	return client, clientErr
}

func UpdateItem(index string, id string, obj interface{}) (*elastic.IndexResponse, error) {
	updatedItem, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(&obj).
		Refresh("true").
		Do(ctx)
	if err != nil {
		panic(err)
	}
	return updatedItem, err
}

func GetItem(index string, id string) (*elastic.GetResult, error) {
	item, err := client.Get().
		Index(index).
		Id(id).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	return item, err
}

func QueryItem(index string, searchSource *elastic.SearchSource) (*elastic.SearchResult, error) {
	result, err := client.Search().
		Index(index).
		SearchSource(searchSource).
		Pretty(true).
		Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}
