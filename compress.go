
func unzipforbuffer(zipfile []byte, savefolder string){
	zipbuf := bytes.NewReader(zipfile) // make io.reader
	archive, err := zip.NewReader(zipbuf,zipbuf.Size()) // io.reader to zip.reader
	if err != nil {
		log.Println("newreader", err)
	} else {
		archive.RegisterDecompressor(zip.Deflate, func(r io.Reader) io.ReadCloser { // register zip decompress (if you used compressor for make "zip")
			return flate.NewReader(r)
		})
		for _, ifile := range archive.File { // zip inside file loop
			pathz := ifile.Name
			modifytime := ifile.Modified
			checkbackslash := strings.HasSuffix(pathz,"\\") // If there's a "\" in the path,
			pathz = strings.Replace(pathz, "\\","/", -1) // backslash to slash
			info := ifile.FileInfo()
			location := path.Join(savefolder,pathz)
			switch {
			case info.IsDir() || checkbackslash : // i've refer to the "github.com\codeclysm\extract"
				if err := os.MkdirAll(location, info.Mode()|os.ModeDir|100); err != nil {
					log.Println("Mkdir Error")
				} else {
					log.Println("folder create", location)
				}
			case info.Mode() & os.ModeSymlink != 0 : //check if symlink
				f, err := ifile.Open()
				if err != nil {
					log.Println("sym err", err)
				} else {
					name, err := ioutil.ReadAll(f)
					if err != nil {
						log.Println("sym read", name)
					} else {
						if err = os.Symlink(location, string(name)); err != nil {
							log.Println("create sym", err)
						} else {
							log.Println("symlink create", location)
						}
					}
				}
			default: // if files
				f, err := ifile.Open() // get io.readcloser for zip inside files
				if err != nil {
					log.Println("openfile", err)
				} else {
					filebufs, err := ioutil.ReadAll(f) // io.readcloser to bytearray
					if err != nil {
						log.Println("fileread", err)
					} else {
						dirpath, _ := path.Split(location) // get directory path (check if dir only)
						_, err:= os.Stat(dirpath)
						if os.IsNotExist(err) { // if folder is not exist. make folder
							err = os.MkdirAll(dirpath, info.Mode()|os.ModeDir|100)
							if err != nil {
								log.Println("mkdir all", err)
							} else {
								log.Println("folder create", dirpath)
							}
						}
						err = ioutil.WriteFile(location, filebufs, ifile.Mode()) // save file
						if err != nil {
							log.Println("file create ", err)
						} else {
							err = os.Chtimes(location, modifytime, modifytime) // edit for save file to file access time // modify time
							if err != nil {
								log.Println("file modifytime ", err)
							} else {
								log.Println("file create", location)
							}
						}
					}
				}

			}
			if pathz == "" {
				continue
			}

		}
	}
}

func Archive(inFilePath []string, writer io.Writer, progress ProgressFunc) error {// i've refer to the "github.com/pierrre/archivefile"
	zipWriter := zip.NewWriter(writer)

	for _,data := range inFilePath{ // support for multiple files
		basePath := filepath.Dir(data)

		err := filepath.Walk(data, func(filePath string, fileInfo os.FileInfo, err error) error { // folder loop
			if err != nil || fileInfo.IsDir() {
				return err
			}

			relativeFilePath, err := filepath.Rel(basePath, filePath)
			if err != nil {
				return err
			}

			archivePath := path.Join(filepath.SplitList(relativeFilePath)...)

			if progress != nil {
				progress(archivePath)
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func() {
				_ = file.Close()
			}()

			zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (closer io.WriteCloser, e error) { // register compressor for reduce size
				return flate.NewWriter(out, flate.BestCompression)
			})
			zipFileWriter, err := zipWriter.CreateHeader(&zip.FileHeader{ // make compress file headers
				Name:               archivePath,
				Method: 			zip.Deflate,
				Modified:           fileInfo.ModTime(),
				NonUTF8: false,


			})
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFileWriter, file) // put file to zip
			return err
		})
		if err != nil {
			return err
		}
	}
	return zipWriter.Close()
}

func ArchiveFile(inFilePath []string, outFilePath string, progress ProgressFunc) error {
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outFile.Close()
	}()

	return Archive(inFilePath, outFile, progress)
}

type ProgressFunc func(archivePath string)
