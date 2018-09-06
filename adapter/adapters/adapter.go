package adapters

import (
	"errors"
	"fmt"
	"github.com/munya/grpc_test.git/pb"
	"golang.org/x/net/context"
	"log"
)

type MessageMap map[string]string

type Sender interface {
	Send(context.Context, *pbservice.Params) (*pbservice.Params, error)
}

type adapter struct {
	serverClient pbservice.ServerClient
}

func (a *adapter) Send(ctx context.Context, phrase *pbservice.Params) (*pbservice.Params, error) {
	res, err := a.serverClient.Send(ctx, phrase)
	return res, err
}

func NewAdapter(sc pbservice.ServerClient) (Sender, error) {
	if sc == nil {
		return nil, errors.New("Can't establish connection to Server. ServerClient is nil.")
	}
	return &adapter{serverClient: sc}, nil
}



type dictInLookupAdapter struct {
	Sender
	inLookup MessageMap
}

func (a *dictInLookupAdapter) Send(ctx context.Context, phrase *pbservice.Params) (*pbservice.Params, error) {
	message, ok := a.inLookup[phrase.Message]
	if !ok {
		log.Printf("missing word in mapping: %s", phrase.Message)
		return nil, fmt.Errorf("missing word '%s' in lookup", phrase.Message)
	}
	return a.Sender.Send(ctx, &pbservice.Params{Message: message})
}

func NewDictInLookupAdapter(f Sender, d MessageMap) (Sender, error) {
	return &dictInLookupAdapter{
		inLookup:  d,
		Sender: f,
	}, nil
}




type dictOutLookupAdapter struct {
	Sender
	outLookup MessageMap
}

func (a *dictOutLookupAdapter) Send(ctx context.Context, phrase *pbservice.Params) (*pbservice.Params, error) {
	res, err := a.Sender.Send(ctx, phrase)
	if err != nil {
		return nil, err
	}
	return &pbservice.Params{Message: a.outLookup[res.Message]}, err
}

func NewDictOutLookupAdapter(f Sender, d MessageMap) (Sender, error) {
	return &dictOutLookupAdapter{
		outLookup: d,
		Sender: f,
	}, nil
}