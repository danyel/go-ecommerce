package handler

import (
	Http "net/http"
	HttpTest "net/http/httptest"
	Testing "testing"
)

func TestParsePaginationRejectsOversizedPage(t *Testing.T) {
	request := HttpTest.NewRequest(Http.MethodGet, "/?page=2&page_size=51", nil)
	if _, _, err := ParsePagination(request, 50); err == nil {
		t.Fatal("expected oversized page_size to be rejected")
	}
}

func TestParsePaginationDefaultsToBoundedPage(t *Testing.T) {
	request := HttpTest.NewRequest(Http.MethodGet, "/", nil)
	page, pageSize, err := ParsePagination(request, 50)
	if err != nil || page != 1 || pageSize != 50 {
		t.Fatalf("unexpected pagination defaults: page=%d page_size=%d err=%v", page, pageSize, err)
	}
}
