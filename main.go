package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	pb "github.com/fffbbbbbb/ocr-grpc-server/ocr"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedOcrServerServer
}

func (s *server) Getcaptcha(ctx context.Context, in *pb.ImageBuffer) (*pb.Captcha, error) {
	1
	image := in.GetImage()
	data, err := getcaptch(&image)
	if err != nil {
		return nil, err
	}
	log.Println(data)
	return &pb.Captcha{Data: data}, nil
}

func getcaptch(image *[]byte) (string, error) {
	for i := 0; i < 30; i++ {
		ul := uuid.NewV4()
		fileName := "./" + ul.String() + ".jepg"
		f, _ := os.Create(fileName)
		defer func() {
			f.Close()
			os.Remove(fileName)
		}()
		_, err := f.Write(*image)
		if err != nil {
			log.Println(err)
			return "", err
		}

		cmd := exec.Command("tesseract", fileName, "stdout")
		// 执行命令，并返回结果
		output, err := cmd.Output()
		if err != nil {
			log.Println(err)
			return "", err
		}
		if len(output) >= 6 {
			return string(output[:4]), nil
		}
	}
	return "", fmt.Errorf("can not ocr image")
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags) // set flags
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOcrServerServer(s, &server{})
	log.Println("service up")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
