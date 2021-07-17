package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/term"
)

var (
	sign   = flag.String("sign", "", "json to sign as claims")
	secret = flag.String("secret", "", "secret to use in signing jwt token")
	exp    = flag.String("exp", "", fmt.Sprintf(
		"expiry - lifetime of secret (optional) (%s)", supportedUnits),
	)
	// color      = flag.Bool("color", false, "colorize output")
	// onlyClaims = flag.Bool("claims", false, "only print claims")

	// support intuitive flags
	_      = flag.Bool("decode", false, "decode json (unused)")
	encode = flag.String("encode", "", "encode jwt token (alias for --sign)")
)

func main() {

	flag.Parse()
	log.SetFlags(0)

	var encodeArg string
	if encodeArg = *sign; encodeArg == "" {
		encodeArg = *encode
	}

	if encodeArg != "" {
		key := *secret
		lifetime := parseExp(*exp)
		if key == "" {
			log.Print("Please enter a secret to sign the JWT (and press Enter)")
			dat, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatalf("Failed to read input: %v", err)
			}
			key = string(dat)
		}
		token, err := generateJWT(encodeArg, key, lifetime)
		if err != nil {
			log.Fatalf("Failed to generate JWT token: %v", err)
		}
		log.Printf("%s", token)
		return
	}
	key := *secret
	if err := decode(flag.Arg(0), key); err != nil {
		log.Fatal(err)
	}
}

func decode(jwtToken, key string) error {
	if jwtToken == "" {
		return errors.New("[jwt] failed to decode, no token specified")
	}
	jwtParts := strings.Split(jwtToken, ".")
	if len(jwtParts) != 3 {
		return errors.New("[jwt] JWT token provided does not look valid")
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(*secret), nil
	})
	validSecret := err == nil

	headerJSON, err := jwt.DecodeSegment(jwtParts[0])
	if err != nil {
		return fmt.Errorf("[jwt] failed to decode JWT header: %v", err)
	}
	claimsJSON, err := jwt.DecodeSegment(jwtParts[1])
	if err != nil {
		return fmt.Errorf("[jwt] failed to decode JWT claims: %v", err)
	}
	var stdClaims jwt.StandardClaims
	if err := json.Unmarshal(claimsJSON, &stdClaims); err != nil {
		return fmt.Errorf("unable to parse standard claims from JWT: %v", err)
	}

	issuedAt := timeFromUnix(stdClaims.IssuedAt)
	expiresAt := timeFromUnix(stdClaims.ExpiresAt)
	notBefore := timeFromUnix(stdClaims.NotBefore)
	now := time.Now()
	var meta string
	if !issuedAt.IsZero() {
		iat := issuedAt.Format(time.RFC822)
		meta += "Issued At: " + iat
		if issuedAt.After(now) {
			meta += "\t// invalid time, issued in future"
		} else {
			meta += fmt.Sprintf("\t// %s ago", readableDuration(time.Since(issuedAt)))
		}
		meta += "\n"
	}
	if !notBefore.IsZero() {
		nbf := notBefore.Format(time.RFC822)
		meta += "Not Before: " + nbf
		if now.Before(notBefore) {
			meta += fmt.Sprintf("\t// cannot be used for %s", readableDuration(notBefore.Sub(now)))
		}
		meta += "\n"
	}
	if !expiresAt.IsZero() {
		exp := expiresAt.Format(time.RFC822)
		meta += "Expires At: " + exp
		if now.After(expiresAt) {
			meta += fmt.Sprintf("\t// expired %s ago", readableDuration(time.Since(expiresAt)))
		} else {
			meta += fmt.Sprintf("\t// expires in %s", readableDuration(expiresAt.Sub(now)))
		}
		meta += "\n"
		var start time.Time
		if start := issuedAt; start.IsZero() {
			start = notBefore
		}
		if !start.IsZero() {
			life := expiresAt.Sub(start)
			meta += "Token Lifetime: " + readableDuration(life)
			meta += "\n"
		}
	}
	if key != "" {
		if validSecret {
			meta += "JWT Token signature verified, key is valid\n"
		} else {
			meta += "JWT Token signature not matching, key is not valid\n"
		}
	}

	log.Printf(`
✻ Header
%s

✻ Claims
%s

Signature: %s

%s
`,
		indentJSON(headerJSON), indentJSON(claimsJSON), jwtParts[2],
		meta,
	)
	return nil
}

func indentJSON(dat []byte) []byte {
	var m map[string]interface{}
	err := json.Unmarshal(dat, &m)
	if err != nil {
		return nil
	}
	out, _ := json.MarshalIndent(m, "", "  ")
	return out
}

func generateJWT(claimsStr, key string, exp time.Duration) ([]byte, error) {
	if key == "" {
		return nil, errors.New("missing signing, please specify with --secret")
	}
	var claims jwt.MapClaims
	err := json.Unmarshal([]byte(claimsStr), &claims)
	if err != nil {
		return nil, err
	}
	if exp != 0 {
		claims["exp"] = time.Now().Add(exp).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	return []byte(tokenStr), nil

}

var (
	errExpDurationError = fmt.Errorf(
		"invalid --exp, please use one of these supported --exp units: %s",
		supportedUnits)
)

const supportedUnits = "yr, mo, w, d, h,m,s"

var units = []struct {
	EndsWith string
	Ms       float64
}{
	{
		EndsWith: "yr",
		Ms:       float64(time.Hour) * 24 * 365,
	},
	{
		EndsWith: "mo",
		Ms:       float64(time.Hour) * 24 * 30,
	},
	{
		EndsWith: "w",
		Ms:       float64(time.Hour) * 24 * 7,
	},
	{
		EndsWith: "d",
		Ms:       float64(time.Hour) * 24,
	},
}

func parseExp(exp string) time.Duration {
	if exp == "" {
		return time.Duration(0)
	}

	for _, unit := range units {
		if !strings.HasSuffix(exp, unit.EndsWith) {
			continue
		}
		val := strings.TrimSuffix(exp, unit.EndsWith)
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Print(errExpDurationError.Error())
			return time.Duration(0)
		}
		total := f * unit.Ms
		return time.Duration(total)
	}
	d, err := time.ParseDuration(exp)
	if err != nil {
		log.Print(errExpDurationError.Error())
	}
	return d
}
func readableDuration(d time.Duration) string {
	f := float64(d.Nanoseconds())
	for _, unit := range units {
		val := f / unit.Ms
		if val > 1 {
			return fmt.Sprintf(`%.0f%s`, val, unit.EndsWith)
		}
	}
	return fmt.Sprint(d)
}

func timeFromUnix(unixTime int64) time.Time {
	if unixTime == 0 {
		return time.Time{}
	}
	return time.Unix(unixTime, 0)
}
