/*
    This is a utility library designed to work with the X Go Binding. This 
    project's main goal is to make various X related tasks easier. For example, 
    binding keys, using the EWMH or ICCCM specs with the window manager, 
    moving/resizing windows, assigning function callbacks to particular events, 
    etc.
*/
package xgbutil

import (
    "fmt"
    "log"
    "code.google.com/p/x-go-binding/xgb"
)

// An XUtil represents the state of xgbutil. It keeps track of the current 
// X connection, the root window, event callbacks, key/mouse bindings, etc.
type XUtil struct {
    conn *xgb.Conn
    root xgb.Id
    atoms map[string]xgb.Id
    atomNames map[xgb.Id]string
    callbacks map[int]map[xgb.Id][]Callback // ev code -> win -> callbacks
}

// Callback is an interface that should be implemented by event callback 
// functions. Namely, to assign a function to a particular event/window
// combination, simply define a function with type '|Event|Fun' (pre-defined
// in xevent/callback.go), and call the 'Connect' method.
// The 'Run' method is used inside the Main event loop, and shouldn't be used
// by the user.
// Also, it is perfectly legitimate to connect to events that don't specify
// a window (like MappingNotify and KeymapNotify). In this case, simply
// use 'xgbutil.NoWindow' as the window id.
type Callback interface {
    Connect(xu *XUtil, win xgb.Id)
    Run(xu *XUtil, ev interface{})
}

// Sometimes we need to specify NO WINDOW when a window is typically
// expected. (Like connecting to MappingNotify or KeymapNotify events.)
// Use this value to do that.
var NoWindow xgb.Id = 0

type XError struct {
    funcName string // some identifier so we know where the error comes from
    err string // free form string explaining the error
    XGBError *xgb.Error // error struct from XGB - to get the raw X error
}

func (xe *XError) Error() string {
    return fmt.Sprintf("%s: %s", xe.funcName, xe.err)
}

// Constructs an error struct from an X error
func Xerr(xgberr interface{}, funcName string, err string,
          params ...interface{}) *XError {
    switch e := xgberr.(type) {
    case *xgb.Error:
        return &XError{
            funcName: funcName,
            err: fmt.Sprintf("%s: %v", fmt.Sprintf(err, params...), e),
            XGBError: e,
        }
    }

    panic(Xuerr("Xerr", "Unsupported error type: %T", err))
}

// Constructs an error struct from an error inside xgbutil (i.e., user error)
func Xuerr(funcName string, err string, params ...interface{}) *XError {
    return &XError{
        funcName: funcName,
        err: fmt.Sprintf(err, params...),
        XGBError: nil,
    }
}

// Dial connects to the X server and creates a new XUtil.
func Dial(display string) (*XUtil, error) {
    c, err := xgb.Dial(display)

    if err != nil {
        return nil, err
    }

    // Initialize our central struct that stores everything.
    xu := &XUtil{
        conn: c,
        root: c.DefaultScreen().Root,
        atoms: make(map[string]xgb.Id, 50), // start with a nice size
        atomNames: make(map[xgb.Id]string, 50),
        callbacks: make(map[int]map[xgb.Id][]Callback, 33),
    }

    // Register the Xinerama extension... because it doesn't cost much.
    err = xu.conn.RegisterExtension("XINERAMA")

    // If we can't register Xinerama, that's okay. Output something
    // and move on.
    if err != nil {
        log.Printf("WARNING: %s\n", err)
        log.Printf("MESSAGE: The 'xinerama' package cannot be used because " +
                   "the XINERAMA extension could not be loaded.")
    }

    return xu, nil
}

// Conn returns the xgb connection object.
func (xu *XUtil) Conn() (*xgb.Conn) {
    return xu.conn
}

// RootWin returns the current root window.
func (xu *XUtil) RootWin() (xgb.Id) {
    return xu.root
}

// SetRootWin will change the current root window to the one provided.
// N.B. This probably shouldn't be used unless you're desperately trying
// to support multiple X screens. (This is *not* the same as Xinerama/RandR or
// TwinView. All of those have a single root window.)
func (xu *XUtil) SetRootWin(root xgb.Id) {
    xu.root = root
}

// AttachCallback associates a (event, window) tuple with an event.
func (xu *XUtil) AttachCallback(evtype int, win xgb.Id, fun Callback) {
    // Do we need to allocate?
    if _, ok := xu.callbacks[evtype]; !ok {
        xu.callbacks[evtype] = make(map[xgb.Id][]Callback, 10)
    }
    if _, ok := xu.callbacks[evtype][win]; !ok {
        xu.callbacks[evtype][win] = make([]Callback, 0)
    }
    xu.callbacks[evtype][win] = append(xu.callbacks[evtype][win], fun)
}

// RunCallbacks executes every callback corresponding to a
// particular event/window tuple.
func (xu *XUtil) RunCallbacks(event interface{}, evtype int, win xgb.Id) {
    for _, cb := range xu.callbacks[evtype][win] {
        cb.Run(xu, event)
    }
}

// DetachWindow removes all callbacks associated with a particular window.
func (xu *XUtil) DetachWindow(win xgb.Id) {
    for evtype, _ := range xu.callbacks {
        delete(xu.callbacks[evtype], win)
    }
}

// GetAtom retrieves an atom identifier from a cache if it exists.
func (xu *XUtil) GetAtom(name string) (aid xgb.Id, ok bool) {
    aid, ok = xu.atoms[name]
    return
}

// GetAtomName retrieves an atom name from a cache if it exists.
func (xu *XUtil) GetAtomName(aid xgb.Id) (name string, ok bool) {
    name, ok = xu.atomNames[aid]
    return
}

// CacheAtom puts an atom into the cache.
func (xu *XUtil) CacheAtom(name string, aid xgb.Id) {
    xu.atoms[name] = aid
    xu.atomNames[aid] = name
}


// True utility/misc functions. Could be factored out at some point.

// BeSafe will recover from any panic produced by xgb or xgbutil and transform
// it into an idiomatic Go error as a second return value.
func BeSafe(err *error) {
    if r := recover(); r != nil {
        // If we get an error that isn't from xgbutil or xgb itself,
        // then let the panic happen.
        var ok bool
        *err, ok = r.(*XError)
        if !ok {
            *err, ok = r.(*xgb.Error)
            if !ok { // some other error, panic!
                panic(r)
            }
        }
    }
}

// put16 adds a 16 bit integer to a byte slice.
// Lifted from the xgb package.
func Put16(buf []byte, v uint16) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
}

// put32 adds a 32 bit integer to a byte slice.
// Lifted from the xgb package.
func Put32(buf []byte, v uint32) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

// get16 extracts a 16 bit integer from a byte slice.
// Lifted from the xgb package.
func Get16(buf []byte) uint16 {
	v := uint16(buf[0])
	v |= uint16(buf[1]) << 8
	return v
}

// get32 extracts a 32 bit integer from a byte slice.
// Lifted from the xgb package.
func Get32(buf []byte) uint32 {
	v := uint32(buf[0])
	v |= uint32(buf[1]) << 8
	v |= uint32(buf[2]) << 16
	v |= uint32(buf[3]) << 24
	return v
}

