package archive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeDuplicate(t *testing.T) {
	type args struct {
		archiveRoot    string
		duplicateFiles []string
	}
	tests := []struct {
		name      string
		args      args
		want      DeDupTask
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "one file not in archiveRoot",
			args: args{
				archiveRoot:    "foo",
				duplicateFiles: []string{"bar/baz.jgp", "foo"},
			},
			want:      DeDupTask{},
			errAssert: assert.Error,
		},
		{
			name: "no file in calendar directory",
			args: args{
				archiveRoot:    "Archive",
				duplicateFiles: []string{"Archive/all/20190417_133044_537842c8.jpg", "Archive/origin/04/20190417_151708_537842c8.jpg"},
			},
			want:      DeDupTask{},
			errAssert: assert.Error,
		},
		{
			name: "deduplicate files in calendar directory",
			args: args{
				archiveRoot:    "Archive",
				duplicateFiles: []string{"Archive/2019/04/20190417_133044_537842c8.jpg", "Archive/2019/04/20190417_151708_537842c8.jpg"},
			},
			want: DeDupTask{
				ToKeep:      "Archive/2019/04/20190417_133044_537842c8.jpg",
				DeleteFiles: []string{"Archive/2019/04/20190417_151708_537842c8.jpg"},
			},
			errAssert: assert.NoError,
		},
		{
			name: "keep only one file in all",
			args: args{
				archiveRoot:    "Archive",
				duplicateFiles: []string{"Archive/2019/04/20190417_133044_537842c8.jpg", "Archive/all/20190417_133044_537842c8.jpg", "Archive/all/20190417_151708_537842c8.jpg"},
			},
			want: DeDupTask{
				ToKeep:        "Archive/2019/04/20190417_133044_537842c8.jpg",
				ReCreateLinks: []string{"Archive/all/20190417_133044_537842c8.jpg"},
				DeleteFiles:   []string{"Archive/all/20190417_151708_537842c8.jpg"},
			},
			errAssert: assert.NoError,
		},
		{
			name: "keep only one file in origin",
			args: args{
				archiveRoot:    "Archive",
				duplicateFiles: []string{"Archive/2019/04/20190417_133044_537842c8.jpg", "Archive/origin/foo/20190417_133044_537842c8.jpg", "Archive/origin/foo/20190417_151708_537842c8.jpg", "Archive/origin/bar/20190417_151708_537842c8.jpg"},
			},
			want: DeDupTask{
				ToKeep:        "Archive/2019/04/20190417_133044_537842c8.jpg",
				ReCreateLinks: []string{"Archive/origin/bar/20190417_151708_537842c8.jpg", "Archive/origin/foo/20190417_133044_537842c8.jpg"},
				DeleteFiles:   []string{"Archive/origin/foo/20190417_151708_537842c8.jpg"},
			},
			errAssert: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeDuplicate(tt.args.archiveRoot, tt.args.duplicateFiles)
			tt.errAssert(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
