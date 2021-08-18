package parser

import (
	"reflect"
	"testing"

	"github.com/thelitdeveloper/imageresizer/pkg/server/entities"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name      string
		qryString string
		want      *entities.ImageProperties
		wantURI   string
	}{
		{
			name:      "Usual request",
			qryString: "crazy/images/blue_marble-500x500.jpg",
			want: &entities.ImageProperties{
				Prefix:   "crazy/images",
				Filename: "blue_marble-500x500.jpg",
				Width:    500,
				Height:   500,
				Convert:  false,
			},
			wantURI: "http://simplytest.s3-website.eu-central-1.amazonaws.com/crazy/images/blue_marble.jpg",
		},
		{
			name:      "jpeg file ext",
			qryString: "crazy/images/blue_marble-640x640.jpeg",
			want: &entities.ImageProperties{
				Prefix:   "crazy/images",
				Filename: "blue_marble-640x640.jpeg",
				Width:    640,
				Height:   640,
				Convert:  false,
			},
			wantURI: "http://simplytest.s3-website.eu-central-1.amazonaws.com/crazy/images/blue_marble.jpeg",
		},
		{
			name:      "Resize with aspect ratio",
			qryString: "crazy/images/blue_marble-500x0.JPG",
			want: &entities.ImageProperties{
				Prefix:   "crazy/images",
				Filename: "blue_marble-500x0.JPG",
				Width:    500,
				Height:   0,
				Convert:  false,
			},
			wantURI: "http://simplytest.s3-website.eu-central-1.amazonaws.com/crazy/images/blue_marble.JPG",
		},
		{
			name:      "PNG",
			qryString: "gopher-800x0.png",
			want: &entities.ImageProperties{
				Prefix:   "",
				Filename: "gopher-800x0.png",
				Width:    800,
				Height:   0,
				Convert:  false,
			},
			wantURI: "http://simplytest.s3-website.eu-central-1.amazonaws.com/gopher.png",
		},
		{
			name:      "Convert png to jpg",
			qryString: "images/gopher-800x0-jpg.png",
			want: &entities.ImageProperties{
				Prefix:   "images",
				Filename: "gopher-800x0-jpg.png",
				Width:    800,
				Height:   0,
				Convert:  true,
			},
			wantURI: "http://simplytest.s3-website.eu-central-1.amazonaws.com/images/gopher.png",
		},
		// Add more ...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, uri := Parse(tt.qryString)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
			if uri != tt.wantURI {
				t.Errorf("Parse() got = %v, want %v", uri, tt.wantURI)
			}
		})
	}
}

func Test_createFallbackKey(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		fname   string
		want    string
		wantErr bool
	}{
		{
			name:    "Without prefix",
			prefix:  "",
			fname:   "gopher-720x0-jpg.png",
			want:    "gopher.png",
			wantErr: false,
		},
		{
			name:    "With prefix",
			prefix:  "images",
			fname:   "theshot-800x0.jpg",
			want:    "images/theshot.jpg",
			wantErr: false,
		},
		{
			name:    "With prefix and jpeg ext",
			prefix:  "dully/crazy/images",
			fname:   "theimage-920x0.jpeg",
			want:    "dully/crazy/images/theimage.jpeg",
			wantErr: false,
		},
		{
			name:    "",
			prefix:  "",
			fname:   "AnotherImage-1024x0-jpg.png",
			want:    "AnotherImage.png",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFallbackKey(tt.prefix, tt.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFallbackKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createFallbackKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeConvertRequestParam(t *testing.T) {
	tests := []struct {
		name  string
		fname string
		want  string
	}{
		{
			name:  "",
			fname: "test-500x500-jpg.png",
			want:  "test-500x500.png",
		},
		{
			name:  "",
			fname: "test-500x500-jpeg.png",
			want:  "test-500x500.png",
		},
		{
			name:  "",
			fname: "test-500x500-JPEG.png",
			want:  "test-500x500.png",
		},
		{
			name:  "",
			fname: "test-500x500-JPG.png",
			want:  "test-500x500.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeConvertRequestParam(tt.fname); got != tt.want {
				t.Errorf("removeConvertRequestParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractConvertRequest(t *testing.T) {
	tests := []struct {
		name  string
		fname string
		want  bool
	}{
		{
			name:  "",
			fname: "test1-jpg.png",
			want:  true,
		},
		{
			name:  "",
			fname: "test1-jpeg.png",
			want:  true,
		},
		{
			name:  "",
			fname: "test1-JPG.png",
			want:  true,
		},
		{
			name:  "",
			fname: "test1-JPEG.png",
			want:  true,
		},
		{
			name:  "",
			fname: "test1-jpg.jpg",
			want:  false,
		},
		{
			name:  "",
			fname: "test1.png",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractConvertRequest(tt.fname); got != tt.want {
				t.Errorf("extractConvertRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	tests := []struct {
		name      string
		qryString string
		want      bool
	}{
		{
			name:      "Empty query",
			qryString: "",
			want:      false,
		},
		{
			name:      "Params missing",
			qryString: "images/theshot.jpg",
			want:      false,
		},
		{
			name:      "One param is missing (height)",
			qryString: "images/theshot-800x.jpeg",
			want:      false,
		},
		{
			name:      "Wrong file extension",
			qryString: "dully/crazy/images/filename-450x450.gif",
			want:      false,
		},
		{
			name:      "Missing file extension",
			qryString: "dully/crazy/images/filename-450x450",
			want:      false,
		},
		{
			name:      "",
			qryString: "images/theshot-800x0.JPEG",
			want:      true,
		},
		{
			name:      "",
			qryString: "gopher-600x0-jpg.png",
			want:      true,
		},
		{
			name:      "",
			qryString: "eat/my/shorts/bart-1024x1024.jpg",
			want:      true,
		},
		{
			name:      "",
			qryString: "moin/moin/Hamburg-0x320.PNG",
			want:      true,
		},
		// Add more ...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validate(tt.qryString); got != tt.want {
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractDimensions(t *testing.T) {
	tests := []struct {
		name         string
		paramSegment string
		want         []int
	}{
		{
			name:         "",
			paramSegment: "filename-500x500-jpg.png",
			want:         []int{500, 500},
		},
		{
			name:         "",
			paramSegment: "filename-480x0.jpg",
			want:         []int{480, 0},
		},
		{
			name:         "",
			paramSegment: "filename-0x824.jpeg",
			want:         []int{0, 824},
		},
		{
			name:         "",
			paramSegment: "filename-600x.png",
			want:         []int{600, 0},
		},
		{
			name:         "",
			paramSegment: "filename-600.png",
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := extractDimensions(tt.paramSegment)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractDimensions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPrefix(t *testing.T) {
	tests := []struct {
		name         string
		pathSegments []string
		want         string
	}{
		{
			name:         "",
			pathSegments: []string{"dully", "crazy", "images"},
			want:         "dully/crazy/images",
		},
		{
			name:         "",
			pathSegments: []string{},
			want:         "",
		},
		{
			name:         "",
			pathSegments: []string{"images"},
			want:         "images",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPrefix(tt.pathSegments); got != tt.want {
				t.Errorf("extractPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
