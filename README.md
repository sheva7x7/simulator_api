# Simulator API

This project provides API for the simulator app

## To run the app

```
go install
go run .
```

### Endpoints

- GET /boundary : gets the boundary of fence
- GET /vehicles/bound : gets the vehicles within the boundary
- POST /vehicle : updates vehicle position
- POST /boundary : update boundary