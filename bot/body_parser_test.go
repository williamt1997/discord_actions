package bot

import (
	"testing"
)

func TestBody_ParserCorrectFormat(t *testing.T) {
	msg, str := Body_Parser("!gobot test test -tester test hello: 'test1, test2, test3' testy -testertwo hi")
	message := string(msg)
	wantone := "{\"tester\":{\"hello\":[\"test1\",\"test2\",\"test3\"]},\"testertwo\":\"hi\"}"
	wanttwo := ""
	if message != wantone {
		t.Errorf("\n\nError: Body Parser Failed To Convert User Command Into JSON:\nWhat We Wanted: %q\nWhat We Got: %q", wantone, string(msg))
	}
	if str != wanttwo {
		t.Errorf("\n\nError: Str Should Remain Empty Only Unill A User Input Error:\nWhat We Wanted: %q\nWhat We Got: %q", wanttwo, str)
	}
}
func TestBody_ParserIncorrectFormatKey(t *testing.T) {
	msg, str := Body_Parser("!gobot test test -tester")
	message := string(msg)
	wantone := ""
	wanttwo := "Invalid JSON: JSON Cannot End With The Key And Only The Key"
	if message != wantone {
		t.Errorf("\n\nError: Body Parser Should Return Empty Byte If Error:\nWhat We Wanted: %q\nWhat We Got: %q", wantone, string(msg))
	}
	if str != wanttwo {
		t.Errorf("\n\nError: Str Should Return Error Message If User Parses A Key Without A Value:\nWhat We Wanted: %q\nWhat We Got: %q", wanttwo, str)
	}
}

func TestBody_ParserIncorrectFormatSubKey(t *testing.T) {
	msg, str := Body_Parser("!gobot test test -tester hello:")
	message := string(msg)
	wantone := ""
	wanttwo := "Invalid JSON: JSON Cannot End With The SubKey And Only The SubKey"
	if message != wantone {
		t.Errorf("\n\nError: Body Parser Should Return Empty Byte If Error:\nWhat We Wanted: %q\nWhat We Got: %q", wantone, string(msg))
	}
	if str != wanttwo {
		t.Errorf("\n\nError: Str Should Return Error Message If User Parses A SubKey Without A Value:\nWhat We Wanted: %q\nWhat We Got: %q", wanttwo, str)
	}
}
