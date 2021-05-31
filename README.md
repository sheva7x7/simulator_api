# Simulator API

This project provides API for the simulator app

## To run the app

```
go install
go run .
```

### Create index for vehicles on initiation of elasticsearch

```
curl --location --request PUT 'localhost:9200/vehicles' \
--header 'Content-Type: application/json' \
--data-raw '{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
            "car_no": {
                "type": "keyword"
            },
            "latlon": {
                "type": "geo_point"
            }
        }
    }
}'
```

```
curl --location --request PUT 'localhost:9200/boundary' \
--header 'Content-Type: application/json' \
--data-raw '{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
            "fence": {
                "type": "geo_shape"
            }
        }
    }
}'
```

### Endpoints

- GET /boundary : gets the boundary of fence
- GET /vehicles/bound : gets the vehicles within the boundary
- POST /vehicle : updates vehicle position
- POST /boundary : update boundary