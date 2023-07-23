package model

import (
	"fmt"
	gen "github.com/alexvishnevskiy/twitter-clone/gen/api/tweets"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MediaToProto converts a Media struct into a
// generated proto counterpart.
func MediaToProto(m *Media) *gen.Media {
	protoTimestamp := timestamppb.New(m.CreatedAt)
	if err := protoTimestamp.CheckValid(); err != nil {
		fmt.Println("Error converting time.Time to timestamppb.Timestamp:", err)
	}

	return &gen.Media{
		Media:     m.Media,
		Content:   m.Content,
		CreatedAt: protoTimestamp,
	}
}

// MediaFromProto converts a proto struct into a
// media counterpart.
func MediaFromProto(m *gen.Media) *Media {
	return &Media{
		Media:     m.Media,
		Content:   m.Content,
		CreatedAt: m.CreatedAt.AsTime(),
	}
}
