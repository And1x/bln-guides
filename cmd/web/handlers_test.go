package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	tests := []struct {
		name string
		tt   time.Time
		want string
	}{
		{
			name: "Local",
			tt:   time.Date(2021, 11, 21, 10, 0, 0, 0, time.Local),
			want: "21 Nov 2021 at 10:00",
		},
		{
			name: "Empty",
			tt:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tt:   time.Date(2021, 11, 21, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "21 Nov 2021 at 10:00",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			hd := humandate(tst.tt)

			if hd != tst.want {
				t.Errorf("want %q; got %q", tst.want, hd)
			}
		})
	}

}
