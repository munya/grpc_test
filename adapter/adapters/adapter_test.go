package adapters

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"time"
	pb "github.com/munya/grpc_test.git/pb"
	"github.com/munya/grpc_test.git/adapter/mocks"
)

var _ = Describe("Test adapter", func() {

	mockCtrl := gomock.NewController(GinkgoT())
	forwarderMock := mocks.NewMockSender(mockCtrl)
	mockCtrl2 := gomock.NewController(GinkgoT())
	serverClientMock := mocks.NewMockServerClient(mockCtrl2)

	outDict := MessageMap{
		"monkey": "marco",
		"follow": "polo",
	}

	inDict := MessageMap{
		"marco": "monkey",
		"polo":  "follow",
	}

	It("should fail if ServerClient is nil", func() {
		a, err := NewAdapter(nil)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring("Can't establish connection to Server. ServerClient is nil."))
		Expect(a).To(BeNil())
	})

	It("DictInLookupAdapter should map marco -> monkey", func() {
		a, _ := NewDictInLookupAdapter(forwarderMock, inDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		forwarderMock.EXPECT().Send(ctx, gomock.Eq(&pb.Params{Message: "monkey"}))
		a.Send(ctx, &pb.Params{Message: "marco"})
	})

	It("DictInLookupAdapter should map polo -> follow", func() {
		a, _ := NewDictInLookupAdapter(forwarderMock, inDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		forwarderMock.EXPECT().Send(ctx, gomock.Eq(&pb.Params{Message: "follow"}))
		a.Send(ctx, &pb.Params{Message: "polo"})
	})

	It("DictInLookupAdapter should return an error if no key found", func() {
		a, _ := NewDictInLookupAdapter(forwarderMock, inDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		_, err := a.Send(ctx, &pb.Params{Message: "yolo"})
		Expect(err).ToNot(BeNil())
	})

	It("dictOutLookupAdapter should map follow -> polo", func() {
		adapter, _ := NewAdapter(serverClientMock)
		a, _ := NewDictOutLookupAdapter(adapter, outDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		serverClientMock.EXPECT().Send(ctx, gomock.Eq(&pb.Params{Message: "monkey"})).Return(&pb.Params{Message: "follow"}, nil)
		res, err := a.Send(ctx, &pb.Params{Message: "monkey"})
		Expect(res.Message).To(Equal("polo"))
		Expect(err).To(BeNil())
	})
	It("dictOutLookupAdapter should map monkey -> marco", func() {
		adapter, _ := NewAdapter(serverClientMock)
		a, _ := NewDictOutLookupAdapter(adapter, outDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		serverClientMock.EXPECT().Send(ctx, gomock.Eq(&pb.Params{Message: "follow"})).Return(&pb.Params{Message: "monkey"}, nil)
		res, err := a.Send(ctx, &pb.Params{Message: "follow"})
		Expect(res.Message).To(Equal("marco"))
		Expect(err).To(BeNil())
	})
	It("dictOutLookupAdapter should return error if no key found", func() {
		adapter, _ := NewAdapter(serverClientMock)
		a, _ := NewDictOutLookupAdapter(adapter, outDict)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		serverClientMock.EXPECT().Send(ctx, gomock.Any()).Return(nil, errors.New("No key found"))
		res, err := a.Send(ctx, &pb.Params{Message: "yolo"})
		Expect(res).To(BeNil())
		Expect(err).NotTo(BeNil())
	})

})
