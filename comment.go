package golexer2

type cLineCommentMatcher struct {
}

func (cLineCommentMatcher) MatchRune(index int, r rune) bool {
	switch index {
	case 0:
		return r == '/'
	case 1:
		return r == '/'
	default:
		return r != '\r' && r != '\n'
	}
}

func (cLineCommentMatcher) TokenType() string {
	return "CLineComment"
}

var _cLineCommentMatcher = new(cLineCommentMatcher)

// C的行注释
func CLineComment() Matcher {
	return _cLineCommentMatcher
}

type unixLineCommentMatcher struct {
}

func (unixLineCommentMatcher) MatchRune(index int, r rune) bool {
	switch index {
	case 0:
		return r == '#'
	default:
		return r != '\r' && r != '\n'
	}
}

func (unixLineCommentMatcher) TokenType() string {
	return "UnixLineComment"
}

var _unixLineCommentMatcher = new(unixLineCommentMatcher)

// Unix的行注释
func UnixLineComment() Matcher {
	return _unixLineCommentMatcher
}

type cBlockCommentMatcher struct {
	endStarIndex int
	blockEnd     bool
}

func (self *cBlockCommentMatcher) MatchRune(index int, r rune) bool {

	if self.blockEnd {
		return false
	}

	switch index {
	case 0:
		return r == '/'
	case 1:
		return r == '*'
	default:

		switch r {
		case '*':
			self.endStarIndex = index
		case '/':
			if index == self.endStarIndex+1 {
				self.blockEnd = true // 需要等下一次再结束，包含/
			}
		default:
			// *只是纯的*
			if self.endStarIndex != -1 {
				self.endStarIndex = -1
			}
		}

		return true
	}
}

func (cBlockCommentMatcher) TokenType() string {
	return "CBlockComment"
}

// C的块注释
func CBlockComment() Matcher {
	return &cBlockCommentMatcher{
		endStarIndex: -1,
	}
}
