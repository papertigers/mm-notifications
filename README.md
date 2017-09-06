# mm-notifications
Run a command on Mattermost mention notifications

## Usage
`go get -u  github.com/papertigers/mm-notifications/cmd/mm-notifications`

##### Config file (~/.mattermost.toml)
```
username =  "mike"
password = "super_secret"
url = "chat.example.com"
```

##### run command
`mm-notifications`


### Working
- [X] Login
- [X] Websocket connection
- [ ] Filter events
- [ ] Run user specified command on event
