package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"google.golang.org/grpc"

	pb "satazor/jwt-with-grpc-performance-tests/pb"
)

var emailsFile, _ = ioutil.ReadFile("./emails.json")
var emails []string
var emailsIdx = 0
var emailsLength = len(emails)

var jwtSecret = []byte(os.Getenv("bmtXKKngXH1HRdshrI7LkJxmyNZyDN1f"))

var publicKey = "{\"x\": \"7-INQ150R-MCWlj5X_wyGLRIRYAA-o8NakJiUq7gOGg\", \"y\": \"dM-GsyJvdDOuALE3l-U9lPL8V3gY_5BPjLH539yTdKU\", \"alg\": \"ES256\", \"crv\": \"P-256\", \"kid\": \"cdd2969c-7e49-4a46-bcbe-e8bbdf74c7f3\", \"kty\": \"EC\"}"
var privateKey = "{\"alg\":\"ES256\",\"crv\":\"P-256\",\"d\":\"h-UIda1elff-qw81gsSQakyzOv8Dozv5RcQqFIV6R1Y\",\"kid\":\"cdd2969c-7e49-4a46-bcbe-e8bbdf74c7f3\",\"kty\":\"EC\",\"x\":\"7-INQ150R-MCWlj5X_wyGLRIRYAA-o8NakJiUq7gOGg\",\"y\":\"dM-GsyJvdDOuALE3l-U9lPL8V3gY_5BPjLH539yTdKU\"}"

var publicKeyJwk, _ = jwk.ParseKey([]byte(publicKey))
var privateKeyJwk, _ = jwk.ParseKey([]byte(privateKey))

var rawPublicKey interface{}
var rawPrivateKey interface{}

type Algorithm interface {
	sign(jwt.MapClaims) string
	verify(string) jwt.MapClaims
}

type HS256 struct{}

type ES256 struct{}

func (HS256) sign(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, _ := token.SignedString(jwtSecret)

	return jwtString
}

func (HS256) verify(jwtString string) jwt.MapClaims {
	token, _ := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims
}

func (ES256) sign(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	jwtString, _ := token.SignedString(rawPrivateKey)

	return jwtString
}

func (ES256) verify(jwtString string) jwt.MapClaims {
	token, _ := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return rawPublicKey, nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	return claims
}

func getAlgorithm(alg string) Algorithm {
	if alg == "HS256" {
		return HS256{}
	}

	return ES256{}
}

// --------------------------------------------------------------------

type JwtServer struct {
	pb.UnimplementedJwtServiceServer
}

func (s *JwtServer) GenerateJwt(ctx context.Context, request *pb.GenerateJwtRequest) (*pb.GenerateJwtResponse, error) {
	var algorithm = getAlgorithm(request.Algorithm)

	nowSeconds := time.Now().UnixMilli() / 1000
	sub := emails[emailsIdx]

	var jwtString = algorithm.sign(jwt.MapClaims{
		"sub": sub,
		"iat": nowSeconds,
		"exp": nowSeconds + 7200,
	})

	var claims = algorithm.verify(jwtString)

	if claims["sub"] != sub {
		return nil, errors.New("sub mismatch")
	}

	emailsIdx += 1

	if emailsIdx >= emailsLength {
		emailsIdx = 0
	}

	return &pb.GenerateJwtResponse{Sub: sub, Jwt: jwtString}, nil
}

func main() {
	json.Unmarshal([]byte(emailsFile), &emails)
	publicKeyJwk.Raw(&rawPublicKey)
	privateKeyJwk.Raw(&rawPrivateKey)

	lis, _ := net.Listen("tcp", ":50052")

	server := grpc.NewServer()

	pb.RegisterJwtServiceServer(server, &JwtServer{})
	log.Printf("server listening at %v", lis.Addr())

	server.Serve(lis)
}
