package httpin

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type UpdateUserProfileInput struct {
	Name   string `in:"form=name"`
	Gender string `in:"form=gender"`
	Avatar File   `in:"form=avatar"`
}

func TestMultipartForm_DecodeFile(t *testing.T) {
	var AvatarBytes = []byte("avatar image content")

	Convey("Upload a file through multipart/form-data requests", t, func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		nameFieldWriter, err := writer.CreateFormField("name")
		So(err, ShouldBeNil)
		nameFieldWriter.Write([]byte("Ggicci T'ang"))

		genderFieldWriter, err := writer.CreateFormField("gender")
		So(err, ShouldBeNil)
		genderFieldWriter.Write([]byte("male"))

		avatarFileWriter, err := writer.CreateFormFile("avatar", "avatar.png")
		So(err, ShouldBeNil)
		_, err = avatarFileWriter.Write(AvatarBytes)
		So(err, ShouldBeNil)

		_ = writer.Close() // error ignored

		r, _ := http.NewRequest("POST", "/", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		// rw := httptest.NewRecorder()
		core, err := New(UpdateUserProfileInput{})
		So(err, ShouldBeNil)
		gotInput, err := core.Decode(r)
		So(err, ShouldBeNil)
		got, ok := gotInput.(*UpdateUserProfileInput)
		So(ok, ShouldBeTrue)
		So(got.Name, ShouldEqual, "Ggicci T'ang")
		So(got.Gender, ShouldEqual, "male")
		So(got.Avatar.Header.Filename, ShouldEqual, "avatar.png")
		So(got.Avatar.Header.Size, ShouldEqual, len(AvatarBytes))
	})

}
