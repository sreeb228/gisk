package tests

import (
	"encoding/json"
	"errors"
	"gitee.com/sreeb/gisk"
	"os"
	"testing"
)

type fileDslGetter struct {
}

func (getter *fileDslGetter) GetDsl(elementType gisk.ElementType, key string, version string) (string, error) {
	bytes, err := os.ReadFile("./dsl/dsl.json")
	if err != nil {
		return "", err
	}

	var dsl map[string]interface{}
	err = json.Unmarshal(bytes, &dsl)
	if err != nil {
		return "", err
	}

	k := string(elementType) + "_" + key + "_" + version

	if v, ok := dsl[k]; ok {
		//v 转成string
		vv, _ := json.Marshal(v)
		return string(vv), nil
	}
	return "", errors.New("not found dsl")
}

func TestRule_Parse(t *testing.T) {
	type fields struct {
		Left     string
		Operator gisk.Operator
		Right    string
	}
	type args struct {
		gisk *gisk.Gisk
	}

	g := gisk.New()

	g.SetDslGetter(&fileDslGetter{})

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantBool bool
		wantErr  bool
	}{
		{name: "test", fields: fields{Left: "variate_age_1", Operator: gisk.IN, Right: "input_18,20,3_array"}, args: args{gisk: g}, wantBool: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &gisk.Rule{
				Left:     tt.fields.Left,
				Operator: tt.fields.Operator,
				Right:    tt.fields.Right,
			}
			gotBool, err := rule.Parse(tt.args.gisk)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBool != tt.wantBool {
				t.Errorf("Parse() gotBool = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}
