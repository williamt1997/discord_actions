package bot

import (
	"encoding/json"
	"fmt"
	"log"
)

// Function To Convert HTTP Json Response Data Into String Format
func BodyReader(body []byte) string {
	var outer string
	var inner string
	var objmap []map[string]interface{}
	//Unmarshalling JSON Data Into A Map: objmap
	if err := json.Unmarshal(body, &objmap); err != nil {
		log.Fatal(err)
	}
	//Looping through each key for the main JSON Items
	for i := range objmap {
		//Looping through nested JSON items and subkeys
		for j := range objmap[i] {
			// Formatting And Appending Each Submap item To Inner String
			inner += fmt.Sprintf("%v", objmap[i][j]) + "\n\n"
		}
		// Formatting And Appending Outer Map With Each Formatted Submap
		outer += inner + "\n\n\n"
		// Resetting inner For Next Submap
		inner = ""
	}
	message := outer
	return message
}
