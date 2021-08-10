# Account service


## Pre requisites

- Docker
- Golang v1.14+
 
 
## Running App 

1. Export DB_PASSWORD
`export DB_PASSWORD="S3cretP@ssw0rd"` 

2. Bring up the mysql container using:

`make infra-local`

3. Run the migrations on the local DB 
 
`make setup`

4. Build and run the app container.

`make app`

5. Inspect logs using docker 

`docker logs account-service-go -f`

## Design
Design an event store
	- credit 2 - expiry t+5
	- credit 2 - expiry t+2
	- debit 3 - 

Check the latest status of the db before the insertion of the DB
- Have some checkpoints of the state of the database.


Interpret the state of the system from the credits with expiry - debits.

### Assume:
	User service exists
		- Prepopulate the user DB with the user data.
	- each user has one account only
	- expiration is done using the minimum expiry in a group of priorities

Run the queries:
```
INSERT INTO public.accounts
(id, user_id, balance, created_at, updated_at)
VALUES(gen_random_uuid(), gen_random_uuid(), 0, timezone('utc'::text, now()), timezone('utc'::text, now()));
```

## Verifying the Functionality

Add Balance  Request
```shell script

curl -X POST \
  http://localhost:8888/account/credit \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{ "userId":"a900c144-\-4324-994f-451a7ac9d46d", "amount":2, "type":"subscription", "priority":2, "expiry":1605688837 }'
```

GET Credit logs request
```shell script
curl -X GET \
  'http://localhost:8888/account/logs?activity=Credit' \
 
```

Debit balance request
```shell script      
curl -X POST \
http://localhost:8888/account/debit \
-H 'content-type: application/json' \
-d '{
    "userId":"a900c144-9f25-4324-994f-451a7ac9d46d",
    "amount":2 
}'
```

get current balance using
```shell script
curl -X GET \
  http://localhost:8888/account/a900c144-9f25-4324-994f-451a7ac9d46d \
  -H 'cache-control: no-cache' \
```

