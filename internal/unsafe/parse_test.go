package unsafe_test

import (
	"testing"

	"github.com/hulo-lang/hulo/internal/unsafe"
)

func TestParse(t *testing.T) {
	unsafe.Parse(`{{ loop item in items }}
	echo {{i}}
{{ end }}`)
}
