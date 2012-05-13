package xgbutil

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
)

// Logger is used through xgbutil when messages need to be emitted to stderr.
var Logger = log.New(os.Stderr, "XGBUTIL: ", 0)

// The current maximum request size. I think we can expand this with
// BigReq, but it probably isn't worth it at the moment.
const MaxReqSize = (1 << 16) * 4

// An XUtil represents the state of xgbutil. It keeps track of the current 
// X connection, the root window, event callbacks, key/mouse bindings, etc.
type XUtil struct {
	// conn is the XGB connection object used to issue protocol requests.
	conn *xgb.Conn

	// quit can be set to true, and the main event loop will finish processing
	// the current event, and gracefully quit afterwards.
	quit bool // when true, the main event loop will stop gracefully

	// setup contains all the setup information retrieved at connection time.
	setup *xproto.SetupInfo

	// screen is a simple alias to the default screen info.
	screen *xproto.ScreenInfo

	// root is an alias to the default root window.
	root xproto.Window

	// atoms is a cache of atom names to resource identifiers. This minimizes
	// round trips to the X server, since atom identifiers never change.
	atoms    map[string]xproto.Atom
	atomsLck *sync.RWMutex

	// atomNames is a cache just like 'atoms', but in the reverse direction.
	atomNames    map[xproto.Atom]string
	atomNamesLck *sync.RWMutex

	// evqueue is the queue that stores the results of xgb.WaitForEvent.
	// Namely, each value is either an Event *or* an Error.
	// I didn't see any particular reason to use a channel for this.
	evqueue []eventOrError

	// callbacks is a map of event numbers to a map of window identifiers
	// to callback functions.
	// This is the data structure that stores all callback functions, where
	// a callback function is always attached to a (event, window) tuple.
	callbacks    map[int]map[xproto.Window][]Callback
	callbacksLck *sync.RWMutex

	// eventTime is the last time recorded by an event. It is automatically
	// updated if xgbutil's main event loop is used.
	eventTime xproto.Timestamp

	// Keymap corresponds to xgbutil's current conception of the keyboard
	// mapping. It is automatically kept up-to-date if xgbutil's event loop
	// is used.
	// It is exported for use in the keybind package. It should not be
	// accessed directly. Instead, use keybind.KeyMapGet.
	Keymap *KeyboardMapping

	// Modmap corresponds to xgbutil's current conception of the modifier key
	// mapping. It is automatically kept up-to-date if xgbutil's event loop
	// is used.
	// It is exported for use in the keybind package. It should not be
	// accessed directly. Instead, use keybind.ModMapGet.
	Modmap *ModifierMapping

	// keyRedirect corresponds to a window identifier that, when set,
	// automatically receives *all* keyboard events. This is a sort-of
	// synthetic grab and is helpful in avoiding race conditions.
	keyRedirect xproto.Window

	// Keybinds is the data structure storing all callbacks for key bindings.
	// This is extremely similar to the general notion of event callbacks,
	// but adds extra support to make handling key bindings easier. (Like
	// specifying human readable key sequences to bind to.)
	// KeyBindKey is a struct representing the 4-tuple
	// (event-type, window-id, modifiers, keycode).
	// It is exported for use in the keybind package. Do not access it directly.
	Keybinds    map[KeyBindKey][]KeyBindCallback
	KeybindsLck *sync.RWMutex

	// Keygrabs is a frequency count of the number of callbacks associated
	// with a particular KeyBindKey. This is necessary because we can only
	// grab a particular key *once*, but we may want to attach several callbacks
	// to a single keypress.
	// It is exported for use in the keybind package. Do not access it directly.
	Keygrabs map[KeyBindKey]int

	// mousebinds is the data structure storing all callbacks for mouse
	// bindings.This is extremely similar to the general notion of event
	// callbacks,but adds extra support to make handling mouse bindings easier.
	// (Like specifying human readable mouse sequences to bind to.)
	// MouseBindKey is a struct representing the 4-tuple
	// (event-type, window-id, modifiers, button).
	mousebinds    map[MouseBindKey][]MouseBindCallback
	mousebindsLck *sync.RWMutex

	// mousegrabs is a frequency count of the number of callbacks associated
	// with a particular MouseBindKey. This is necessary because we can only
	// grab a particular mouse button *once*, but we may want to attach
	// several callbacks to a single button press.
	mousegrabs map[MouseBindKey]int

	// mouseDrag is true if a drag is currently in progress.
	mouseDrag bool

	// mouseDragStep is the function executed for each step (i.e., pointer
	// movement) in the current mouse drag. Note that this is nil when a drag
	// is not in progress.
	mouseDragStep MouseDragFun

	// mouseDragEnd is the function executed at the end of the current
	// mouse drag. This is nil when a drag is not in progress.
	mouseDragEnd MouseDragFun

	// gc is a general purpose graphics context; used to paint images.
	// Since we don't do any real X drawing, we don't really care about the
	// particulars of our graphics context.
	gc xproto.Gcontext

	// dummy is a dummy window used for mouse/key GRABs.
	// Basically, whenever a grab is instituted, mouse and key events are
	// redirected to the dummy the window.
	dummy xproto.Window

	// errorHandler is the function that handles errors *in the event loop*.
	// By default, it simply emits them to stderr.
	errorHandler ErrorHandlerFun
}

