package special

import (
	"1bwiki/tmpl/layout"
	"bytes"
)

func Register() string {
	var _buffer bytes.Buffer
	_buffer.WriteString("\n\n\n\n<form action=\"/special/register\" method=\"post\">\n\n\t<label>Username</label><br>\n\n\t<input type=\"text\" name=\"username\"><br>\n\n\t<label>Password</label><br>\n\n\t<input type=\"text\" name=\"password\"><br>\n\n\t<label>Confirm Password</label><br>\n\n\t<input type=\"text\" name=\"passwordConfirm\"><br><br>\n\n\t<!-- include stuff here to show password when hidden -->\n\n\t<button type=\"submit\">Submit</button>\n\n</form>")

	return layout.Base(_buffer.String())
}
