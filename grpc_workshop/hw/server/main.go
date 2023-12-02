package main

import (
	"google.golang.org/grpc"
	"grpc_workshop/hw/service"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type server struct {
	file_transfer.UnimplementedFileTransferServiceServer
	filesPath string
	batchSize int
}

func (s *server) GetFileData(req *file_transfer.FileName, stream file_transfer.FileTransferService_GetFileDataServer) error {
	file, err := os.Open(s.filesPath + req.Name)
	if err != nil {
		return err
	}

	buf := make([]byte, s.batchSize)

	for {
		count, err := file.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&file_transfer.FileData{Data: buf[:count]}); err != nil {
			return err
		}
	}
}

func (s *server) GetAllFileNames(req *file_transfer.Empty, stream file_transfer.FileTransferService_GetAllFileNamesServer) error {
	files, err := ioutil.ReadDir(s.filesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := stream.Send(&file_transfer.FileName{Name: file.Name()}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) GetFileMetaData(req *file_transfer.FileName, stream file_transfer.FileTransferService_GetFileMetaDataServer) error {
	fileStat, err := os.Stat(s.filesPath + req.Name)
	if err != nil {
		return err
	}

	if err := stream.Send(&file_transfer.FileMetaData{Name: fileStat.Name(), Size: fileStat.Size(), IsDir: fileStat.IsDir()}); err != nil {
		return err
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	file_transfer.RegisterFileTransferServiceServer(s, &server{filesPath: "/Users/nikunis/GolandProjects/-ab-hw9/grpc_workshop/hw/files/", batchSize: 1024})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