// NewConn connects to the X server using the DISPLAY environment variable
// and creates a new XUtil.
func NewConn() (*XUtil, error) {
	return NewConnDisplay("")
}

// NewConnDisplay connects to the X server and creates a new XUtil.
// If 'display' is empty, the DISPLAY environment variable is used. Otherwise
// there are several different display formats supported:
//
// NewConn(":1") -> net.Dial("unix", "", "/tmp/.X11-unix/X1")
// NewConn("/tmp/launch-123/:0") -> net.Dial("unix", "", "/tmp/launch-123/:0")
// NewConn("hostname:2.1") -> net.Dial("tcp", "", "hostname:6002")
// NewConn("tcp/hostname:1.0") -> net.Dial("tcp", "", "hostname:6001")
func NewConnDisplay(display string) (*XUtil, error) {
	c, err := xgb.NewConnDisplay(display)

	if err != nil {
		return nil, err
	}

	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)

	// Initialize our central struct that stores everything.
	xu := &XUtil{
		conn:          c,
		quit:          false,
		evqueue:       make([]eventOrError, 0),
		setup:         setup,
		screen:        screen,
		root:          screen.Root,
		eventTime:     xproto.Timestamp(0), // last event time
		atoms:         make(map[string]xproto.Atom, 50),
		atomsLck:      &sync.RWMutex{},
		atomNames:     make(map[xproto.Atom]string, 50),
		atomNamesLck:  &sync.RWMutex{},
		callbacks:     make(map[int]map[xproto.Window][]Callback, 33),
		callbacksLck:  &sync.RWMutex{},
		Keymap:        nil, // we don't have anything yet
		Modmap:        nil,
		keyRedirect:   0,
		Keybinds:      make(map[KeyBindKey][]KeyBindCallback, 10),
		KeybindsLck:   &sync.RWMutex{},
		Keygrabs:      make(map[KeyBindKey]int, 10),
		mousebinds:    make(map[MouseBindKey][]MouseBindCallback, 10),
		mousebindsLck: &sync.RWMutex{},
		mousegrabs:    make(map[MouseBindKey]int, 10),
		mouseDrag:     false,
		mouseDragStep: nil,
		mouseDragEnd:  nil,
		errorHandler:  defaultErrorHandler,
	}

	// Create a general purpose graphics context
	xu.gc, err = xproto.NewGcontextId(xu.conn)
	if err != nil {
		return nil, err
	}
	xproto.CreateGC(xu.conn, xu.gc, xproto.Drawable(xu.root),
		xproto.GcForeground, []uint32{xu.screen.WhitePixel})

	// Create a dummy window
	xu.dummy, err = xproto.NewWindowId(xu.conn)
	if err != nil {
		return nil, err
	}
	xproto.CreateWindow(xu.conn, xu.Screen().RootDepth, xu.dummy, xu.RootWin(),
		-1000, -1000, 1, 1, 0,
		xproto.WindowClassInputOutput, xu.Screen().RootVisual,
		xproto.CwEventMask|xproto.CwOverrideRedirect,
		[]uint32{1, xproto.EventMaskPropertyChange})
	xproto.MapWindow(xu.conn, xu.dummy)

	// Register the Xinerama extension... because it doesn't cost much.
	err = xinerama.Init(xu.conn)

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

