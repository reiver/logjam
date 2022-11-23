# logjam

Private repo for logjam.

**DO NOT SHARE THIS!**

_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-

# Run Backend
```bash
cd backend
go run . -vvvvvv -http-port=8080
```

# Endpoints
To broadcast
```
http://localhost:8090/files/frontend/?role=broadcast
```

To be audience
```
http://localhost:8090/files/frontend/?role=audience
```


If you want to use it in a specific `room`  then you need to pass the `room` as a query param.
