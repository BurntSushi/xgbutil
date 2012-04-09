/*
   This is a utility library designed to work with the X Go Binding. This 
   project's main goal is to make various X related tasks easier. For example, 
   binding keys, using the EWMH or ICCCM specs with the window manager, 
   moving/resizing windows, assigning function callbacks to particular events, 
   etc.
*/
package xgbutil

import (
	"code.google.com/p/jamslam-x-go-binding/xgb"
	"fmt"
	"log"
)

const MAX_REQ_SIZE = (1 << 16) * 4

// An XUtil represents the state of xgbutil. It keeps track of the current 
// X connection, the root window, event callbacks, key/mouse bindings, etc.
type XUtil struct {
	conn      *xgb.Conn
	quit      bool // when true, the main event loop will stop gracefully
	evqueue   []xgb.Event
	screen    *xgb.ScreenInfo
	root      xgb.Id
	eventTime xgb.Timestamp
	atoms     map[string]xgb.Id
	atomNames map[xgb.Id]string
	callbacks map[int]map[xgb.Id][]Callback // ev code -> win -> callbacks

	keymap *KeyboardMapping
	modmap *ModifierMapping

	keyRedirect xgb.Id
	keybinds    map[KeyBindKey][]KeyBindCallback
	keygrabs    map[KeyBindKey]int

	mousebinds    map[MouseBindKey][]MouseBindCallback
	mousegrabs    map[MouseBindKey]int
	mouseDrag     bool
	mouseDragStep MouseDragFun
	mouseDragEnd  MouseDragFun

	gc xgb.Id // a general purpose graphics context; used to paint images

	dummy xgb.Id // a dummy window used for mouse/key GRABs

	ignoreWindowErrors map[xgb.Id]bool // when true, errors don't go to stderr
}

type MouseDragFun func(xu *XUtil, rootX, rootY, eventX, eventY int)
type MouseDragBeginFun func(xu *XUtil, rootX, rootY,
	eventX, eventY int) (bool, xgb.Id)

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

// XError encapsulates any error returned by xgbutil.
type XError struct {
	funcName string     // some identifier so we know where the error comes from
	err      string     // free form string explaining the error
	XGBError *xgb.Error // error struct from XGB - to get the raw X error
}

// Error turns values of type *XError into a nice string.
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
			err:      fmt.Sprintf("%s: %v", fmt.Sprintf(err, params...), e),
			XGBError: e,
		}
	}

	panic(Xuerr("Xerr", "Unsupported error type: %T", err))
}

// Constructs an error struct from an error inside xgbutil (i.e., user error)
func Xuerr(funcName string, err string, params ...interface{}) *XError {
	return &XError{
		funcName: funcName,
		err:      fmt.Sprintf(err, params...),
		XGBError: nil,
	}
}

// IgnoreMods is a list of X modifiers that we don't want interfering
// with our mouse or key bindings. In particular, for each mouse or key binding 
// issued, there is a seperate mouse or key binding made for each of the 
// following modifiers.
var IgnoreMods []uint16 = []uint16{
	0,
	xgb.ModMaskLock,                // Num lock
	xgb.ModMask2,                   // Caps lock
	xgb.ModMaskLock | xgb.ModMask2, // Caps and Num lock
}