// Quit elegantly exits out of the main event loop.
func (xu *XUtil) Quit() {
	xu.quit = true
}

// Quitting returns whether it's time to quit.
// This is only used in the main event loop in xevent.
func (xu *XUtil) Quitting() bool {
	return xu.quit
}

// ExtInitialized returns true if an extension has been initialized.
// This is useful for determining whether an extension is available or not.
func (xu *XUtil) ExtInitialized(extName string) bool {
	_, ok := xu.Conn().Extensions[extName]
	return ok
}

// Sync forces XGB to catch up with all events/requests and synchronize.
// This is done by issuing a benign round trip request to X.
func (xu *XUtil) Sync() {
	xproto.GetInputFocus(xu.Conn()).Reply()
}

// Setup returns the setup information retrieved during connection time.
func (xu *XUtil) Setup() *xproto.SetupInfo {
	return xu.setup
}

// Screen returns the default screen
func (xu *XUtil) Screen() *xproto.ScreenInfo {
	return xu.screen
}

// RootWin returns the current root window.
func (xu *XUtil) RootWin() xproto.Window {
	return xu.root
}

// RootWinSet will change the current root window to the one provided.
// N.B. This probably shouldn't be used unless you're desperately trying
// to support multiple X screens. (This is *not* the same as Xinerama/RandR or
// TwinView. All of those have a single root window.)
func (xu *XUtil) RootWinSet(root xproto.Window) {
	xu.root = root
}

// TimeGet gets the most recent time seen by an event.
func (xu *XUtil) TimeGet() xproto.Timestamp {
	return xu.eventTime
}

// TimeSet sets the most recent time seen by an event.
func (xu *XUtil) TimeSet(t xproto.Timestamp) {
	xu.eventTime = t
}

// GC gets a general purpose graphics context that is typically used to simply
// paint images.
func (xu *XUtil) GC() xproto.Gcontext {
	return xu.gc
}

// Dummy gets the id of the dummy window.
func (xu *XUtil) Dummy() xproto.Window {
	return xu.dummy
}

// AttachCallback associates a (event, window) tuple with an event.
// This function should not be used. It is exported for use in the xevent
// package.
// See the Callback type for an example of attaching event handlers.
func (xu *XUtil) AttachCallback(evtype int, win xproto.Window, fun Callback) {
	xu.callbacksLck.Lock()
	defer xu.callbacksLck.Unlock()

	if _, ok := xu.callbacks[evtype]; !ok {
		xu.callbacks[evtype] = make(map[xproto.Window][]Callback, 20)
	}
	if _, ok := xu.callbacks[evtype][win]; !ok {
		xu.callbacks[evtype][win] = make([]Callback, 0)
	}
	xu.callbacks[evtype][win] = append(xu.callbacks[evtype][win], fun)
}

// RunCallbacks executes every callback corresponding to a
// particular event/window tuple.
// This function should not be used. It is exported for use in the xevent
// package.
func (xu *XUtil) RunCallbacks(event interface{}, evtype int,
	win xproto.Window) {

	xu.callbacksLck.RLock()
	defer xu.callbacksLck.RUnlock()

	for _, cb := range xu.callbacks[evtype][win] {
		cb.Run(xu, event)
	}
}

// DetachWindow removes all callbacks associated with a particular window.
// This function should not be used, since it only cleans up Callback and not
// key and mouse bindings. Instead, use xevent.Detach instead.
func (xu *XUtil) DetachWindow(win xproto.Window) {
	xu.callbacksLck.Lock()
	defer xu.callbacksLck.Unlock()

	for evtype, _ := range xu.callbacks {
		delete(xu.callbacks[evtype], win)
	}
}

// RedirectKeyEvents, when set to a window id (greater than 0), will force
// *all* Key{Press,Release} to callbacks attached to the specified window.
// This is close to emulating a Keyboard grab without the racing.
// To stop redirecting key events, use window identifier '0'.
func (xu *XUtil) RedirectKeyEvents(wid xproto.Window) {
	xu.keyRedirect = wid
}

