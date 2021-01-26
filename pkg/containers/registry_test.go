package containers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"tuber/testutilities"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestToken(t *testing.T) {
	testCases := []struct {
		name, path, scope, username, password, domain string
		expectedToken                                 string
		expectedError                                 error
		testStatusCode                                int
		testStatusResponse, testBody                  string
		testError                                     error
	}{
		{
			name:               "token found",
			path:               "success/path",
			domain:             "foo.com",
			expectedToken:      "successToken",
			testStatusCode:     http.StatusOK,
			testStatusResponse: "200 OKAY",
			testBody:           `{"token": "successToken"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			mockClient := testutilities.NewTestClient(func(req *http.Request) (*http.Response, error) {
				reqURL := fmt.Sprintf("https://%s/v2/token?scope=repository:success/path:pull", tc.domain)
				require.Equal(tt, reqURL, req.URL.String())

				fmt.Println()
				return &http.Response{
					StatusCode: tc.testStatusCode,
					Status:     tc.testStatusResponse,
					Header:     make(http.Header),
					Body:       ioutil.NopCloser(strings.NewReader(tc.testBody)),
				}, tc.testError
			})

			r := NewRegistry(tc.domain, tc.password, mockClient)

			actualToken, actualErr := r.token(tc.path)
			testutilities.CheckError(tt, tc.expectedError, actualErr)
			assert.Equal(tt, tc.expectedToken, actualToken)
		})
	}
}
