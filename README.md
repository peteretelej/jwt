# jwt

A JWT CLI written in go.

Also includes:
- A single html file web app for quickly decoding tokens: available online at [https://jwt.etelej.com/](https://jwt.etelej.com/))
  - File available at [./webapp](./tree/main/webapp) - uses [petite-vue](https://github.com/vuejs/petite-vue)

## Installation

Download executable for your OS from [Releases](./releases)
- Works on most platforms: Windows, iOS, Linux etc

Alternative installation:
```
go get github.com/peteretelej/jwt
```

## Usage

Decode token and print details in human readable format
- Color sections 
- Add headers
- Parse timestamps
```
./jwt TOKEN
```

Verify if token is signed by a secret
```
./jwt --secret SECRET TOKEN 
```

Generate JWT token
```
./jwt --sign '{"user": "John Doe"}' --secret demopass
```

Specify an expiry period for the generated token
```
./jwt --sign '{"user": "John Doe"}' --secret demopass --exp 1y
```
- Supports durations (eg year to second) such as `yr, mo, w, d, h,m,s`
- examples: `--exp 6mo` (6 months), `--exp 2w` (2 weeks)



### Tips
If you don't want to paste secret as plain text in command line:
  - do not specify the `--secret` argument, it will them prompt for pasting or typing secret

```
./jwt --sign '{"user": "Jane Doe"}'
Please enter a secret to sign the JWT (and press Enter)
```
