package httpin

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
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

type UpdateGitHubIssueInput struct {
	Title       string `in:"form=title"`
	Attachments []File `in:"form=attachment"`
}

func TestMultipartForm_DecodeFile_WithInvalidFileHeaders(t *testing.T) {
	Convey("Decode file with nil header", t, func() {
		gotInput, err := decodeFile(nil)
		So(errors.Is(err, ErrNilFile), ShouldBeTrue)
		got, ok := gotInput.(File)
		So(ok, ShouldBeTrue)
		So(got.Valid, ShouldBeFalse)
	})

	Convey("Decode file with invalid (broken) header", t, func() {
		fileHeader := &multipart.FileHeader{
			Filename: "avatar.png",
			Size:     10,
		}
		gotInput, err := decodeFile(fileHeader)
		So(err, ShouldBeError)
		got, ok := gotInput.(File)
		So(ok, ShouldBeTrue)
		So(got.Valid, ShouldBeFalse)
	})
}

func TestMultipartForm_UploadSingleFile(t *testing.T) {
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

		core, err := New(UpdateUserProfileInput{})
		So(err, ShouldBeNil)
		gotInput, err := core.Decode(r)
		So(err, ShouldBeNil)
		got, ok := gotInput.(*UpdateUserProfileInput)
		So(ok, ShouldBeTrue)
		So(got.Name, ShouldEqual, "Ggicci T'ang")
		So(got.Gender, ShouldEqual, "male")
		So(got.Avatar.Valid, ShouldBeTrue)
		So(got.Avatar.Header.Filename, ShouldEqual, "avatar.png")
		So(got.Avatar.Header.Size, ShouldEqual, len(AvatarBytes))
		uploadedContent, err := ioutil.ReadAll(got.Avatar.File)
		So(err, ShouldBeNil)
		So(uploadedContent, ShouldResemble, AvatarBytes)
	})

	Convey("No files uploaded", t, func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		nameFieldWriter, err := writer.CreateFormField("name")
		So(err, ShouldBeNil)
		nameFieldWriter.Write([]byte("Ggicci T'ang"))

		_ = writer.Close() // error ignored

		r, _ := http.NewRequest("POST", "/", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())
		core, err := New(UpdateUserProfileInput{})
		So(err, ShouldBeNil)
		gotInput, err := core.Decode(r)
		So(err, ShouldBeNil)
		got, ok := gotInput.(*UpdateUserProfileInput)
		So(ok, ShouldBeTrue)
		So(got.Name, ShouldEqual, "Ggicci T'ang")
		So(got.Avatar.Valid, ShouldBeFalse)
		So(got.Avatar.File, ShouldBeNil)
		So(got.Avatar.Header, ShouldBeNil)
	})

	Convey("Broken boundaries should cause server to fail", t, func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		avatarFileWriter, err := writer.CreateFormFile("avatar", "avatar.png")
		So(err, ShouldBeNil)
		_, err = avatarFileWriter.Write(AvatarBytes)
		So(err, ShouldBeNil)
		writer.Close() // error ignored

		raw := body.Bytes()
		var brokenBody = bytes.NewBuffer(raw[:len(raw)-10])
		brokenBody.Write([]byte("xxx")) // break the boundary

		r, _ := http.NewRequest("POST", "/", brokenBody)
		r.Header.Set("Content-Type", writer.FormDataContentType())
		core, err := New(UpdateUserProfileInput{})
		So(err, ShouldBeNil)

		gotInput, err := core.Decode(r)
		So(gotInput, ShouldBeNil)
		So(err, ShouldBeError)
	})
}

func TestMultipartForm_UploadMultiFiles(t *testing.T) {
	var Attachments = [][]byte{
		[]byte("attachment #1"),
		[]byte("attachment #2"),
		[]byte("attachment #3"),
	}

	Convey("Upload multiple files through multipart/form-data requests", t, func() {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		title := "feature-request: integrate with open-telemetry"
		titleFieldWriter, err := writer.CreateFormField("title")
		So(err, ShouldBeNil)
		titleFieldWriter.Write([]byte(title))

		for i, attContent := range Attachments {
			filename := fmt.Sprintf("attachment-%d.txt", i+1)
			attachmentFileWriter, err := writer.CreateFormFile("attachment", filename)
			So(err, ShouldBeNil)
			_, err = attachmentFileWriter.Write(attContent)
			So(err, ShouldBeNil)
		}
		_ = writer.Close() // error ignored

		r, _ := http.NewRequest("POST", "/", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		core, err := New(UpdateGitHubIssueInput{})
		So(err, ShouldBeNil)
		gotInput, err := core.Decode(r)
		So(err, ShouldBeNil)
		got, ok := gotInput.(*UpdateGitHubIssueInput)
		So(ok, ShouldBeTrue)

		So(got.Title, ShouldEqual, title)
		So(got.Attachments, ShouldHaveLength, len(Attachments))
		for i, att := range got.Attachments {
			So(att.Valid, ShouldBeTrue)
			So(att.Header.Filename, ShouldEqual, fmt.Sprintf("attachment-%d.txt", i+1))
			So(att.Header.Size, ShouldEqual, len(Attachments[i]))
			uploadedContent, err := ioutil.ReadAll(att.File)
			So(err, ShouldBeNil)
			So(uploadedContent, ShouldResemble, Attachments[i])
		}
	})
}
