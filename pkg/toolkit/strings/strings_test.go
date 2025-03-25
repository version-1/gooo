package strings

import "testing"

func TestSnakeCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "snake case from snake case",
			args: args{s: "test_snake_case"},
			want: "test_snake_case",
		},
		{
			name: "snake case from camel case",
			args: args{s: "TestSnakeCase"},
			want: "test_snake_case",
		},
		{
			name: "snake case from upper case",
			args: args{s: "UPPERCASE"},
			want: "uppercase",
		},
		{
			name: "snake case from kebab case",
			args: args{s: "kebab-case-example"},
			want: "kebab_case_example",
		},
		{
			name: "snake case from upper snake case",
			args: args{s: "UPPER_SNAKE_CASE"},
			want: "upper_snake_case",
		},
		{
			name: "snake case from train case",
			args: args{s: "Train-Case-Example"},
			want: "train_case_example",
		},
		{
			name: "snake case from dot notation",
			args: args{s: "dot.notation.example"},
			want: "dot_notation_example",
		},
		{
			name: "snake case from space delimited",
			args: args{s: "Space Delimited Example"},
			want: "space_delimited_example",
		},
		{
			name: "snake case from symbols",
			args: args{s: " !@#^&*()-_=+[]{} |;:',<.>/?`"},
			want: "_!@#^&*()__=+[]{}_|;:',<_>/?`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.args.s); got != tt.want {
				t.Errorf("snakeCase() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "camel case from camel case",
			args: args{s: "TestCamelCase"},
			want: "TestCamelCase",
		},
		{
			name: "snake case from upper case",
			args: args{s: "UPPERCASE"},
			want: "Uppercase",
		},
		{
			name: "camel case from kebab case",
			args: args{s: "kebab-case-example"},
			want: "KebabCaseExample",
		},
		{
			name: "camel case from upper snake case",
			args: args{s: "UPPER_SNAKE_CASE"},
			want: "UpperSnakeCase",
		},
		{
			name: "camel case from train case",
			args: args{s: "Train-Case-Example"},
			want: "TrainCaseExample",
		},
		{
			name: "camel case from dot notation",
			args: args{s: "dot.notation.example"},
			want: "DotNotationExample",
		},
		{
			name: "camel case from space delimited",
			args: args{s: "Space Delimited Example"},
			want: "SpaceDelimitedExample",
		},
		{
			name: "camel case from symbols",
			args: args{s: " !@#^&*()-_=+[]{} |;:',<.>/?`"},
			want: "!@#^&*()=+[]{}|;:',<>/?`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.args.s); got != tt.want {
				t.Errorf("camelCase() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestKebabCase(t *testing.T) {

}

func TestPascalCase(t *testing.T) {

}
