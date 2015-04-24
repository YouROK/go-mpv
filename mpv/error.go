package mpv

//#include <mpv/client.h>
import "C"

import (
	"fmt"
	"strconv"
)

type MpvError struct {
	errcode   int
	errstring string
}

//const char *mpv_error_string(int error);
func NewMpvError(errcode C.int) *MpvError {
	if int(errcode) == ERROR_SUCCESS {
		return nil
	}
	err := MpvError{int(errcode), ""}
	err.errstring = C.GoString(C.mpv_error_string(errcode))
	return &err
}

func (m *MpvError) Error() string {
	return fmt.Sprintln("Mpv error", strconv.Itoa(m.errcode), m.errstring)
}
