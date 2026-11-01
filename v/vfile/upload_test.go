package vfile

import "testing"

func TestValidateUploadExt(t *testing.T) {
	for _, name := range []string{"a.jpg", "b.PNG", "doc.pdf", "x.zip"} {
		if _, err := ValidateUploadExt(name); err != nil {
			t.Fatalf("allowed file rejected: %s", name)
		}
	}
	for _, name := range []string{"x.html", "x.svg", "x.js", "x.php", "x.exe", "noext", "x.jpg.html"} {
		if _, err := ValidateUploadExt(name); err == nil {
			t.Fatalf("dangerous file accepted: %s", name)
		}
	}
	ext, err := ValidateUploadExt("Photo.JPEG")
	if err != nil || ext != ".jpeg" {
		t.Fatalf("ext normalization failed: %q %v", ext, err)
	}
}
