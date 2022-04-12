/*
Copyright © 2021-2022 Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package billing

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	pb "github.com/slntopp/nocloud/pkg/billing/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func MakeBillingServiceClientOrFail() (context.Context, pb.BillingServiceClient) {
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opt := grpc.WithTransportCredentials(creds)
	insecure := viper.GetBool("insecure")
	if insecure {
		opt = grpc.WithInsecure()
	}
	conn, err := grpc.Dial(host.(string), opt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic(err)
	}

	token := viper.Get("token")
	if token == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Token is unset")
	}

	client := pb.NewBillingServiceClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer " + token.(string))
	return ctx, client
}

func PrintPlan(p *pb.Plan) {
	fmt.Println("UUID:", p.Uuid)
	fmt.Printf("%s | %s", p.Title, p.Type)
	if p.Public {
		fmt.Print("| Public")
	}
	fmt.Println("\nResources:")
	for _, c := range p.Resources {
		fmt.Printf("  %s: ( %f / %d ) NCU/sec if ", c.Key, c.Price, c.Period)
		on := make([]string, len(c.On))
		for i, e := range c.On {
			on[i] = e.String()
		}
		if c.Except {
			fmt.Print("not ")
		}
		fmt.Println(strings.Join(on, ", "))
	}
}