package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_workshop/hw/service"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(
		insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := file_transfer.NewFileTransferServiceClient(conn)

	// Read file
	log.Println("Read file")
	stream1, err := c.GetFileData(context.Background(), &file_transfer.FileName{Name: "1.txt"})
	if err != nil {
		log.Fatalf("Error on stream messages: %v", err)
	}

	for {
		msg, err := stream1.Recv()
		if err == io.EOF {
			// End of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving message: %v", err)
		}

		log.Printf("Received message: %s", msg.Data)
	}

	// Get all file names
	log.Println("Get all file names")
	stream2, err := c.GetAllFileNames(context.Background(), &file_transfer.Empty{})
	if err != nil {
		log.Fatalf("Error on stream messages: %v", err)
	}

	for {
		msg, err := stream2.Recv()
		if err == io.EOF {
			// End of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving message: %v", err)
		}

		log.Printf("Received message: %s", msg.Name)
	}

	// Get file metadata
	log.Println("Get file metadata")
	stream3, err := c.GetFileMetaData(context.Background(), &file_transfer.FileName{Name: "1.txt"})
	if err != nil {
		log.Fatalf("Error on stream messages: %v", err)
	}

	for {
		msg, err := stream3.Recv()
		if err == io.EOF {
			// End of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving message: %v", err)
		}

		log.Printf("Received message: %s %d %t", msg.Name, msg.Size, msg.IsDir)
	}
}
