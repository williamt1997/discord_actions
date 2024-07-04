package bot

import (
	"encoding/json"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var Key string
var SubKey string
var Num int

// Function To Transform The Users Inputted Request Body Into JSON Format
func BodyParser(input string) ([]byte, string) {
	var str = ""
	//Regular Rexpression To Split Strings Into A Strings Slice Via Whitespace, Strings That Are Encapsulated Via Single Quotes And Strings that Are Encapsulated Via Square Brackets
	r := regexp.MustCompile(`(\[[^\]]+\]|'[^']+'|\S+)`)
	inputs := r.FindAllString(input, -1)
	//Creating A Map To Store JSON Body
	bodyMap := make(map[string]interface{})
	//Creating A Map To Store Nested JSON Objects: I.E. Key{subkey:value}
	subMap := make(map[string]interface{})

	//For loop that will itterate through the inputs slice
	for i := 0; i < len(inputs); i++ {
		//If Text Starts With A Dash/Hyphen: The Dash Represents The A Key In The Main Json Body
		if strings.HasPrefix(inputs[i], "-") {
			//Handle Errors When User Tries To Enter Key Without Value
			if i+1 >= len(inputs) {
				str = "Invalid JSON: JSON Cannot End With The Key And Only The Key"
				var fullBodyJson []byte
				return fullBodyJson, str
			}
			Key = strings.TrimPrefix(inputs[i], "-")
			Num = i + 1
			// Checking If The Next Item In The Slice Is A Sub Key
			if Num < len(inputs) && strings.HasSuffix(inputs[Num], ":") {
				// Resetting Nested Json Map To Allow For A New Subkey/Nest
				subMap = make(map[string]interface{})
			} else if Num < len(inputs) {
				//Checking If The Next Item In The Slice Is Encapsulated With Square Brackets
				// Square Brackets Are Used To Allow Users To Value Pairs With More Than One Word, Allowing For Complex Sentence Structures
				if strings.HasPrefix(inputs[Num], "[") && strings.HasSuffix(inputs[Num], "]") {
					// Trimming Brackets From Text And Assigning Values To With The Current Key In The Main Body Map
					bodyMap[Key] = strings.Trim(inputs[Num], "[]")
					//Checking if the next item in the slice is a comma seperated array which is encapsulated with single quotes
				} else if strings.Contains(inputs[Num], ",") && strings.Contains(inputs[Num], "'") {
					// Trimming Single Quotes From slice item and splitting text into a slice using a comma and a space.
					bodyList := strings.Split(strings.Trim(inputs[Num], "'"), ", ")
					// Assigning Values To With The Current Key In The Main Body Map
					bodyMap[Key] = bodyList
				} else {
					bodyMap[Key] = strings.Trim(inputs[Num], "', ")
				}
				i++
			}
		}
		//If Text Ends With A Colon: The Colon Represents The A Subkey Of A Nested JSON Object
		if strings.HasSuffix(inputs[i], ":") {
			//Handle Errors When User Tries To Enter SubKey Without Value
			if i+1 >= len(inputs) {
				str = "Invalid JSON: JSON Cannot End With The SubKey And Only The SubKey"
				var fullBodyJson []byte
				return fullBodyJson, str
			}
			SubKey = strings.TrimSuffix(inputs[i], ":")
			Num = i + 1
			// Check if there is a value for the sub-key
			if Num < len(inputs) {
				// Stop Adding Values To Subkey If Next Slice Item Is A Main Body Key
				if strings.Contains(inputs[Num], "-") {
					continue
					//Checking If The Next Item In The Slice Is Encapsulated With Square Brackets
					// Square Brackets Are Used To Allow Users To Value Pairs With More Than One Word, Allowing For Complex Sentence Structures
				} else if strings.HasPrefix(inputs[Num], "[") && strings.HasSuffix(inputs[Num], "]") {
					// Trimming Brackets From Text And Assigning Values To With The Current SubKey In The Nested JSON Object / Map
					subMap[SubKey] = strings.Trim(inputs[Num], "[]")
					//Checking if the next item in the slice is a comma seperated array which is encapsulated with single quotes
				} else if strings.Contains(inputs[Num], ",") && strings.Contains(inputs[Num], "'") {
					// Trimming Single Quotes From slice item and splitting text into a slice using a comma and a space.
					bodyList := strings.Split(strings.Trim(inputs[Num], "'"), ", ")
					//Assigning Values To With The Current SubKey In The Nested JSON Object / Map
					subMap[SubKey] = bodyList
				} else {
					subMap[SubKey] = strings.Trim(inputs[Num], "', ")
				}
				i++
			}
			//Assigning The Nested JSON Object To The Current Body With Under The Current Body Key
			bodyMap[Key] = subMap
		}
	}
	//Converting Map To JSON
	fullBodyJson, err := json.Marshal(bodyMap)
	//Handling Errors Where JSON Conversion Fails
	if err != nil {
		str = "Invalid JSON: JSON Cannot End With The SubKey And Only The SubKey"
		zap.L().Error("Error:", zap.Error(err))
	}

	return fullBodyJson, str
}
