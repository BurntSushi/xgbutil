/*
   Package xcursor.go facilitates the use of different cursors with X.
   Please see 'cursors.go' for a list of all available cursors.
*/
package xcursor

import "github.com/BurntSushi/xgb"
import "github.com/BurntSushi/xgbutil"

// CreateCursor sets some default colors for nice and easy cursor creation.
// Just supply a cursor constant from 'xcursor/cursors.go'.
func CreateCursor(xu *xgbutil.XUtil, cursor uint16) xgb.Id {
	return CreateCursorExtra(xu, cursor, 0, 0, 0, 0xffff, 0xffff, 0xffff)
}

// CreateCursorExtra features all available parameters to creating a cursor.
func CreateCursorExtra(xu *xgbutil.XUtil, cursor, foreRed, foreGreen,
	foreBlue, backRed, backGreen, backBlue uint16) xgb.Id {

	fontId := xu.Conn().NewId()
	cursorId := xu.Conn().NewId()

	xu.Conn().OpenFont(fontId, "cursor")
	xu.Conn().CreateGlyphCursor(cursorId, fontId, fontId,
		cursor, cursor+1,
		foreRed, foreGreen, foreBlue,
		backRed, backGreen, backBlue)
	xu.Conn().CloseFont(fontId)

	return cursorId
}
