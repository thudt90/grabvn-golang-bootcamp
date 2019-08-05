package main

import (
	pb "../passengerfeedback"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":8080"
)

// PassengerFeedbackServerImp implements server
type PassengerFeedbackServerImp struct {
	FeedbackMap map[string]*pb.PassengerFeedback
}

//
var SuccessResponse = &pb.ResponseCode{
	Code:    0,
	Message: "Success",
}

//
var NotFoundFeedbackResponse = &pb.ResponseCode{
	Code:    1,
	Message: "Not found Feedback",
}

//
var ExistsFeedbackResponse = &pb.ResponseCode{
	Code:    2,
	Message: "Exists Feedback",
}

// AddPassengerFeedback add feedback
func (s *PassengerFeedbackServerImp) AddPassengerFeedback(ctx context.Context, pf *pb.PassengerFeedback) (*pb.AddPassengerFeedbackResult, error) {
	var responseCode = SuccessResponse

	if _, ok := s.FeedbackMap[pf.BookingCode]; ok {
		responseCode = ExistsFeedbackResponse
	} else {
		s.FeedbackMap[pf.BookingCode] = pf
	}

	return &pb.AddPassengerFeedbackResult{
		ResponseCode: responseCode,
	}, nil
}

// GetPassengerFeedbackByBookingCode is get feedback
func (s *PassengerFeedbackServerImp) GetPassengerFeedbackByBookingCode(ctx context.Context, pfr *pb.GetPassengerFeedbackRequest) (*pb.GetPassengerFeedbackCodeResponse, error) {
	var responseCode = NotFoundFeedbackResponse
	var feedback *pb.PassengerFeedback

	if v, ok := s.FeedbackMap[pfr.BookingCode]; ok {
		responseCode = SuccessResponse
		feedback = v
	}

	return &pb.GetPassengerFeedbackCodeResponse{
		ResponseCode:      responseCode,
		PassengerFeedback: feedback,
	}, nil
}

// func GetPassengerFeedbackByPassengerId should be GetPassengerFeedbackByPassengerID
func (s *PassengerFeedbackServerImp) GetPassengerFeedbackByPassengerId(ctx context.Context, pfr *pb.GetPassengerFeedbackRequest) (*pb.GetPassengerFeedbackIDResponse, error) {
	var responseCode = SuccessResponse
	var feedbacks []*pb.PassengerFeedback

	for _, v := range s.FeedbackMap {
		if v.PassengerID == pfr.PassengerID {
			feedbacks = append(feedbacks, v)
		}
	}

	if len(feedbacks) == 0 {
		responseCode = NotFoundFeedbackResponse
	}

	return &pb.GetPassengerFeedbackIDResponse{
		ResponseCode:       responseCode,
		PassengerFeedbacks: feedbacks,
	}, nil
}

// DeletePassengerFeedbackPassengerId is delete passenger deleteFeedbackByPassengerID
func (s *PassengerFeedbackServerImp) DeletePassengerFeedbackPassengerId(ctx context.Context, pfr *pb.DeletePassengerFeedbackRequest) (*pb.DeletePassengerFeedbackResponse, error) {
	var bookingCodes []string

	for _, v := range s.FeedbackMap {
		if v.PassengerID == pfr.PassengerID {
			bookingCodes = append(bookingCodes, v.BookingCode)
			//fmt.Println("Add ", v.BookingCode)
		}
	}

	for _, bookingCode := range bookingCodes {
		delete(s.FeedbackMap, bookingCode)
	}

	return &pb.DeletePassengerFeedbackResponse{
		ResponseCode: SuccessResponse,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPassengerFeedbackServiceServer(s, &PassengerFeedbackServerImp{
		FeedbackMap: make(map[string]*pb.PassengerFeedback),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
