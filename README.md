# Go_zip
i've refer "github.com/pierrre/archivefile" //  "github.com\codeclysm\extract"
두개의 library를 참조했습니다.

## 사용방법(how to used) :
#### unzipforbuffer(zipfile([]byte), savefolder(string))
##### example :
```
    file,err := ioutil.ReadFile("test.zip")
    if err != nil {
    // error processing
    } else {
      unzipforbuffer(file,"tmp")
      //tmp폴더에 저장하기 
      //saved tmp folder
    }
```
##### advenced example :
need resty ("github.com/go-resty/resty")
```
Requestz = resty.New()
r, e := Requestz.R().EnableTrace().Get(url)
if e != nil {
	unzipforbuffer(r.Body(),"tmp") // do not need saved zip file
}
```
    
#### ArchiveFile(file//dir lists, savefilename,progress)
 ##### example:
  ```
    ArchiveFile([]string{"test.txt","testfolder"},"test.zip",nil) 
  ```
  ##### if used progress :
  ```
    ArchiveFile(files,filename, func(archivePath string) {
		  log.Println(archivePath) // archivePath is compress file list
	  })
  ```
