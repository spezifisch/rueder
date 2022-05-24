package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cockroachdb/copyist"
	"github.com/gobuffalo/pop/v6"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestAPIUser_First(t *testing.T) {
	defer copyist.Open(t).Close()
	firstRunTest(t)
}

func TestAPIUser_Folders(t *testing.T) {
	defer copyist.Open(t).Close()
	beforeEachTest(t)

	testPop, err := pop.Connect("test")
	assert.NoError(t, err)

	// good claims
	mockClaims := helpers.AuthClaims{
		ID:         testUserID,
		Origin:     "somewhere",
		Name:       "someone",
		OriginName: "somewhere:someone",
	}

	// some bad claims
	invalidClaims1 := mockClaims
	invalidClaims1.Name = ""

	invalidClaims2 := mockClaims
	invalidClaims2.Origin = ""
	invalidClaims2.Name = ""
	invalidClaims2.OriginName = ":"

	badUserClaims1 := mockClaims
	badUserClaims1.ID = nonExistentUserID

	type fields struct {
		pop                  *pop.Connection
		folderCountLimit     int
		folderFeedCountLimit int
	}
	type args struct {
		claims *helpers.AuthClaims
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantRet         []controller.Folder
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "repo nil",
			fields: fields{
				pop: nil,
			},
			args: args{
				claims: &mockClaims,
			},
			wantErr:         true,
			wantErrContains: "invalid repository",
		},
		{
			name: "claims nil",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims: nil,
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "invalid claims 1",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims: &invalidClaims1,
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "invalid claims 2",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims: &invalidClaims2,
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "bad user claims 1",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims: &badUserClaims1,
			},
			wantErr:         true,
			wantErrContains: "user doesn't exist",
		},
		{
			name: "folders empty",
			fields: fields{
				pop:                  testPop,
				folderCountLimit:     5,
				folderFeedCountLimit: 5,
			},
			args: args{
				claims: &mockClaims,
			},
			wantErr: false,
			wantRet: []controller.Folder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &APIPopRepository{
				pop:                  tt.fields.pop,
				folderCountLimit:     tt.fields.folderCountLimit,
				folderFeedCountLimit: tt.fields.folderFeedCountLimit,
			}
			gotRet, err := r.Folders(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("APIPopRepository.Folders() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr && tt.wantErrContains != "" {
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ChangeFolders() error = %v, wantErr %v doesn't contain %s", err, tt.wantErr, tt.wantErrContains)
				}
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("APIPopRepository.Folders() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestAPIUser_ChangeFolders(t *testing.T) {
	defer copyist.Open(t).Close()
	beforeEachTest(t)

	testPop, err := pop.Connect("test")
	assert.NoError(t, err)

	// good claims
	mockClaims := helpers.AuthClaims{
		ID:         testUserID,
		Origin:     "somewhere",
		Name:       "someone",
		OriginName: "somewhere:someone",
	}

	// some bad claims
	invalidClaims1 := mockClaims
	invalidClaims1.Name = ""

	invalidClaims2 := mockClaims
	invalidClaims2.Origin = ""
	invalidClaims2.Name = ""
	invalidClaims2.OriginName = ":"

	badUserClaims1 := mockClaims
	badUserClaims1.ID = nonExistentUserID

	type fields struct {
		pop                  *pop.Connection
		folderCountLimit     int
		folderFeedCountLimit int
	}
	type args struct {
		claims  *helpers.AuthClaims
		folders []controller.Folder
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "repo nil",
			fields: fields{
				pop: nil,
			},
			args: args{
				claims:  &mockClaims,
				folders: []controller.Folder{},
			},
			wantErr:         true,
			wantErrContains: "invalid repository",
		},
		{
			name: "folders nil",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims:  &mockClaims,
				folders: nil,
			},
			wantErr:         true,
			wantErrContains: "invalid folders",
		},
		{
			name: "claims nil",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims:  nil,
				folders: []controller.Folder{},
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "invalid claims 1",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims:  &invalidClaims1,
				folders: []controller.Folder{},
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "invalid claims 2",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims:  &invalidClaims2,
				folders: []controller.Folder{},
			},
			wantErr:         true,
			wantErrContains: "invalid claims",
		},
		{
			name: "bad user claims 1",
			fields: fields{
				pop: testPop,
			},
			args: args{
				claims:  &badUserClaims1,
				folders: []controller.Folder{},
			},
			wantErr:         true,
			wantErrContains: "user doesn't exist",
		},
		{
			name: "folders empty",
			fields: fields{
				pop:                  testPop,
				folderCountLimit:     5,
				folderFeedCountLimit: 5,
			},
			args: args{
				claims:  &mockClaims,
				folders: []controller.Folder{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &APIPopRepository{
				pop: tt.fields.pop,
			}
			if err := a.ChangeFolders(tt.args.claims, tt.args.folders); (err != nil) != tt.wantErr {
				t.Errorf("ChangeFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr && tt.wantErrContains != "" {
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Errorf("ChangeFolders() error = %v, wantErr %v doesn't contain %s", err, tt.wantErr, tt.wantErrContains)
				}
			}
		})
	}
}
