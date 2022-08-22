// Package searcher parse conditions and pagination
// nolint
package searcher

import (
	"encoding/json"
	"testing"

	"gitee.com/szxjyt/filbox-backend/db/mysql"
)

func Test_parseOrder(t *testing.T) {
	type args struct {
		order string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				order: "+created_at",
			},
			want: "created_at",
		},
		{
			name: "t1",
			args: args{
				order: "-created_at",
			},
			want: "created_at desc",
		}, {
			name: "t3",
			args: args{
				order: "-created_at, -id, +name",
			},
			want: "created_at desc, id desc, name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOrder(tt.args.order); got != tt.want {
				t.Errorf("parseOrder() = [%v], want %v", got, tt.want)
			}
		})
	}
}

func Test_condition(t *testing.T) {
	c := mysql.QueryCondition{
		[]*mysql.Condition{
			&mysql.Condition{
				"name",
				"like",
				"a",
			}, &mysql.Condition{
				"sex",
				"=",
				"man",
			},
		},
	}
	d, _ := json.Marshal(c)
	t.Logf("%+v", string(d))
}
