package path

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	box := NewMockBox(ctrl)
	box.EXPECT().Value().Return(true).AnyTimes()

	scope := NewMockScope(ctrl)
	scope.EXPECT().GetIdentValue("aaa").Return(box, nil).AnyTimes()
	scope.EXPECT().GetIdentValue("aaa.bbb").Return(box, nil).AnyTimes()

	res, err := ioutil.ReadFile("./testfiles/success")
	if err != nil {
		t.Fatal(err)
	}

	buf := bufio.NewReader(bytes.NewBuffer(res))
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}

		t.Run(string(line), func(t *testing.T) {
			query, err := Parse(string(line))
			if err != nil {
				t.Fatal(err)
			}

			done, err := query.Run(scope)
			if err != nil {
				t.Fatal(err)
			}
			if expected, actual := true, done; expected != actual {
				t.Errorf("expected %v, actual: %v", expected, actual)
			}
		})
	}
}
