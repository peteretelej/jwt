# jwt

A JWT CLI written in Go.

Also includes:
- A single html file web app for quickly decoding tokens: available online at [https://jwt.etelej.com/](https://jwt.etelej.com/))
  - File available at [./webapp](https://github.com/peteretelej/jwt/tree/main/webapp) - uses [petite-vue](https://github.com/vuejs/petite-vue)

## Installation

Download executable for your OS from [Releases](https://github.com/peteretelej/jwt/releases/latest)
- Works on most platforms: Windows, iOS, Linux etc

Alternative installation:
```
go get github.com/peteretelej/jwt
```

## Usage

Decode token and print details in human readable format
- Add headers to sections
- Parse timestamps
- Provide human readable durations eg ( expires in 1yr )
```
./jwt TOKEN
```

Verify if token is signed by a secret
```
./jwt --secret SECRET TOKEN 
```

Generate JWT token
```
./jwt --secret demopass --sign '{"user": "John Doe"}' 
```

Specify an expiry period for the generated token
```
./jwt  --secret demopass --exp 1y --sign '{"user": "John Doe"}'
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

<details>
 <summary>Usage Examples</summary>

Generating a signed JWT token
```
./jwt --sign '{"name": "John Doe"}' --secret demopass --exp 2w
```
![image](https://user-images.githubusercontent.com/2271973/126047461-08ee52b5-88e3-404c-98f4-77d992c14ec1.png)

Decode JWT token
```
./jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjc3NTkyNzUsIm5hbWUiOiJKb2huIERvZSJ9.SF2XbD6QpxxcV95Oa_AC1oXysfWcF9gmyMEaNAHagP0
```
![image](https://user-images.githubusercontent.com/2271973/126047690-cdfd72f3-f6bb-4423-903e-3e33a9bcab14.png)

Verifying a token's signing key 
```
./jwt --secret demopass eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjc3NTkyNzUsIm5hbWUiOiJKb2huIERvZSJ9.SF2XbD6QpxxcV95Oa_AC1oXysfWcF9gmyMEaNAHagP0
```
![image](https://user-images.githubusercontent.com/2271973/126047696-6b3676b7-1050-4faa-8b7e-57f1e0777761.png)


Decoding a token with multiple standard claims
```
./jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNjIxMjY5Mjg2LCJleHAiOjE2Mzk3NTg4ODYsIm5iZiI6MTYyMzk0NzY4Nn0.GR5pGSiJZk3Ls0A429K3HIZfsmQqGnyIhPusDT5F5BU
```
![image](https://user-images.githubusercontent.com/2271973/126047835-e6dc6ae3-cbd0-4a24-a851-b56a148ec994.png)


</details>

License: **MIT**
