# logjam

Private repo for logjam.

**DO NOT SHARE THIS!**

# Endpoints
To broadcast
```
/files/frontend/?role=broadcast
```

To be audience
```
/files/frontend/?role=audience
```

View logs
```
/logs
```

If you want to use it in a specific `room`  then you need to pass the `room` as a query param.

# Setting and retrieving background URL

```javascript
sparkRTC.socket.send(JSON.stringify({type:'metadata-get'}));
sparkRTC.socket.send(JSON.stringify({type:'metadata-set', data: '{"backgroundUrl": "whatever"}'}));

sparkRtc.metaData.backgroundUrl;
```
