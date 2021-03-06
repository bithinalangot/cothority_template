package service

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"time"

	"github.com/dedis/cothority_template"
	"github.com/dedis/cothority_template/protocol"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

// Used for tests
var templateID onet.ServiceID

func init() {
	var err error
	templateID, err = onet.RegisterNewService(template.ServiceName, newService)
	log.ErrFatal(err)
}

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor
	path string
	// Count holds the number of calls to 'ClockRequest'
	Count int
}

// ClockRequest starts a template-protocol and returns the run-time.
func (s *Service) ClockRequest(req *template.ClockRequest) (*template.ClockResponse, onet.ClientError) {
	s.Count++
	tree := req.Roster.GenerateNaryTreeWithRoot(2, s.ServerIdentity())
	pi, err := s.CreateProtocol(protocol.Name, tree)
	if err != nil {
		return nil, onet.NewClientError(err)
	}
	start := time.Now()
	pi.Start()
	resp := &template.ClockResponse{
		Children: <-pi.(*protocol.Template).ChildCount,
	}
	resp.Time = time.Now().Sub(start).Seconds()
	return resp, nil
}

// CountRequest returns the number of instantiations of the protocol.
func (s *Service) CountRequest(req *template.CountRequest) (*template.CountResponse, onet.ClientError) {
	return &template.CountResponse{s.Count}, nil
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolOnet, this will not be called, as the Onet will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	log.Lvl3("Not templated yet")
	return nil, nil
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *onet.Context) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	if err := s.RegisterHandlers(s.ClockRequest, s.CountRequest); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}
