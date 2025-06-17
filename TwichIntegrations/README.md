A set of utilities I use for my streams that connect to twich API

# Usage

- Set up .env file. Template:
```
TWICH_API_KEY=<API key of your app>
TWICH_API_SECRET=<API secret of your app>
BROADCASTER_LOGIN=<Tartget chat owner>
MODERATOR_LOGIN=<Token issuer name>
```
- go run .

# Features

## Twich TTS
Connects to twich chat and reads messages. Uses espeak as a driver

### Dependencies
- espeak

### Params
For every chatter there will be speech params:
| Param | Description     | Range     | espeak flag |
|-------|-----------------|-----------|-------------|
| Voice | Accent          | Table     | --voices    |
| Speed | Words per min   | [125-225] | -s          |
| Pitch | Voice pitch     | [1-100]   | -p          |
| Cap   | Capitak letters | [1-200]   | -k          |

## Random chatter selection
Impliments system like DougDoug's chatter selector

### Requirements
- You have to be an moderator on target streamer's channel
- Set EnableRandomChatter in dotenv to true

### Usage
