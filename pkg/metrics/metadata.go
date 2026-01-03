package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"
)

const (
	MetadataRequestID  = "x-request-id"
	MetadataPlatform   = "x-platform"
	MetadataAppVersion = "x-app-version"
)

var AllMetadataFields = []string{
	MetadataRequestID,
	MetadataPlatform,
	MetadataAppVersion,
}

func ExtractLabelsFromMetadata(ctx context.Context) prometheus.Labels {
	labels := prometheus.Labels{
		MetadataRequestID:  "",
		MetadataPlatform:   "",
		MetadataAppVersion: "",
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return labels
	}

	if values := md.Get(MetadataRequestID); len(values) > 0 {
		labels[MetadataRequestID] = values[0]
	}
	if values := md.Get(MetadataPlatform); len(values) > 0 {
		labels[MetadataPlatform] = values[0]
	}
	if values := md.Get(MetadataAppVersion); len(values) > 0 {
		labels[MetadataAppVersion] = values[0]
	}

	return labels
}