// RedirectKeyGet gets the window that key events are being redirected to.
// If 0, then no redirection occurs.
func (xu *XUtil) RedirectKeyGet() xproto.Window {
	return xu.keyRedirect
}

// AtomGet retrieves an atom identifier from a cache if it exists.
// This function should not be used. It is exported for use in the xprop
// package. Instead, to intern an atom, use xprop.Atom or xprop.Atm.
func (xu *XUtil) AtomGet(name string) (xproto.Atom, bool) {
	xu.atomsLck.RLock()
	defer xu.atomsLck.RUnlock()

	aid, ok := xu.atoms[name]
	return aid, ok
}

// AtomNameGet retrieves an atom name from a cache if it exists.
// This function should not be used. It is exported for use in the xprop
// package. Instead, to get the name of an atom, use xprop.AtomName.
func (xu *XUtil) AtomNameGet(aid xproto.Atom) (string, bool) {
	name, ok := xu.atomNames[aid]

	xu.atomNamesLck.RLock()
	defer xu.atomNamesLck.RUnlock()

	return name, ok
}

// CacheAtom puts an atom into the cache.
// This function should not be used. It is exported for use in the xprop
// package.
func (xu *XUtil) CacheAtom(name string, aid xproto.Atom) {
	xu.atomsLck.Lock()
	xu.atomNamesLck.Lock()
	defer xu.atomsLck.Unlock()
	defer xu.atomNamesLck.Unlock()

	xu.atoms[name] = aid
	xu.atomNames[aid] = name
}

// Grabs the server. Everything becomes synchronous.
func (xu *XUtil) Grab() {
	xproto.GrabServer(xu.Conn())
}

// Ungrabs the server.
func (xu *XUtil) Ungrab() {
	xproto.UngrabServer(xu.Conn())
}

// Enqueue queues up an event read from X.
// Note that an event read may return an error, in which case, this queue
// entry will be an error and not an event.
func (xu *XUtil) Enqueue(everr eventOrError) {
	xu.evqueue = append(xu.evqueue, everr)
}

// Dequeue pops an event/error from the queue and returns it.
func (xu *XUtil) Dequeue() eventOrError {
	everr := xu.evqueue[0]
	xu.evqueue = xu.evqueue[1:]
	return everr
}

// DequeueAt removes a particular item from the queue
// This is primarily used in the main event loop when compressing events like
// MotionNotify.
func (xu *XUtil) DequeueAt(i int) {
	xu.evqueue = append(xu.evqueue[:i], xu.evqueue[i+1:]...)
}

// QueueEmpty returns whether the event queue is empty or not.
func (xu *XUtil) QueueEmpty() bool {
	return len(xu.evqueue) == 0
}

// QueuePeek returns the current queue so we can examine it.
// This can be useful when trying to determine if a particular kind of
// event will be processed in the future.
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
// The default error handler just emits them to stderr.
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
// This function should not be used. It is exported for use in the xevent
// package.
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
//
// Example to respond to ConfigureNotify events on window 0x1
//
//     xevent.ConfigureNotifyFun(
//		func(X *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
//			fmt.Printf("(%d, %d) %dx%d\n", e.X, e.Y, e.Width, e.Height)
//		}).Connect(X, 0x1)
type Callback interface {
	Connect(xu *XUtil, win xproto.Window)
	Run(xu *XUtil, ev interface{})
}

// Sometimes we need to specify NO WINDOW when a window is typically
// expected. (Like connecting to MappingNotify or KeymapNotify events.)
// Use this value to do that.
var NoWindow xproto.Window = 0

// IgnoreMods is a list of X modifiers that we don't want interfering
// with our mouse or key bindings. In particular, for each mouse or key binding 
// issued, there is a seperate mouse or key binding made for each of the 
// modifiers specified.
var IgnoreMods []uint16 = []uint16{
	0,
	xproto.ModMaskLock,                   // Num lock
	xproto.ModMask2,                      // Caps lock
	xproto.ModMaskLock | xproto.ModMask2, // Caps and Num lock
}
