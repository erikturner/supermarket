# Supermarket

RESTful API to add, delete, and fetch produce from system

## Install:

docker desktop

## Running the code
```
docker pull erikturner/supermarket:1
```

### Spin up the service

```
# start docker container 
docker run --rm -p 8080:8080 erikturner/supermarket:1
```

### Making Requests
#### GET Request
```
curl -X GET -H "Content-Type: application/json" http://localhost:8080/produce/fetch 

Rest Client:
http://localhost:8080/produce/fetch

```
#### POST Request
```
curl -X POST -H "Content-Type: application/json" http://localhost:8080/produce/add -d '[{"produceCode":"TR32-YUT7-93WE-290K","name":"Grapes","unitPrice":1.00}]'

Rest Client:
http://localhost:8080/produce/add

[
	{
		"produceCode":"TR32-YUT7-93WE-290K",
		"name":"Grapes",
		"unitPrice":1.00
	}
]

```
#### DELETE Request
```
curl -X DELETE -H "Content-Type: application/json" http://localhost:8080/produce/TR32-YUT7-93WE-290K/remove

Rest Client:
http://localhost:8080/produce/TR32-YUT7-93WE-290K/remove

```
