# Random chatters

## API specs

### View socket
The WebView socket is used to show the current character in OBS
To use, connect to ```/api/rnd/ws```. This socket impliments following schemas:

- Ping
Send any mesage to socket and it will reply with this exact payload. Can be used to check connection status

- Message event
A message is sent to connection when the active chatter sends something in chat. The schema is as follows:
```
{
	"type": "message",
    "message": <message>
}
```

- Connection event
A new chatter is connectd. The schema is as follows:
```
{
	"type": "conenct",
    "username": string, // chatter's username
    "userpfp": string  // chatter's user profile
}
```

- Disconnect event
A current chatter is disconnectd. The schema is as follows:
```
{
	"type": "disconnect"
}
```

### Control flow
To control the flow app has WebView at ``` /rnd/control ```.
It impliments following APIs:

- ```/api/rnd/connect```
Selects a random not-ignored user from chat and sends "connect" message to a web view

- ```/api/rnd/disconnect```
Disconnects currently selected chatter and sends "disconnect" message to web view

- ```/api/rnd/bannedusers```
Returns a list of all ignored chatters. Schema is as follows:
```
{
    "chatters": string[]
}
```
This endpoint is not protected

- ```/api/rnd/ban?user=<username>```
Adds a user to blocked list. Uses displayed username

- ```/api/rnd/pardon?user=<username>```
Removes a user from blocked list. Uses username from list

 - ```/api/rnd/currnetchatter``` 
 Returns the name if the current user. Schema is as follows:
```
{
    "username": string
}
This endpoint is not protected

### Authentication
To access the control API endpoints you have to proveide an app token. It is accessible in .env file
You have to append it to your HTTP request via URL query like that:
``` /api/end/connect?token=<your token> ```
