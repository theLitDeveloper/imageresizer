package aws3

import (
	"testing"

	"github.com/thelitdeveloper/imageresizer/pkg/server/entities"
)

func Test_detectContentType(t *testing.T) {
	tests := []struct {
		name  string
		props *entities.ImageProperties
		want  string
	}{
		{
			name: "",
			props: &entities.ImageProperties{
				Filename: "blue_marble-500x500.jpg",
				Convert:  false,
			},
			want: "image/jpeg",
		},
		{
			name: "",
			props: &entities.ImageProperties{
				Filename: "gopher-720x0.png",
				Convert:  false,
			},
			want: "image/png",
		},
		{
			name: "",
			props: &entities.ImageProperties{
				Filename: "gopher-480x0-jpg.jpg",
				Convert:  true,
			},
			want: "image/jpeg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectContentType(tt.props); got != tt.want {
				t.Errorf("detectContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
