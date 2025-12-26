package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"

	pb "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/file/v1"
	"github.com/kitanoyoru/kgym/internal/apps/file/internal/service"
	"github.com/kitanoyoru/kgym/internal/apps/user/pkg/metrics"
	"github.com/kitanoyoru/kgym/pkg/metrics/prometheus"
	"go.uber.org/multierr"
)

const (
	GRPCServiceMetricsPrefix = "kgym.file.api.grpc"
)

type FileServiceServer struct {
	pb.UnimplementedFileServiceServer

	service service.IService
}

func NewFileService(service service.IService) (*FileServiceServer, error) {
	methods := []string{
		"UploadUserAvatar",
		"GetFileURL",
		"DeleteFile",
	}

	for _, method := range methods {
		if err := metrics.GlobalRegistry.RegisterMetric(prometheus.MetricConfig{
			Name: fmt.Sprintf("%s.%s", GRPCServiceMetricsPrefix, method),
			Type: prometheus.Counter,
		}); err != nil {
			return nil, err
		}
	}

	return &FileServiceServer{
		service: service,
	}, nil
}

func (s *FileServiceServer) UploadUserAvatar(stream pb.FileService_UploadUserAvatarServer) error {
	ctx := stream.Context()

	uploadRequest, doneChan, err := s.streamToPipe(ctx, stream)
	if err != nil {
		return err
	}

	resp, err := s.service.Upload(ctx, uploadRequest)
	if err != nil {
		return err
	}

	streamErr := <-doneChan
	if streamErr != nil {
		return streamErr
	}

	return stream.SendAndClose(&pb.UploadUserAvatar_Response{
		File: &pb.File{
			Metadata: &pb.Metadata{
				Id: resp.ID,
			},
		},
	})
}

func (s *FileServiceServer) GetFileURL(ctx context.Context, req *pb.GetFileURL_Request) (*pb.GetFileURL_Response, error) {
	url, err := s.service.GetURL(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetFileURL_Response{
		Url: url,
	}, nil
}

func (s *FileServiceServer) DeleteFile(ctx context.Context, req *pb.DeleteFile_Request) (*pb.DeleteFile_Response, error) {
	err := s.service.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteFile_Response{}, nil
}

func (s *FileServiceServer) streamToPipe(ctx context.Context, stream pb.FileService_UploadUserAvatarServer) (service.UploadRequest, chan error, error) {
	var (
		uploadRequest          = service.UploadRequest{}
		pipeReader, pipeWriter = io.Pipe()
		doneChan               = make(chan error, 1)
		firstRunFlag           = true
		firstRunChan           = make(chan struct{}, 1)
	)

	go func() {
		defer func() {
			err := multierr.Combine(
				pipeWriter.Close(),
				pipeReader.Close(),
			)
			if err != nil {
				select {
				case doneChan <- err:
				default:
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				doneChan <- ctx.Err()
				return
			default:
				req, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						doneChan <- nil
						return
					}
					doneChan <- err
					return
				}

				if firstRunFlag {
					firstRunFlag = false
					uploadRequest.Name = req.Metadata.Name
					uploadRequest.ContentType = req.Metadata.ContentType
					firstRunChan <- struct{}{}
				}

				if _, err := pipeWriter.Write(req.Data); err != nil {
					doneChan <- err
					return
				}
			}
		}
	}()

	<-firstRunChan
	uploadRequest.Reader = pipeReader

	return uploadRequest, doneChan, nil
}