// Dial connects to the X server and creates a new XUtil.
func Dial(display string) (*XUtil, error) {
	c, err := xgb.Dial(display)

	if err != nil {
		return nil, err
	}

	// Initialize our central struct that stores everything.
	xu := &XUtil{
		conn:               c,
		quit:               false,
		evqueue:            make([]xgb.Event, 0),
		screen:             c.DefaultScreen(),
		root:               c.DefaultScreen().Root,
		eventTime:          xgb.Timestamp(0), // last event time
		atoms:              make(map[string]xgb.Id, 50),
		atomNames:          make(map[xgb.Id]string, 50),
		callbacks:          make(map[int]map[xgb.Id][]Callback, 33),
		keymap:             nil, // we don't have anything yet
		modmap:             nil,
		keyRedirect:        0,
		keybinds:           make(map[KeyBindKey][]KeyBindCallback, 10),
		keygrabs:           make(map[KeyBindKey]int, 10),
		mousebinds:         make(map[MouseBindKey][]MouseBindCallback, 10),
		mousegrabs:         make(map[MouseBindKey]int, 10),
		mouseDrag:          false,
		mouseDragStep:      nil,
		mouseDragEnd:       nil,
		ignoreWindowErrors: make(map[xgb.Id]bool, 10),
	}

	// Create a general purpose graphics context
	xu.gc = xu.conn.NewId()
	xu.conn.CreateGC(xu.gc, xu.root, xgb.GCForeground,
		[]uint32{xu.screen.WhitePixel})

	// Create a dummy window
	xu.dummy = xu.conn.NewId()
	xu.conn.CreateWindow(xu.Screen().RootDepth, xu.dummy, xu.RootWin(),
		-1000, -1000, 1, 1, 0,
		xgb.WindowClassInputOutput, xu.Screen().RootVisual,
		xgb.CWEventMask|xgb.CWOverrideRedirect,
		[]uint32{1, xgb.EventMaskPropertyChange})
	xu.conn.MapWindow(xu.dummy)

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
func (xu *XUtil) Conn() *xgb.Conn {
	return xu.conn
}

// Die forcefully shuts everything down.
func (xu *XUtil) Die() {
	xu.Conn().Close()
}

// Quit elegantly exits out of the main event loop.
func (xu *XUtil) Quit() {
	xu.quit = true
}

// Quitting returns whether it's time to quit.
func (xu *XUtil) Quitting() bool {
	return xu.quit
}

// Forces XGB to catch up with all events and synchronize.
func (xu *XUtil) Flush() {
	xu.conn.GetInputFocus()
}

// Enqueue queues up an event read from X.
func (xu *XUtil) Enqueue(ev xgb.Event) {
	xu.evqueue = append(xu.evqueue, ev)
}

// Dequeue pops an event from the queue and returns it.
func (xu *XUtil) Dequeue() xgb.Event {
	ev := xu.evqueue[0]
	xu.evqueue = xu.evqueue[1:]
	return ev
}

// DequeueAt removes a particular item from the queue
func (xu *XUtil) DequeueAt(i int) {
	xu.evqueue = append(xu.evqueue[:i], xu.evqueue[i+1:]...)
}

// QueueEmpty returns whether the event queue is empty or not.
func (xu *XUtil) QueueEmpty() bool {
	return len(xu.evqueue) == 0
}

// QueuePeek returns the current queue so we can examine it
func (xu *XUtil) QueuePeek() []xgb.Event {
	return xu.evqueue
}

// Screen returns the default screen
func (xu *XUtil) Screen() *xgb.ScreenInfo {
	return xu.screen
}

// RootWin returns the current root window.
func (xu *XUtil) RootWin() xgb.Id {
	return xu.root
}

// SetRootWin will change the current root window to the one provided.
// N.B. This probably shouldn't be used unless you're desperately trying
// to support multiple X screens. (This is *not* the same as Xinerama/RandR or
// TwinView. All of those have a single root window.)
func (xu *XUtil) SetRootWin(root xgb.Id) {
	xu.root = root
}

// GetTime gets the most recent time seen by an event.
func (xu *XUtil) GetTime() xgb.Timestamp {
	return xu.eventTime
}

// SetTime sets the most recent time seen by an event.
func (xu *XUtil) SetTime(t xgb.Timestamp) {
	xu.eventTime = t
}

// GC gets a general purpose graphics context that is typically used to simply
// paint images.
func (xu *XUtil) GC() xgb.Id {
	return xu.gc
}

// Dummy gets the id of the dummy window.
func (xu *XUtil) Dummy() xgb.Id {
	return xu.dummy
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

// Grabs the server. Everything becomes synchronous.
func (xu *XUtil) Grab() {
	xu.conn.GrabServer()
}

// Ungrabs the server.
func (xu *XUtil) Ungrab() {
	xu.conn.UngrabServer()
}

// IgnoredWindow returns true if the given client id's errors should be
// exempt from stderr and false otherwise.
func (xu *XUtil) IgnoredWindow(id xgb.Id) bool {
	ignored, ok := xu.ignoreWindowErrors[id]
	return ok && ignored
}

// IgnoreWindowErrorsAdd ignores all errors generated by the client ID
// provided. "ignore" means to not emit to stderr.
func (xu *XUtil) IgnoreWindowErrorsAdd(id xgb.Id) {
	xu.ignoreWindowErrors[id] = true
}

// IgnoreWindowErrorsRemove stops ignoring all errors generated by the client 
// ID provided. "ignore" means to not emit to stderr.
func (xu *XUtil) IgnoreWindowErrorsRemove(id xgb.Id) {
	delete(xu.ignoreWindowErrors, id)
}

// True utility/misc functions. Could be factored out to another package at 
// some point.

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
