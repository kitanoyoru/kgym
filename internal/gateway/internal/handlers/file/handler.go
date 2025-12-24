package file

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pbFile "github.com/kitanoyoru/kgym/contracts/protobuf/gen/go/file/v1"
)

type Handler struct {
	mux *runtime.ServeMux

	grpcFileServiceClient pbFile.FileServiceClient

	bodyLimit, chunkSize int
}

type Config struct {
	GRPCEndpoint    string
	GRPCDialOptions []grpc.DialOption

	BodyLimit, ChunkSize int
}

func New(ctx context.Context, mux *runtime.ServeMux, cfg Config) (*Handler, error) {
	conn, err := grpc.NewClient(cfg.GRPCEndpoint, cfg.GRPCDialOptions...)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()

		if cerr := conn.Close(); cerr != nil {
			grpclog.Errorf("Failed to close conn to %s: %v", cfg.GRPCEndpoint, cerr)
		}
	}()

	return &Handler{
		mux:                   mux,
		grpcFileServiceClient: pbFile.NewFileServiceClient(conn),
		bodyLimit:             cfg.BodyLimit,
		chunkSize:             cfg.ChunkSize,
	}, nil
}

func (h *Handler) UploadUserAvatar() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(h.mux, r)

		md := map[string]string{
			"x-request-id":  r.Header.Get("X-Request-ID"),
			"x-platform":    r.Header.Get("X-Platform"),
			"x-app-version": r.Header.Get("X-App-Version"),
		}
		authorization := r.Header.Get("Authorization")
		if authorization != "" {
			if token, err := r.Cookie("access_token"); err == nil {
				md["authorization"] = "Bearer " + token.Value
			}
		}
		md["authorization"] = authorization

		ctx := metadata.NewOutgoingContext(r.Context(), metadata.New(md))

		err := r.ParseMultipartForm(int64(h.bodyLimit))
		if err != nil {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
			return
		}
		if r.MultipartForm == nil {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, status.Errorf(codes.InvalidArgument, "multipart form is nil"))
		}
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				grpclog.Errorf("Failed to remove multipart form: %v", err)
			}
		}()

		files := r.MultipartForm.File["file"]
		if len(files) == 0 {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, status.Errorf(codes.InvalidArgument, "file is required"))
			return
		}

		file := files[0]

		f, err := file.Open()
		if err != nil {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				grpclog.Errorf("Failed to close file: %v", err)
			}
		}()

		stream, err := h.grpcFileServiceClient.UploadUserAvatar(ctx)
		if err != nil {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
			return
		}

		buf := make([]byte, 0, h.chunkSize)
		for {
			n, err := f.Read(buf[:cap(buf)])
			buf = buf[:n]
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
				return
			}

			err = stream.Send(&pbFile.UploadUserAvatar_Request{
				// TODO: add metadata
				Data: buf,
			})
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
				return
			}

			if n < h.chunkSize {
				break
			}
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			runtime.HTTPError(ctx, h.mux, outboundMarshaler, w, r, err)
			return
		}

		runtime.ForwardResponseMessage(ctx, h.mux, inboundMarshaler, w, r, res, h.mux.GetForwardResponseOptions()...)
	}
}
