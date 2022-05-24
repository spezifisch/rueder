package api

import (
	"reflect"
	"testing"

	"github.com/cockroachdb/copyist"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/stretchr/testify/assert"
)

func TestAPI_First(t *testing.T) {
	defer copyist.Open(t).Close()
	firstRunTest(t)
}

func TestAPI_Feeds(t *testing.T) {
	defer copyist.Open(t).Close()
	beforeEachTest(t)

	// this connection is reused between tests in this test table
	testPop, err := pop.Connect("test")
	assert.NoError(t, err)

	type fields struct {
		pop                  *pop.Connection
		folderCountLimit     int
		folderFeedCountLimit int
	}
	tests := []struct {
		name      string
		fields    fields
		wantFeeds []controller.Feed
		wantErr   bool
	}{
		{
			name: "no feeds",
			fields: fields{
				pop: testPop,
			},
			wantFeeds: []controller.Feed{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &APIPopRepository{
				pop:                  tt.fields.pop,
				folderCountLimit:     tt.fields.folderCountLimit,
				folderFeedCountLimit: tt.fields.folderFeedCountLimit,
			}
			gotFeeds, err := r.Feeds()
			if (err != nil) != tt.wantErr {
				t.Errorf("APIPopRepository.Feeds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFeeds, tt.wantFeeds) {
				t.Errorf("APIPopRepository.Feeds() = %v, want %v", gotFeeds, tt.wantFeeds)
			}
		})
	}
}

func TestAPI_GetArticle(t *testing.T) {
	defer copyist.Open(t).Close()
	beforeEachTest(t)

	// this connection is reused between tests in this test table
	testPop, err := pop.Connect("test")
	assert.NoError(t, err)

	type fields struct {
		pop                  *pop.Connection
		folderCountLimit     int
		folderFeedCountLimit int
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRet controller.Article
		wantErr bool
	}{
		{
			name: "non-existent article",
			fields: fields{
				pop: testPop,
			},
			args: args{
				id: nonExistentArticleID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &APIPopRepository{
				pop:                  tt.fields.pop,
				folderCountLimit:     tt.fields.folderCountLimit,
				folderFeedCountLimit: tt.fields.folderFeedCountLimit,
			}
			gotRet, err := r.GetArticle(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIPopRepository.GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("APIPopRepository.GetArticle() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
