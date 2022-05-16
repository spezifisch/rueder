package pop

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/spezifisch/rueder3/backend/pkg/api/controller"
	"github.com/spezifisch/rueder3/backend/pkg/helpers"
)

func TestAPIPopRepository_ChangeFolders(t *testing.T) {
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
				folders: nil,
			},
			wantErr: true,
		},
		{
			name: "folders empty",
			args: args{
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
