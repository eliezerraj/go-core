package gprc

import (
	"time"
	"context"
	"github.com/rs/zerolog/log"

	"github.com/golang/protobuf/jsonpb"	
	pb "github.com/golang/protobuf/proto"
	
	"google.golang.org/grpc"	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var childLogger = log.With().Str("component","go-core").Str("package", "grpc").Logger()

type GrpcClientWorker struct {
	GrcpClient		*grpc.ClientConn
}

// About convert proto to json
func (s *GrpcClientWorker) ProtoToJSON(msg pb.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	return marshaler.MarshalToString(msg)
}

// About convert json to proto
func (s *GrpcClientWorker) JSONToProto(data string, msg pb.Message) error {
	return jsonpb.UnmarshalString(data, msg)
}

// About start a grpc client
func (s *GrpcClientWorker) StartGrpcClient(host string) (*GrpcClientWorker, error){
	childLogger.Debug().Str("func","StartGrpcClient").Interface("host",host).Send()

	// Prepare options
	var opts []grpc.DialOption
	opts = append(opts, grpc.FailOnNonTempDialError(true)) // Wait for ready
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)) // 
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(5*time.Second))
	opts = append(opts, grpc.WithBlock()) // Wait for ready
	
	// Dail a server
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
	  childLogger.Error().Err(err).Msg("erro connect to grpc server")
	  return nil, err
	}

	return &GrpcClientWorker{
		GrcpClient : conn,
	}, nil
}

// About get connection
func (s *GrpcClientWorker) TestConnection(ctx context.Context) (error) {
	childLogger.Debug().Str("func","TestConnection").Send()
	
	if (s.GrcpClient == nil){
		return status.Errorf(codes.Internal, "client grpc is nul")
	}

	client := grpc_health_v1.NewHealthClient(s.GrcpClient)
	_, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
	if err != nil {
		childLogger.Error().Err(err).Msg("failed to gPRC test connection")
		return err
	}

	return nil
}

// About close connection
func (s *GrpcClientWorker) CloseConnection() () {
	childLogger.Debug().Str("func","CloseConnection").Send()

	if err := s.GrcpClient.Close(); err != nil {
		childLogger.Error().Err(err).Msg("failed to close gPRC connection")
	}
}