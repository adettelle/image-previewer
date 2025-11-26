package internalhttp

/*
import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adettelle/image-previewer/internal/mocks"
	"github.com/adettelle/image-previewer/internal/previewservice"
	"github.com/c2fo/testify/require"
	"github.com/golang/mock/gomock"
)

func TestProxingErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockPreviewer := mocks.NewMockPreviewer(ctrl)
	// ih := &ImageHandler{
	// 	PreviewServise: mockPreviewer,
	// }
	// m := ih.PreviewServise.(*mocks.MockPreviewer)

	mockDownloader := mocks.NewMockDownloader(ctrl)

	ds := previewservice.DownloadService{}

	// reqURL := "/fill/{width}/{height}/*"
	reqURL := "/fill/400/300/raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/image-examples/Minla_strigula_4415x2943.jpg"
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	require.NoError(t, err)

	// request.SetPathValue("userid", "1")

	response := httptest.NewRecorder()

	imageAddr := "https://raw.githubusercontent.com/adettelle/image-previewer/refs/heads/create_api/image-examples/Minla_strigula_4415x2943.jpg"
	originalImageName := base64.StdEncoding.EncodeToString([]byte(imageAddr))
	// resizedImageName := originalImageName + "_" + "400" + "_" + "300"
	pathToOriginal := "/tmp/"
	pathToOriginalFile := pathToOriginal + originalImageName

	m.EXPECT().GeneratePreview(400, 300, imageAddr, "scale").Return(previewservice.ResizedImage, nil)
	// m.EXPECT().DownloadFile(pathToOriginalFile, imageAddr).Return(nil)

	ih.preview(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	// require.Contains(t, response.Header().Get("Content-Type"), "application/json")

}
*/
