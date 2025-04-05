package url_download

type DlProgress struct {
	filename  string
	filepath  string
	totalSize int64
	written   uint64
	progPct   uint8
}

type DlResult struct {
	filename string
	filepath string
	filesize int64
}

type DlProgressWriter struct {
	total    uint64
	written  uint64
	lastPct  uint8
	filename string
	filepath string
	progChan chan<- DlProgress
}

func (pw *DlProgressWriter) Write(chunk []byte) (int, error) {
	chunkSize := len(chunk)
	pw.written += uint64(chunkSize)

	if pw.total > 0 {
		pct := uint8(float64(pw.written) / float64(pw.total) * 100)
		if (pct != pw.lastPct) && (pct%5 == 0) {
			pw.progChan <- DlProgress{
				filename:  pw.filename,
				filepath:  pw.filepath,
				totalSize: int64(pw.total),
				written:   pw.written,
				progPct:   pct,
			}
			pw.lastPct = pct
		}
	}

	return chunkSize, nil
}
