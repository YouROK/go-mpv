package mpv

/*
#include <mpv/client.h>
#include <stdlib.h>
#cgo LDFLAGS: -lmpv

char** makeCharArray1(int size) {
    return calloc(sizeof(char*), size);
}
void setArrayString1(char** a, int i, char* s) {
    a[i] = s;
}

*/
import "C"
import (
	"unsafe"
)

type Mpv struct {
	handle              *C.mpv_handle
	wakeup_callbackVar  interface{}
	wakeup_callbackFunc func(d interface{})
}

func Create() *Mpv {
	return &Mpv{C.mpv_create(), nil, nil}
}

func (m *Mpv) ClientName() string {
	return C.GoString(C.mpv_client_name(m.handle))
}

func (m *Mpv) Initialize() int {
	return int(C.mpv_initialize(m.handle))
}

func (m *Mpv) DetachDestroy() {
	C.mpv_detach_destroy(m.handle)
}

func (m *Mpv) TerminateDestroy() {
	C.mpv_terminate_destroy(m.handle)
}

func (m *Mpv) CreateClient(name string) *Mpv {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cmpv := C.mpv_create_client(m.handle, cname)
	if cmpv != nil {
		return &Mpv{cmpv, nil, nil}
	}
	return nil
}

func (m *Mpv) LoadConfigFile(fileName string) int {
	cfn := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfn))
	return int(C.mpv_load_config_file(m.handle, cfn))
}

func (m *Mpv) Suspend() {
	C.mpv_suspend(m.handle)
}

func (m *Mpv) Resume() {
	C.mpv_resume(m.handle)
}

func (m *Mpv) GetTimeUS() int64 {
	return int64(C.mpv_get_time_us(m.handle))
}

func (m *Mpv) SetOption(name string, format Format, data interface{}) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ptr := data2Ptr(format, data)
	return NewMpvError(C.mpv_set_option(m.handle, cname, C.mpv_format(format), ptr))
}

func (m *Mpv) SetOptionString(name, data string) int {
	cname := C.CString(name)
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cdata))
	return int(C.mpv_set_option_string(m.handle, cname, cdata))
}

func (m *Mpv) Command(command []string) int {
	cArray := C.makeCharArray1(C.int(len(command) + 1))
	if cArray == nil {
		panic("got NULL from calloc")
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString1(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	return int(C.mpv_command(m.handle, cArray))
	return -1
}

func (m *Mpv) CommandNode(command []string) int {
	//int mpv_command_node(mpv_handle *ctx, mpv_node *args, mpv_node *result);
	//TODO
	panic("Not supported command")
	return -1
}

func (m *Mpv) CommandString(command string) int {
	ccmd := C.CString(command)
	defer C.free(unsafe.Pointer(ccmd))
	return int(C.mpv_command_string(m.handle, ccmd))
}

func (m *Mpv) CommandAsync(replyUserdata uint64, command []string) int {
	cArray := C.makeCharArray1(C.int(len(command) + 1))
	if cArray == nil {
		panic("got NULL from calloc")
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString1(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	return int(C.mpv_command_async(m.handle, C.uint64_t(replyUserdata), cArray))
}

func (m *Mpv) CommandNodeAsync(command []string) int {
	//int mpv_command_node_async(mpv_handle *ctx, uint64_t reply_userdata, mpv_node *args);
	//TODO
	panic("Not supported command")
	return -1
}

func (m *Mpv) SetProperty(name string, format Format, data interface{}) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ptr := data2Ptr(format, data)
	return NewMpvError(C.mpv_set_property(m.handle, cname, C.mpv_format(format), ptr))
}

func (m *Mpv) SetPropertyString(name, data string) int {
	cname := C.CString(name)
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cdata))
	return int(C.mpv_set_property_string(m.handle, cname, cdata))
}

func (m *Mpv) SetPropertyAsync(name string, replyUserdata uint64, format Format, data interface{}) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	ptr := data2Ptr(format, data)
	return NewMpvError(C.mpv_set_property_async(m.handle,C.uint64_t(replyUserdata),cname,C.mpv_format(format),ptr))
}

func (m *Mpv) GetProperty(name string, format Format) (interface{}, error) {
	//int mpv_get_property(mpv_handle *ctx, const char *name, mpv_format format, void *data);
	//TODO
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var data unsafe.Pointer

	err := NewMpvError(C.mpv_get_property(m.handle,cname,C.mpv_format(format),data))

	return nil, err
}

