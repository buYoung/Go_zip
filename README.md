# Go_zip


사용방법 :
how to used:
unzipforbuffer(zipfile([]byte), savefolder(string))
//
  example :
    file,err := ioutil.ReadFile("test.zip")
    if err != nil {
    // error processing
    } else {
      unzipforbuffer(file,"tmp")
      //tmp폴더에 저장하기 
      //saved tmp folder
    }
    
ArchiveFile(file or dir, savefilename,progress)
//
  example:
    ArchiveFile([]string{"test.txt","testfolder"},"test.zip",nil) 
  if used progress :
    ArchiveFile(files,filename, func(archivePath string) {
		  log.Println(archivePath) // archivePath is compress file list
	  })
