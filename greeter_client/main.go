/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"flag"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	clientCert := flag.String("client-cert", "", "client-cert path")
	clientKey := flag.String("client-key", "", "client-key path")
	caPath := flag.String("ca-path", "", "ca-cert path")
	flag.Parse()

	certs, capool, err := loadCredentials(
		// certFile,
		// "./cfssl-certs/client.pem",
		*clientCert,
		// keyFile,
		// "./cfssl-certs/client-key.pem",
		*clientKey,
		// caFile,
		// "./cfssl-certs/ca.pem",
		*caPath,
	)
	if err != nil {
		log.Fatalf("failed to load certs: %v", err)
	}

	cred := credentials.NewTLS(&tls.Config{
		Certificates: certs,
		RootCAs:      capool,
	})
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(
		cred,
		// credentials.NewTLS(tlsConfig),
	))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

}

func loadCredentials(certFile, keyFile, caFile string) ([]tls.Certificate, *x509.CertPool, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, nil, err
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, nil, err
	}
	return []tls.Certificate{cert}, certPool, nil
}
