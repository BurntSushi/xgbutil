package xgbutil

import (
	"log"
	"os"

	"github.com/BurntSushi/xgb"
)

// Logger is used through xgbutil when messages need to be emitted to stderr.
var Logger = log.New(os.Stderr, "XGBUTIL: ", 0)

// The current maximum request size. I think we can expand this with
// BigReq, but it probably isn't worth it at the moment.
const MaxReqSize = (1 << 16) * 4

// An XUtil represents the state of xgbutil. It keeps track of the current 
// X connection, the root window, event callbacks, key/mouse bindings, etc.
type XUtil struct {
	conn      *xgb.Conn
	quit      bool // when true, the main event loop will stop gracefully
	evqueue   []eventOrError
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

	errorHandler ErrorHandlerFun
}

// NewConn connects to the X server using the DISPLAY environment variable
// and creates a new XUtil.
func NewConn() (*XUtil, error) {
	return NewConnDisplay("")
}

// NewConnDisplay connects to the X server and creates a new XUtil.
// See the XGB documentation for xgb.NewConnDisplay for which values of
// 'display' are supported.
func NewConnDisplay(display string) (*XUtil, error) {
	c, err := xgb.NewConnDisplay(display)

	if err != nil {
		return nil, err
	}

	// Initialize our central struct that stores everything.
	xu := &XUtil{
		conn:          c,
		quit:          false,
		evqueue:       make([]eventOrError, 0),
		screen:        c.DefaultScreen(),
		root:          c.DefaultScreen().Root,
		eventTime:     xgb.Timestamp(0), // last event time
		atoms:         make(map[string]xgb.Id, 50),
		atomNames:     make(map[xgb.Id]string, 50),
		callbacks:     make(map[int]map[xgb.Id][]Callback, 33),
		keymap:        nil, // we don't have anything yet
		modmap:        nil,
		keyRedirect:   0,
		keybinds:      make(map[KeyBindKey][]KeyBindCallback, 10),
		keygrabs:      make(map[KeyBindKey]int, 10),
		mousebinds:    make(map[MouseBindKey][]MouseBindCallback, 10),
		mousegrabs:    make(map[MouseBindKey]int, 10),
		mouseDrag:     false,
		mouseDragStep: nil,
		mouseDragEnd:  nil,
		errorHandler:  defaultErrorHandler,
	}

	// Create a general purpose graphics context
	xu.gc, err = xu.conn.NewId()
	if err != nil {
		return nil, err
	}
	xu.conn.CreateGC(xu.gc, xu.root, xgb.GcForeground,
		[]uint32{xu.screen.WhitePixel})

	// Create a dummy window
	xu.dummy, err = xu.conn.NewId()
	if err != nil {
		return nil, err
	}
	xu.conn.CreateWindow(xu.Screen().RootDepth, xu.dummy, xu.RootWin(),
		-1000, -1000, 1, 1, 0,
		xgb.WindowClassInputOutput, xu.Screen().RootVisual,
		xgb.CwEventMask|xgb.CwOverrideRedirect,
		[]uint32{1, xgb.EventMaskPropertyChange})
	xu.conn.MapWindow(xu.dummy)

	// Register the Xinerama extension... because it doesn't cost much.
	err = xu.conn.XineramaInit()

	// If we can't register Xinerama, that's okay. Output something
	// and move on.
	if err != nil {
		Logger.Printf("WARNING: %s\n", err)
		Logger.Printf("MESSAGE: The 'xinerama' package cannot be used " +
			"because the XINERAMA extension could not be loaded.")
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

// Enqueue queues up an event read from X.
func (xu *XUtil) Enqueue(everr eventOrError) {
	xu.evqueue = append(xu.evqueue, everr)
}

// Dequeue pops an event from the queue and returns it.
func (xu *XUtil) Dequeue() eventOrError {
	everr := xu.evqueue[0]
	xu.evqueue = xu.evqueue[1:]
	return everr
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
func (xu *XUtil) QueuePeek() []eventOrError {
	return xu.evqueue
}

// ErrorHandlerFun is the type of function required to handle errors that
// come in through the main event loop.
type ErrorHandlerFun func(err xgb.Error)

// ErrorHandlerSet sets the default error handler for errors that come
// into the main event loop. (This may be removed in the future in favor
// of a particular callback interface like events, but these sorts of errors
// aren't handled often in practice, so maybe not.)
// This is only called for errors returned from unchecked (asynchronous error
// handling) requests.
func (xu *XUtil) ErrorHandlerSet(fun ErrorHandlerFun) {
	xu.errorHandler = fun
}

// ErrorHandlerGet retrieves the default error handler.
func (xu *XUtil) ErrorHandlerGet() ErrorHandlerFun {
	return xu.errorHandler
}

// defaultErrorHandler just emits errors to stderr.
func defaultErrorHandler(err xgb.Error) {
	Logger.Println(err)
}

// eventOrError is a struct that contains either an event value or an error
// value. It is an error to contain both. Containing neither indicates an
// error too.
type eventOrError struct {
	Event xgb.Event
	Err   xgb.Error
}

// newEventOrError creates a new eventOrError value.
func NewEventOrError(event xgb.Event, err xgb.Error) eventOrError {
	return eventOrError{
		Event: event,
		Err:   err,
	}
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
