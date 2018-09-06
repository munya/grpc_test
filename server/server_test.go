package main

import (
	"context"
	pb "github.com/munya/grpc_test.git/pb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	dict := map[string]string{
		"monkey": "follow",
		"follow": "monkey",
	}
	mps := &Server{dict}
	It("should return follow for monkey", func() {
		phrase := &pb.Params{Message: "monkey"}
		result, err := mps.Send(context.Background(), phrase)
		Expect(err).To(BeNil())
		Expect(result.GetMessage()).To(Equal("follow"))
	})
	It("should return monkey for follow", func() {
		phrase := &pb.Params{Message: "follow"}
		result, err := mps.Send(context.Background(), phrase)
		Expect(err).To(BeNil())
		Expect(result.GetMessage()).To(Equal("monkey"))
	})
	It("should return empty string for random word", func() {
		phrase := &pb.Params{Message: "random word"}
		result, err := mps.Send(context.Background(), phrase)
		Expect(err).To(BeNil())
		Expect(result.GetMessage()).To(Equal(""))
	})
})