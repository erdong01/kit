package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rxt/cmd/report/proto/report"
)

const (
	address = "localhost:50001"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := report.NewReportRpcClient(conn)
	r, err := c.Show(context.Background(), &report.ReportRequest{ExamId: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.ExamId)
}
