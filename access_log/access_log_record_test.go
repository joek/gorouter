package access_log_test

import (
	. "github.com/cloudfoundry/gorouter/access_log"

	router_http "github.com/cloudfoundry/gorouter/common/http"
	"github.com/cloudfoundry/gorouter/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/url"
	"time"
)

var _ = Describe("AccessLogRecord", func() {
	defaultLogHeaders := make([]string, 0)

	It("Makes a record with all values", func() {
		record := CompleteAccessLogRecord()

		recordString := "FakeRequestHost - " +
			"[01/01/2000:00:00:00 +0000] " +
			"\"FakeRequestMethod http://example.com/request FakeRequestProto\" " +
			"200 " +
			"23 " +
			"\"FakeReferer\" " +
			"\"FakeUserAgent\" " +
			"FakeRemoteAddr " +
			"x_forwarded_for:\"FakeProxy1, FakeProxy2\" " +
			"vcap_request_id:abc-123-xyz-pdq " +
			"response_time:60.000000000 " +
			"app_id:FakeApplicationId\n"
		Expect(record.LogMessage(&defaultLogHeaders)).To(Equal(recordString))
	})

	It("Makes a record with values missing", func() {
		record := AccessLogRecord{
			Request: &http.Request{
				Host:   "FakeRequestHost",
				Method: "FakeRequestMethod",
				Proto:  "FakeRequestProto",
				URL: &url.URL{
					Opaque: "http://example.com/request",
				},
				Header:     http.Header{},
				RemoteAddr: "FakeRemoteAddr",
			},
			RouteEndpoint: &route.Endpoint{
				ApplicationId: "FakeApplicationId",
			},
			StartedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		}

		recordString := "FakeRequestHost - " +
			"[01/01/2000:00:00:00 +0000] " +
			"\"FakeRequestMethod http://example.com/request FakeRequestProto\" " +
			"MissingResponseStatusCode " +
			"0 " +
			"\"-\" " +
			"\"-\" " +
			"FakeRemoteAddr " +
			"x_forwarded_for:\"-\" " +
			"vcap_request_id:- " +
			"response_time:MissingFinishedAt " +
			"app_id:FakeApplicationId\n"

		Expect(record.LogMessage(&defaultLogHeaders)).To(Equal(recordString))
	})

	It("does not create a log message when route endpoint missing", func() {
		record := AccessLogRecord{}
		Expect(record.LogMessage(&defaultLogHeaders)).To(Equal(""))
	})

	It("Makes a record with custom headers", func() {
		header := http.Header{}
		header.Set("Foo", "FOO")
		header.Set("Bar", "BAR")

		record := AccessLogRecord{
			Request: &http.Request{
				Host:   "FakeRequestHost",
				Method: "FakeRequestMethod",
				Proto:  "FakeRequestProto",
				URL: &url.URL{
					Opaque: "http://example.com/request",
				},
				RemoteAddr: "FakeRemoteAddr",
			},
			RouteEndpoint: &route.Endpoint{
				ApplicationId: "FakeApplicationId",
			},
			StartedAt:      time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			ResponseHeader: header,
		}

		recordString := "FakeRequestHost - " +
			"[01/01/2000:00:00:00 +0000] " +
			"\"FakeRequestMethod http://example.com/request FakeRequestProto\" " +
			"MissingResponseStatusCode " +
			"0 " +
			"\"-\" " +
			"\"-\" " +
			"FakeRemoteAddr " +
			"x_forwarded_for:\"-\" " +
			"vcap_request_id:- " +
			"foo:\"FOO\" " +
			"bar:\"BAR\" " +
			"response_time:MissingFinishedAt " +
			"app_id:FakeApplicationId\n"

		logHeaders := []string{"foo", "bar"}
		Expect(record.LogMessage(&logHeaders)).To(Equal(recordString))
	})

	It("creates a log message if headers are missing", func() {
		record := CompleteAccessLogRecord()

		recordString := "FakeRequestHost - " +
			"[01/01/2000:00:00:00 +0000] " +
			"\"FakeRequestMethod http://example.com/request FakeRequestProto\" " +
			"200 " +
			"23 " +
			"\"FakeReferer\" " +
			"\"FakeUserAgent\" " +
			"FakeRemoteAddr " +
			"x_forwarded_for:\"FakeProxy1, FakeProxy2\" " +
			"vcap_request_id:abc-123-xyz-pdq " +
			"foo:\"-\" " +
			"bar:\"-\" " +
			"response_time:60.000000000 " +
			"app_id:FakeApplicationId\n"

		logHeaders := []string{"foo", "bar"}

		Expect(record.LogMessage(&logHeaders)).To(Equal(recordString))
	})

})

func CompleteAccessLogRecord() AccessLogRecord {
	return AccessLogRecord{
		Request: &http.Request{
			Host:   "FakeRequestHost",
			Method: "FakeRequestMethod",
			Proto:  "FakeRequestProto",
			URL: &url.URL{
				Opaque: "http://example.com/request",
			},
			Header: http.Header{
				"Referer":                       []string{"FakeReferer"},
				"User-Agent":                    []string{"FakeUserAgent"},
				"X-Forwarded-For":               []string{"FakeProxy1, FakeProxy2"},
				router_http.VcapRequestIdHeader: []string{"abc-123-xyz-pdq"},
			},
			RemoteAddr: "FakeRemoteAddr",
		},
		BodyBytesSent: 23,
		StatusCode:    200,
		RouteEndpoint: &route.Endpoint{
			ApplicationId: "FakeApplicationId",
		},
		StartedAt:  time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		FinishedAt: time.Date(2000, time.January, 1, 0, 1, 0, 0, time.UTC),
	}
}
