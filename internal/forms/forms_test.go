package forms

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid form when should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("test", "not_present")
	if form.Valid() {
		t.Error("got valid form when required fields missing")
	}
}

func TestForm_Has(t *testing.T) {
	r, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.Has("notPresent", r) {
		t.Error("form shows field present that is missing")
	}

	postedData := url.Values{}
	postedData.Add("test", "present")
	form = New(postedData)

	if !form.Has("test", r) {
		t.Error("form shows does not have field that is present")
	}
}

func TestForm_MinLength(t *testing.T) {
	r, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("test", 3, r)
	if form.Valid() {
		t.Error("got valid form when field does not satisfy minlength")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid form when invalid email")
	}
}
