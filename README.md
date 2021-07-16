# jwt

A JWT CLI written in go.

Also includes:
- A single html file web app for quickly decoding tokens: available online at [https://jwt.etelej.com/](https://jwt.etelej.com/))
  - File available at [./webapp](./tree/main/webapp) - uses [petite-vue](https://github.com/vuejs/petite-vue)

## Installation

Download executable for your OS from [Releases](./releases)
- Works on most platforms: Windows, iOS, Linux etc

## Usage

Decode token and print details in human readable format
- Color sections 
- Add headers
- Parse timestamps
```
./jwt TOKEN
```

Decode token, and print value as json 
```
./jwt -json TOKEN
```

Verify if token is signed by a secret
```
./jwt TOKEN --secret SECRET
```



Create a token
- Create a jwt token signed by a secret
```
./jwt --sign '{"user": "john doe"}' --secret "SECRET"
```


### Tips
If you don't want to paste secret in command line:
  - if you just specify `--secret` at the end, it will prompt for the secret which you can paste

```
./jwt TOKEN --secret

# output (prompts for secret, does not print)
Please paste secret and press Enter to continue
```
