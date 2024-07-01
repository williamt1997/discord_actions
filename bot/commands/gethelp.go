package commands

import (
	"bytes"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// Get Help Function When A User Successfully Enters A Valid Microservice Name But An Incorrect Endpoint.
// I.E. !gobot <microservice_name> <microservice_endpoint: bad_endpoint324432432dsa>
func Get_Help(url string) ([]byte, string) {
	var msg []byte
	var str string
	// Preparing To Make A POST Request To The Microservice ENDPOINT Of Help
	//Setting POST Request BODY As Blank Byte
	body := new(bytes.Buffer)
	//Sending a POST request with an empty body to the microservice's help endpoint to check its availability And To Get Help Response Data From Endpoint
	resp, err := http.Post(url, "application/json", body)
	// Check And Handle Errors Whilst Making The Post Request
	if err != nil {
		str = "Microservice Connection Issue! Report This to An Admin"
		zap.L().Error("Error", zap.Error(err))
		return msg, str
	} else {
		// Check if the response status code is less than 400: Successful Request
		if resp.StatusCode < 400 {
			// Closing The Response Body After Reading
			defer resp.Body.Close()
			// Reading The Response Body
			body, err := io.ReadAll(resp.Body)
			// Check And Handle Errors When Reading The Response Body Fails
			if err != nil {
				str = "Error Reading Response Body! Report This to An Admin"
				zap.L().Error("Response Read Error", zap.Error(err))
				return msg, str
			} else {
				msg := body
				str = ""
				return msg, str
			}
			// If HTTP Request Returns A Status Code Greater Or Equal To 400 For The Help Endpoint: Unuccessful Request
		} else {
			str = "Help Endpoint Not found Either! Report This to An Admin"
			return msg, str
		}
	}
}
