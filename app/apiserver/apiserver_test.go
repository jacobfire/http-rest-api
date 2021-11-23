package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/jacobfire/http-rest-api/app/store"
	"github.com/jacobfire/http-rest-api/configs"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAPIServer_configureStore(t *testing.T) {
	type fields struct {
		config *configs.Config
		logger *logrus.Logger
		router *mux.Router
		store  *store.Store
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &APIServer{
				//config: tt.fields.config,
				logger: tt.fields.logger,
				router: tt.fields.router,
				store:  tt.fields.store,
			}
			if err := s.configureStore(); (err != nil) != tt.wantErr {
				t.Errorf("configureStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaginate(t *testing.T) {
	type args struct {
		total int
		page  int
		size  int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Paginate(tt.args.total, tt.args.page, tt.args.size)
			if got != tt.want {
				t.Errorf("Paginate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Paginate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSum(t *testing.T) {
	t.Run("sum 1+2", func(t *testing.T) {
		got := Sum(1, 2)
		if got != 3 {
			t.Errorf("Paginate() got = %v, want %v", got, 4)
		}
	})

	t.Run("sum 1+2+1", func(t *testing.T) {
		got := Sum(1, 2)
		if got != 3 {
			t.Errorf("Paginate() got = %v, want %v", got, 3)
		}
	})
}