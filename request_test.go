package f

import (
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestRequest(t *testing.T) {

	var req *Request

	BeforeEach(func() {
		req = CreateRequestMock(nil)
	})

	Report(t)
}
