package file

type PreviewService struct {
	//cache
	//pathToSaveIncommingImages
}

func New() *PreviewService {
	return &PreviewService{}
}

// returns pathToSave
func (ps *PreviewService) GeneratePreview(outWidth int, outHeight int, imageAddr string) (string, error) {
	// fmt.Println("Got from cache: " + imageAddr)
	return "", nil
}

func countActualSize() {

}
