package main

import (
	"bufio"
	"context"
	// "errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	// "net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "../passengerfeedback"
)

const (
	address = "localhost:8080"
)

func readNumber() (int32, error) {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	text = strings.TrimSuffix(text, "\n")
	i64, err := strconv.ParseInt(text, 10, 0)
	if err != nil {
		return -1, err
	}

	return int32(i64), nil
}

func readText() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(text, "\n"), nil
}

func selectMenu() int32 {
	fmt.Println("Actions Menu:")
	fmt.Println(" 1. Add feedback")
	fmt.Println(" 2. Get feedback by booking code")
	fmt.Println(" 3. Get feedback by passenger id")
	fmt.Println(" 4. Delete feedback by passenger id")
	fmt.Println(" 5. Exit")
	fmt.Print("Please select your action: ")

	num, err := readNumber()
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return -1
	}

	return num
}

func addPassengerFeedback(client pb.PassengerFeedbackServiceClient) {
	var feedback *pb.PassengerFeedback

	fmt.Println("Please enter feedback information")
	fmt.Print("Passenger ID: ")
	passengerID, errID := readNumber()
	if errID != nil {
		log.Fatalf("Couldn't add your ID %v", errID)
		return
	}

	fmt.Print("Booking Code: ")
	bookingCode, errCode := readText()
	if errCode != nil {
		log.Fatalf("Couldn't add your code %v", errCode)
		return
	}

	fmt.Print("Feedback: ")
	feedbackText, errFeedback := readText()
	if errFeedback != nil {
		log.Fatalf("Couldn't add your feedback %v", errFeedback)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	feedback = &pb.PassengerFeedback{
		PassengerID: passengerID,
		BookingCode: bookingCode,
		Feedback:    feedbackText,
	}
	response, err := client.AddPassengerFeedback(ctx, feedback)
	if err != nil {
		log.Fatalf("Couldn't add feedback: %v", err)
	}

	fmt.Println(response.ResponseCode.Message+".", "Please press enter to return actions menu.")
	fmt.Scanln()
}

func getFeedbackByBookingCode(client pb.PassengerFeedbackServiceClient) {
	var requestData *pb.GetPassengerFeedbackRequest

	fmt.Print("Please enter booking code: ")

	bookingCode, errCode := readText()
	if errCode != nil {
		log.Fatalf("Couldn't add your code %v", errCode)
		return
	}

	// another way to add into request
	requestData = &pb.GetPassengerFeedbackRequest{
		BookingCode: bookingCode,
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetPassengerFeedbackByBookingCode(ctx, requestData)
	if err != nil {
		log.Fatalf("Couldn't get feedback by booking code: %v", err)
	}

	if response.ResponseCode.Code != 0 {
		fmt.Println(response.ResponseCode.Message+".", "Please press enter to return actions menu.")
	} else {
		var feedback = response.PassengerFeedback
		fmt.Println("Results: ")
		fmt.Printf("\t Passenger ID: %d - Booking code: %s - Feedback: %s\n", feedback.PassengerID, feedback.BookingCode, feedback.Feedback)
		fmt.Println("Please press enter to return actions menu.")
	}

	fmt.Scanln()
}

func getFeedbackByPassengerID(client pb.PassengerFeedbackServiceClient) {
	var err error
	var passengerID int32
	var requestData *pb.GetPassengerFeedbackRequest

	passengerID, errID := readNumber()
	if errID != nil {
		log.Fatalf("Couldn't add your ID %v", errID)
		return
	}

	requestData = &pb.GetPassengerFeedbackRequest{
		PassengerID: passengerID,
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetPassengerFeedbackByPassengerId(ctx, requestData)

	if err != nil {
		log.Fatalf("Couldn't get feedback by passenger id: %v", err)
	}

	if response.ResponseCode.Code != 0 {
		fmt.Println(response.ResponseCode.Message+".", "Please press enter to return actions menu.")
	} else {
		var feedbacks = response.PassengerFeedbacks
		fmt.Println("Results: ")

		for _, feedback := range feedbacks {
			fmt.Printf("\t Passenger ID: %d - Booking code: %s - Feedback: %s\n", feedback.PassengerID, feedback.BookingCode, feedback.Feedback)
		}
		fmt.Println("Please press enter to return actions menu.")
	}

	fmt.Scanln()
}

func deleteFeedbackByPassengerID(client pb.PassengerFeedbackServiceClient) {
	var requestData *pb.DeletePassengerFeedbackRequest

	fmt.Print("Passenger ID: ")
	passengerID, errID := readNumber()

	if errID != nil {
		log.Fatalf("Couldn't add your ID %v", errID)
		return
	}

	requestData = &pb.DeletePassengerFeedbackRequest{
		PassengerID: passengerID,
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.DeletePassengerFeedbackPassengerId(ctx, requestData)
	if err != nil {
		log.Fatalf("Couldn't delete feedback: %v", err)
	}

	fmt.Println(response.ResponseCode.Message+".", "Please press enter to return actions menu.")
	fmt.Scanln()
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	var passengerFeedbackClient = pb.NewPassengerFeedbackServiceClient(conn)

	for {
		switch action := selectMenu(); action {
		case 1:
			addPassengerFeedback(passengerFeedbackClient)
		case 2:
			getFeedbackByBookingCode(passengerFeedbackClient)
		case 3:
			getFeedbackByPassengerID(passengerFeedbackClient)
		case 4:
			deleteFeedbackByPassengerID(passengerFeedbackClient)
		case 5:
			fmt.Println("Exit...")
			return
		default: // Cover -1 return
			fmt.Println("Action isn't exists. Please press enter to return actions menu.")
			fmt.Scanln()
		}

	}
}
