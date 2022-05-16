package pop

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

func TestAPIPopRepository_ChangeFolders(t *testing.T) {
	mockClaims := helpers.AuthClaims{
		Origin:     "somewhere",
		Name:       "someone",
		OriginName: "somewhere:someone",
	}
	mockClaims.ID, _ = uuid.NewV4()

	badClaims1 := mockClaims
	badClaims1.Name = ""

	badClaims2 := mockClaims
	badClaims2.Origin = ""
	badClaims2.Name = ""
	badClaims2.OriginName = ":"

	type fields struct {
		pop *pop.Connection
	}
	type args struct {
		claims  *helpers.AuthClaims
		folders []controller.Folder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "folders nil",
			args: args{
				claims:  &mockClaims,
				folders: nil,
			},
			wantErr: true,
		},
		{
			name: "claims nil",
			args: args{
				claims:  nil,
				folders: []controller.Folder{},
			},
			wantErr: true,
		},
		{
			name: "bad claims 1",
			args: args{
				claims:  &badClaims1,
				folders: []controller.Folder{},
			},
			wantErr: true,
		},
		{
			name: "bad claims 2",
			args: args{
				claims:  &badClaims2,
				folders: []controller.Folder{},
			},
			wantErr: true,
		},
		{
			name: "folders empty",
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
				t.Errorf("APIPopRepository.ChangeFolders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
