// Code generated by "stringer -type=TokenTyp"; DO NOT EDIT.

package lisp

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tokenError-0]
	_ = x[tokenEOF-1]
	_ = x[tokenComment-2]
	_ = x[tokenLParen-3]
	_ = x[tokenRParen-4]
	_ = x[tokenQuote-5]
	_ = x[tokenAtom-6]
	_ = x[tokenNumber-7]
}

const _TokenTyp_name = "tokenErrortokenEOFtokenCommenttokenLParentokenRParentokenQuotetokenAtomtokenNumber"

var _TokenTyp_index = [...]uint8{0, 10, 18, 30, 41, 52, 62, 71, 82}

func (i TokenTyp) String() string {
	if i < 0 || i >= TokenTyp(len(_TokenTyp_index)-1) {
		return "TokenTyp(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenTyp_name[_TokenTyp_index[i]:_TokenTyp_index[i+1]]
}
