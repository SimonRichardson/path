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

	bad := NewMockScope(ctrl)

	subsubchild := NewMockScope(ctrl)
	subsubchild.EXPECT().GetIdentValue("eee").Return(bad, nil).AnyTimes()

	subchild := NewMockScope(ctrl)
	subchild.EXPECT().GetIdentValue("ccc").Return(bad, nil).AnyTimes()
	subchild.EXPECT().GetIdentValue("ddd").Return(subsubchild, nil).AnyTimes()

	child := NewMockScope(ctrl)
	child.EXPECT().GetIdentValue("bbb").Return(subchild, nil).AnyTimes()
	child.EXPECT().GetIdentValue("xxx").Return(nil, ErrNotFound).AnyTimes()
	child.EXPECT().GetIdentValue("yyy").Return(bad, nil).AnyTimes()
	child.EXPECT().GetAllIdents().Return([]string{"bbb"}).AnyTimes()

	root := NewMockScope(ctrl)
	root.EXPECT().GetIdentValue("aaa").Return(child, nil).AnyTimes()

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

			done, err := query.Run(root)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(done)
		})
	}
}