func (m *Mpv) GetPropertyString(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(C.mpv_get_property_string(m.handle, cname))
}

func (m *Mpv) GetPropertyOsdString(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(C.mpv_get_property_osd_string(m.handle, cname))
}

func (m *Mpv) GetPropertyAsync(name string, format Format) (interface{}, int) {
	//int mpv_get_property_async(mpv_handle *ctx, uint64_t reply_userdata, const char *name, mpv_format format);
	//TODO
	return nil, -1
}

func (m *Mpv) ObserveProperty(replyUserdata uint64, name string, format Format) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.mpv_observe_property(m.handle, C.uint64_t(replyUserdata), cname, C.mpv_format(format)))
}

func (m *Mpv) UnObserveProperty(registeredReplyUserdata uint64) int {
	return int(C.mpv_unobserve_property(m.handle, C.uint64_t(registeredReplyUserdata)))
}

func (m *Mpv) RequestEvent(event EventId, enable bool) int {
	var en C.int = 0
	if enable {
		en = 1
	}
	return int(C.mpv_request_event(m.handle, C.mpv_event_id(event), en))
}

func (m *Mpv) RequestLogMessages(minLevel string) int {
	clevel := C.CString(minLevel)
	defer C.free(unsafe.Pointer(clevel))
	return int(C.mpv_request_log_messages(m.handle, clevel))
}

func (m *Mpv) WaitEvent(timeout float32) /*Event*/ {
	//TODO
	//mpv_event *mpv_wait_event(mpv_handle *ctx, double timeout);
}

func (m *Mpv) Wakeup() {
	C.mpv_wakeup(m.handle)
}

func (m *Mpv) SetWakeupCallback(callback func(d interface{}), d interface{}) {
	/*callbackFunc = callback
	callbackVar = d*/
	//	C.mpv_set_wakeup_callback(m.handle,,unsafe.Pointer(d))
	//TODO void mpv_set_wakeup_callback(mpv_handle *ctx, void (*cb)(void *d), void *d);
}

func (m *Mpv) GetWakeupPipe() int {
	return int(C.mpv_get_wakeup_pipe(m.handle))
}

func (m *Mpv) WaitAsyncRequests() {
	C.mpv_wait_async_requests(m.handle)
}

func (m *Mpv) GetSubApi(api SubApi) unsafe.Pointer {
	return unsafe.Pointer(C.mpv_get_sub_api(m.handle, C.mpv_sub_api(api)))
}

/*

void mpv_free(void *data);
void mpv_free_node_contents(mpv_node *node);
*/

func data2Ptr(format Format, data interface{}) unsafe.Pointer {
	var ptr unsafe.Pointer = nil
	switch format {
	case FORMAT_STRING, FORMAT_OSD_STRING:
		{
			val := C.CString(data.(string))
			ptr = unsafe.Pointer(&val)
			defer C.free(unsafe.Pointer(val))
		}
	case FORMAT_INT64:
		{
			val := C.int64_t(data.(int64))
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_DOUBLE:
		{
			val := C.double(data.(float32))
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_FLAG:
		{
			val := C.int(0)
			if data.(bool) {
				val = 1
			}
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_NODE, FORMAT_NODE_ARRAY, FORMAT_NODE_MAP, FORMAT_NONE:
		{
			//TODO
			panic("Not supported property")
		}
	}
	return ptr
}

func ptr2Data(format Format, data unsafe.Pointer) interface{} {
	var ptr unsafe.Pointer = nil
	switch format {
	case FORMAT_STRING, FORMAT_OSD_STRING:
		{

			val := C.CString(data.(string))
			ptr = unsafe.Pointer(&val)
			defer C.free(unsafe.Pointer(val))
		}
	case FORMAT_INT64:
		{
			val := C.int64_t(data.(int64))
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_DOUBLE:
		{
			val := C.double(data.(float32))
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_FLAG:
		{
			val := C.int(0)
			if data.(bool) {
				val = 1
			}
			ptr = unsafe.Pointer(&val)
		}
	case FORMAT_NODE, FORMAT_NODE_ARRAY, FORMAT_NODE_MAP, FORMAT_NONE:
		{
			//TODO
			panic("Not supported property")
		}
	}
	return ptr
}
