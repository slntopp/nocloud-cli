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
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	pb "github.com/slntopp/nocloud/pkg/billing/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func MakeBillingServiceClientOrFail() (context.Context, pb.BillingServiceClient) {
	host := viper.Get("nocloud")
	if host == nil {
		fmt.Fprintln(os.Stderr, "Error setting connection up")
		panic("Host is unset")
	}

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	insec := viper.GetBool("insecure")
	if insec {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.Dial(host.(string), grpc.WithTransportCredentials(creds))
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
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token.(string))
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

var processedLabels = map[bool]string{
	true:  "V",
	false: "X",
}

func PrintRecords(pool []*pb.Record) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Instance", "Recource", "Start", "End", "Total NCU", "Processed"})

	for _, r := range pool {
		t.AppendRow(table.Row{
			r.Instance, r.Resource,
			time.Unix(r.Start, 0).Format("2006-01-02 15:04:05"),
			time.Unix(r.End, 0).Format("2006-01-02 15:04:05"),
			r.Total, processedLabels[r.Processed],
		})
	}

	t.Render()
}

func PrintTransactions(pool []*pb.Transaction, meta bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if meta {
		t.AppendHeader(table.Row{"UUID", "Account", "Service", "Timestamp", "Total NCU", "Meta"})
	} else {
		t.AppendHeader(table.Row{"UUID", "Account", "Service", "Timestamp", "Total NCU"})
	}

	rowFunc := MakeTrRow
	if meta {
		rowFunc = MakeTrRowWithMeta
	}

	sort.Slice(pool, func(i, j int) bool {
		return (pool[i].Exec - pool[j].Exec) < 0
	})
	for _, tr := range pool {
		t.AppendRow(rowFunc(tr))
	}

	t.Render()
}

func MakeTrRow(t *pb.Transaction) table.Row {
	ts := t.Exec
	if t.Processed {
		ts = t.Proc
	}
	return table.Row{
		t.Uuid,
		t.Account,
		t.Service,
		time.Unix(ts, 0).Format("2006-01-02 15:04:05"),
		t.Total,
	}
}

func MakeTrRowWithMeta(t *pb.Transaction) table.Row {
	ts := t.Exec
	if t.Processed {
		ts = t.Proc
	}
	meta, _ := json.Marshal(t.Meta)
	return table.Row{
		t.Uuid,
		t.Account,
		t.Service,
		ts,
		t.Total,
		string(meta),
	}
}
